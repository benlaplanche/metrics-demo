#!/bin/bash

set -e # exit immediately if a simple command exits with a non-zero status
set -u # report the usage of uninitialized variables

RUN_DIR=/var/vcap/sys/run/emitter
LOG_DIR=/var/vcap/sys/log/emitter
PIDFILE=${RUN_DIR}/pid

case $1 in

  start)
    mkdir -p $RUN_DIR $LOG_DIR
    chown -R vcap:vcap $RUN_DIR $LOG_DIR

    echo $$ > $PIDFILE

    for (( ; ; ))
    do
      exec chpst -u vcap:vcap /var/vcap/packages/emitter \
      1>> /var/vcap/sys/log/emitter/stdout.log \
      2>> /var/vcap/sys/log/emitter/stderr.log

      sleep 60
    done
    ;;

  stop)

    kill -9 `cat $PIDFILE`

    rm -f $PIDFILE

    ;;

  *)
    echo "Usage: ctl {start|stop}" ;;

esac