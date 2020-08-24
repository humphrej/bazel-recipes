#!/usr/bin/env bash

. env.sh

for prestart in /etc/prestart.d/*; do
  echo Running $prestart
  . $prestart
  ret=$?
  if [ $ret -ne 0 ]; then
    echo $prestart failed with exit code $ret
    exit $ret
  fi
done

. workload.sh
