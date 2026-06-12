#Requires -Version 5.1
<#
.SYNOPSIS
    Build / test the go-yespower cgo bindings on Windows.

    Compatible with Windows PowerShell 5.1 and PowerShell 7+.

.DESCRIPTION
    The cgo path (Hash / HashVersion) needs a C compiler. This script makes the
    MSYS2 UCRT64 gcc available to the `go` toolchain, then runs the test suite
    and/or the example.

    A C compiler is only located if one is not already on PATH. When it is
    needed, the UCRT64 bin directory is searched in this order:

        1. $env:UCRT64                  (direct path to the ucrt64 dir)
        2. $env:MSYS2_ROOT\ucrt64       (MSYS2 install root)
        3. C:\msys64\ucrt64            (default install location)
        4. Registry                     (MSYS2 64bit uninstall InstallLocation)

.PARAMETER Test
    Run `go test -v ./...`. Enabled by default (also when no flag is given).

.PARAMETER Example
    Run `go run ./example`.

.EXAMPLE
    .\build.ps1                 # runs the tests (default)

.EXAMPLE
    .\build.ps1 -Example        # runs only the example

.EXAMPLE
    .\build.ps1 -Test -Example  # runs both
#>
[CmdletBinding()]
param(
    [switch]$Test,
    [switch]$Example
)

$ErrorActionPreference = 'Stop'

# Default action: if neither flag is given, run the tests.
if (-not $Test -and -not $Example) {
    $Test = $true
}

# Run from the directory this script lives in (the go module root).
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Push-Location $ScriptDir

try {
    # ------------------------------------------------------------------
    # Ensure a C compiler is available (cgo requirement).
    # ------------------------------------------------------------------
    if (Get-Command gcc -ErrorAction SilentlyContinue) {
        Write-Host 'gcc already on PATH; skipping MSYS2 lookup.' -ForegroundColor DarkGray
    } else {
        Write-Host 'gcc not found on PATH; locating MSYS2 UCRT64...' -ForegroundColor Cyan

        # Build the ordered list of candidate ucrt64\bin directories.
        $candidates = @()

        if ($env:UCRT64) {
            $candidates += (Join-Path $env:UCRT64 'bin')
        }
        if ($env:MSYS2_ROOT) {
            $candidates += (Join-Path $env:MSYS2_ROOT 'ucrt64\bin')
        }
        $candidates += 'C:\msys64\ucrt64\bin'

        # Registry fallback: the MSYS2 installer records its InstallLocation.
        $regPaths = @(
            'HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\MSYS2 64bit',
            'HKLM:\SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall\MSYS2 64bit'
        )
        foreach ($rp in $regPaths) {
            try {
                $loc = (Get-ItemProperty -Path $rp -Name InstallLocation -ErrorAction Stop).InstallLocation
                if ($loc) {
                    $candidates += (Join-Path $loc 'ucrt64\bin')
                }
            } catch {
                # key absent — ignore and try the next source
            }
        }

        # Pick the first candidate that actually contains gcc.exe.
        $ucrtBin = $null
        foreach ($c in $candidates) {
            if ($c -and (Test-Path (Join-Path $c 'gcc.exe'))) {
                $ucrtBin = $c
                break
            }
        }

        if (-not $ucrtBin) {
            throw @"
Could not locate an MSYS2 UCRT64 gcc.exe.
Searched:
$($candidates -join "`n")

Install MSYS2 (https://www.msys2.org) and the toolchain:
    pacman -S mingw-w64-ucrt-x86_64-gcc
Then re-run, optionally setting `$env:MSYS2_ROOT` or `$env:UCRT64`.
"@
        }

        Write-Host "Using gcc from: $ucrtBin" -ForegroundColor Green
        $env:PATH = "$ucrtBin;$env:PATH"
    }

    # cgo must be enabled to compile the C sources in ../yespower-c.
    $env:CGO_ENABLED = '1'

    if ($Test) {
        Write-Host 'Running: go test -v ./...' -ForegroundColor Cyan
        & go test -v ./...
        if ($LASTEXITCODE -ne 0) { throw "go test failed (exit $LASTEXITCODE)" }
    }

    if ($Example) {
        Write-Host 'Running: go run ./example' -ForegroundColor Cyan
        & go run ./example
        if ($LASTEXITCODE -ne 0) { throw "go run ./example failed (exit $LASTEXITCODE)" }
    }
}
finally {
    Pop-Location
}
