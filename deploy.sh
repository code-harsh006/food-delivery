#!/bin/bash

# Deployment script for multiple platforms

echo "🚀 Food Delivery Backend Deployment Script"
echo "==========================================="

# Check if platform is specified
if [ -z "$1" ]; then
    echo "Usage: ./deploy.sh [platform]"
    echo "Platforms: vercel, railway, render, heroku, docker"
    exit 1
fi

PLATFORM=$1

case $PLATFORM in
    "vercel")
        echo "📦 Deploying to Vercel..."
        
        # Use the correct vercel.json
        cp vercel.json vercel.json.backup 2>/dev/null || true
        cp vercel-simple.json vercel.json
        
        # Deploy
        if command -v vercel &> /dev/null; then
            vercel --prod
        else
            echo "❌ Vercel CLI not found. Install with: npm i -g vercel"
            exit 1
        fi
        
        # Restore backup
        mv vercel.json.backup vercel.json 2>/dev/null || true
        ;;
        
    "railway")
        echo "🚂 Deploying to Railway..."
        if command -v railway &> /dev/null; then
            railway up
        else
            echo "❌ Railway CLI not found. Install with: npm i -g @railway/cli"
            exit 1
        fi
        ;;
        
    "render")
        echo "🎨 Deploying to Render..."
        echo "Please connect your GitHub repository to Render dashboard"
        echo "Render will automatically deploy using render.yaml configuration"
        ;;
        
    "heroku")
        echo "🟣 Deploying to Heroku..."
        if command -v heroku &> /dev/null; then
            heroku stack:set container
            git push heroku main
        else
            echo "❌ Heroku CLI not found. Install from: https://devcenter.heroku.com/articles/heroku-cli"
            exit 1
        fi
        ;;
        
    "docker")
        echo "🐳 Building Docker image..."
        docker build -f Dockerfile.production -t food-delivery-backend .
        echo "✅ Docker image built successfully!"
        echo "Run with: docker run -p 8080:8080 food-delivery-backend"
        ;;
        
    *)
        echo "❌ Unknown platform: $PLATFORM"
        echo "Supported platforms: vercel, railway, render, heroku, docker"
        exit 1
        ;;
esac

echo "✅ Deployment process completed!"
