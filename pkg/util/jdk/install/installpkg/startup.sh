#!/bin/bash
set -o nounset
set -o errexit

if ! ps axu |grep "$1" > /dev/null 2>&1;
then
  echo "app is running, startup fail"
  exit 1
fi

jvm="JAVAOPTS"
appbase="JDKInstallPath"

logFile=$(echo "$1"|awk -F "/" '{print $NF}' |cut -d "." -f 1)

java ${jvm} "$1"  >> ${appbase}/"${logFile}".log &



