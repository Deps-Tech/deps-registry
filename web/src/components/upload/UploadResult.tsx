import { AlertCircle, CheckCircle2, ExternalLink } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import type { Manifest } from '@/lib/manifest';

interface UploadResultProps {
  result?: { manifest: Manifest; prUrl: string } | null;
  error?: string | null;
}

export function UploadResult({ result, error }: UploadResultProps) {
  if (error) {
    return (
      <Card className="border-destructive">
        <CardContent className="pt-6">
          <div className="flex items-start gap-3">
            <AlertCircle className="h-5 w-5 text-destructive mt-0.5" />
            <div>
              <h3 className="font-semibold text-destructive">Ошибка</h3>
              <p className="text-sm text-muted-foreground mt-1">{error}</p>
            </div>
          </div>
        </CardContent>
      </Card>
    );
  }

  if (result) {
    return (
      <Card className="border-green-500">
        <CardContent className="pt-6">
          <div className="flex items-start gap-3">
            <CheckCircle2 className="h-5 w-5 text-green-500 mt-0.5" />
            <div className="flex-1">
              <h3 className="font-semibold text-green-500">Успешно!</h3>
              <p className="text-sm text-muted-foreground mt-1">
                Pull request создан успешно
              </p>
              <Button
                variant="outline"
                size="sm"
                className="mt-3"
                onClick={() => window.open(result.prUrl, '_blank')}
              >
                <ExternalLink className="h-4 w-4 mr-2" />
                Открыть Pull Request
              </Button>

              <div className="mt-4 p-4 bg-muted rounded-lg">
                <h4 className="text-sm font-semibold mb-2">Сгенерированный манифест:</h4>
                <pre className="text-xs overflow-x-auto">
                  {JSON.stringify(result.manifest, null, 2)}
                </pre>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    );
  }

  return null;
}

