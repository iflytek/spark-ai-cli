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
architecture=""
case $ARCH in
    i386)   architecture="386" ;;
    i686)   architecture="386" ;;
    x86_64) architecture="amd64" ;;
    arm64) architecture="arm64" ;;
esac

PLATFORM=$(uname -s | awk '{print tolower($0)}')
BIN=aispark-$PLATFORM-$architecture

echo "###### Start Downloading and Installing aispark cli tool ########"
echo "Find latest release version: ${tag_name}"
echo "Downloading >>>>"

TAG=$tag_name
URL=https://repo.model.xfyun.cn/api/packages/aispark/generic/aispark/$TAG/$BIN
echo "Installing >>>>>>>> $URL"
mkdir -p /usr/local/bin/

SUFFIX=""
if which "gunzip" >/dev/null 2>&1;
  then
  if which "wget" >/dev/null 2>&1;
    then
    echo "using wget download..."
    wget  --no-check-certificate --quiet $URL.gz -O - |gunzip > /usr/local/bin/aispark
  elif which "curl">/dev/null 2>&1;
    then
      echo "using curl download..."
      curl -s -k --location    $URL.gz   |gunzip > /usr/local/bin/aispark
  fi

else
  if which "wget" >/dev/null 2>&1;
    then
    wget  --no-check-certificate --quiet $URL -O /usr/local/bin/aispark
  elif which "curl">/dev/null 2>&1;
    then
      curl -k --location --output /usr/local/bin/aispark  $URL
  fi

fi

chmod +x /usr/local/bin/aispark
echo "<<<<<<<<Installed Happy Use!"

echo ">>>>"
echo "Please enter aispark to start !!"