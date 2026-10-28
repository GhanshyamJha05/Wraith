# Wraith Shell Hook for Zsh
# Rule 1: Preserves $? and $status exactly
# Rule 2: Read-only observation
# Rule 3: Non-blocking execution within budget

_wraith_precmd() {
    local saved_exit=$?
    local last_cmd
    last_cmd=$(fc -ln -1 | sed -e 's/^[[:space:]]*//')

    if [[ -n "$last_cmd" && ! "$last_cmd" =~ ^[[:space:]]*wraith ]]; then
        "__WRAITH_BIN__" hook ingest --shell zsh --last-exit "$saved_exit" --cmd "$last_cmd" --port "__WRAITH_PORT__" --budget-ms "__WRAITH_BUDGET_MS__" 2>/dev/null || true
    fi

    return $saved_exit
}

autoload -Uz add-zsh-hook 2>/dev/null || true
if (( $+functions[add-zsh-hook] )); then
    add-zsh-hook precmd _wraith_precmd
else
    precmd_functions+=(_wraith_precmd)
fi
