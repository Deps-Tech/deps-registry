'use client';

import { useState, useEffect } from 'react';
import { Card, CardContent } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Package, Search, AlertCircle } from 'lucide-react';
import { PackageCard } from '@/components/browse/PackageCard';

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

export default function BrowsePage() {
  const [scripts, setScripts] = useState<PackageInfo[]>([]);
  const [deps, setDeps] = useState<PackageInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [activeTab, setActiveTab] = useState('all');

  useEffect(() => {
    fetchPackages();
  }, []);

  const fetchPackages = async () => {
    try {
      const response = await fetch('/api/packages');
      if (!response.ok) throw new Error('Ошибка загрузки пакетов');
      
      const data = await response.json();
      setScripts(data.scripts || []);
      setDeps(data.deps || []);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const filteredScripts = scripts.filter(pkg =>
    pkg.id.toLowerCase().includes(searchQuery.toLowerCase()) ||
    pkg.name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
    pkg.tags?.some(tag => tag.toLowerCase().includes(searchQuery.toLowerCase()))
  );

  const filteredDeps = deps.filter(pkg =>
    pkg.id.toLowerCase().includes(searchQuery.toLowerCase()) ||
    pkg.name?.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const allPackages = [...filteredScripts, ...filteredDeps];

  if (loading) {
    return (
      <div className="container mx-auto py-8">
        <div className="flex items-center justify-center h-64">
          <p className="text-muted-foreground">Загрузка пакетов...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mx-auto py-8">
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
      </div>
    );
  }

  return (
    <div className="container mx-auto py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-2">Обзор пакетов</h1>
        <p className="text-muted-foreground mb-6">
          {scripts.length} скриптов и {deps.length} библиотек
        </p>

        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Поиск пакетов..."
            className="pl-10"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-3 mb-6">
          <TabsTrigger value="all">
            Все ({allPackages.length})
          </TabsTrigger>
          <TabsTrigger value="scripts">
            Скрипты ({filteredScripts.length})
          </TabsTrigger>
          <TabsTrigger value="deps">
            Библиотеки ({filteredDeps.length})
          </TabsTrigger>
        </TabsList>

        <TabsContent value="all" className="space-y-4">
          {allPackages.length === 0 ? (
            <Card>
              <CardContent className="pt-6 text-center">
                <Package className="h-12 w-12 mx-auto mb-4 text-muted-foreground" />
                <p className="text-muted-foreground">Пакеты не найдены</p>
              </CardContent>
            </Card>
          ) : (
            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
              {filteredScripts.map(pkg => <PackageCard key={`s-${pkg.id}`} pkg={pkg} type="script" />)}
              {filteredDeps.map(pkg => <PackageCard key={`d-${pkg.id}`} pkg={pkg} type="dep" />)}
            </div>
          )}
        </TabsContent>

        <TabsContent value="scripts" className="space-y-4">
          {filteredScripts.length === 0 ? (
            <Card>
              <CardContent className="pt-6 text-center">
                <Package className="h-12 w-12 mx-auto mb-4 text-muted-foreground" />
                <p className="text-muted-foreground">Скрипты не найдены</p>
              </CardContent>
            </Card>
          ) : (
            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
              {filteredScripts.map(pkg => <PackageCard key={pkg.id} pkg={pkg} type="script" />)}
            </div>
          )}
        </TabsContent>

        <TabsContent value="deps" className="space-y-4">
          {filteredDeps.length === 0 ? (
            <Card>
              <CardContent className="pt-6 text-center">
                <Package className="h-12 w-12 mx-auto mb-4 text-muted-foreground" />
                <p className="text-muted-foreground">Библиотеки не найдены</p>
              </CardContent>
            </Card>
          ) : (
            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
              {filteredDeps.map(pkg => <PackageCard key={pkg.id} pkg={pkg} type="dep" />)}
            </div>
          )}
        </TabsContent>
      </Tabs>
    </div>
  );
}
