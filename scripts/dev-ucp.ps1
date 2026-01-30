$ErrorActionPreference = "Stop"

param(
  [int]$MockWebhookFailStatus = 0,
  [int]$SleepSeconds = 3,
  [switch]$StopExisting,
  [switch]$StopOnly,
  [switch]$NoVerify
)

$root = "C:\Users\1000067\Documents\gitea\meowucp"
Set-Location $root

if ($StopExisting) {
  Write-Host "Stopping existing processes on ports 8081, 9091, 9092..."
  $ports = @(8081, 9091, 9092)
  foreach ($port in $ports) {
    $pids = (netstat -ano | Select-String ":$port" | ForEach-Object {
      ($_ -split "\s+")[-1]
    }) | Where-Object { $_ -match "^\d+$" } | Sort-Object -Unique
    foreach ($pid in $pids) {
      try {
        taskkill /F /PID $pid | Out-Null
      } catch {
        Write-Host "Failed to stop PID $pid" -ForegroundColor Yellow
      }
    }
  }

  Write-Host "Stopping existing worker processes..."
  $workerPids = (wmic process where "CommandLine like '%cmd/worker/main.go%' or CommandLine like '%cmd\\worker\\main.go%'" get ProcessId | ForEach-Object { $_.Trim() }) | Where-Object { $_ -match "^\d+$" } | Sort-Object -Unique
  foreach ($pid in $workerPids) {
    try {
      taskkill /F /PID $pid | Out-Null
    } catch {
      Write-Host "Failed to stop worker PID $pid" -ForegroundColor Yellow
    }
  }

  Write-Host "Stopping existing mock services..."
  $mockPids = (wmic process where "CommandLine like '%cmd/mock-jwks/main.go%' or CommandLine like '%cmd\\mock-jwks\\main.go%' or CommandLine like '%cmd/mock-webhook/main.go%' or CommandLine like '%cmd\\mock-webhook\\main.go%'" get ProcessId | ForEach-Object { $_.Trim() }) | Where-Object { $_ -match "^\d+$" } | Sort-Object -Unique
  foreach ($pid in $mockPids) {
    try {
      taskkill /F /PID $pid | Out-Null
    } catch {
      Write-Host "Failed to stop mock PID $pid" -ForegroundColor Yellow
    }
  }

  Write-Host "Stopping existing API processes..."
  $apiPids = (wmic process where "CommandLine like '%cmd/api/main.go%' or CommandLine like '%cmd\\api\\main.go%'" get ProcessId | ForEach-Object { $_.Trim() }) | Where-Object { $_ -match "^\d+$" } | Sort-Object -Unique
  foreach ($pid in $apiPids) {
    try {
      taskkill /F /PID $pid | Out-Null
    } catch {
      Write-Host "Failed to stop API PID $pid" -ForegroundColor Yellow
    }
  }

  Write-Host "Stopping existing send-webhook processes..."
  $sendPids = (wmic process where "CommandLine like '%cmd/send-webhook/main.go%' or CommandLine like '%cmd\\send-webhook\\main.go%'" get ProcessId | ForEach-Object { $_.Trim() }) | Where-Object { $_ -match "^\d+$" } | Sort-Object -Unique
  foreach ($pid in $sendPids) {
    try {
      taskkill /F /PID $pid | Out-Null
    } catch {
      Write-Host "Failed to stop send-webhook PID $pid" -ForegroundColor Yellow
    }
  }
}

if ($StopOnly) {
  Write-Host "StopOnly set, exiting without start."
  exit 0
}

Write-Host "[1/4] Running migrations..."
go run scripts/migrate.go

Write-Host "[2/4] Starting mock JWKS, mock webhook, API, worker..."
if ($MockWebhookFailStatus -ne 0) {
  $env:MOCK_WEBHOOK_FAIL_STATUS = $MockWebhookFailStatus
} else {
  Remove-Item Env:MOCK_WEBHOOK_FAIL_STATUS -ErrorAction SilentlyContinue
}

Start-Process -FilePath go -ArgumentList "run","cmd/mock-jwks/main.go" -WorkingDirectory $root
Start-Process -FilePath go -ArgumentList "run","cmd/mock-webhook/main.go" -WorkingDirectory $root
Start-Process -FilePath go -ArgumentList "run","cmd/api/main.go" -WorkingDirectory $root
Start-Process -FilePath go -ArgumentList "run","cmd/worker/main.go" -WorkingDirectory $root

Write-Host "[3/4] Waiting $SleepSeconds seconds..."
Start-Sleep -Seconds $SleepSeconds

if (-not $NoVerify) {
  Write-Host "[4/4] Sending webhook and inspecting state..."
  go run cmd/send-webhook/main.go
  go run scripts/inspect_webhook_state.go
} else {
  Write-Host "[4/4] NoVerify set, skipping webhook/inspect."
}
