<#
Wraith Shell Hook for PowerShell 7 and Windows PowerShell 5.1
Rule 1: Preserves $LASTEXITCODE and $? exactly
Rule 2: Read-only observation
Rule 3: Non-blocking execution within budget
#>
function Invoke-WraithHook {
    # CRITICAL: Capture exit state before anything else executes
    $savedExit = $LASTEXITCODE
    $savedSuccess = $?

    try {
        # Determine the last command executed
        $lastCmd = ""
        try {
            $history = Get-History -Count 1 -ErrorAction SilentlyContinue
            if ($null -ne $history) {
                $lastCmd = $history.CommandLine
            }
        } catch {}

        # If no command line or internal wraith command, do nothing
        if ([string]::IsNullOrWhiteSpace($lastCmd) -or $lastCmd -match '^\s*wraith(\.exe)?(\s|$)') {
            return
        }

        # Auto-start daemon in background if not responding
        $wraithBin = '__WRAITH_BIN__'
        $wraithPort = '__WRAITH_PORT__'
        $budgetMs = '__WRAITH_BUDGET_MS__'

        # Invoke wraith hook directly with timeout; prints whisper if triggered
        # Output is received directly onto the terminal
        & "$wraithBin" hook ingest --shell pwsh --last-exit $savedExit --cmd "$lastCmd" --port $wraithPort --budget-ms $budgetMs 2>$null
    } catch {
        # Fail silent: never break prompt
    } finally {
        # CRITICAL RULE 1: Restore exact exit code and status
        if ($null -ne $savedExit) {
            $global:LASTEXITCODE = $savedExit
        }
        # Force $? to match previous success state
        if (-not $savedSuccess) {
            Write-Error "WraithExitCodeRestore" -ErrorAction SilentlyContinue 2>$null
            if ($null -ne $savedExit) {
                $global:LASTEXITCODE = $savedExit
            }
        }
    }
}

# Attach hook to PowerShell prompt without destroying existing prompt
if (Test-Path Function:\prompt) {
    $script:OriginalWraithPrompt = (Get-Item Function:\prompt).ScriptBlock
} else {
    $script:OriginalWraithPrompt = { "PS $($executionContext.SessionState.Path.CurrentLocation)$('>' * ($nestedPromptLevel + 1)) " }
}

function Global:prompt {
    # 1. Run Wraith hook
    Invoke-WraithHook

    # 2. Call original prompt
    & $script:OriginalWraithPrompt
}
