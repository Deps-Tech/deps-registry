import { Package, ExternalLink, User, Download } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';

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

interface PackageCardProps {
  pkg: PackageInfo;
  type: 'script' | 'dep';
}

export function PackageCard({ pkg, type }: PackageCardProps) {
  return (
    <Card className="hover:shadow-lg transition-shadow">
      <CardHeader>
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-1">
              <Package className="h-5 w-5 text-primary" />
              <CardTitle className="text-lg">{pkg.name || pkg.id}</CardTitle>
              <Badge variant="outline">{type === 'script' ? 'скрипт' : 'библиотека'}</Badge>
            </div>
            {pkg.id !== pkg.name && (
              <p className="text-xs text-muted-foreground mb-2">{pkg.id}</p>
            )}
            <CardDescription>
              Версия: {pkg.latestVersion}
              {pkg.versions.length > 1 && ` (+${pkg.versions.length - 1})`}
            </CardDescription>
          </div>
        </div>
      </CardHeader>
      <CardContent className="space-y-3">
        {pkg.tags && pkg.tags.length > 0 && (
          <div className="flex flex-wrap gap-1">
            {pkg.tags.map((tag, idx) => (
              <Badge key={idx} variant="secondary" className="text-xs">
                {tag}
              </Badge>
            ))}
          </div>
        )}

        {pkg.dependencies && pkg.dependencies.length > 0 && (
          <div>
            <p className="text-xs font-semibold mb-1">Зависимости:</p>
            <div className="flex flex-wrap gap-1">
              {pkg.dependencies.slice(0, 5).map((dep, idx) => (
                <Badge key={idx} variant="outline" className="text-xs">
                  {dep}
                </Badge>
              ))}
              {pkg.dependencies.length > 5 && (
                <Badge variant="outline" className="text-xs">
                  +{pkg.dependencies.length - 5}
                </Badge>
              )}
            </div>
          </div>
        )}

        {pkg.security && (
          <div className="flex flex-wrap gap-1">
            {pkg.security.networkAccess && (
              <Badge variant="destructive" className="text-xs">Сеть</Badge>
            )}
            {pkg.security.usesFFI && (
              <Badge variant="destructive" className="text-xs">FFI</Badge>
            )}
            {pkg.security.fileAccess && pkg.security.fileAccess.length > 0 && (
              <Badge variant="destructive" className="text-xs">Файлы</Badge>
            )}
          </div>
        )}

        <div className="flex items-center justify-between pt-2">
          {pkg.metadata?.uploadedBy && (
            <div className="flex items-center gap-1 text-xs text-muted-foreground">
              <User className="h-3 w-3" />
              <span>{pkg.metadata.uploadedBy}</span>
            </div>
          )}
          <div className="flex gap-2">
            {pkg.metadata?.sourceUrl && (
              <Button
                variant="ghost"
                size="sm"
                onClick={() => window.open(pkg.metadata?.sourceUrl, '_blank')}
              >
                <ExternalLink className="h-4 w-4" />
              </Button>
            )}
            <Button
              variant="outline"
              size="sm"
              onClick={() => {
                navigator.clipboard.writeText(`require '${pkg.id}'`);
              }}
            >
              <Download className="h-4 w-4 mr-1" />
              Копировать
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

