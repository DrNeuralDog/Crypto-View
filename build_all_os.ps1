Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$ProjectRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$BinDir = Join-Path $ProjectRoot "bin"

if (-not (Test-Path $BinDir)) {
    New-Item -ItemType Directory -Path $BinDir | Out-Null
}

$targets = @(
    @{ GOOS = "windows"; GOARCH = "amd64"; Output = "cryptoview-windows-amd64.exe"; Tags = ""; CGO = "" },
    @{ GOOS = "linux"; GOARCH = "amd64"; Output = "cryptoview-linux-amd64"; Tags = "ci"; CGO = "0" },
    @{ GOOS = "darwin"; GOARCH = "amd64"; Output = "cryptoview-darwin-amd64"; Tags = ""; CGO = "1" }
)

foreach ($target in $targets) {
    $outputPath = Join-Path $BinDir $target.Output
    Write-Host "Building $($target.GOOS)/$($target.GOARCH) -> $outputPath"

    $env:GOOS = $target.GOOS
    $env:GOARCH = $target.GOARCH
    if ($target.CGO -ne "") {
        $env:CGO_ENABLED = $target.CGO
    } else {
        Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue
    }

    if ($target.Tags -ne "") {
        go build -tags $target.Tags -o $outputPath "$ProjectRoot\cmd\cryptoview"
    } else {
        go build -o $outputPath "$ProjectRoot\cmd\cryptoview"
    }
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Build failed for $($target.GOOS)/$($target.GOARCH)"
        exit $LASTEXITCODE
    }
}

Remove-Item Env:GOOS -ErrorAction SilentlyContinue
Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue

Write-Host "Cross-platform build completed successfully."
