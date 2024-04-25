param(
    [string]$repoowner = "iflytek",
    [string]$reponame = "spark-ai-cli",
    [string]$toolname = "aispark",
    [string]$toolsymlink = "aispark",
    [string]$region = "cn",
    [switch]$help
)

# Check if running as Administrator and restart with elevated permissions if not
function Test-Admin {
    $user = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($user)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

if (-not (Test-Admin)) {
    Write-Host "Restarting script with administrator privileges..."
    $script = $MyInvocation.MyCommand.Definition
    $scriptPath = Split-Path -Parent $script
    Start-Process PowerShell -ArgumentList "-NoProfile -ExecutionPolicy Bypass -File `"$script`"" -Verb RunAs -WorkingDirectory $scriptPath
    exit
}

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

# Define the temporary directory for all intermediate files
$tempDir = Join-Path $env:TEMP "aispark"
New-Item -ItemType Directory -Path $tempDir -Force | Out-Null


function Expand-GZipFile {
    param(
        [string]$inputPath,    # Input .gz file path
        [string]$outputPath    # Output file path after decompression
    )

    # Build the complete input and output paths
    $fullInputPath = Join-Path $tempDir $inputPath
    $fullOutputPath = Join-Path $tempDir $outputPath

    # Ensure the output directory exists, create it if it doesn't
    $outputDir = Split-Path -Path $fullOutputPath -Parent
    if (-Not (Test-Path -Path $outputDir)) {
        New-Item -ItemType Directory -Path $outputDir -Force | Out-Null
        Write-Host "Output directory created: $outputDir"
    }

    try {
        # Open the .gz file stream
        if (-Not (Test-Path $fullInputPath)) {
            Write-Error "Input file does not exist: $fullInputPath"
            return
        }

        $inputStream = [System.IO.File]::OpenRead($fullInputPath)
        $outputStream = [System.IO.File]::Create($fullOutputPath)
        $gzipStream = New-Object System.IO.Compression.GZipStream($inputStream, [System.IO.Compression.CompressionMode]::Decompress)

        # Copy the compressed data to the output file
        $gzipStream.CopyTo($outputStream)

        # Close all streams
        $gzipStream.Dispose()
        $outputStream.Dispose()
        $inputStream.Dispose()

        Write-Output "File has been decompressed to '$outputPath'."
    }
    catch {
        Write-Error "An error occurred: $_"
    }
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

if ($env:OS -notmatch "Windows") {
    Write-Host "You are running the powershell script on a non-windows platform. Please use the install.sh script instead."
}

# Fetch the latest release tag from GitHub API
$API_URL = "https://api.github.com/repos/$repoowner/$reponame/releases/latest"
$LATEST_TAG = (Invoke-RestMethod -Uri $API_URL).tag_name

# Set the download URL based on the platform and latest release tag
if ($region -eq "cn"){
    $DOWNLOAD_URL = "https://521github.com/$repoowner/$reponame/releases/download/$LATEST_TAG/${toolname}-${OS}-${ARCH}.exe.gz"
}else{
    $DOWNLOAD_URL = "https://github.com/$repoowner/$reponame/releases/download/$LATEST_TAG/${toolname}-${OS}-${ARCH}.exe.gz"
}

Write-Host $DOWNLOAD_URL


# Download the file
$downloadedFilePath = Join-Path $tempDir "${toolname}.exe.gz"
Invoke-WebRequest -Uri $DOWNLOAD_URL -OutFile $downloadedFilePath

# Extract the GZip file
$extractedFilePath = "${toolname}.exe"
Expand-GZipFile -inputPath "${toolname}.exe.gz" -outputPath $extractedFilePath

# check if the file already exists
$toolDir = "C:\Program Files\aispark"
$toolPath = ${toolDir}+"\${toolsymlink}.exe"
if (Test-Path $toolDir) {
    if (Test-Path $toolPath){
        Write-Host "delete old file"
        Remove-Item $toolPath
    }
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
Move-Item "${tempDir}/${toolname}.exe" $toolPath

# Clean up
Remove-Item -Recurse -Force "${tempDir}"

# Print success message
Write-Host "The $toolname has been installed successfully (version: $LATEST_TAG)."