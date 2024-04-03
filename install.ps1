param(
    [string]$repoowner = "iflytek",
    [string]$reponame = "spark-ai-cli",
    [string]$toolname = "aispark",
    [string]$toolsymlink = "aispark",
    [switch]$region = "cn"
    [switch]$help
)

if ($help) {
    Write-Host "aispark Installer Help!"
    Write-Host " Usage: "
    Write-Host "    aispark -help <Shows this message>"
    Write-Host "    aispark -repoowner <Owner of the repo>"
    Write-Host "    aispark -reponame <Set the repository name we will look for>"
    Write-Host "    aispark -toolname <Set the name of the tool (inside the .zip build)>"
    Write-Host "    aispark -toolsymlink <Set name of the local executable>"

    exit 0
}

# if user not admin then quit
function IsUserAdministrator {
    $user = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($user)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

if (-not (IsUserAdministrator)) {
    Write-Host "Please run as administrator"
    exit 1
}

# Detect the platform (architecture and OS)
$ARCH = $null
$OS = "windows"


if ($env:PROCESSOR_ARCHITECTURE -eq "AMD64") {
    $ARCH = "amd64"
} elseif ($env:PROCESSOR_ARCHITECTURE -eq "arm64") {
    $ARCH = "arm64"
} else {
    $ARCH = "i386"
}

if ($env:OS -notmatch "Windows") {
    Write-Host "You are running the powershell script on a non-windows platform. Please use the install.sh script instead."
}

# Fetch the latest release tag from GitHub API
$API_URL = "https://api.github.com/repos/$repoowner/$reponame/releases/latest"
$LATEST_TAG = (Invoke-RestMethod -Uri $API_URL).tag_name

# Set the download URL based on the platform and latest release tag
if ($region -eq "cn"){
    $DOWNLOAD_URL = "https://521github.com/$repoowner/$reponame/releases/download/$LATEST_TAG/${toolname}-${OS}-${ARCH}.exe"
}else{
    $DOWNLOAD_URL = "https://github.com/$repoowner/$reponame/releases/download/$LATEST_TAG/${toolname}-${OS}-${ARCH}.exe"
}

Write-Host $DOWNLOAD_URL

# Download the file
Invoke-WebRequest -Uri $DOWNLOAD_URL -OutFile "${toolname}.exe"

# Extract the ZIP file
# $extractedDir = "${toolname}-temp"
# Expand-Archive -Path "${toolname}.zip" -DestinationPath $extractedDir -Force

# check if the file already exists
$toolPath = "C:\Program Files\aispark\${toolsymlink}.exe"
if (Test-Path $toolPath) {
    Remove-Item $toolPath
} else {
    New-Item -ItemType Directory -Path "C:\Program Files\aispark\"
}

# Add the file to path
$currentPath = [System.Environment]::GetEnvironmentVariable("PATH", "User")

# Append the desired path to the current PATH value if it's not already present
if (-not ($currentPath -split ";" | Select-String -SimpleMatch "C:\Program Files\aispark\")) {
    $updatedPath = $currentPath + ";" + "C:\Program Files\aispark\"

    # Set the updated PATH value
    [System.Environment]::SetEnvironmentVariable("PATH", $updatedPath, "User")   # Use "User" instead of "Machine" for user-level PATH

    Write-Host "The path has been added to the PATH variable. You may need to restart applications to see the changes." -ForegroundColor Red
}

# Make the binary executable
Move-Item "${toolname}.exe" $toolPath
Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy Unrestricted

# Clean up

# Print success message
Write-Host "The $toolname has been installed successfully (version: $LATEST_TAG)."