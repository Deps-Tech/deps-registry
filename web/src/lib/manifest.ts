import { z } from 'zod';

export const FileInfoSchema = z.object({
  sha256: z.string(),
  size: z.number(),
});

export const SecuritySchema = z.object({
  networkAccess: z.boolean().optional(),
  fileAccess: z.array(z.string()).optional(),
  usesFFI: z.boolean().optional(),
});

export const MetadataSchema = z.object({
  sourceUrl: z.string().optional(),
  tags: z.array(z.string()).optional(),
  deprecated: z.boolean().optional(),
});

export const ManifestSchema = z.object({
  manifestVersion: z.string(),
  id: z.string(),
  name: z.string().optional(),
  version: z.string(),
  files: z.record(FileInfoSchema),
  dependencies: z.record(z.string()).optional(),
  security: SecuritySchema.optional(),
  metadata: MetadataSchema.optional(),
});

export type FileInfo = z.infer<typeof FileInfoSchema>;
export type Security = z.infer<typeof SecuritySchema>;
export type Metadata = z.infer<typeof MetadataSchema>;
export type Manifest = z.infer<typeof ManifestSchema>;

export async function generateSHA256(content: string): Promise<string> {
  const encoder = new TextEncoder();
  const data = encoder.encode(content);
  const hashBuffer = await crypto.subtle.digest('SHA-256', data);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  return hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
}

export async function createManifest(
  files: { name: string; content: string }[],
  metadata: {
    id: string;
    name: string;
    version: string;
    dependencies: string[];
    security: {
      usesNetwork: boolean;
      usesFFI: boolean;
      filePaths: string[];
    };
    tags: string[];
    sourceUrl?: string;
  },
  depVersions: Record<string, string>
): Promise<Manifest> {
  const fileMap: Record<string, FileInfo> = {};

  for (const file of files) {
    const sha256 = await generateSHA256(file.content);
    const size = new TextEncoder().encode(file.content).length;
    fileMap[file.name] = { sha256, size };
  }

  const deps: Record<string, string> = {};
  for (const dep of metadata.dependencies) {
    if (depVersions[dep]) {
      deps[dep] = depVersions[dep];
    }
  }

  return {
    manifestVersion: '1.0',
    id: metadata.id,
    name: metadata.name,
    version: metadata.version,
    files: fileMap,
    dependencies: Object.keys(deps).length > 0 ? deps : undefined,
    security: {
      networkAccess: metadata.security.usesNetwork || undefined,
      fileAccess: metadata.security.filePaths.length > 0 ? metadata.security.filePaths : undefined,
      usesFFI: metadata.security.usesFFI || undefined,
    },
    metadata: {
      tags: metadata.tags.length > 0 ? metadata.tags : undefined,
      sourceUrl: metadata.sourceUrl,
    },
  };
}

