param(
    [switch]$SkipDocker
)

$root = Split-Path -Parent $PSScriptRoot
$scriptsDir = Join-Path $root "scripts"

New-Item -ItemType Directory -Force -Path $scriptsDir | Out-Null

$pidFile = Join-Path $scriptsDir "pids.json"
if (Test-Path $pidFile) {
    try {
        $content = Get-Content $pidFile -Raw | ConvertFrom-Json
        foreach ($key in $content.PSObject.Properties.Name) {
            $id = $content.$key
            if ($id) {
                Write-Host "Stopping $key (PID $id) ..."
                Stop-Process -Id $id -Force -ErrorAction SilentlyContinue
            }
        }
        Remove-Item $pidFile -ErrorAction SilentlyContinue
        Write-Host "Stopped background processes."
    } catch {
        Write-Warning "Failed to read/stop processes from $pidFile: $_"
    }
} else {
    Write-Host "PID file not found: $pidFile. No background PIDs to stop."
}

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
    try {
        Invoke-DockerCompose -ActionArgs @('down')
    } catch {
        Write-Warning "Failed to stop docker compose services: $_"
    }
}

Write-Host "Done."