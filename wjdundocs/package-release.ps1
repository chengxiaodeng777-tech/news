$ErrorActionPreference = 'Stop'
Add-Type -AssemblyName System.IO.Compression
Add-Type -AssemblyName System.IO.Compression.FileSystem

$packageRoot = (Resolve-Path -LiteralPath $PSScriptRoot).Path
$archivePath = Join-Path $packageRoot 'backend-article-upload-linux-amd64.zip'
$files = @(
  Get-Item -LiteralPath (Join-Path $packageRoot 'go-admin-linux-amd64')
  Get-Item -LiteralPath (Join-Path $packageRoot 'README.md')
  Get-ChildItem -LiteralPath (Join-Path $packageRoot 'source') -Recurse -File
)

if (Test-Path -LiteralPath $archivePath) {
  Remove-Item -LiteralPath $archivePath -Force
}
$archive = [System.IO.Compression.ZipFile]::Open($archivePath, [System.IO.Compression.ZipArchiveMode]::Create)
try {
  foreach ($file in $files) {
    $relativePath = $file.FullName.Substring($packageRoot.Length + 1).Replace('\', '/')
    [System.IO.Compression.ZipFileExtensions]::CreateEntryFromFile(
      $archive, $file.FullName, $relativePath,
      [System.IO.Compression.CompressionLevel]::Optimal
    ) | Out-Null
  }
} finally {
  $archive.Dispose()
}

Write-Output "Created $archivePath with $($files.Count) files"
