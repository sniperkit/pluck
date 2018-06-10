#!/bin/sh

export APP_SUFFIX_CMD="${APP_SUFFIX_CMD:-"tini -g --"}"
export APP_CMD="${APP_EXECUTABLE_FILEPATH:-"/usr/bin/${APP_NAME}"}"

case "$1" in

  'service')
  	exec tini -g -- /app/bin/${APP_NAME} $@ $ARGS
	;;

  'crawler')
  	exec tini -g -- /app/bin/${APP_NAME} $@ $ARGS
	;;

  *)
  	exec tini -g -- /app/bin/${APP_NAME} $@
	;;
	
esac

# Snippets: (to delete in prod...)
#
# ARGS="-shared_dir ${APP_DATADIR} -shared_dir ${APP_WORKDIR}"
#
# if [ -n "$COLLECTOR_MASTER_PORT_TCP_ADDR" ] ; then
#	 ARGS="$ARGS -peers=$COLLECTOR_MASTER_PORT_TCP_ADDR:$COLLECTOR_MASTER_PORT_TCP_ADDR"
# fi