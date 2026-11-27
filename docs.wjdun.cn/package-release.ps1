$ErrorActionPreference = 'Stop'
Add-Type -AssemblyName System.IO.Compression.FileSystem
Add-Type -AssemblyName System.IO.Compression

$siteRoot = (Resolve-Path -LiteralPath $PSScriptRoot).Path
$zipPath = Join-Path $siteRoot 'frontend-article-publish.zip'
$rootFiles = @(
  'index.html', '404.html', 'favicon.ico', '.htaccess', '.user.ini',
  '_app-config-5.7.0-7362080d.js', '_app-config-5.7.0-3afe33e7.js'
)
$assetDirectories = @('.well-known', 'css', 'js', 'jse')

$files = @(
  foreach ($relative in $rootFiles) {
    $target = Join-Path $siteRoot $relative
    if (Test-Path -LiteralPath $target -PathType Leaf) { Get-Item -LiteralPath $target -Force }
  }
  foreach ($relative in $assetDirectories) {
    Get-ChildItem -LiteralPath (Join-Path $siteRoot $relative) -Recurse -File -Force
  }
)

if (Test-Path -LiteralPath $zipPath) { Remove-Item -LiteralPath $zipPath -Force }
$archive = [System.IO.Compression.ZipFile]::Open($zipPath, [System.IO.Compression.ZipArchiveMode]::Create)
try {
  foreach ($file in $files) {
    $entryName = $file.FullName.Substring($siteRoot.Length + 1).Replace('\', '/')
    [System.IO.Compression.ZipFileExtensions]::CreateEntryFromFile(
      $archive, $file.FullName, $entryName,
      [System.IO.Compression.CompressionLevel]::Optimal
    ) | Out-Null
  }
} finally {
  $archive.Dispose()
}

Write-Output "Created $zipPath with $($files.Count) files"
