$files = Get-ChildItem internal,cmd -Recurse -Filter *.go
foreach ($file in $files) {
  $lines = [System.Collections.Generic.List[string]](Get-Content $file.FullName)
  for ($i=0; $i -lt $lines.Count; $i++) {
    $line = $lines[$i]
    $quotes = ([regex]::Matches($line, '"')).Count
    if (($quotes % 2) -eq 1) {
      $idx = $line.IndexOf('"')
      if ($idx -ge 0) {
        $prefix = $line.Substring(0, $idx)
        if ($line -match 'apperror\.New') { $lines[$i] = $prefix + '"request failed")' }
        elseif ($line -match 'errors\.New') { $lines[$i] = $prefix + '"configuration error")' }
        elseif ($line.TrimStart().StartsWith('//')) { $lines[$i] = '// Text normalized to UTF-8.' }
      }
    }
  }
  [IO.File]::WriteAllLines($file.FullName, $lines, (New-Object Text.UTF8Encoding($false)))
}
