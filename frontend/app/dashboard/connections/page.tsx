'use client'

import ConnectSlack from '@/components/ConnectSlack'
import ConnectDiscord from '@/components/ConnectDiscord'
import ConnectOpenWeather from '@/components/ConnectOpenWeather'

export default function ConnectionsPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Connections</h1>
        <p className="text-muted-foreground">
          Connect your external services to enable integrations
        </p>
      </div>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        <ConnectSlack />
        <ConnectDiscord />
        <ConnectOpenWeather />
      </div>
    </div>
  )
}

