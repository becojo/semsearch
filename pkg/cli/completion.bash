#!/bin/bash

# Function to extract metavariables from command line arguments
_extract_metavariables() {
    local metavars=()
    local i=0

    # Look through all command words for pattern arguments
    for ((i=1; i<${#COMP_WORDS[@]}; i++)); do
        local word="${COMP_WORDS[i]}"
        local prev_word=""
        if ((i > 0)); then
            prev_word="${COMP_WORDS[i-1]}"
        fi

        # Check if this word is a pattern argument value
        case "${prev_word}" in
            --pattern|-p|--pattern-inside|-pi|--pattern-not|-pn|--pattern-not-inside|-pni|--pattern-regex|-pr|--pattern-not-regex|-pnr)
                # Extract metavariables from the pattern (supports $IDENTIFIER and $...IDENTIFIER)
                local vars=$(echo "${word}" | grep -Eo '\$([.]{3})?[A-Z_][A-Z0-9_]*' | sed 's/\$//' | sort -u)
                if [[ -n "${vars}" ]]; then
                    metavars+=(${vars})
                fi
                ;;
        esac
    done

    # Remove duplicates and return
    printf '%s\n' "${metavars[@]}" | sort -u
}

_semsearch_completion() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    # Flags that don't take arguments
    local flags0="--autofix --debug --export --pattern-either --pattern-sinks --pattern-sources --patterns --pop --rule --semgrep --verbose"

    # Flags that take arguments
    local flags1="--config --eval --fix --fix-regex --focus-metavariable --format --id --language --message --metadata --metavariable-pattern --metavariable-regex --option --path --path-exclude --path-include --pattern --pattern-inside --pattern-not --pattern-not-inside --pattern-not-regex --pattern-regex --severity"

    # Format options
    local formats="yaml json sarif text emacs vim github-actions gitlab-sast gitlab-secrets junit-xml"

    # Language options (common ones)
    local languages="bash c cpp csharp dockerfile generic go java javascript json php python ruby rust scala terraform typescript yaml"

    # Severity options
    local severities="INFO WARNING ERROR"

    # Rule options
    local option_keys="generic_ellipsis_max_span"

    # Short flags and their expansions
    local shortcuts=(
        "-af" "--autofix"
        "-c" "--config"
        "-e" "--eval"
        "-f" "--format"
        "-fm" "--focus-metavariable"
        "-fr" "--fix-regex"
        "-fx" "--fix"
        "-i" "--path"
        "-l" "--language"
        "-m" "--message"
        "-mp" "--metavariable-pattern"
        "-mr" "--metavariable-regex"
        "-p" "--pattern"
        "-pe" "--pattern-either"
        "-pi" "--pattern-inside"
        "-pn" "--pattern-not"
        "-pni" "--pattern-not-inside"
        "-pnr" "--pattern-not-regex"
        "-pr" "--pattern-regex"
        "-ps" "--patterns"
        "-psk" "--pattern-sinks"
        "-pso" "--pattern-sources"
        "-sv" "--severity"
    )

    # Check if the previous word expects a value
    case "${prev}" in
        --config|--path|-i|--path-include|--path-exclude)
            # Complete with file/directory paths
            COMPREPLY=( $(compgen -f "${cur}") )
            return 0
            ;;
        --format|-f)
            COMPREPLY=( $(compgen -W "${formats}" -- ${cur}) )
            return 0
            ;;
        --language|-l)
            COMPREPLY=( $(compgen -W "${languages}" -- ${cur}) )
            return 0
            ;;
        --severity|-sv)
            COMPREPLY=( $(compgen -W "${severities}" -- ${cur}) )
            return 0
            ;;
        --pattern|-p|--pattern-inside|-pi|--pattern-not|-pn|--pattern-not-inside|-pni|--pattern-regex|-pr|--pattern-not-regex|-pnr|--eval|-e|--fix|-fx|--fix-regex|-fr|--id|--message|-m)
            # These expect pattern/code strings - no completion
            return 0
            ;;
        --metavariable-pattern|-mp|--focus-metavariable|-fm)
            # Complete with metavariables extracted from patterns (names only)
            local metavars=($(_extract_metavariables))
            if [[ ${#metavars[@]} -gt 0 ]]; then
                COMPREPLY=( $(compgen -W "${metavars[*]}" -- ${cur}) )
            fi
            return 0
            ;;
        --metavariable-regex|-mr)
            # Complete with metavariables extracted from patterns
            local metavars=($(_extract_metavariables))
            if [[ ${#metavars[@]} -gt 0 ]]; then
                # Add = to each metavariable for key=value completion
                local completions=()
                for metavar in "${metavars[@]}"; do
                    completions+=("${metavar}=")
                done
                COMPREPLY=( $(compgen -W "${completions[*]}" -- ${cur}) )
            fi
            return 0
            ;;
        --metadata)
            # These expect key=value pairs - no completion
            return 0
            ;;
        --option)
            # Complete with common option keys
            local completions=()
            for option_key in ${option_keys}; do
                completions+=("${option_key}=")
            done
            COMPREPLY=( $(compgen -W "${completions[*]}" -- ${cur}) )
            return 0
            ;;
    esac

    # If current word starts with -, complete with flags
    if [[ ${cur} == -* ]]; then
        local all_flags="${flags0} ${flags1}"

        # Add short flags
        for ((i=0; i<${#shortcuts[@]}; i+=2)); do
            all_flags="${all_flags} ${shortcuts[i]}"
        done

        COMPREPLY=( $(compgen -W "${all_flags}" -- ${cur}) )
        return 0
    fi

    # Special commands
    if [[ ${cur} == h* ]]; then
        COMPREPLY=( $(compgen -W "help" -- ${cur}) )
        return 0
    fi

    # Default to file completion for non-flag arguments
    COMPREPLY=( $(compgen -f "${cur}") )
    return 0
}

# Register the completion function
complete -o nospace -F _semsearch_completion semsearch