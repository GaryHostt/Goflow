'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { isAuthenticated } from '@/lib/api'

export default function Home() {
  const router = useRouter()

  useEffect(() => {
    if (isAuthenticated()) {
      router.push('/dashboard')
    } else {
      router.push('/login')
    }
  }, [router])

  return (
    <div className="flex min-h-screen items-center justify-center">
      <p>Loading...</p>
    </div>
  )
}

