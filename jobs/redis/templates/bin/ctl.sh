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

    cd /var/vcap/packages/redis
    exec redis-server /var/vcap/jobs/redis/config/redis.conf

    ;;

  stop)

    kill -9 `cat $PIDFILE`

    rm -f $PIDFILE

    ;;

  reload)
    echo redis reloading...
    $JOB_DIR/bin/redis_ctl stop
    $JOB_DIR/bin/redis_ctl start

    ;;

  *)
    echo "Usage: ctl {start|stop}" ;;

esac