import { Octokit } from '@octokit/rest';
import type { Manifest } from './manifest';

export interface GitHubConfig {
  owner: string;
  repo: string;
  token: string;
}

export class GitHubService {
  private octokit: Octokit;
  private owner: string;
  private repo: string;

  constructor(config: GitHubConfig) {
    this.octokit = new Octokit({ auth: config.token });
    this.owner = config.owner;
    this.repo = config.repo;
  }

  async ensureFork(userLogin: string): Promise<{ owner: string; repo: string }> {
    if (userLogin.toLowerCase() === this.owner.toLowerCase()) {
      return { owner: this.owner, repo: this.repo };
    }

    try {
      await this.octokit.repos.get({
        owner: userLogin,
        repo: this.repo,
      });
      return { owner: userLogin, repo: this.repo };
    } catch (error: any) {
      if (error.status === 404) {
        try {
          await this.octokit.repos.get({
            owner: this.owner,
            repo: this.repo,
          });
        } catch (repoError: any) {
          if (repoError.status === 404) {
            throw new Error(`Repository ${this.owner}/${this.repo} not found or not accessible`);
          }
          throw repoError;
        }

        try {
          await this.octokit.repos.createFork({
            owner: this.owner,
            repo: this.repo,
          });
          
          await new Promise(resolve => setTimeout(resolve, 5000));
          return { owner: userLogin, repo: this.repo };
        } catch (forkError: any) {
          throw new Error(`Failed to create fork: ${forkError.message}`);
        }
      }
      throw error;
    }
  }

  async getLatestCommitSha(branch: string = 'main'): Promise<string> {
    const { data } = await this.octokit.repos.getBranch({
      owner: this.owner,
      repo: this.repo,
      branch,
    });
    return data.commit.sha;
  }

  async createBranch(
    userLogin: string,
    branchName: string,
    baseSha: string
  ): Promise<void> {
    await this.octokit.git.createRef({
      owner: userLogin,
      repo: this.repo,
      ref: `refs/heads/${branchName}`,
      sha: baseSha,
    });
  }

  async createOrUpdateFile(
    userLogin: string,
    branch: string,
    path: string,
    content: string,
    message: string
  ): Promise<void> {
    const encodedContent = Buffer.from(content).toString('base64');

    try {
      const { data: existingFile } = await this.octokit.repos.getContent({
        owner: userLogin,
        repo: this.repo,
        path,
        ref: branch,
      });

      if ('sha' in existingFile) {
        await this.octokit.repos.createOrUpdateFileContents({
          owner: userLogin,
          repo: this.repo,
          path,
          message,
          content: encodedContent,
          branch,
          sha: existingFile.sha,
        });
      }
    } catch (error: any) {
      if (error.status === 404) {
        await this.octokit.repos.createOrUpdateFileContents({
          owner: userLogin,
          repo: this.repo,
          path,
          message,
          content: encodedContent,
          branch,
        });
      } else {
        throw error;
      }
    }
  }

  async createPullRequest(
    userLogin: string,
    branch: string,
    title: string,
    body: string
  ): Promise<string> {
    const { data: pr } = await this.octokit.pulls.create({
      owner: this.owner,
      repo: this.repo,
      title,
      body,
      head: `${userLogin}:${branch}`,
      base: 'main',
    });

    return pr.html_url;
  }

  async addPackage(
    userLogin: string,
    packageType: 'scripts' | 'deps',
    manifest: Manifest,
    files: { name: string; content: string }[]
  ): Promise<string> {
    const forkInfo = await this.ensureFork(userLogin);
    
    const baseSha = await this.getLatestCommitSha();
    const branchName = `add-${packageType}-${manifest.id}-${Date.now()}`;
    
    await this.createBranch(forkInfo.owner, branchName, baseSha);

    const basePath = `${packageType}/${manifest.id}/${manifest.version}`;

    for (const file of files) {
      await this.createOrUpdateFile(
        forkInfo.owner,
        branchName,
        `${basePath}/${file.name}`,
        file.content,
        `feat(${packageType}): add ${manifest.id} ${file.name}`
      );
    }

    await this.createOrUpdateFile(
      forkInfo.owner,
      branchName,
      `${basePath}/dep.json`,
      JSON.stringify(manifest, null, 2),
      `feat(${packageType}): add ${manifest.id} manifest`
    );

    const prBody = `
## Add ${manifest.name || manifest.id} v${manifest.version}

**Package Type:** ${packageType}

**Dependencies:** ${manifest.dependencies ? Object.keys(manifest.dependencies).join(', ') : 'None'}

**Security:**
- Network Access: ${manifest.security?.networkAccess ? '⚠️ Yes' : '✅ No'}
- FFI Usage: ${manifest.security?.usesFFI ? '⚠️ Yes' : '✅ No'}
- File Access: ${manifest.security?.fileAccess?.length ? `⚠️ ${manifest.security.fileAccess.length} file(s)` : '✅ No'}

**Files:**
${files.map(f => `- ${f.name}`).join('\n')}
    `.trim();

    return await this.createPullRequest(
      forkInfo.owner,
      branchName,
      `feat(${packageType}): add ${manifest.id} v${manifest.version}`,
      prBody
    );
  }
}

export async function getDependencyVersions(
  owner: string,
  repo: string,
  token: string
): Promise<Record<string, string>> {
  const octokit = new Octokit({ auth: token });
  const versions: Record<string, string> = {};

  try {
    const { data: deps } = await octokit.repos.getContent({
      owner,
      repo,
      path: 'deps',
    });

    if (Array.isArray(deps)) {
      for (const dep of deps) {
        if (dep.type === 'dir') {
          try {
            const { data: versionDirs } = await octokit.repos.getContent({
              owner,
              repo,
              path: `deps/${dep.name}`,
            });

            if (Array.isArray(versionDirs)) {
              const sorted = versionDirs
                .filter(v => v.type === 'dir')
                .map(v => v.name)
                .sort()
                .reverse();
              
              if (sorted.length > 0) {
                versions[dep.name] = sorted[0];
              }
            }
          } catch (err) {
            console.error(`Error fetching versions for ${dep.name}:`, err);
          }
        }
      }
    }
  } catch (error: any) {
    if (error.status === 404) {
      console.warn('deps directory not found, using empty dependency versions');
    } else {
      console.error('Error fetching dependency versions:', error);
    }
  }

  return versions;
}

