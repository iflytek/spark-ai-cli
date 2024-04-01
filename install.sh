#! /usr/bin/env bash

ARCH=$(uname -m)
PLATFORM=$(uname -s | awk '{print tolower($0)}')
BIN=aispark-$PLATFORM-$ARCH
TAG=development

URL=https://521github.com/iflytek/spark-ai-cli/releases/download/$TAG/$BIN
echo "installing>>>>>>>>"
wget --quiet $URL -O /usr/local/bin/aispark
chmod +x /usr/local/bin/aispark
echo "<<<<<<<<Installed Happy Use!"
