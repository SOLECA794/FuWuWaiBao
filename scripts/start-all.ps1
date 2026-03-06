param(
    [switch]$SkipDocker
)

$root = Split-Path -Parent $PSScriptRoot
New-Item -ItemType Directory -Force -Path "$root\logs" | Out-Null
New-Item -ItemType Directory -Force -Path "$root\tmp" | Out-Null

Write-Host "Root: $root"

if (-not $SkipDocker) {
    Write-Host "Starting docker-compose (backend/docker-compose.yml) ..."
    docker-compose -f "$root\backend\docker-compose.yml" up -d
}

$pids = @{}

function Start-Background($file,$args,$workdir,$out,$err,$name){
    Write-Host "Starting $name..."
    $proc = Start-Process -FilePath $file -ArgumentList $args -WorkingDirectory $workdir -NoNewWindow -PassThru -RedirectStandardOutput $out -RedirectStandardError $err
    Start-Sleep -Milliseconds 500
    if ($proc -and $proc.Id) { $pids[$name] = $proc.Id; Write-Host "$name started (PID $($proc.Id))" } else { Write-Warning "Failed to start $name" }
}

# Start AI engine (in conda env 'fuww_ai')
Start-Background "conda" @('run','-n','fuww_ai','python','ai_engine/main.py') "$root" "$root\logs\ai_engine.log" "$root\logs\ai_engine.err.log" "ai_engine"

# Start backend (Go)
Start-Background "go" @('run','./api/main.go') "$root\backend" "$root\logs\backend.log" "$root\logs\backend.err.log" "backend"

# Start student frontend
Start-Background "npm" @('run','dev') "$root\frontend\student" "$root\logs\student_frontend.log" "$root\logs\student_frontend.err.log" "student_frontend"

# Start teacher frontend
Start-Background "npm" @('run','dev') "$root\frontend\teacher" "$root\logs\teacher_frontend.log" "$root\logs\teacher_frontend.err.log" "teacher_frontend"

# Save pids
$pids | ConvertTo-Json | Out-File -Encoding utf8 "$root\scripts\pids.json"

Write-Host "Started components. PID file: $root\scripts\pids.json"
Write-Host "Logs: $root\logs"
Write-Host "If something fails, check corresponding log files and run scripts/stop-all.ps1 to stop processes."

Write-Host "Done."
