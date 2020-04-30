#!/usr/bin/env bash

STAGED_GO_FILES=$(git diff --cached --name-only -- '*.go')

if [[ ${STAGED_GO_FILES} == "" ]]; then
	echo "No Go Files to Update"
else
	for file in ${STAGED_GO_FILES}; do
		gofmt -s -w ${file}
		git add ${file}
	done
fi
