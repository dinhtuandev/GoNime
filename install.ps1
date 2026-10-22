# Install GoNime on Windows without cloning the repository.
#   irm https://raw.githubusercontent.com/dinhtuandev/GoNime/main/install.ps1 | iex
# Options: $env:VERSION = "v1.2.3"; $env:DIR = "C:\Tools\GoNime"
$ErrorActionPreference = "Stop"

$repo = "dinhtuandev/GoNime"
$bin = "gonime"

switch ($env:PROCESSOR_ARCHITECTURE) {
    "AMD64" { $arch = "amd64" }
    "ARM64" { $arch = "arm64" }
    default { throw "Unsupported architecture: $env:PROCESSOR_ARCHITECTURE" }
}

$version = $env:VERSION
if (-not $version) {
    $latest = Invoke-RestMethod -Uri "https://api.github.com/repos/$repo/releases/latest" -UseBasicParsing
    $version = $latest.tag_name
}
if (-not $version) { throw "Could not determine the latest release" }

$asset = "${bin}_${version}_windows_${arch}.zip"
$base = "https://github.com/$repo/releases/download/$version"
$target = if ($env:DIR) { $env:DIR } else { Join-Path $env:LOCALAPPDATA "Programs\GoNime" }

$tmp = Join-Path ([IO.Path]::GetTempPath()) ("gonime-" + [Guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Path $tmp | Out-Null

try {
    Write-Host "Downloading $asset ($version)..."
    Invoke-WebRequest -Uri "$base/$asset" -OutFile (Join-Path $tmp $asset) -UseBasicParsing
    Invoke-WebRequest -Uri "$base/checksums.txt" -OutFile (Join-Path $tmp "checksums.txt") -UseBasicParsing

    $entry = Select-String -Path (Join-Path $tmp "checksums.txt") -Pattern " $asset$" | Select-Object -First 1
    if (-not $entry) { throw "No checksum entry for $asset" }
    $expected = ($entry.Line -split "\s+")[0].ToLower()
    $actual = (Get-FileHash -Path (Join-Path $tmp $asset) -Algorithm SHA256).Hash.ToLower()
    if ($expected -ne $actual) { throw "Checksum mismatch for $asset" }

    Expand-Archive -Path (Join-Path $tmp $asset) -DestinationPath $tmp -Force
    $exe = Join-Path $tmp "$bin.exe"
    if (-not (Test-Path $exe)) {
        $exe = (Get-ChildItem -Path $tmp -Filter "*.exe" -Recurse | Select-Object -First 1).FullName
    }
    if (-not $exe) { throw "No executable found in $asset" }

    New-Item -ItemType Directory -Path $target -Force | Out-Null
    Copy-Item -Path $exe -Destination (Join-Path $target "$bin.exe") -Force
} finally {
    Remove-Item -Recurse -Force $tmp -ErrorAction SilentlyContinue
}

$dest = Join-Path $target "$bin.exe"
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$target*") {
    [Environment]::SetEnvironmentVariable("Path", "$userPath;$target", "User")
    Write-Host "Added $target to your user PATH (restart the terminal)."
}

Write-Host "Installed to $dest"
& $dest version
