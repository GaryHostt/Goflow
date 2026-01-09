'use client'

import { useState } from 'react'
import { credentials } from '@/lib/api'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

export default function ConnectDiscord() {
  const [webhookUrl, setWebhookUrl] = useState('')
  const [loading, setLoading] = useState(false)
  const [success, setSuccess] = useState(false)
  const [error, setError] = useState('')

  const handleConnect = async () => {
    setLoading(true)
    setError('')
    setSuccess(false)

    try {
      await credentials.create('discord', webhookUrl)
      setSuccess(true)
      setWebhookUrl('')
    } catch (err) {
      setError('Failed to save Discord connection')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>Connect Discord</CardTitle>
        <CardDescription>
          Paste your Discord Webhook URL from your server settings.
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {success && (
          <div className="rounded-md bg-green-50 p-3 text-sm text-green-800">
            Discord connected successfully!
          </div>
        )}
        {error && (
          <div className="rounded-md bg-destructive/15 p-3 text-sm text-destructive">
            {error}
          </div>
        )}
        <Input
          placeholder="https://discord.com/api/webhooks/..."
          value={webhookUrl}
          onChange={(e) => setWebhookUrl(e.target.value)}
        />
        <Button onClick={handleConnect} disabled={loading || !webhookUrl} className="w-full">
          {loading ? 'Saving...' : 'Save Connection'}
        </Button>
      </CardContent>
    </Card>
  )
}

