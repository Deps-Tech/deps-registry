import { NextResponse } from 'next/server';

const CDN_URL = 'https://cdn.depscian.tech/index.json';

interface PackageInfo {
  id: string;
  name?: string;
  versions: string[];
  latestVersion: string;
  tags?: string[];
  dependencies?: string[];
  security?: {
    networkAccess?: boolean;
    usesFFI?: boolean;
    fileAccess?: string[];
  };
  metadata?: {
    uploadedBy?: string;
    sourceUrl?: string;
  };
}

function transformCDNData(cdnData: any): PackageInfo[] {
  const packages: PackageInfo[] = [];
  
  for (const [id, pkgData] of Object.entries<any>(cdnData || {})) {
    const versions = Object.keys(pkgData.versions || {}).sort().reverse();
    const latestVersion = pkgData.latest || versions[0];
    const latestManifest = pkgData.versions?.[latestVersion]?.manifest;

    if (latestManifest) {
      packages.push({
        id,
        name: latestManifest.name || id,
        versions,
        latestVersion,
        tags: latestManifest.metadata?.tags || latestManifest.tags,
        dependencies: latestManifest.dependencies 
          ? Object.keys(latestManifest.dependencies)
          : undefined,
        security: latestManifest.security,
        metadata: latestManifest.metadata,
      });
    }
  }
  
  return packages;
}

export async function GET(request: Request) {
  try {
    const { searchParams } = new URL(request.url);
    const type = searchParams.get('type') as 'scripts' | 'deps' | undefined;

    const response = await fetch(CDN_URL, { next: { revalidate: 300 } });
    
    if (!response.ok) {
      throw new Error('Failed to fetch from CDN');
    }

    const data = await response.json();
    const scripts = transformCDNData(data.scripts);
    const deps = transformCDNData(data.dependencies);

    if (type === 'scripts') {
      return NextResponse.json({ packages: scripts });
    } else if (type === 'deps') {
      return NextResponse.json({ packages: deps });
    }

    return NextResponse.json({ scripts, deps });
  } catch (error: any) {
    console.error('Error fetching packages:', error);
    return NextResponse.json(
      { error: error.message || 'Failed to fetch packages' },
      { status: 500 }
    );
  }
}

export const revalidate = 300;
