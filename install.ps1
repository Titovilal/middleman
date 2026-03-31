$ErrorActionPreference = "Stop"

$Repo = "Titovilal/context0"
$Binary = "ctx.exe"
$Asset = "ctx-windows-amd64.exe"

# Get latest release tag
$Release = Invoke-RestMethod "https://api.github.com/repos/$Repo/releases/latest"
$Tag = $Release.tag_name
if (-not $Tag) {
    Write-Error "Could not find latest release"
    exit 1
}

$Url = "https://github.com/$Repo/releases/download/$Tag/$Asset"
$InstallDir = "$env:LOCALAPPDATA\ctx"
$InstallPath = "$InstallDir\$Binary"

Write-Host "Downloading ctx $Tag (windows/amd64)..."
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
Invoke-WebRequest -Uri $Url -OutFile $InstallPath

# Add to PATH if not already there
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
    Write-Host "Added $InstallDir to your PATH."
}

Write-Host "Done. Restart your terminal and run 'ctx' to get started."
