import { auth, signIn } from '@/lib/auth';
import { Button } from '@/components/ui/button';
import { Package, Upload, Github } from 'lucide-react';
import Link from 'next/link';

export default async function Home() {
  const session = await auth();

  return (
    <div className="h-screen flex items-center justify-center container mx-auto">
      <div className="text-center mb-16">
        <h1 className="text-5xl font-bold mb-4">Catalyst Registry</h1>
        <p className="text-xl text-muted-foreground mb-8">
          Реестр пакетов для Arizona Catalyst от сообщества
        </p>
        {session?.user ? (
          <div className="flex gap-4 justify-center">
            <Button asChild size="lg">
              <Link href="/add">
                <Upload className="mr-2 h-5 w-5" />
                Загрузить пакет
              </Link>
            </Button>
            <Button asChild variant="outline" size="lg">
              <Link href="/browse">
                <Package className="mr-2 h-5 w-5" />
                Обзор пакетов
              </Link>
            </Button>
          </div>
        ) : (
          <form
            action={async () => {
              'use server';
              await signIn('github', { redirectTo: '/add' });
            }}
          >
            <Button type="submit" size="lg">
              <Github className="mr-2 h-5 w-5" />
              Войти через GitHub
            </Button>
          </form>
        )}
      </div>
    </div>
  );
}
