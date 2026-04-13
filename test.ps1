# Quality-check entrypoint for CryptoView
Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$ProjectRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
Push-Location $ProjectRoot
try {
    $targets = @("cmd", "internal", "resources")
    $unformatted = @()
    foreach ($target in $targets) {
        if (Test-Path $target) {
            $unformatted += @(gofmt -l $target)
        }
    }
    $unformatted = @($unformatted | Where-Object { -not [string]::IsNullOrWhiteSpace($_) })
    if ($unformatted.Count -gt 0) {
        Write-Host "The following files are not gofmt-formatted:"
        $unformatted | ForEach-Object { Write-Host "  $_" }
        throw "gofmt check failed"
    }

    go vet ./...
    if ($LASTEXITCODE -ne 0) {
        throw "go vet failed with exit code $LASTEXITCODE"
    }

    go test ./...
    if ($LASTEXITCODE -ne 0) {
        throw "go test failed with exit code $LASTEXITCODE"
    }

    Write-Host "`nQuality checks passed."
}
finally {
    Pop-Location
}
