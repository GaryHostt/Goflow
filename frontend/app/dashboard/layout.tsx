'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import ProtectedLayout from '@/components/ProtectedRoute'
import { clearToken } from '@/lib/api'
import { Button } from '@/components/ui/button'

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const router = useRouter()
  const [activeTab, setActiveTab] = useState('dashboard')

  useEffect(() => {
    const path = window.location.pathname
    if (path.includes('workflows')) setActiveTab('workflows')
    else if (path.includes('connections')) setActiveTab('connections')
    else if (path.includes('logs')) setActiveTab('logs')
    else setActiveTab('dashboard')
  }, [])

  const handleLogout = () => {
    clearToken()
    router.push('/login')
  }

  return (
    <ProtectedLayout>
      <div className="flex min-h-screen bg-gray-50">
        {/* Sidebar */}
        <aside className="w-64 bg-white border-r">
          <div className="p-6">
            <h1 className="text-2xl font-bold text-primary">iPaaS</h1>
            <p className="text-sm text-muted-foreground">Integration Platform</p>
          </div>
          
          <nav className="px-4 space-y-2">
            <Link href="/dashboard">
              <div className={`px-4 py-2 rounded-md transition-colors cursor-pointer ${
                activeTab === 'dashboard' 
                  ? 'bg-primary text-white' 
                  : 'hover:bg-gray-100'
              }`}>
                Dashboard
              </div>
            </Link>
            
            <Link href="/dashboard/workflows">
              <div className={`px-4 py-2 rounded-md transition-colors cursor-pointer ${
                activeTab === 'workflows' 
                  ? 'bg-primary text-white' 
                  : 'hover:bg-gray-100'
              }`}>
                Workflows
              </div>
            </Link>
            
            <Link href="/dashboard/connections">
              <div className={`px-4 py-2 rounded-md transition-colors cursor-pointer ${
                activeTab === 'connections' 
                  ? 'bg-primary text-white' 
                  : 'hover:bg-gray-100'
              }`}>
                Connections
              </div>
            </Link>
            
            <Link href="/dashboard/logs">
              <div className={`px-4 py-2 rounded-md transition-colors cursor-pointer ${
                activeTab === 'logs' 
                  ? 'bg-primary text-white' 
                  : 'hover:bg-gray-100'
              }`}>
                Logs
              </div>
            </Link>
          </nav>

          <div className="absolute bottom-0 w-64 p-4 border-t">
            <Button onClick={handleLogout} variant="outline" className="w-full">
              Logout
            </Button>
          </div>
        </aside>

        {/* Main content */}
        <main className="flex-1 p-8">
          {children}
        </main>
      </div>
    </ProtectedLayout>
  )
}

