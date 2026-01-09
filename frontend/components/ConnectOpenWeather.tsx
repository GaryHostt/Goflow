'use client'

import { useState } from 'react'
import { credentials } from '@/lib/api'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

export default function ConnectOpenWeather() {
  const [apiKey, setApiKey] = useState('')
  const [loading, setLoading] = useState(false)
  const [success, setSuccess] = useState(false)
  const [error, setError] = useState('')

  const handleConnect = async () => {
    setLoading(true)
    setError('')
    setSuccess(false)

    try {
      await credentials.create('openweather', apiKey)
      setSuccess(true)
      setApiKey('')
    } catch (err) {
      setError('Failed to save OpenWeather API key')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>Connect OpenWeather</CardTitle>
        <CardDescription>
          Enter your OpenWeather API key. Get one free at openweathermap.org.
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {success && (
          <div className="rounded-md bg-green-50 p-3 text-sm text-green-800">
            OpenWeather API key saved successfully!
          </div>
        )}
        {error && (
          <div className="rounded-md bg-destructive/15 p-3 text-sm text-destructive">
            {error}
          </div>
        )}
        <Input
          placeholder="Your API Key"
          value={apiKey}
          onChange={(e) => setApiKey(e.target.value)}
          type="password"
        />
        <Button onClick={handleConnect} disabled={loading || !apiKey} className="w-full">
          {loading ? 'Saving...' : 'Save API Key'}
        </Button>
      </CardContent>
    </Card>
  )
}

