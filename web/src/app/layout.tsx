import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { auth, signOut } from "@/lib/auth";
import { Button } from "@/components/ui/button";
import { ThemeProvider } from "@/components/theme-provider";
import { ThemeToggle } from "@/components/theme-toggle";
import Link from "next/link";
import { Package, LogOut } from "lucide-react";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Catalyst Registry",
  description: "Реестр пакетов для Arizona Catalyst",
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const session = await auth();

  return (
    <html lang="en" suppressHydrationWarning>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
          <header className="border-b">
            <div className="container mx-auto px-4 py-4 flex items-center justify-between">
              <Link href="/" className="flex items-center gap-2 font-bold text-xl">
                <Package className="h-6 w-6" />
                Catalyst Registry
              </Link>
              <nav className="flex items-center gap-4">
                {session?.user ? (
                  <>
                    <Link href="/browse">
                      <Button variant="ghost">Обзор</Button>
                    </Link>
                    <Link href="/add">
                      <Button variant="default">Загрузить</Button>
                    </Link>
                    <ThemeToggle />
                    <div className="flex items-center gap-2">
                      <span className="text-sm text-muted-foreground">
                        {session.user.name || session.user.email}
                      </span>
                      <form
                        action={async () => {
                          "use server";
                          await signOut();
                        }}
                      >
                        <Button type="submit" variant="ghost" size="icon">
                          <LogOut className="h-4 w-4" />
                        </Button>
                      </form>
                    </div>
                  </>
                ) : (
                  <ThemeToggle />
                )}
              </nav>
            </div>
          </header>
          <main>{children}</main>
          <footer className="border-t mt-16">
            <div className="container mx-auto px-4 py-8 text-center text-sm text-muted-foreground">
              <p>Catalyst Registry - реестр пакетов для Arizona Catalyst</p>
            </div>
          </footer>
        </ThemeProvider>
      </body>
    </html>
  );
}
