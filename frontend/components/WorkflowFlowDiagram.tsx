"use client";

import { ArrowRight, Zap, Send, MessageSquare, Phone, Cloud, Wifi, Code, Database, Star, Building2 } from "lucide-react";

interface FlowDiagramProps {
  triggerType: string;
  actionType: string;
}

// Connector metadata with icons and colors
const connectorData: Record<string, { name: string; icon: any; color: string; bgColor: string }> = {
  slack_message: {
    name: "Slack",
    icon: MessageSquare,
    color: "text-purple-600",
    bgColor: "bg-purple-50 border-purple-200",
  },
  discord_post: {
    name: "Discord",
    icon: Send,
    color: "text-indigo-600",
    bgColor: "bg-indigo-50 border-indigo-200",
  },
  twilio_sms: {
    name: "Twilio SMS",
    icon: Phone,
    color: "text-red-600",
    bgColor: "bg-red-50 border-red-200",
  },
  weather_check: {
    name: "OpenWeather",
    icon: Cloud,
    color: "text-blue-600",
    bgColor: "bg-blue-50 border-blue-200",
  },
  news_fetch: {
    name: "News API",
    icon: Wifi,
    color: "text-orange-600",
    bgColor: "bg-orange-50 border-orange-200",
  },
  cat_fetch: {
    name: "Cat API",
    icon: Star,
    color: "text-pink-600",
    bgColor: "bg-pink-50 border-pink-200",
  },
  fakestore_fetch: {
    name: "Fake Store",
    icon: Database,
    color: "text-green-600",
    bgColor: "bg-green-50 border-green-200",
  },
  soap_call: {
    name: "SOAP Bridge",
    icon: Code,
    color: "text-gray-600",
    bgColor: "bg-gray-50 border-gray-200",
  },
  swapi_fetch: {
    name: "SWAPI",
    icon: Star,
    color: "text-yellow-600",
    bgColor: "bg-yellow-50 border-yellow-200",
  },
  salesforce: {
    name: "Salesforce",
    icon: Building2,
    color: "text-cyan-600",
    bgColor: "bg-cyan-50 border-cyan-200",
  },
};

// Trigger type metadata
const triggerData: Record<string, { name: string; icon: any; color: string; bgColor: string }> = {
  webhook: {
    name: "Webhook",
    icon: Zap,
    color: "text-emerald-600",
    bgColor: "bg-emerald-50 border-emerald-200",
  },
  schedule: {
    name: "Schedule",
    icon: Cloud,
    color: "text-blue-600",
    bgColor: "bg-blue-50 border-blue-200",
  },
};

export function WorkflowFlowDiagram({ triggerType, actionType }: FlowDiagramProps) {
  const trigger = triggerData[triggerType];
  const action = connectorData[actionType];

  // Don't show if either is not selected
  if (!trigger || !action) {
    return (
      <div className="flex items-center justify-center h-full text-gray-400 text-sm">
        Select a trigger and action to visualize your workflow
      </div>
    );
  }

  const TriggerIcon = trigger.icon;
  const ActionIcon = action.icon;

  return (
    <div className="flex flex-col items-center justify-center h-full p-6 space-y-6">
      {/* Title */}
      <div className="text-center">
        <h3 className="text-sm font-semibold text-gray-700 mb-1">Workflow Flow</h3>
        <p className="text-xs text-gray-500">Visual representation of your integration</p>
      </div>

      {/* Flow Diagram */}
      <div className="flex items-center space-x-4">
        {/* Trigger Box */}
        <div className="flex flex-col items-center">
          <div
            className={`border-2 rounded-lg p-4 ${trigger.bgColor} transition-all hover:shadow-md`}
          >
            <TriggerIcon className={`w-8 h-8 ${trigger.color}`} />
          </div>
          <span className="text-xs font-medium text-gray-600 mt-2">{trigger.name}</span>
          <span className="text-xs text-gray-400">Trigger</span>
        </div>

        {/* Arrow */}
        <div className="flex flex-col items-center">
          <ArrowRight className="w-6 h-6 text-gray-400" />
          <span className="text-xs text-gray-400 mt-1">When</span>
        </div>

        {/* GoFlow Engine */}
        <div className="flex flex-col items-center">
          <div className="border-2 border-dashed border-gray-300 rounded-lg p-4 bg-white">
            <div className="flex items-center space-x-2">
              <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-purple-600 rounded flex items-center justify-center">
                <span className="text-white text-xs font-bold">GF</span>
              </div>
            </div>
          </div>
          <span className="text-xs font-medium text-gray-600 mt-2">GoFlow</span>
          <span className="text-xs text-gray-400">Engine</span>
        </div>

        {/* Arrow */}
        <div className="flex flex-col items-center">
          <ArrowRight className="w-6 h-6 text-gray-400" />
          <span className="text-xs text-gray-400 mt-1">Execute</span>
        </div>

        {/* Action Box */}
        <div className="flex flex-col items-center">
          <div
            className={`border-2 rounded-lg p-4 ${action.bgColor} transition-all hover:shadow-md`}
          >
            <ActionIcon className={`w-8 h-8 ${action.color}`} />
          </div>
          <span className="text-xs font-medium text-gray-600 mt-2">{action.name}</span>
          <span className="text-xs text-gray-400">Action</span>
        </div>
      </div>

      {/* Description */}
      <div className="text-center max-w-md">
        <p className="text-xs text-gray-600">
          <span className="font-semibold">{trigger.name}</span> triggers →{" "}
          <span className="font-semibold">GoFlow</span> processes →{" "}
          <span className="font-semibold">{action.name}</span> executes
        </p>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-3 gap-4 w-full max-w-md">
        <div className="bg-gray-50 rounded-lg p-3 text-center border border-gray-200">
          <div className="text-lg font-bold text-gray-700">~50ms</div>
          <div className="text-xs text-gray-500">Avg Latency</div>
        </div>
        <div className="bg-gray-50 rounded-lg p-3 text-center border border-gray-200">
          <div className="text-lg font-bold text-gray-700">99.9%</div>
          <div className="text-xs text-gray-500">Uptime</div>
        </div>
        <div className="bg-gray-50 rounded-lg p-3 text-center border border-gray-200">
          <div className="text-lg font-bold text-gray-700">10</div>
          <div className="text-xs text-gray-500">Workers</div>
        </div>
      </div>
    </div>
  );
}

