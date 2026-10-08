#!/usr/bin/env bash

# Read stdin
CURRENT_RELEASE=$(cat -)

# Output 1 if there's no existing release, and increment
# release number if releases already exist
if [ "x$CURRENT_RELEASE" == "x" ]; then
  echo "1"
else
  NEW_RELEASE=$(expr $CURRENT_RELEASE + 1)
  echo "$NEW_RELEASE"
fi
