#!/bin/bash

# Deployment verification script

if [ -z "$1" ]; then
    echo "Usage: ./verify-deployment.sh <your-vercel-url>"
    echo "Example: ./verify-deployment.sh https://your-app.vercel.app"
    exit 1
fi

BASE_URL=$1

echo "🧪 Testing deployment at: $BASE_URL"
echo "=================================="

# Test health endpoint
echo "1. Testing health endpoint..."
HEALTH_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/health")
if [ "$HEALTH_RESPONSE" = "200" ]; then
    echo "✅ Health check passed"
else
    echo "❌ Health check failed (HTTP $HEALTH_RESPONSE)"
fi

# Test root endpoint
echo "2. Testing root endpoint..."
ROOT_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/")
if [ "$ROOT_RESPONSE" = "200" ]; then
    echo "✅ Root endpoint passed"
else
    echo "❌ Root endpoint failed (HTTP $ROOT_RESPONSE)"
fi

# Test auth endpoints
echo "3. Testing auth endpoints..."
AUTH_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/api/v1/auth/login" -X POST -H "Content-Type: application/json" -d '{}')
if [ "$AUTH_RESPONSE" = "200" ]; then
    echo "✅ Auth endpoint passed"
else
    echo "❌ Auth endpoint failed (HTTP $AUTH_RESPONSE)"
fi

# Test products endpoint
echo "4. Testing products endpoint..."
PRODUCTS_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/api/v1/products")
if [ "$PRODUCTS_RESPONSE" = "200" ]; then
    echo "✅ Products endpoint passed"
else
    echo "❌ Products endpoint failed (HTTP $PRODUCTS_RESPONSE)"
fi

# Test cart endpoint
echo "5. Testing cart endpoint..."
CART_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/api/v1/cart")
if [ "$CART_RESPONSE" = "200" ]; then
    echo "✅ Cart endpoint passed"
else
    echo "❌ Cart endpoint failed (HTTP $CART_RESPONSE)"
fi

echo ""
echo "🎉 Deployment verification completed!"
echo "Visit your app at: $BASE_URL"
