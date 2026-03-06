$root = Split-Path -Parent $PSScriptRoot
$pidFile = "$root\scripts\pids.json"
if (Test-Path $pidFile) {
    try {
        $content = Get-Content $pidFile -Raw | ConvertFrom-Json
        foreach ($key in $content.PSObject.Properties.Name) {
            $id = $content.$key
            Write-Host "Stopping $key (PID $id) ..."
            Stop-Process -Id $id -ErrorAction SilentlyContinue
        }
        Remove-Item $pidFile -ErrorAction SilentlyContinue
        Write-Host "Stopped background processes."
    } catch {
        Write-Warning "Failed to read/stop processes from $pidFile: $_"
    }
} else {
    Write-Host "PID file not found: $pidFile. No background PIDs to stop."
}

Write-Host "Stopping docker-compose services (backend/docker-compose.yml) ..."
docker-compose -f "$root\backend\docker-compose.yml" down

Write-Host "Done."
