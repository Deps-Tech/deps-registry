import { auth } from '@/lib/auth';
import { redirect } from 'next/navigation';
import { UploadForm } from '@/components/upload/UploadForm';

export default async function AddPage() {
  const session = await auth();

  if (!session?.user) {
    redirect('/');
  }

  return (
    <div className="container mx-auto py-8 max-w-4xl">
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-2">Добавить пакет</h1>
        <p className="text-muted-foreground">
          Загрузите Lua скрипты или зависимости в Catalyst Registry
        </p>
      </div>
      <UploadForm />
    </div>
  );
}

