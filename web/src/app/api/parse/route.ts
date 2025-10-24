import { NextRequest, NextResponse } from 'next/server';
import { extractMetadata, analyzeLua } from '@/lib/parser';
import { getDependencyVersions } from '@/lib/github';

const REPO_OWNER = process.env.GITHUB_REPO_OWNER || 'Deps-Tech';
const REPO_NAME = process.env.GITHUB_REPO_NAME || 'deps-registry';

export async function POST(request: NextRequest) {
  try {
    const formData = await request.formData();
    const files = formData.getAll('files') as File[];

    if (!files || files.length === 0) {
      return NextResponse.json(
        { error: 'No files provided' },
        { status: 400 }
      );
    }

    const parsedFiles = await Promise.all(
      files.map(async (file) => ({
        name: file.name,
        content: await file.text(),
      }))
    );

    const mainFile = parsedFiles[0];
    const metadata = extractMetadata(mainFile.content, mainFile.name);

    const depVersions = await getDependencyVersions(
      REPO_OWNER,
      REPO_NAME,
      ''
    );

    const availableDeps = new Set(Object.keys(depVersions));
    const analysis = analyzeLua(mainFile.content, metadata.id, availableDeps);

    return NextResponse.json({
      id: metadata.id,
      name: metadata.name,
      version: metadata.version,
      author: metadata.author,
      dependencies: analysis.dependencies,
      security: {
        usesNetwork: analysis.usesNetwork,
        usesFFI: analysis.usesFFI,
        filePaths: analysis.filePaths,
      },
    });
  } catch (error: any) {
    console.error('Parse error:', error);
    return NextResponse.json(
      { error: error.message || 'Internal server error' },
      { status: 500 }
    );
  }
}

