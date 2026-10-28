# Wraith Shell Hook for Bash
# Rule 1: Preserves $? exactly
# Rule 2: Read-only observation
# Rule 3: Non-blocking execution within budget

_wraith_precmd() {
    local saved_exit=$?
    local last_cmd
    last_cmd=$(HISTTIMEFORMAT= history 1 | sed -e 's/^[ ]*[0-9]*[ ]*//')

    if [[ -n "$last_cmd" && ! "$last_cmd" =~ ^[[:space:]]*wraith ]]; then
        "__WRAITH_BIN__" hook ingest --shell bash --last-exit "$saved_exit" --cmd "$last_cmd" --port "__WRAITH_PORT__" --budget-ms "__WRAITH_BUDGET_MS__" 2>/dev/null || true
    fi

    return $saved_exit
}

if [[ ! "$PROMPT_COMMAND" =~ _wraith_precmd ]]; then
    if [[ -z "$PROMPT_COMMAND" ]]; then
        PROMPT_COMMAND="_wraith_precmd"
    else
        PROMPT_COMMAND="_wraith_precmd; $PROMPT_COMMAND"
    fi
fi
