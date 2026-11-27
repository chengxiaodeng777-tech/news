$ErrorActionPreference = 'Stop'
Add-Type -AssemblyName System.IO.Compression
Add-Type -AssemblyName System.IO.Compression.FileSystem

$workspace = (Resolve-Path -LiteralPath $PSScriptRoot).Path
$sourceRoot = Join-Path $workspace 'backend-go-admin'
$zipPath = Join-Path $workspace 'backend-article-publish-linux-amd64.zip'
$sourceFiles = @(
  foreach ($name in @('go.mod', 'go.sum')) {
    Get-Item -LiteralPath (Join-Path $sourceRoot $name)
  }
  foreach ($name in @('cmd', 'ent', 'internal')) {
    Get-ChildItem -LiteralPath (Join-Path $sourceRoot $name) -Recurse -File -Force
  }
)

if (Test-Path -LiteralPath $zipPath) { Remove-Item -LiteralPath $zipPath -Force }
$archive = [System.IO.Compression.ZipFile]::Open($zipPath, [System.IO.Compression.ZipArchiveMode]::Create)
try {
  function Add-File([System.IO.FileInfo]$file, [string]$entryName) {
    [System.IO.Compression.ZipFileExtensions]::CreateEntryFromFile(
      $archive, $file.FullName, $entryName.Replace('\', '/'),
      [System.IO.Compression.CompressionLevel]::Optimal
    ) | Out-Null
  }

  Add-File (Get-Item -LiteralPath (Join-Path $sourceRoot 'release\go-admin-linux-amd64')) 'go-admin-linux-amd64'
  foreach ($file in $sourceFiles) {
    $relative = $file.FullName.Substring($sourceRoot.Length + 1)
    Add-File $file ("source/" + $relative)
  }
  Add-File (Get-Item -LiteralPath (Join-Path $workspace 'BACKEND-README.md')) 'README.md'
} finally {
  $archive.Dispose()
}

Write-Output "Created $zipPath with $($sourceFiles.Count) Go source files"
