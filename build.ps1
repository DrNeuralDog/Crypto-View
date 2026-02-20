param(
    [string]$Output = "cryptoview.exe"
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$ProjectRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$BinDir = Join-Path $ProjectRoot "bin"
$OutputPath = Join-Path $BinDir $Output

if (-not (Test-Path $BinDir)) {
    New-Item -ItemType Directory -Path $BinDir | Out-Null
}

Write-Host "Building CryptoView -> $OutputPath"
go build -o $OutputPath "$ProjectRoot\cmd\cryptoview"

if ($LASTEXITCODE -ne 0) {
    Write-Error "Build failed with exit code $LASTEXITCODE"
    exit $LASTEXITCODE
}

Write-Host "Build completed successfully."
