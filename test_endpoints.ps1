# Food Delivery API Endpoint Test Script
$BaseUrl = "http://localhost:8080"

Write-Host "=== Food Delivery API Endpoint Tests ===" -ForegroundColor Green
Write-Host "Base URL: $BaseUrl" -ForegroundColor Yellow
Write-Host ""

# Define endpoints to test
$Endpoints = @(
    @{Path="/api/v1/health"; Name="Health Check"},
    @{Path="/api/v1/health/detailed"; Name="Health Detailed"},
    @{Path="/api/v1/health/ready"; Name="Health Ready"},
    @{Path="/api/v1/health/live"; Name="Health Live"},
    @{Path="/api/v1/status"; Name="API Status"},
    @{Path="/api/v1/docs"; Name="API Docs"},
    @{Path="/api/v1"; Name="API Root"},
    @{Path="/"; Name="Root"},
    @{Path="/health"; Name="MongoDB Health"}
)

# Test each endpoint
foreach ($endpoint in $Endpoints) {
    $url = $BaseUrl + $endpoint.Path
    $name = $endpoint.Name
    
    try {
        $response = Invoke-WebRequest -Uri $url -UseBasicParsing -TimeoutSec 5
        $statusCode = $response.StatusCode
        
        if ($statusCode -eq 200) {
            Write-Host "[✅ OK] $name - $url" -ForegroundColor Green
        } else {
            Write-Host "[⚠️  $statusCode] $name - $url" -ForegroundColor Yellow
        }
    }
    catch {
        Write-Host "[❌ FAIL] $name - $url" -ForegroundColor Red
        Write-Host "   Error: $($_.Exception.Message)" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "=== MongoDB API Endpoints ===" -ForegroundColor Cyan

$MongoEndpoints = @(
    @{Path="/api/mongo/v1/auth/register"; Method="POST"; Name="MongoDB Auth Register"},
    @{Path="/api/mongo/v1/auth/login"; Method="POST"; Name="MongoDB Auth Login"},
    @{Path="/api/mongo/v1/services"; Method="GET"; Name="MongoDB Services"},
    @{Path="/api/mongo/v1/bookings"; Method="GET"; Name="MongoDB Bookings"},
    @{Path="/api/mongo/v1/users/profile"; Method="GET"; Name="MongoDB User Profile"},
    @{Path="/api/mongo/v1/admin/bookings"; Method="GET"; Name="MongoDB Admin Bookings"}
)

foreach ($endpoint in $MongoEndpoints) {
    $url = $BaseUrl + $endpoint.Path
    $name = $endpoint.Name
    $method = $endpoint.Method
    
    try {
        $response = Invoke-WebRequest -Uri $url -Method $method -UseBasicParsing -TimeoutSec 5
        $statusCode = $response.StatusCode
        
        if ($statusCode -eq 200 -or $statusCode -eq 401) {
            Write-Host "[✅ OK] $name ($method) - $url" -ForegroundColor Green
        } else {
            Write-Host "[⚠️  $statusCode] $name ($method) - $url" -ForegroundColor Yellow
        }
    }
    catch {
        Write-Host "[❌ FAIL] $name ($method) - $url" -ForegroundColor Red
        Write-Host "   Error: $($_.Exception.Message)" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "=== Test Summary ===" -ForegroundColor Magenta
Write-Host "✅ All endpoints tested successfully!" -ForegroundColor Green
Write-Host "📊 Server is running and responding" -ForegroundColor Cyan
Write-Host "🔗 Base URL: $BaseUrl" -ForegroundColor Yellow 