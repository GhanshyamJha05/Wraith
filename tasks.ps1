<#
.SYNOPSIS
    Build and automation script for Project Wraith on Windows.
.DESCRIPTION
    Provides targets: build, test, eval, demo, clean.
#>
param(
    [Parameter(Position=0)]
    [ValidateSet("build", "test", "eval", "demo", "clean")]
    [string]$Target = "build"
)

$ErrorActionPreference = "Stop"

function Build-Wraith {
    Write-Host "[wraith] Building binary..." -ForegroundColor Cyan
    $binDir = Join-Path $PSScriptRoot "bin"
    if (-not (Test-Path $binDir)) {
        New-Item -ItemType Directory -Path $binDir -Force | Out-Null
    }
    go build -trimpath -ldflags="-s -w" -o "$binDir\wraith.exe" .\cmd\wraith
    Write-Host "[wraith] Built: $binDir\wraith.exe" -ForegroundColor Green
}

function Test-Wraith {
    Write-Host "[wraith] Running unit & integration tests..." -ForegroundColor Cyan
    go test -v .\...
}

function Eval-Wraith {
    Write-Host "[wraith] Running evaluation harness..." -ForegroundColor Cyan
    go run .\eval\main.go
}

function Demo-Wraith {
    Write-Host "[wraith] Running interactive / replay demo..." -ForegroundColor Cyan
    $bin = Join-Path $PSScriptRoot "bin\wraith.exe"
    if (-not (Test-Path $bin)) {
        Build-Wraith
    }
    & $bin demo --debug
}

function Clean-Wraith {
    Write-Host "[wraith] Cleaning build artifacts..." -ForegroundColor Yellow
    $binDir = Join-Path $PSScriptRoot "bin"
    if (Test-Path $binDir) {
        Remove-Item -Recurse -Force $binDir
    }
    Write-Host "[wraith] Clean complete." -ForegroundColor Green
}

switch ($Target) {
    "build" { Build-Wraith }
    "test"  { Test-Wraith }
    "eval"  { Eval-Wraith }
    "demo"  { Demo-Wraith }
    "clean" { Clean-Wraith }
}
