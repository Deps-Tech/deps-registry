'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { FileDropzone } from './FileDropzone';
import { MetadataEditor } from './MetadataEditor';
import { UploadResult } from './UploadResult';
import type { Manifest } from '@/lib/manifest';

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

export function UploadForm() {
  const [files, setFiles] = useState<File[]>([]);
  const [tags, setTags] = useState('');
  const [sourceUrl, setSourceUrl] = useState('');
  const [packageType, setPackageType] = useState<'scripts' | 'deps'>('scripts');
  const [loading, setLoading] = useState(false);
  const [parsing, setParsing] = useState(false);
  const [result, setResult] = useState<{ manifest: Manifest; prUrl: string } | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [parsedData, setParsedData] = useState<ParsedData | null>(null);
  const [editableData, setEditableData] = useState<ParsedData | null>(null);
  const [existingVersions, setExistingVersions] = useState<string[]>([]);
  const [isUpdate, setIsUpdate] = useState(false);

  const handleFilesChange = async (acceptedFiles: File[]) => {
    setFiles(acceptedFiles);
    setError(null);
    setResult(null);
    setParsing(true);

    try {
      const formData = new FormData();
      acceptedFiles.forEach(file => formData.append('files', file));

      const response = await fetch('/api/parse', {
        method: 'POST',
        body: formData,
      });

      const data = await response.json();

      if (response.ok) {
        setParsedData(data);
        setEditableData(data);
        
        const pkgResponse = await fetch(`/api/packages?type=${packageType}`);
        if (pkgResponse.ok) {
          const pkgData = await pkgResponse.json();
          const existing = pkgData.packages?.find((p: any) => p.id === data.id);
          if (existing && existing.versions) {
            setExistingVersions(existing.versions);
            setIsUpdate(true);
          }
        }
      } else {
        throw new Error(data.error || 'Ошибка парсинга файлов');
      }
    } catch (err: any) {
      setError(err.message);
    } finally {
      setParsing(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (files.length === 0 || !editableData) {
      setError('Выберите хотя бы один .lua файл');
      return;
    }

    setLoading(true);
    setError(null);
    setResult(null);

    try {
      const formData = new FormData();
      files.forEach(file => formData.append('files', file));
      formData.append('tags', tags);
      formData.append('type', packageType);
      formData.append('metadata', JSON.stringify(editableData));
      if (sourceUrl) formData.append('sourceUrl', sourceUrl);

      const response = await fetch('/api/upload', {
        method: 'POST',
        body: formData,
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Ошибка загрузки');
      }

      setResult(data);
      setFiles([]);
      setTags('');
      setSourceUrl('');
      setParsedData(null);
      setEditableData(null);
      setExistingVersions([]);
      setIsUpdate(false);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Загрузить пакет</CardTitle>
          <CardDescription>
            Загрузите Lua скрипты или зависимости в реестр
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-6">
            <Tabs value={packageType} onValueChange={(v) => setPackageType(v as any)}>
              <TabsList className="grid w-full grid-cols-2">
                <TabsTrigger value="scripts">Скрипт</TabsTrigger>
                <TabsTrigger value="deps">Зависимость</TabsTrigger>
              </TabsList>
            </Tabs>

            <FileDropzone 
              files={files} 
              onFilesChange={handleFilesChange}
              parsing={parsing}
            />

            {editableData && (
              <MetadataEditor
                data={editableData}
                onChange={setEditableData}
                isUpdate={isUpdate}
                existingVersions={existingVersions}
              />
            )}

            <div className="space-y-2">
              <Label htmlFor="tags">Теги (через запятую)</Label>
              <Input
                id="tags"
                placeholder="helper, rp, automation"
                value={tags}
                onChange={(e) => setTags(e.target.value)}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="sourceUrl">Ссылка на источник (опционально)</Label>
              <Input
                id="sourceUrl"
                type="url"
                placeholder="https://github.com/..."
                value={sourceUrl}
                onChange={(e) => setSourceUrl(e.target.value)}
              />
            </div>

            <Button 
              type="submit" 
              className="w-full" 
              disabled={loading || files.length === 0 || !editableData}
            >
              {loading ? 'Загрузка...' : isUpdate ? 'Обновить пакет и создать PR' : 'Загрузить и создать PR'}
            </Button>
          </form>
        </CardContent>
      </Card>

      <UploadResult result={result} error={error} />
    </div>
  );
}

