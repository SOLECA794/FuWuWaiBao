param(
    [switch]$SkipDocker
)

$root = Split-Path -Parent $PSScriptRoot
$logsDir = Join-Path $root "logs"
$tmpDir = Join-Path $root "tmp"
$scriptsDir = Join-Path $root "scripts"

New-Item -ItemType Directory -Force -Path $logsDir | Out-Null
New-Item -ItemType Directory -Force -Path $tmpDir | Out-Null
New-Item -ItemType Directory -Force -Path $scriptsDir | Out-Null

Write-Host "Root: $root"

function Invoke-DockerCompose {
    param(
        [Parameter(Mandatory)]
        [string[]]$ActionArgs
    )

    $composeFile = Join-Path $root "backend\docker-compose.yml"
    $variants = @(
        @{ Cmd = "docker"; Args = @("compose","-f",$composeFile) + $ActionArgs },
        @{ Cmd = "docker-compose"; Args = @("-f",$composeFile) + $ActionArgs }
    )

    foreach ($variant in $variants) {
        if (Get-Command $variant.Cmd -ErrorAction SilentlyContinue) {
            Write-Host "Running $($variant.Cmd) $($variant.Args -join ' ') ..."
            & $variant.Cmd @($variant.Args)
            if ($LASTEXITCODE -eq 0) {
                return
            }
            throw "Command '$($variant.Cmd) $($variant.Args -join ' ')' failed with exit code $LASTEXITCODE."
        }
    }

    throw "Neither 'docker' nor 'docker-compose' is available in PATH. Please install Docker CLI before running this script."
}

if (-not $SkipDocker) {
    Invoke-DockerCompose -ActionArgs @('up','-d')
}

$pids = @{}

function Start-Background {
    param(
        [Parameter(Mandatory)][string]$Name,
        [Parameter(Mandatory)][string]$FilePath,
        [Parameter(Mandatory)][string[]]$ArgumentList,
        [Parameter(Mandatory)][string]$WorkingDirectory,
        [Parameter(Mandatory)][string]$StdOut,
        [Parameter(Mandatory)][string]$StdErr
    )

    Write-Host "Starting $Name..."
    $proc = Start-Process -FilePath $FilePath -ArgumentList $ArgumentList -WorkingDirectory $WorkingDirectory -NoNewWindow -PassThru -RedirectStandardOutput $StdOut -RedirectStandardError $StdErr
    Start-Sleep -Milliseconds 500

    if ($proc -and $proc.Id) {
        $pids[$Name] = $proc.Id
        Write-Host "$Name started (PID $($proc.Id))"
    } else {
        Write-Warning "Failed to start $Name"
    }
}

Start-Background -Name "ai_engine" -FilePath "conda" -ArgumentList @('run','-n','fuww_ai','python','ai_engine/main.py') -WorkingDirectory $root -StdOut "$logsDir\ai_engine.log" -StdErr "$logsDir\ai_engine.err.log"
Start-Background -Name "backend" -FilePath "go" -ArgumentList @('run','./api/main.go') -WorkingDirectory (Join-Path $root "backend") -StdOut "$logsDir\backend.log" -StdErr "$logsDir\backend.err.log"
Start-Background -Name "student_frontend" -FilePath "npm" -ArgumentList @('run','serve') -WorkingDirectory (Join-Path $root "frontend\student") -StdOut "$logsDir\student_frontend.log" -StdErr "$logsDir\student_frontend.err.log"
Start-Background -Name "teacher_frontend" -FilePath "npm" -ArgumentList @('run','dev','--','--host','0.0.0.0','--port','5173') -WorkingDirectory (Join-Path $root "frontend\teacher") -StdOut "$logsDir\teacher_frontend.log" -StdErr "$logsDir\teacher_frontend.err.log"

$pidFile = Join-Path $scriptsDir "pids.json"
$pids | ConvertTo-Json | Out-File -Encoding utf8 $pidFile

Write-Host "Started components. PID file: $pidFile"
Write-Host "Logs: $logsDir"
Write-Host "If something fails, check the log files and run scripts/stop-all.ps1 to stop the background processes."
Write-Host "Done."