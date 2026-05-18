import "./globals.css";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "RuntimeWall Dashboard",
  description: "Security-first runtime and governance for autonomous AI agents",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className="h-full">
      <body className="h-full bg-bg text-slate-200 font-mono antialiased">
        {children}
      </body>
    </html>
  );
}
