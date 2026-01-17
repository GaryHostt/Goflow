"use client";

import { ArrowRight, Zap, Send, MessageSquare, Phone, Cloud, Wifi, Code, Database, Star, Building2, TestTube, Calendar } from "lucide-react";
import { useState } from "react";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { Label } from "./ui/label";
import { Textarea } from "./ui/textarea";

interface FlowDiagramProps {
  triggerType: string;
  actionType: string;
  config: any;
  webhookPath?: string;
  webhookPayload?: string;
  onTriggerTypeChange: (type: string) => void;
  onActionTypeChange: (type: string) => void;
  onWebhookPathChange?: (path: string) => void;
  onWebhookPayloadChange?: (payload: string) => void;
  onConfigChange: (config: any) => void;
}

// Connector metadata with icons and colors
const connectorData: Record<string, { name: string; icon: any; color: string; bgColor: string }> = {
  slack_message: { name: "Slack", icon: MessageSquare, color: "text-purple-600", bgColor: "bg-purple-50 border-purple-200" },
  discord_post: { name: "Discord", icon: Send, color: "text-indigo-600", bgColor: "bg-indigo-50 border-indigo-200" },
  twilio_sms: { name: "Twilio SMS", icon: Phone, color: "text-red-600", bgColor: "bg-red-50 border-red-200" },
  weather_check: { name: "OpenWeather", icon: Cloud, color: "text-blue-600", bgColor: "bg-blue-50 border-blue-200" },
  news_fetch: { name: "News API", icon: Wifi, color: "text-orange-600", bgColor: "bg-orange-50 border-orange-200" },
  cat_fetch: { name: "Cat API", icon: Star, color: "text-pink-600", bgColor: "bg-pink-50 border-pink-200" },
  fakestore_fetch: { name: "Fake Store", icon: Database, color: "text-green-600", bgColor: "bg-green-50 border-green-200" },
  soap_call: { name: "SOAP Bridge", icon: Code, color: "text-gray-600", bgColor: "bg-gray-50 border-gray-200" },
  swapi_fetch: { name: "SWAPI", icon: Star, color: "text-yellow-600", bgColor: "bg-yellow-50 border-yellow-200" },
  salesforce: { name: "Salesforce", icon: Building2, color: "text-cyan-600", bgColor: "bg-cyan-50 border-cyan-200" },
  testing: { name: "Testing", icon: TestTube, color: "text-emerald-600", bgColor: "bg-emerald-50 border-emerald-200" },
};

const triggerData: Record<string, { name: string; icon: any; color: string; bgColor: string }> = {
  webhook: { name: "Webhook", icon: Zap, color: "text-emerald-600", bgColor: "bg-emerald-50 border-emerald-200" },
  schedule: { name: "Schedule", icon: Calendar, color: "text-blue-600", bgColor: "bg-blue-50 border-blue-200" },
};

const actionOptions = [
  { value: "slack_message", label: "Slack Message" },
  { value: "discord_post", label: "Discord Post" },
  { value: "twilio_sms", label: "Twilio SMS" },
  { value: "testing", label: "Testing/Mock" },
  { value: "weather_check", label: "Weather Check" },
  { value: "news_fetch", label: "News API" },
  { value: "salesforce", label: "Salesforce" },
];

