'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import Image from 'next/image'
import { auth, setToken, api } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { XCircle, Zap } from 'lucide-react'

const DEV_MODE = process.env.NEXT_PUBLIC_DEV_MODE === 'true'

export default function LoginPage() {
  const router = useRouter()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const response = await auth.login(email, password)
      
      if (response.token || response.data?.token) {
        const token = response.token || response.data.token
        setToken(token)
        router.push('/dashboard')
      } else {
        setError('Login failed. Please check your credentials.')
      }
    } catch (err: any) {
      // Display the specific error message
      setError(err.message || 'An error occurred. Please try again.')
      console.error('Login error:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleDevLogin = async () => {
    setError('')
    setLoading(true)

    try {
      const response = await api.post('/auth/dev-login')
      
      if (response.token) {
        setToken(response.token)
        router.push('/dashboard')
      } else {
        setError('Dev login failed.')
      }
    } catch (err: any) {
      setError(err.message || 'Dev login failed.')
      console.error('Dev login error:', err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50">
      <Card className="w-[400px]">
        <CardHeader className="text-center">
          <div className="flex justify-center mb-4">
            <Image
              src="/goflow-logo.png"
              alt="GoFlow Logo"
              width={120}
              height={120}
              priority
            />
          </div>
          <CardTitle>Welcome to GoFlow</CardTitle>
          <CardDescription>Sign in to your integration platform</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            {error && (
              <Alert variant="destructive">
                <XCircle className="h-4 w-4" />
                <AlertTitle>Error</AlertTitle>
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}
            
            <div className="space-y-2">
              <label htmlFor="email" className="text-sm font-medium">
                Email
              </label>
              <Input
                id="email"
                type="email"
                placeholder="you@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                autoComplete="email"
                pattern="[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$"
                title="Please enter a valid email address"
              />
            </div>

            <div className="space-y-2">
              <label htmlFor="password" className="text-sm font-medium">
                Password
              </label>
              <Input
                id="password"
                type="password"
                placeholder="••••••••"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                minLength={6}
                maxLength={128}
                autoComplete="current-password"
                title="Password must be at least 6 characters"
              />
            </div>

            <Button type="submit" className="w-full" disabled={loading}>
              {loading ? 'Signing in...' : 'Sign In'}
            </Button>

            {DEV_MODE && (
              <div className="relative">
                <div className="absolute inset-0 flex items-center">
                  <span className="w-full border-t" />
                </div>
                <div className="relative flex justify-center text-xs uppercase">
                  <span className="bg-white px-2 text-muted-foreground">
                    Development Mode
                  </span>
                </div>
              </div>
            )}

            {DEV_MODE && (
              <Button
                type="button"
                variant="outline"
                className="w-full border-orange-500 text-orange-600 hover:bg-orange-50"
                onClick={handleDevLogin}
                disabled={loading}
              >
                <Zap className="mr-2 h-4 w-4" />
                {loading ? 'Logging in...' : 'Skip Login - Dev Mode'}
              </Button>
            )}

            <p className="text-center text-sm text-muted-foreground">
              Don't have an account?{' '}
              <Link href="/register" className="text-primary hover:underline">
                Create account
              </Link>
            </p>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}

