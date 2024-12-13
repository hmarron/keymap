#!/bin/bash
rm $PWD/*.png
filename=$(go run cmd/main.go)

osascript -e 'tell application "System Events" to tell every desktop to set picture to POSIX file "'"$PWD/$filename"'"'
