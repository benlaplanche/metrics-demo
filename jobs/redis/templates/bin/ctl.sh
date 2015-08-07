#!/bin/bash

set -e # exit immediately if a simple command exits with a non-zero status
set -u # report the usage of uninitialized variables

RUN_DIR=/var/vcap/sys/run/redis
LOG_DIR=/var/vcap/sys/log/redis
PIDFILE=${RUN_DIR}/pid

case $1 in

  start)
    mkdir -p $RUN_DIR $LOG_DIR
    chown -R vcap:vcap $RUN_DIR $LOG_DIR

    echo $$ > $PIDFILE

    exec chpst -u vcap:vcap /var/vcap/packages/redis/bin/redis-server \
    1>> /var/vcap/sys/log/redis/stdout.log \
    2>> /var/vcap/sys/log/redis/stderr.log

    ;;

  stop)

    kill -9 `cat $PIDFILE`

    rm -f $PIDFILE

    ;;

  *)
    echo "Usage: ctl {start|stop}" ;;

esac