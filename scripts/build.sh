#!/bin/bash

# exit if fails
set -e

####
# function
####

xbuild() {
  # build windows amd64
  env GOOS=windows GOARCH=amd64 go build "$CMD_PATH"/"$FOLDER"/"$FOLDER".go
  mv "$FOLDER".exe "$BIN_PATH"/"$FOLDER"_amd64.exe

  echo Built "$BIN_PATH"/"$FOLDER"_amd64.exe

  # build macos arm64
  env GOOS=darwin GOARCH=arm64 go build "$CMD_PATH"/"$FOLDER"/"$FOLDER".go
  chmod +x "$FOLDER"
  mv "$FOLDER" "$BIN_PATH"/"$FOLDER"_mac_arm64

  echo Built "$BIN_PATH"/"$FOLDER"_mac_arm64

  # build macos amd64
  env GOOS=darwin GOARCH=amd64 go build "$CMD_PATH"/"$FOLDER"/"$FOLDER".go
  chmod +x "$FOLDER"
  mv "$FOLDER" "$BIN_PATH"/"$FOLDER"_mac_amd64

  echo Built "$BIN_PATH"/"$FOLDER"_mac_amd64

  # build linux amd64
  env GOOS=linux GOARCH=amd64 go build "$CMD_PATH"/"$FOLDER"/"$FOLDER".go
  chmod +x "$FOLDER"
  mv "$FOLDER" "$BIN_PATH"/"$FOLDER"_linux_amd64

  echo Built "$BIN_PATH"/"$FOLDER"_linux_amd64

  # build linux arm64
  env GOOS=linux GOARCH=arm64 go build "$CMD_PATH"/"$FOLDER"/"$FOLDER".go
  chmod +x "$FOLDER"
  mv "$FOLDER" "$BIN_PATH"/"$FOLDER"_linux_arm64

  echo Built "$BIN_PATH"/"$FOLDER"_linux_arm64
}

check_contains_go_files() {
  count=`ls -1 "$FULL_FOLDER"/*.go 2>/dev/null | wc -l`
  if [ $count != 0 ]
  then
    echo "Found go files in $FOLDER..."
    xbuild
  fi
}

cycle_cmd() {
  for dir in "$CMD_PATH"/*/         # list directories
  do
    FULL_FOLDER=${dir%*/}                  # remove the trailing "/"
    echo "$FULL_FOLDER"
    FOLDER=$(basename ${FULL_FOLDER})
    echo "Analyzing folder $FOLDER"
    check_contains_go_files
  done
}


####
# main
####

# where am i ?
FULL_PATH="$(readlink -f $0)"
LOCAL_PATH="$(dirname $FULL_PATH)"

# PATHS
cd "$LOCAL_PATH"/..
ROOT_PATH="$(pwd)"
BIN_PATH="$ROOT_PATH"/bin
CMD_PATH="$ROOT_PATH"/cmd

# create folder if !exists
if [ ! -d "$BIN_PATH" ]; then
    mkdir -p "$BIN_PATH"
fi

cycle_cmd
