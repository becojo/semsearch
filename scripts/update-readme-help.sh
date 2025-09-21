#!/usr/bin/env bash

alias semsearch="go run ./cmd/semsearch"

function matchInside() {
  echo '<!-- help start -->'
  echo '...'
  echo '<!-- help end -->'
}

function match() {
  echo '```'
  echo '...'
  echo '```'
}

function replace() {
  echo '```'
  go run ./cmd/semsearch help
  echo '```'
}

semsearch -i README.md \
  --option generic_ellipsis_max_span=500 \
  --pattern-inside "$(matchInside)" \
  --pattern "$(match)" \
  --fix "$(replace)" \
  --autofix
