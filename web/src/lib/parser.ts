export interface Analysis {
  dependencies: string[];
  filePaths: string[];
  usesNetwork: boolean;
  usesFFI: boolean;
}

export interface LuaMetadata {
  id: string;
  name: string;
  version: string;
  author?: string;
}

const REQUIRE_REGEX = /require\s*\(?\s*["']([\w\.-]+)["']\s*\)?/g;
const NETWORK_REGEX = /(http\.request|socket\.tcp|socket\.connect)/;
const FFI_REGEX = /\brequire\s*\(?\s*["']ffi["']\s*\)?/;
const FILE_ACCESS_REGEX = /io\.open\s*\(\s*["']([^"']+)["']/g;

const BUILTIN_MODULES = new Set([
  'bit', 'math', 'string', 'table', 'os', 'io', 'debug', 
  'coroutine', 'package', 'utf8', 'bit32'
]);

export function extractMetadata(content: string, filename: string): LuaMetadata {
  const nameMatch = content.match(/script_name\s*\(\s*["'](.+?)["']\s*\)/);
  const versionMatch = content.match(/script_version\s*\(\s*["'](.+?)["']\s*\)/);
  const authorMatch = content.match(/script_author\s*\(\s*["'](.+?)["']\s*\)/);

  const name = nameMatch?.[1] || filename.replace(/\.lua$/, '');
  const version = versionMatch?.[1] || '1.0.0';
  const author = authorMatch?.[1];

  const id = name
    .toLowerCase()
    .replace(/[^a-z0-9-]+/g, '-')
    .replace(/^-+|-+$/g, '');

  return { id, name, version, author };
}

export function analyzeLua(
  content: string,
  excludeID: string,
  availableDeps: Set<string>
): Analysis {
  const dependencies: string[] = [];
  const filePaths: string[] = [];
  const depSet = new Set<string>();

  const requireMatches = content.matchAll(REQUIRE_REGEX);
  for (const match of requireMatches) {
    let dep = match[1].replace(/^lib\./, '');
    
    const parts = dep.split('.');
    const rootDep = parts[0];
    const lastPart = parts[parts.length - 1];
    
    if (BUILTIN_MODULES.has(rootDep)) {
      continue;
    }

    if (excludeID && rootDep.toLowerCase() === excludeID.toLowerCase()) {
      continue;
    }

    const depsToCheck = [
      dep.replace(/\./g, '-'),
      rootDep,
      lastPart,
    ];

    let found = false;
    for (const depToCheck of depsToCheck) {
      if (availableDeps.has(depToCheck.toLowerCase())) {
        depSet.add(depToCheck.toLowerCase());
        found = true;
        break;
      }
    }

    if (!found && parts.length > 1) {
      const combinations = [
        parts.join('-'),
        parts.slice(1).join('-'),
        lastPart,
      ];
      
      for (const combo of combinations) {
        if (availableDeps.has(combo.toLowerCase())) {
          depSet.add(combo.toLowerCase());
          break;
        }
      }
    }
  }

  const fileMatches = content.matchAll(FILE_ACCESS_REGEX);
  for (const match of fileMatches) {
    filePaths.push(match[1]);
  }

  return {
    dependencies: Array.from(depSet),
    filePaths,
    usesNetwork: NETWORK_REGEX.test(content),
    usesFFI: FFI_REGEX.test(content),
  };
}

export async function parseFiles(
  files: File[]
): Promise<{ content: string; filename: string }[]> {
  const results = await Promise.all(
    files.map(async (file) => ({
      content: await file.text(),
      filename: file.name,
    }))
  );
  return results;
}

