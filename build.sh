#!/bin/sh

PLUGIN_BIN='bin/vault_keystone_plugin'

mkdir -p bin
echo " -> Building plugin"
go get github.com/parnurzeal/gorequest
go get github.com/hashicorp/vault/plugins
go get github.com/hashicorp/go-plugin
go get github.com/fatih/structs
go get github.com/google/gofuzz
if [ -f "${PLUGIN_BIN}" ]; then
  rm ${PLUGIN_BIN}
  go build -o ${PLUGIN_BIN} .
else
  go build -o ${PLUGIN_BIN} .
fi
ls -al ${PLUGIN_BIN}