export function WorkflowFlowDiagram(props: FlowDiagramProps) {
  const [showTriggerConfig, setShowTriggerConfig] = useState(false);
  const [showActionConfig, setShowActionConfig] = useState(false);
  const [tempWebhookPath, setTempWebhookPath] = useState(props.webhookPath || '/integration/my-endpoint');
  const [tempWebhookPayload, setTempWebhookPayload] = useState(props.webhookPayload || '{\n  "user": {\n    "name": ""\n  }\n}');
  const [tempConfig, setTempConfig] = useState(props.config);
  const [tempInterval, setTempInterval] = useState(props.config.interval || 10);

  const trigger = triggerData[props.triggerType];
  const action = connectorData[props.actionType];

  if (!trigger || !action) {
    return (
      <div className="flex items-center justify-center h-full text-gray-400 text-sm py-20">
        Select a trigger and action to visualize your workflow
      </div>
    );
  }

  const TriggerIcon = trigger.icon;
  const ActionIcon = action.icon;

  const handleTriggerClick = () => {
    setTempWebhookPath(props.webhookPath || '/integration/my-endpoint');
    setTempWebhookPayload(props.webhookPayload || '{\n  "user": {\n    "name": ""\n  }\n}');
    setTempInterval(props.config.interval || 10);
    setShowTriggerConfig(true);
  };

  const handleActionClick = () => {
    setTempConfig(props.config);
    setShowActionConfig(true);
  };

  const saveTriggerConfig = () => {
    if (props.triggerType === 'webhook') {
      if (props.onWebhookPathChange) props.onWebhookPathChange(tempWebhookPath);
      if (props.onWebhookPayloadChange) props.onWebhookPayloadChange(tempWebhookPayload);
    } else if (props.triggerType === 'schedule') {
      props.onConfigChange({ ...props.config, interval: tempInterval });
    }
    setShowTriggerConfig(false);
  };

  const saveActionConfig = () => {
    props.onConfigChange(tempConfig);
    setShowActionConfig(false);
  };

  // Parse webhook payload to extract field paths
  const getFieldPaths = (obj: any, prefix = ''): string[] => {
    let paths: string[] = [];
    try {
      const parsed = typeof obj === 'string' ? JSON.parse(obj) : obj;
      Object.keys(parsed).forEach(key => {
        const path = prefix ? `${prefix}.${key}` : key;
        paths.push(path);
        if (typeof parsed[key] === 'object' && parsed[key] !== null && !Array.isArray(parsed[key])) {
          paths = paths.concat(getFieldPaths(parsed[key], path));
        }
      });
    } catch (e) {
      // Invalid JSON, return empty
    }
    return paths;
  };

  const availableFields = props.triggerType === 'webhook' ? getFieldPaths(props.webhookPayload || '{}') : [];

  return (
    <div className="space-y-8 py-6">
      {/* Flow Diagram */}
      <div className="flex items-center justify-center space-x-6">
        {/* Trigger Box */}
        <button
          type="button"
          onClick={handleTriggerClick}
          className={`flex flex-col items-center group cursor-pointer`}
        >
          <div className={`border-2 rounded-lg p-6 ${trigger.bgColor} transition-all hover:shadow-lg hover:scale-105 relative`}>
            <TriggerIcon className={`w-12 h-12 ${trigger.color}`} />
            <div className="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity bg-black/5 rounded-lg">
              <span className="text-xs font-medium">Click to configure</span>
            </div>
          </div>
          <span className="text-sm font-medium text-gray-700 mt-3">{trigger.name}</span>
          <span className="text-xs text-gray-500">Trigger</span>
        </button>

        <ArrowRight className="w-8 h-8 text-gray-400" />

        {/* Engine */}
        <div className="flex flex-col items-center">
          <div className="border-2 border-dashed border-gray-300 rounded-lg p-6 bg-white">
            <div className="w-12 h-12 bg-gradient-to-br from-blue-500 to-purple-600 rounded flex items-center justify-center">
              <span className="text-white text-sm font-bold">GF</span>
            </div>
          </div>
          <span className="text-sm font-medium text-gray-700 mt-3">GoFlow</span>
          <span className="text-xs text-gray-500">Engine</span>
        </div>

        <ArrowRight className="w-8 h-8 text-gray-400" />

        {/* Action Box */}
        <button
          type="button"
          onClick={handleActionClick}
          className={`flex flex-col items-center group cursor-pointer`}
        >
          <div className={`border-2 rounded-lg p-6 ${action.bgColor} transition-all hover:shadow-lg hover:scale-105 relative`}>
            <ActionIcon className={`w-12 h-12 ${action.color}`} />
            <div className="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity bg-black/5 rounded-lg">
              <span className="text-xs font-medium">Click to configure</span>
            </div>
          </div>
          <span className="text-sm font-medium text-gray-700 mt-3">{action.name}</span>
          <span className="text-xs text-gray-500">Action</span>
        </button>
      </div>

      {/* Trigger Configuration Panel */}
      {showTriggerConfig && (
        <div className="border rounded-lg p-6 bg-gray-50 space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-semibold">Configure {trigger.name} Trigger</h3>
            <button onClick={() => setShowTriggerConfig(false)} className="text-gray-500 hover:text-gray-700">✕</button>
          </div>

          {/* Trigger Type Selector */}
          <div className="space-y-2">
            <Label>Trigger Type</Label>
            <select
              className="w-full h-10 rounded-md border border-gray-300 px-3"
              value={props.triggerType}
              onChange={(e) => props.onTriggerTypeChange(e.target.value)}
            >
              <option value="webhook">Webhook</option>
              <option value="schedule">Schedule</option>
            </select>
          </div>

          {props.triggerType === 'webhook' && (
            <>
              <div className="space-y-2">
                <Label>Webhook Endpoint Path</Label>
                <Input
                  value={tempWebhookPath}
                  onChange={(e) => setTempWebhookPath(e.target.value)}
                  placeholder="/integration/my-endpoint"
                />
                <p className="text-xs text-gray-500">
                  Full URL: {`http://localhost:8080/webhook${tempWebhookPath}`}
                </p>
              </div>
              <div className="space-y-2">
                <Label>Expected Payload Schema (JSON)</Label>
                <Textarea
                  value={tempWebhookPayload}
                  onChange={(e) => setTempWebhookPayload(e.target.value)}
                  rows={8}
                  className="font-mono text-sm"
                  placeholder='{\n  "user": {\n    "name": ""\n  }\n}'
                />
                <p className="text-xs text-gray-500">
                  Define the structure of incoming webhook data for field mapping
                </p>
              </div>
            </>
          )}

          {props.triggerType === 'schedule' && (
            <div className="space-y-2">
              <Label>Interval (minutes)</Label>
              <Input
                type="number"
                min="1"
                value={tempInterval}
                onChange={(e) => setTempInterval(parseInt(e.target.value))}
              />
            </div>
          )}

          <div className="flex justify-end gap-2 pt-4">
            <Button variant="outline" onClick={() => setShowTriggerConfig(false)}>Cancel</Button>
            <Button onClick={saveTriggerConfig}>Save Configuration</Button>
          </div>
        </div>
      )}

      {/* Action Configuration Panel */}
      {showActionConfig && (
        <div className="border rounded-lg p-6 bg-gray-50 space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-semibold">Configure {action.name} Action</h3>
            <button onClick={() => setShowActionConfig(false)} className="text-gray-500 hover:text-gray-700">✕</button>
          </div>

          {/* Action Type Selector */}
          <div className="space-y-2">
            <Label>Action Type</Label>
            <select
              className="w-full h-10 rounded-md border border-gray-300 px-3"
              value={props.actionType}
              onChange={(e) => props.onActionTypeChange(e.target.value)}
            >
              {actionOptions.map(opt => (
                <option key={opt.value} value={opt.value}>{opt.label}</option>
              ))}
            </select>
          </div>

          {/* Available Fields (from webhook) */}
          {availableFields.length > 0 && (
            <div className="p-3 bg-blue-50 border border-blue-200 rounded">
              <p className="text-sm font-medium text-blue-900 mb-2">Available Fields:</p>
              <div className="flex flex-wrap gap-2">
                {availableFields.map(field => (
                  <code key={field} className="text-xs bg-white px-2 py-1 rounded border">
                    {`{{${field}}}`}
                  </code>
                ))}
              </div>
            </div>
          )}

          {/* Action-specific configuration */}
          {props.actionType === 'slack_message' && (
            <div className="space-y-2">
              <Label>Message</Label>
              <Textarea
                value={tempConfig.slack_message || ''}
                onChange={(e) => setTempConfig({ ...tempConfig, slack_message: e.target.value })}
                placeholder="Hello {{user.name}}!"
                rows={4}
              />
            </div>
          )}

          {props.actionType === 'testing' && (
            <div className="space-y-2">
              <Label>Response JSON</Label>
              <Textarea
                value={tempConfig.testing_response_json || '{"message": "success"}'}
                onChange={(e) => setTempConfig({ ...tempConfig, testing_response_json: e.target.value })}
                rows={6}
                className="font-mono text-sm"
              />
            </div>
          )}

          <div className="flex justify-end gap-2 pt-4">
            <Button variant="outline" onClick={() => setShowActionConfig(false)}>Cancel</Button>
            <Button onClick={saveActionConfig}>Save Configuration</Button>
          </div>
        </div>
      )}
    </div>
  );
}
