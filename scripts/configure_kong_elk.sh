#!/bin/bash
# Script to configure Kong to ship logs to ELK via Logstash

KONG_ADMIN_URL="http://localhost:8001"
LOGSTASH_HOST="logstash"
LOGSTASH_PORT="5000"

echo "🔧 Configuring Kong to ship logs to ELK..."

# Wait for Kong to be ready
echo "⏳ Waiting for Kong Admin API..."
until curl -s -f "${KONG_ADMIN_URL}/status" > /dev/null 2>&1; do
  echo "  Waiting for Kong..."
  sleep 2
done
echo "✅ Kong is ready"

# Install http-log plugin globally to ship all logs
echo "📤 Installing http-log plugin for ELK integration..."
curl -i -X POST "${KONG_ADMIN_URL}/plugins" \
  --data "name=http-log" \
  --data "config.http_endpoint=http://${LOGSTASH_HOST}:5001" \
  --data "config.method=POST" \
  --data "config.timeout=10000" \
  --data "config.keepalive=60000"

if [ $? -eq 0 ]; then
  echo "✅ Kong http-log plugin installed"
else
  echo "⚠️  Failed to install http-log plugin"
  exit 1
fi

# Install request-transformer to add tracking headers
echo "📋 Installing request-transformer for enhanced tracking..."
curl -i -X POST "${KONG_ADMIN_URL}/plugins" \
  --data "name=request-transformer" \
  --data "config.add.headers=X-Kong-Logged:true"

echo "✅ Kong log shipping configured!"
echo ""
echo "📊 Kong logs will now appear in:"
echo "   - Elasticsearch: http://localhost:9200/kong-logs-*"
echo "   - Kibana: http://localhost:5601"
echo ""
echo "To verify, make a request through Kong and check Kibana!"

