import './globals.css'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'iPaaS - Integration Platform',
  description: 'Simple Integration Platform as a Service',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  )
}

