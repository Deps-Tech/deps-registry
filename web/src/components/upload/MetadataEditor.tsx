import { Edit2 } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';

interface ParsedData {
  id: string;
  name: string;
  version: string;
  author?: string;
  dependencies: string[];
  security: {
    usesNetwork: boolean;
    usesFFI: boolean;
    filePaths: string[];
  };
}

interface MetadataEditorProps {
  data: ParsedData;
  onChange: (data: ParsedData) => void;
  isUpdate?: boolean;
  existingVersions?: string[];
}

export function MetadataEditor({ data, onChange, isUpdate, existingVersions }: MetadataEditorProps) {
  return (
    <div className="space-y-4 p-4 border rounded-lg bg-muted/50">
      <div className="flex items-center gap-2 mb-2">
        <Edit2 className="h-4 w-4" />
        <span className="font-semibold">Метаданные пакета</span>
        {isUpdate && <Badge variant="secondary">Обновление</Badge>}
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label htmlFor="edit-id">ID</Label>
          <Input
            id="edit-id"
            value={data.id}
            onChange={(e) => onChange({ ...data, id: e.target.value })}
          />
        </div>

        <div className="space-y-2">
          <Label htmlFor="edit-version">Версия</Label>
          <Input
            id="edit-version"
            value={data.version}
            onChange={(e) => onChange({ ...data, version: e.target.value })}
          />
          {existingVersions && existingVersions.length > 0 && (
            <p className="text-xs text-muted-foreground">
              Существующие: {existingVersions.join(', ')}
            </p>
          )}
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="edit-name">Название</Label>
        <Input
          id="edit-name"
          value={data.name}
          onChange={(e) => onChange({ ...data, name: e.target.value })}
        />
      </div>

      {data.author && (
        <div className="space-y-2">
          <Label htmlFor="edit-author">Автор</Label>
          <Input
            id="edit-author"
            value={data.author}
            onChange={(e) => onChange({ ...data, author: e.target.value })}
          />
        </div>
      )}

      {data.dependencies.length > 0 && (
        <div className="space-y-2">
          <Label>Зависимости</Label>
          <div className="flex flex-wrap gap-2">
            {data.dependencies.map((dep, idx) => (
              <Badge key={idx} variant="outline">{dep}</Badge>
            ))}
          </div>
        </div>
      )}

      {(data.security.usesNetwork || data.security.usesFFI || data.security.filePaths.length > 0) && (
        <div className="space-y-2">
          <Label>Предупреждения безопасности</Label>
          <div className="flex flex-wrap gap-2">
            {data.security.usesNetwork && <Badge variant="destructive">Сеть</Badge>}
            {data.security.usesFFI && <Badge variant="destructive">FFI</Badge>}
            {data.security.filePaths.length > 0 && <Badge variant="destructive">Файлы</Badge>}
          </div>
        </div>
      )}
    </div>
  );
}

