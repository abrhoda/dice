#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
if [[ "${TRACE-0}" == "1" ]]; then
    set -o xtrace
fi

if [[ "${1-}" =~ ^-*h(elp)?$ ]]; then
    echo 'Usage: bash ./coverage.bash <path to coverage.out file> <coverage threshold>

This script returns a 0 or 1 depending on if the current unit test coverage is greater than a threshold.
'
    exit
fi

cd "$(dirname "$0")"

main() {
  cd "$(pwd)/.."

  # TOTAL will be a decimal %
  TOTAL=$(go tool cover -func="$1" | grep total | awk '{print substr($3, 1, length($3)-1)}')
  

  if [[ `echo "$TOTAL $2" | awk '{print ($1 > $2)}'` == 1 ]];then
    echo "$TOTAL is greater than $2"
    return 0
  else
    echo "$TOTAL is less than $2"
    return 1
  fi
}

main "$@"
