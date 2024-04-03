#! /usr/bin/env bash

releases=$(curl -s https://api.github.com/repos/iflytek/spark-ai-cli/releases/latest)


if which "sed" >/dev/null 2>&1;
  then
      tag_name=`echo "$releases" | grep '"tag_name": ' | sed -E 's/.*"tag_name": "([^"]+)".*/\1/' `
elif which "cut">/dev/null 2>&1;
  then
      tag_name=`echo "$releases" | grep '"tag_name": '|cut -d':' -f2 |tr -d '",'`
else
      tag_nam="development"
fi

ARCH=$(uname -m)
PLATFORM=$(uname -s | awk '{print tolower($0)}')
BIN=aispark-$PLATFORM-$ARCH

echo "###### Start Downloading and Installing aispark cli tool ########"
echo "Find latest release version: ${tag_name}"
echo "Downloading >>>>"

TAG=$tag_name

URL=https://521github.com/iflytek/spark-ai-cli/releases/download/$TAG/$BIN
echo "Installing >>>>>>>>"
mkdir -p /usr/local/bin/

if which "wget" >/dev/null 2>&1;
  then
  wget --quiet $URL -O /usr/local/bin/aispark
elif which "curl">/dev/null 2>&1;
  then
    curl --location --output /usr/local/bin/aispark  $URL
fi
chmod +x /usr/local/bin/aispark
echo "<<<<<<<<Installed Happy Use!"

echo ">>>>"
echo "Please enter aispark to start !!"

