param(
 [Parameter(Mandatory=$true)][ValidatePattern('^v?[0-9]+\.[0-9]+\.[0-9]+(?:-[A-Za-z0-9.-]+)?$')][string]$Version,
 [string]$Executable = 'build/bin/ElmoArchive.exe'
)
$ErrorActionPreference='Stop'
$repoRoot=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
Set-Location -LiteralPath $repoRoot
$releaseRoot=Join-Path $repoRoot "build/release/$Version"
if (Test-Path -LiteralPath $releaseRoot) { throw 'This release folder already exists. Choose a new version; no files were replaced.' }
$exePath=(Resolve-Path -LiteralPath $Executable).Path
New-Item -ItemType Directory -Path $releaseRoot | Out-Null
$package=Join-Path $releaseRoot 'ElmoArchive'
New-Item -ItemType Directory -Path $package | Out-Null
Copy-Item -LiteralPath $exePath -Destination (Join-Path $package 'ElmoArchive.exe')
Copy-Item -LiteralPath LICENSE,README.md -Destination $package
Copy-Item -LiteralPath docs -Destination (Join-Path $package 'docs') -Recurse
$source=Join-Path $releaseRoot 'source'
node scripts/release-preflight.mjs --export $source
if ($LASTEXITCODE -ne 0) { throw 'Source preflight failed; do not distribute this folder.' }
Compress-Archive -LiteralPath $package -DestinationPath (Join-Path $releaseRoot "ElmoArchive-$Version-windows-x64.zip")
# Compress-Archive can omit dotfiles: use .NET to include .github and .gitignore.
Add-Type -AssemblyName System.IO.Compression.FileSystem
[System.IO.Compression.ZipFile]::CreateFromDirectory($source,(Join-Path $releaseRoot "ElmoArchive-$Version-source.zip"))
Get-ChildItem -LiteralPath $releaseRoot -Filter '*.zip' | ForEach-Object {
 $hash=Get-FileHash -LiteralPath $_.FullName -Algorithm SHA256
 "$($hash.Hash.ToLowerInvariant())  $($_.Name)"
} | Set-Content -LiteralPath (Join-Path $releaseRoot 'SHA256SUMS.txt')
Write-Output "Release candidate ready: $releaseRoot"

