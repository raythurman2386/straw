#Requires -Version 5.1
<#
.SYNOPSIS
    Installs Straw — File Automation System on Windows.

.DESCRIPTION
    Downloads the latest (or specified) release of Straw from GitHub,
    extracts the binaries, and optionally adds them to the user PATH.

.PARAMETER Version
    Specific version tag to install (e.g. "v0.1.0"). Defaults to latest.

.PARAMETER InstallDir
    Directory to install binaries into. Defaults to "$env:LOCALAPPDATA\straw\bin".

.PARAMETER NoPath
    Skip adding the install directory to the user PATH.

.EXAMPLE
    irm https://raw.githubusercontent.com/raythurman2386/straw/main/install.ps1 | iex

.EXAMPLE
    .\install.ps1 -Version v0.2.0

.EXAMPLE
    .\install.ps1 -InstallDir "C:\tools\straw" -NoPath
#>

param(
    [string]$Version = "",
    [string]$InstallDir = "",
    [switch]$NoPath
)

$ErrorActionPreference = "Stop"

$repo = "raythurman2386/straw"

if (-not $InstallDir) {
    $InstallDir = Join-Path $env:LOCALAPPDATA "straw\bin"
}

# --- Detect architecture ---
$arch = switch ($env:PROCESSOR_ARCHITECTURE) {
    "AMD64"   { "amd64" }
    "x86"     { "386" }
    "ARM64"   { "arm64" }
    default   {
        Write-Error "Unsupported architecture: $env:PROCESSOR_ARCHITECTURE"
        exit 1
    }
}

Write-Host "Detected: windows $arch" -ForegroundColor Cyan

# --- Resolve version ---
if (-not $Version) {
    Write-Host "Fetching latest release..."
    try {
        $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$repo/releases/latest" -UseBasicParsing
        $Version = $release.tag_name
    } catch {
        # Fall back to pre-releases
        try {
            $releases = Invoke-RestMethod -Uri "https://api.github.com/repos/$repo/releases" -UseBasicParsing
            $Version = $releases[0].tag_name
        } catch {
            Write-Error "Could not determine latest version. Use -Version to specify one."
            exit 1
        }
    }
}

$versionNum = $Version.TrimStart("v")
Write-Host "Installing straw $Version..." -ForegroundColor Green

# --- Download ---
$fileName = "straw_${versionNum}_windows_${arch}.zip"
$url = "https://github.com/$repo/releases/download/$Version/$fileName"
$tmpDir = Join-Path $env:TEMP "straw_install_$(Get-Random)"
$zipPath = Join-Path $tmpDir $fileName

New-Item -ItemType Directory -Path $tmpDir -Force | Out-Null

Write-Host "Downloading from GitHub releases..."
try {
    Invoke-WebRequest -Uri $url -OutFile $zipPath -UseBasicParsing
} catch {
    Write-Error "Download failed. Check that version '$Version' exists for windows/$arch at:`n  $url"
    Remove-Item -Recurse -Force $tmpDir -ErrorAction SilentlyContinue
    exit 1
}

# --- Extract ---
Write-Host "Extracting..."
Expand-Archive -Path $zipPath -DestinationPath $tmpDir -Force

# Verify binaries exist
$strawExe = Get-ChildItem -Path $tmpDir -Recurse -Filter "straw.exe" | Select-Object -First 1
$strawdExe = Get-ChildItem -Path $tmpDir -Recurse -Filter "strawd.exe" | Select-Object -First 1

if (-not $strawExe -or -not $strawdExe) {
    Write-Error "Expected binaries (straw.exe, strawd.exe) not found in archive."
    Remove-Item -Recurse -Force $tmpDir -ErrorAction SilentlyContinue
    exit 1
}

# --- Install binaries ---
New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
Copy-Item -Path $strawExe.FullName -Destination $InstallDir -Force
Copy-Item -Path $strawdExe.FullName -Destination $InstallDir -Force

Write-Host "Binaries installed to $InstallDir" -ForegroundColor Green

# --- Add to PATH ---
if (-not $NoPath) {
    $userPath = [System.Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -notlike "*$InstallDir*") {
        [System.Environment]::SetEnvironmentVariable("Path", "$userPath;$InstallDir", "User")
        $env:Path = "$env:Path;$InstallDir"
        Write-Host "Added $InstallDir to user PATH" -ForegroundColor Green
    } else {
        Write-Host "$InstallDir is already in PATH" -ForegroundColor Yellow
    }
}

# --- Setup config ---
$configDir = Join-Path $env:APPDATA "straw"
$configFile = Join-Path $configDir "config.toml"

if (-not (Test-Path $configFile)) {
    New-Item -ItemType Directory -Path $configDir -Force | Out-Null

    # Check if config.example.toml was in the archive
    $exampleConfig = Get-ChildItem -Path $tmpDir -Recurse -Filter "config.example.toml" | Select-Object -First 1
    if ($exampleConfig) {
        Copy-Item -Path $exampleConfig.FullName -Destination $configFile
        Write-Host "Created default config at $configFile" -ForegroundColor Green
    } else {
        # Create a minimal config
        $downloadsPath = Join-Path $env:USERPROFILE "Downloads"
        $minimalConfig = @"
# Straw Configuration
# See https://github.com/raythurman2386/straw for documentation

[tui]
theme = "default"

[[watch]]
path = "~/Downloads"
recursive = true

[[rules]]
name = "Organize PDFs"
enabled = true
[rules.match]
extension = ".pdf"
[[rules.actions]]
type = "move"
target = "~/Documents/PDFs"
"@
        Set-Content -Path $configFile -Value $minimalConfig
        Write-Host "Created minimal config at $configFile" -ForegroundColor Green
    }
    Write-Host "Edit your config: $configFile" -ForegroundColor Yellow
} else {
    Write-Host "Config already exists at $configFile" -ForegroundColor Yellow
}

# --- Cleanup ---
Remove-Item -Recurse -Force $tmpDir -ErrorAction SilentlyContinue

# --- Done ---
Write-Host ""
Write-Host "---------------------------------------------------" -ForegroundColor Cyan
Write-Host "  Installation complete!" -ForegroundColor Green
Write-Host "---------------------------------------------------" -ForegroundColor Cyan
Write-Host ""
Write-Host "Binaries installed:"
Write-Host "  - straw.exe  (TUI client)"
Write-Host "  - strawd.exe (daemon)"
Write-Host ""
Write-Host "Configuration: $configFile"
Write-Host ""
Write-Host "To get started:"
Write-Host "  1. Edit your config:  notepad $configFile"
Write-Host "  2. Start the daemon:  strawd"
Write-Host "     (To run at startup, add a shortcut to existing 'strawd.exe' in your 'Startup' folder)"
Write-Host "  3. Start the TUI:     straw"

if (-not $NoPath) {
    Write-Host ""
    Write-Host "NOTE: You may need to restart your terminal for PATH changes to take effect." -ForegroundColor Yellow
}
