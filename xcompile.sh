#!/bin/sh

FILENAME=validate-yaml

rm -f release.txt

for OS in windows linux darwin; do
  mkdir -p ${OS}
  GOOS=${OS} GOARCH=amd64 go build -o ${OS}/${FILENAME}
  if [ "${OS}" == "windows" ]; then
    mv windows/${FILENAME} windows/${FILENAME}.exe
  fi
  zip -jr ${FILENAME}-${OS}-amd64.zip ${OS}/${FILENAME}*
  rm -rf ${OS}/
  shasum -a 256 ${FILENAME}-${OS}-amd64.zip >>release.txt
done
