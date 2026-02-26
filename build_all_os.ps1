Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"
$HostIsWindows = ($env:OS -eq "Windows_NT")
$SkippedTargets = New-Object System.Collections.Generic.List[string]

function Invoke-NativeCommand {
    param(
        [Parameter(Mandatory = $true)]
        [string]$FilePath,
        [Parameter(Mandatory = $true)]
        [string[]]$Arguments
    )

    & $FilePath @Arguments
    if ($LASTEXITCODE -ne 0) {
        throw "Command failed with exit code ${LASTEXITCODE}: $FilePath $($Arguments -join ' ')"
    }
}

function Convert-PngToIco {
    param(
        [Parameter(Mandatory = $true)]
        [string]$PngPath,
        [Parameter(Mandatory = $true)]
        [string]$IcoPath
    )

    $pngBytes = [System.IO.File]::ReadAllBytes($PngPath)
    if ($pngBytes.Length -lt 24) {
        throw "PNG file is too small: $PngPath"
    }

    $expectedSig = [byte[]](137, 80, 78, 71, 13, 10, 26, 10)
    for ($i = 0; $i -lt $expectedSig.Length; $i++) {
        if ($pngBytes[$i] -ne $expectedSig[$i]) {
            throw "Invalid PNG signature: $PngPath"
        }
    }

    $width = [System.BitConverter]::ToUInt32([byte[]]($pngBytes[19], $pngBytes[18], $pngBytes[17], $pngBytes[16]), 0)
    $height = [System.BitConverter]::ToUInt32([byte[]]($pngBytes[23], $pngBytes[22], $pngBytes[21], $pngBytes[20]), 0)
    $iconWidth = if ($width -ge 256) { [byte]0 } else { [byte]$width }
    $iconHeight = if ($height -ge 256) { [byte]0 } else { [byte]$height }

    $stream = New-Object System.IO.MemoryStream
    $writer = New-Object System.IO.BinaryWriter($stream)
    try {
        $writer.Write([UInt16]0)
        $writer.Write([UInt16]1)
        $writer.Write([UInt16]1)
        $writer.Write($iconWidth)
        $writer.Write($iconHeight)
        $writer.Write([byte]0)
        $writer.Write([byte]0)
        $writer.Write([UInt16]1)
        $writer.Write([UInt16]32)
        $writer.Write([UInt32]$pngBytes.Length)
        $writer.Write([UInt32]22)
        $writer.Write($pngBytes)
        [System.IO.File]::WriteAllBytes($IcoPath, $stream.ToArray())
    }
    finally {
        $writer.Dispose()
        $stream.Dispose()
    }
}

$ProjectRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$AppSourceDir = Join-Path $ProjectRoot "cmd\cryptoview"
$BinDir = Join-Path $ProjectRoot "bin"
$WindowsIconPath = Join-Path $ProjectRoot "resources\Logo\CryptoView Icon.png"
$TempIcoPath = Join-Path $env:TEMP "CryptoView_build_icon_cross.ico"

$goCommand = Get-Command go -ErrorAction SilentlyContinue
if ($null -eq $goCommand) {
    throw "go command not found in PATH."
}
$GoExe = $goCommand.Source

$goBinPath = if ([string]::IsNullOrWhiteSpace($env:GOBIN)) {
    Join-Path $env:USERPROFILE "go\bin"
} else {
    $env:GOBIN
}
$RsrcExe = Join-Path $goBinPath "rsrc.exe"

if (-not (Test-Path $BinDir)) {
    New-Item -ItemType Directory -Path $BinDir | Out-Null
}

if (-not (Test-Path $WindowsIconPath)) {
    throw "Windows icon PNG not found: $WindowsIconPath"
}

if (-not (Test-Path $RsrcExe)) {
    Write-Host "Installing rsrc tool (github.com/akavel/rsrc)..."
    Invoke-NativeCommand -FilePath $GoExe -Arguments @("install", "github.com/akavel/rsrc@latest")
}

$targets = @(
    @{ GOOS = "windows"; GOARCH = "amd64"; Output = "cryptoview-windows-amd64.exe"; Tags = ""; CGO = "" },
    @{ GOOS = "linux"; GOARCH = "amd64"; Output = "cryptoview-linux-amd64"; Tags = "ci"; CGO = "0" },
    @{ GOOS = "darwin"; GOARCH = "amd64"; Output = "cryptoview-darwin-amd64"; Tags = ""; CGO = "1" }
)

foreach ($target in $targets) {
    $outputPath = Join-Path $BinDir $target.Output
    Write-Host "Building $($target.GOOS)/$($target.GOARCH) -> $outputPath"

    if ($HostIsWindows -and $target.GOOS -eq "darwin" -and $target.CGO -eq "1") {
        $SkippedTargets.Add("$($target.GOOS)/$($target.GOARCH)") | Out-Null
        Write-Warning "Skipping darwin build on Windows: Fyne macOS build requires CGO and a macOS cross-toolchain (osxcross/clang). Windows/Linux artifacts were built."
        continue
    }

    $env:GOOS = $target.GOOS
    $env:GOARCH = $target.GOARCH
    if ($target.CGO -ne "") {
        $env:CGO_ENABLED = $target.CGO
    } else {
        Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue
    }

    $sysoPath = $null
    if ($target.GOOS -eq "windows") {
        $sysoPath = Join-Path $AppSourceDir ("rsrc_windows_{0}.syso" -f $target.GOARCH)
        Convert-PngToIco -PngPath $WindowsIconPath -IcoPath $TempIcoPath
        Invoke-NativeCommand -FilePath $RsrcExe -Arguments @("-ico", $TempIcoPath, "-arch", $target.GOARCH, "-o", $sysoPath)
    }

    try {
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
    finally {
        if ($null -ne $sysoPath -and (Test-Path $sysoPath)) {
            Remove-Item $sysoPath -Force
        }
    }
}

if (Test-Path $TempIcoPath) {
    Remove-Item $TempIcoPath -Force
}

Remove-Item Env:GOOS -ErrorAction SilentlyContinue
Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue

if ($SkippedTargets.Count -gt 0) {
    Write-Host "Cross-platform build completed with skipped targets: $($SkippedTargets -join ', ')"
} else {
    Write-Host "Cross-platform build completed successfully."
}
