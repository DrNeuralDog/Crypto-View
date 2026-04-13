param(
    [string]$Output = "cryptoview.exe"
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

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
$OutputPath = Join-Path $BinDir $Output
$WindowsIconPath = Join-Path $ProjectRoot "resources\Logo\CryptoView Icon.png"
$TempIcoPath = Join-Path $env:TEMP "CryptoView_build_icon.ico"

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

$goArch = (& $GoExe env GOARCH).Trim()
if (-not $goArch) {
    throw "Unable to determine GOARCH."
}
$sysoPath = Join-Path $AppSourceDir ("rsrc_windows_{0}.syso" -f $goArch)

Write-Host "Building CryptoView -> $OutputPath"
Convert-PngToIco -PngPath $WindowsIconPath -IcoPath $TempIcoPath
try {
    Invoke-NativeCommand -FilePath $RsrcExe -Arguments @("-ico", $TempIcoPath, "-arch", $goArch, "-o", $sysoPath)
    Invoke-NativeCommand -FilePath $GoExe -Arguments @("build", "-o", $OutputPath, "$ProjectRoot\cmd\cryptoview")
}
finally {
    if (Test-Path $TempIcoPath) {
        Remove-Item $TempIcoPath -Force
    }
    if (Test-Path $sysoPath) {
        Remove-Item $sysoPath -Force
    }
}

Write-Host "Build completed successfully."
