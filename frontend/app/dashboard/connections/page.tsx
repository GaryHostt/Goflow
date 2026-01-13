'use client'

import { useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  MessageSquare, 
  MessageCircle, 
  Cloud, 
  Phone, 
  Newspaper, 
  Cat, 
  ShoppingBag,
  Zap,
  Star,
  Briefcase,
  Gamepad2,
  Hash,
  Rocket,
  Globe,
  Dog,
  CheckCircle2,
  AlertCircle
} from 'lucide-react'
import { credentials } from '@/lib/api'

type ConnectorConfig = {
  id: string
  name: string
  description: string
  icon: any
  fields: Array<{
    key: string
    label: string
    type: string
    placeholder: string
    required?: boolean
  }>
  category: 'messaging' | 'data' | 'fun' | 'enterprise'
  color: string
}

const CONNECTORS: ConnectorConfig[] = [
  // Messaging
  {
    id: 'slack',
    name: 'Slack',
    description: 'Send messages to Slack channels via webhooks',
    icon: MessageSquare,
    fields: [
      { key: 'webhook_url', label: 'Webhook URL', type: 'url', placeholder: 'https://hooks.slack.com/services/...', required: true }
    ],
    category: 'messaging',
    color: 'text-purple-600'
  },
  {
    id: 'discord',
    name: 'Discord',
    description: 'Send messages to Discord channels via webhooks',
    icon: MessageCircle,
    fields: [
      { key: 'webhook_url', label: 'Webhook URL', type: 'url', placeholder: 'https://discord.com/api/webhooks/...', required: true }
    ],
    category: 'messaging',
    color: 'text-indigo-600'
  },
  {
    id: 'twilio',
    name: 'Twilio',
    description: 'Send SMS messages via Twilio',
    icon: Phone,
    fields: [
      { key: 'account_sid', label: 'Account SID', type: 'text', placeholder: 'ACxxxxxxxxxxxxx', required: true },
      { key: 'auth_token', label: 'Auth Token', type: 'password', placeholder: 'Your auth token', required: true },
      { key: 'from_number', label: 'From Number', type: 'tel', placeholder: '+1234567890', required: true }
    ],
    category: 'messaging',
    color: 'text-red-600'
  },
  
  // Data APIs
  {
    id: 'openweather',
    name: 'OpenWeather',
    description: 'Get weather data and forecasts',
    icon: Cloud,
    fields: [
      { key: 'api_key', label: 'API Key', type: 'password', placeholder: 'Your OpenWeather API key', required: true }
    ],
    category: 'data',
    color: 'text-blue-600'
  },
  {
    id: 'newsapi',
    name: 'NewsAPI',
    description: 'Fetch latest news articles',
    icon: Newspaper,
    fields: [
      { key: 'api_key', label: 'API Key', type: 'password', placeholder: 'Your NewsAPI key', required: true }
    ],
    category: 'data',
    color: 'text-orange-600'
  },
  {
    id: 'nasa',
    name: 'NASA API',
    description: 'Access NASA data and imagery (use DEMO_KEY for testing)',
    icon: Rocket,
    fields: [
      { key: 'api_key', label: 'API Key', type: 'text', placeholder: 'DEMO_KEY or your API key', required: false }
    ],
    category: 'data',
    color: 'text-blue-700'
  },
  {
    id: 'restcountries',
    name: 'REST Countries',
    description: 'Get information about countries (no API key needed)',
    icon: Globe,
    fields: [],
    category: 'data',
    color: 'text-green-600'
  },
  
  // Fun APIs
  {
    id: 'catapi',
    name: 'The Cat API',
    description: 'Random cat images and facts',
    icon: Cat,
    fields: [],
    category: 'fun',
    color: 'text-pink-600'
  },
  {
    id: 'dogapi',
    name: 'Dog CEO API',
    description: 'Random dog images by breed',
    icon: Dog,
    fields: [],
    category: 'fun',
    color: 'text-amber-600'
  },
  {
    id: 'pokeapi',
    name: 'PokeAPI',
    description: 'Pokémon data and information',
    icon: Gamepad2,
    fields: [],
    category: 'fun',
    color: 'text-yellow-600'
  },
  {
    id: 'boredapi',
    name: 'Bored API',
    description: 'Random activity suggestions',
    icon: Zap,
    fields: [],
    category: 'fun',
    color: 'text-purple-500'
  },
  {
    id: 'numbersapi',
    name: 'Numbers API',
    description: 'Interesting facts about numbers',
    icon: Hash,
    fields: [],
    category: 'fun',
    color: 'text-teal-600'
  },
  {
    id: 'fakestore',
    name: 'Fake Store API',
    description: 'Mock e-commerce product data',
    icon: ShoppingBag,
    fields: [],
    category: 'enterprise',
    color: 'text-emerald-600'
  },
  
  // Enterprise
  {
    id: 'swapi',
    name: 'SWAPI',
    description: 'Star Wars API - movies, characters, planets',
    icon: Star,
    fields: [],
    category: 'fun',
    color: 'text-yellow-500'
  },
  {
    id: 'salesforce',
    name: 'Salesforce',
    description: 'CRM data and operations',
    icon: Briefcase,
    fields: [
      { key: 'instance_url', label: 'Instance URL', type: 'url', placeholder: 'https://yourinstance.salesforce.com', required: true },
      { key: 'access_token', label: 'Access Token', type: 'password', placeholder: 'Your OAuth access token', required: true }
    ],
    category: 'enterprise',
    color: 'text-blue-500'
  },
]

