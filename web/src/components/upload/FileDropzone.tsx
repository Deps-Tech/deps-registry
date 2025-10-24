import { useDropzone } from 'react-dropzone';
import { Upload, FileCode } from 'lucide-react';
import { Label } from '@/components/ui/label';

interface FileDropzoneProps {
  files: File[];
  onFilesChange: (files: File[]) => void;
  parsing?: boolean;
}

export function FileDropzone({ files, onFilesChange, parsing }: FileDropzoneProps) {
  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    accept: { 'text/x-lua': ['.lua'] },
    maxSize: 10 * 1024 * 1024,
    onDrop: onFilesChange,
  });

  return (
    <div>
      <Label>Файл</Label>
      <div
        {...getRootProps()}
        className={`border-2 border-dashed rounded-lg p-8 mt-2 text-center cursor-pointer transition-colors ${
          isDragActive ? 'border-primary bg-primary/5' : 'border-muted-foreground/25'
        }`}
      >
        <input {...getInputProps()} />
        <Upload className="mx-auto h-12 w-12 text-muted-foreground mb-4" />
        {files.length > 0 ? (
          <div className="space-y-2">
            {files.map((file, idx) => (
              <div key={idx} className="flex items-center justify-center gap-2 text-sm">
                <FileCode className="h-4 w-4" />
                <span>{file.name}</span>
                <span className="text-muted-foreground">
                  ({(file.size / 1024).toFixed(2)} KB)
                </span>
              </div>
            ))}
          </div>
        ) : (
          <div>
            <p className="text-sm text-muted-foreground">
              Переместите сюда ваш .lua файл, или нажмите для выбора
            </p>
            <p className="text-xs text-muted-foreground mt-2">
              Максимальный размер файла: 10 МБ
            </p>
          </div>
        )}
      </div>
      {parsing && (
        <p className="text-sm text-muted-foreground mt-2">Анализируем файлы...</p>
      )}
    </div>
  );
}

