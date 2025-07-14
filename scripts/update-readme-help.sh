#!/usr/bin/env bash

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

go run ./cmd/semsearch -i README.md \
  --option generic_ellipsis_max_span=500 \
  --pattern-inside "$(matchInside)" \
  --pattern "$(match)" \
  --fix "$(replace)" \
  --autofix
