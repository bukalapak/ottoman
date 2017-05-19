#!/bin/bash

OUTFILE=$1
SKIPPED=(
  github.com/bukalapak/ottoman/qtest
  github.com/bukalapak/ottoman
)

_SKIPPED=$(printf "|%s" "${SKIPPED[@]}")
PACKAGES=`go list ./... | grep -v /vendor/ | grep -E -v "${_SKIPPED:1}$"`
EXITCODE=0

echo "mode: atomic" > $OUTFILE

for PKG in $PACKAGES; do
  echo ======= $PKG

  go test -race -v -coverprofile=profile.out -covermode=atomic $PKG; __EXITCODE__=$?

  if [ "$__EXITCODE__" -ne "0" ]; then
    EXITCODE=$__EXITCODE__
  fi

  if [ -f profile.out ]; then
    tail -n +2 profile.out >> $OUTFILE; rm profile.out
  fi
done

exit $EXITCODE