export default function ConnectionsPage() {
  const [selectedConnector, setSelectedConnector] = useState<string | null>(null)
  const [formData, setFormData] = useState<Record<string, string>>({})
  const [loading, setLoading] = useState(false)
  const [success, setSuccess] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)

  const connector = CONNECTORS.find(c => c.id === selectedConnector)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!connector) return

    setLoading(true)
    setError(null)
    setSuccess(null)

    try {
      // Save credentials to backend
      const credentialData = {
        service_name: connector.id,
        api_key: JSON.stringify(formData) // Store all fields as JSON
      }
      
      await credentials.create(credentialData.service_name, credentialData.api_key)
      
      setSuccess(`${connector.name} connected successfully!`)
      setFormData({})
      setTimeout(() => {
        setSelectedConnector(null)
        setSuccess(null)
      }, 2000)
    } catch (err: any) {
      setError(err.message || 'Failed to save connection')
    } finally {
      setLoading(false)
    }
  }

  const categories = [
    { id: 'all', label: 'All', icon: null },
    { id: 'messaging', label: 'Messaging', icon: MessageSquare },
    { id: 'data', label: 'Data', icon: Cloud },
    { id: 'fun', label: 'Fun', icon: Gamepad2 },
    { id: 'enterprise', label: 'Enterprise', icon: Briefcase },
  ]

  const [activeCategory, setActiveCategory] = useState('all')

  const filteredConnectors = activeCategory === 'all' 
    ? CONNECTORS 
    : CONNECTORS.filter(c => c.category === activeCategory)

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Connections</h1>
        <p className="text-muted-foreground">
          Connect your external services to enable integrations
        </p>
      </div>

      <Tabs value={activeCategory} onValueChange={setActiveCategory}>
        <TabsList>
          {categories.map(cat => (
            <TabsTrigger key={cat.id} value={cat.id}>
              {cat.icon && <cat.icon className="w-4 h-4 mr-2" />}
              {cat.label}
              <Badge variant="secondary" className="ml-2">
                {cat.id === 'all' ? CONNECTORS.length : CONNECTORS.filter(c => c.category === cat.id).length}
              </Badge>
            </TabsTrigger>
          ))}
        </TabsList>
      </Tabs>

      {selectedConnector && connector ? (
        <Card>
          <CardHeader>
            <div className="flex items-start justify-between">
              <div className="flex items-center gap-3">
                <connector.icon className={`w-8 h-8 ${connector.color}`} />
                <div>
                  <CardTitle>{connector.name}</CardTitle>
                  <CardDescription>{connector.description}</CardDescription>
                </div>
              </div>
              <Button variant="ghost" onClick={() => setSelectedConnector(null)}>
                Close
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            {success && (
              <Alert className="mb-4 border-green-200 bg-green-50">
                <CheckCircle2 className="h-4 w-4 text-green-600" />
                <AlertDescription className="text-green-800">{success}</AlertDescription>
              </Alert>
            )}
            
            {error && (
              <Alert className="mb-4" variant="destructive">
                <AlertCircle className="h-4 w-4" />
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}

            {connector.fields.length === 0 ? (
              <Alert>
                <AlertDescription>
                  This connector doesn't require any API keys or credentials. You can use it directly in your workflows!
                </AlertDescription>
              </Alert>
            ) : (
              <form onSubmit={handleSubmit} className="space-y-4">
                {connector.fields.map(field => (
                  <div key={field.key} className="space-y-2">
                    <Label htmlFor={field.key}>
                      {field.label}
                      {field.required && <span className="text-red-500 ml-1">*</span>}
                    </Label>
                    <Input
                      id={field.key}
                      type={field.type}
                      placeholder={field.placeholder}
                      value={formData[field.key] || ''}
                      onChange={(e) => setFormData({ ...formData, [field.key]: e.target.value })}
                      required={field.required}
                    />
                  </div>
                ))}

                <Button type="submit" disabled={loading} className="w-full">
                  {loading ? 'Connecting...' : `Connect ${connector.name}`}
                </Button>
              </form>
            )}
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {filteredConnectors.map(conn => (
            <Card 
              key={conn.id} 
              className="cursor-pointer hover:shadow-lg transition-shadow"
              onClick={() => setSelectedConnector(conn.id)}
            >
              <CardHeader>
                <div className="flex items-start gap-3">
                  <conn.icon className={`w-6 h-6 ${conn.color}`} />
                  <div className="flex-1">
                    <CardTitle className="text-lg">{conn.name}</CardTitle>
                    <CardDescription className="text-sm mt-1">
                      {conn.description}
                    </CardDescription>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <div className="flex items-center justify-between">
                  <Badge variant="outline" className="text-xs">
                    {conn.category}
                  </Badge>
                  {conn.fields.length === 0 && (
                    <Badge variant="secondary" className="text-xs">
                      No setup needed
                    </Badge>
                  )}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {filteredConnectors.length === 0 && (
        <div className="text-center py-12">
          <p className="text-muted-foreground">No connectors found in this category</p>
        </div>
      )}
    </div>
  )
}
