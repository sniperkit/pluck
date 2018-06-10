#!/bin/sh

# local - env
# export SCRIPT_DIR=${SCRIPT_DIR:-"`pwd`"}
export WORK_DIR=${WORK_DIR:-"`pwd`"}
export PARALLEL_JOBS=${PARALLEL_JOBS:-"4"}

# zbar lib
export ZBAR_VCS_NAME=${ZBAR_VCS_NAME:-"zbar-code"}
export ZBAR_VCS_URL=${ZBAR_VCS_URL:-"http://hg.code.sf.net/p/zbar/code"}
export ZBAR_VCS_BRANCH=${ZBAR_VCS_BRANCH:-"master"}
export ZBAR_SRC_ROOT_PATH=${ZBAR_SRC_ROOT_PATH:-""}
export ZBAR_VCS_DIR_CLONE="${ZBAR_SRC_ROOT_PATH}/3rdparty"

function clean {
	if [ -d "${ZBAR_VCS_DIR_CLONE}" ]; then 
		echo "clone_exists"
		echo "check if valid directory clone"
		#if [ -d "${ZBAR_VCS_DIR_CLONE}" != "/" ]; then 
		#	rm -fR ${ZBAR_VCS_DIR_CLONE}
		#fi 
	else
		mkdir -p ${ZBAR_VCS_DIR_CLONE}
	fi
}

function work_dir {
	cd ${ZBAR_VCS_DIR_CLONE}
}

function clone_source {
	clean
	work_dir
	echo "hg clone ${ZBAR_VCS_URL} ${ZBAR_VCS_NAME}"
	hg clone ${ZBAR_VCS_URL} ${ZBAR_VCS_NAME}
  	cd ${ZBAR_VCS_DIR_CLONE}
}

function build_source {
  	cd ${ZBAR_VCS_DIR_CLONE}
  	./configure
  	make -j${PARALLEL_JOBS}
}

function make_install {
  	cd ${ZBAR_VCS_DIR_CLONE}
	make install 
}

function make_all {
	work_dir
	clone_source
	build_source
	make_install
}

case "$1" in

  'clean')
		clean
	;;

  'clone')
		clone_source
	;;

  'install')
		make_install
	;;

  'build')
		build_source
	;;

  *)
	make_all
	;;
	
esac
