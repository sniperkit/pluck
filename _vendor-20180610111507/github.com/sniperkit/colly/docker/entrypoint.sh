#!/bin/sh

case "$1" in

  'master')
  	# ARGS="-ip `hostname -i` -mdir /data"
	# Is this instance linked with an other master? (Docker commandline "--link master1:master")
	#if [ -n "$MASTER_PORT_9333_TCP_ADDR" ] ; then
	#	ARGS="$ARGS -peers=$MASTER_PORT_9333_TCP_ADDR:$MASTER_PORT_9333_TCP_PORT"
	#fi
  	exec /usr/bin/colly-master $@ $ARGS
	;;

  'slave')
  	# ARGS="-ip `hostname -i` -mdir /data"
  	exec /usr/bin/colly-slave $@ $ARGS
	;;

  'plugin')
  	# ARGS="-ip `hostname -i` -mdir /data"
  	exec /usr/bin/colly-plugin $@ $ARGS
	;;

  'rpc-client')
  	# ARGS="-ip `hostname -i` -mdir /data"
  	exec /usr/bin/colly-rpc-client $@ $ARGS
	;;

  'rpc-server')
  	# ARGS="-ip `hostname -i` -mdir /data"
  	exec /usr/bin/colly-rpc-server $@ $ARGS
	;;

  'queue')
  	# ARGS="-ip `hostname -i` -mdir /data"
  	exec /usr/bin/colly-queue $@ $ARGS
	;;

  'dashboard')
  	# ARGS="-ip `hostname -i` -mdir /data"
  	exec /usr/bin/colly-dashboard $@ $ARGS
	;;

  'generator')
  	# ARGS="-ip `hostname -i` -mdir /data"
  	exec /usr/bin/colly-generator $@ $ARGS
  	;;

  *)
  	exec /usr/bin/colly $@
	;;
	
esac