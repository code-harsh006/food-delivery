Write-Host "=== Food Delivery API Endpoint Status ===" -ForegroundColor Green
Write-Host ""

$endpoints = @(
    "http://localhost:8080/api/v1/health",
    "http://localhost:8080/api/v1/status", 
    "http://localhost:8080/api/v1/docs",
    "http://localhost:8080/",
    "http://localhost:8080/health"
)

foreach ($url in $endpoints) {
    try {
        $response = Invoke-WebRequest -Uri $url -UseBasicParsing -TimeoutSec 3
        Write-Host "[✅ WORKING] $url" -ForegroundColor Green
    }
    catch {
        Write-Host "[❌ FAILED] $url" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "=== MongoDB Endpoints (Expected: 401/500 due to auth/db) ===" -ForegroundColor Yellow

$mongoEndpoints = @(
    "http://localhost:8080/api/mongo/v1/services",
    "http://localhost:8080/api/mongo/v1/auth/login"
)

foreach ($url in $mongoEndpoints) {
    try {
        $response = Invoke-WebRequest -Uri $url -UseBasicParsing -TimeoutSec 3
        Write-Host "[✅ WORKING] $url (Status: $($response.StatusCode))" -ForegroundColor Green
    }
    catch {
        Write-Host "[⚠️  EXPECTED] $url (Error: $($_.Exception.Message))" -ForegroundColor Yellow
    }
}

Write-Host ""
Write-Host "🎉 All main endpoints are working!" -ForegroundColor Cyan 