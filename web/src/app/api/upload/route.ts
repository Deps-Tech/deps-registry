import { NextRequest, NextResponse } from 'next/server';
import { auth } from '@/lib/auth';
import { extractMetadata, analyzeLua } from '@/lib/parser';
import { createManifest } from '@/lib/manifest';
import { GitHubService, getDependencyVersions } from '@/lib/github';

const MAX_FILE_SIZE = 10 * 1024 * 1024;
const REPO_OWNER = process.env.GITHUB_REPO_OWNER || 'Deps-Tech';
const REPO_NAME = process.env.GITHUB_REPO_NAME || 'deps-registry';

export async function POST(request: NextRequest) {
  try {
    const session = await auth();
    
    if (!session?.user) {
      return NextResponse.json(
        { error: 'Unauthorized' },
        { status: 401 }
      );
    }

    const formData = await request.formData();
    const files = formData.getAll('files') as File[];
    const tags = formData.get('tags') as string;
    const packageType = (formData.get('type') as 'scripts' | 'deps') || 'scripts';
    const sourceUrl = formData.get('sourceUrl') as string;
    const metadataStr = formData.get('metadata') as string;

    if (!files || files.length === 0) {
      return NextResponse.json(
        { error: 'No files provided' },
        { status: 400 }
      );
    }

    for (const file of files) {
      if (file.size > MAX_FILE_SIZE) {
        return NextResponse.json(
          { error: `File ${file.name} exceeds 10MB limit` },
          { status: 400 }
        );
      }

      if (!file.name.endsWith('.lua')) {
        return NextResponse.json(
          { error: `File ${file.name} is not a .lua file` },
          { status: 400 }
        );
      }
    }

    const parsedFiles = await Promise.all(
      files.map(async (file) => ({
        name: file.name,
        content: await file.text(),
      }))
    );

    const accessToken = (session as any).accessToken;
    const depVersions = await getDependencyVersions(
      REPO_OWNER,
      REPO_NAME,
      accessToken
    );

    let metadata;
    let dependencies;
    let security;

    if (metadataStr) {
      const editedMetadata = JSON.parse(metadataStr);
      metadata = {
        id: editedMetadata.id,
        name: editedMetadata.name,
        version: editedMetadata.version,
        author: editedMetadata.author,
      };
      dependencies = editedMetadata.dependencies || [];
      security = editedMetadata.security || {
        usesNetwork: false,
        usesFFI: false,
        filePaths: [],
      };
    } else {
      const mainFile = parsedFiles[0];
      metadata = extractMetadata(mainFile.content, mainFile.name);
      const availableDeps = new Set(Object.keys(depVersions));
      const analysis = analyzeLua(mainFile.content, metadata.id, availableDeps);
      dependencies = analysis.dependencies;
      security = {
        usesNetwork: analysis.usesNetwork,
        usesFFI: analysis.usesFFI,
        filePaths: analysis.filePaths,
      };
    }

    const tagList = tags
      ? tags.split(',').map(t => t.trim()).filter(Boolean)
      : [];

    const octokit = new (await import('@octokit/rest')).Octokit({ auth: accessToken });
    const { data: user } = await octokit.users.getAuthenticated();
    const userLogin = user.login;

    const manifest = await createManifest(
      parsedFiles,
      {
        id: metadata.id,
        name: metadata.name,
        version: metadata.version,
        dependencies,
        security,
        tags: tagList,
        sourceUrl,
      },
      depVersions
    );

    (manifest.metadata as any).uploadedBy = userLogin;

    const githubService = new GitHubService({
      owner: REPO_OWNER,
      repo: REPO_NAME,
      token: accessToken,
    });
    
    const prUrl = await githubService.addPackage(
      userLogin,
      packageType,
      manifest,
      parsedFiles
    );

    return NextResponse.json({
      success: true,
      manifest,
      prUrl,
    });
  } catch (error: any) {
    console.error('Upload error:', error);
    return NextResponse.json(
      { error: error.message || 'Internal server error' },
      { status: 500 }
    );
  }
}

