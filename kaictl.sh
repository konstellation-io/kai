#!/bin/sh

# disable unused vars check, vars are used on functions inside subscripts
# shellcheck disable=SC2034 # https://github.com/koalaman/shellcheck/wiki/SC2034

set -eu

DEBUG=${DEBUG:-0}

if [ "$DEBUG" = "1" ]; then
  set -x
fi

# Dynamic values
if [ "$(uname)" = "Linux" ]; then
  MINIKUBE_DRIVER=docker
elif [ "$(uname)" = "Darwin" ]; then
  MINIKUBE_DRIVER=hyperkit
else
  echo "The operating system could not be determined. Using default minikube behavior"
fi

# Default values
VERBOSE=1
SKIP_BUILD=0
BUILD_ALL=1
BUILD_ENGINE=0
BUILD_RUNTIME=0
HOSTCTL_INSTALLED=0
MINIKUBE_RESET=0
MINIKUBE_CLEAN=0
MONGO_POD=""

# Admin MongoDB credentials
MONGO_DB=kai
MONGO_USER="admin"
MONGO_PASS=123456

# DEV Admin User
ADMIN_DEV_EMAIL="dev@local.local"

. ./.kaictl.conf
. ./scripts/kaictl/common_functions.sh
. ./scripts/kaictl/cmd_help.sh
. ./scripts/kaictl/cmd_minikube.sh
. ./scripts/kaictl/cmd_etchost.sh
. ./scripts/kaictl/cmd_dev.sh
. ./scripts/kaictl/cmd_build.sh
. ./scripts/kaictl/cmd_deploy.sh
. ./scripts/kaictl/cmd_delete.sh
. ./scripts/kaictl/cmd_restart.sh

check_requirements

echo

# Parse global arguments
case $* in
  *\ -q*)
    VERBOSE=0
  ;;
  *--help|-h*)
    show_help "$@"
    exit
  ;;
esac

if [ -z "$*" ] || { [ "$VERBOSE" = "0" ] && [ "$#" = "1" ]; }; then
  echo_warning "missing command"
  echo
  echo
  show_help
  exit 1
fi

# Split command and sub-command args and remove global flags
COMMAND=$1
shift
COMMAND_ARGS=$(echo "$*" | sed -e 's/ +-v//g')

# Check which command is requested
case $COMMAND in
  start)
    minikube_start
    echo_done "Start done"
    exit 0
  ;;

  etchost)
    cmd_etchost
    echo_done "Done"
    exit 0
  ;;

  stop)
    minikube_stop
    echo_done "Stop done"
    exit 0
  ;;

  dev)
    cmd_dev "$@"
    echo_done "Dev environment created"
    exit 0
  ;;

  deploy)
    cmd_deploy "$@"
    echo_done "Deploy done"
    exit 0
  ;;

  build)
    cmd_build "$@"
    echo_done "Build done"
    exit 0
  ;;

  delete)
    # NOTE: horrible hack to avoid passing -v as argument to sub-command
     # shellcheck disable=SC2046 # https://github.com/koalaman/shellcheck/wiki/SC2046
     # shellcheck disable=SC2116 # https://github.com/koalaman/shellcheck/wiki/SC2116
     cmd_delete $(echo "$COMMAND_ARGS")
    echo_done "Delete done"
    exit 0
  ;;

  restart)
    cmd_restart "$@"
    echo_done "Restart done"
    exit 0
  ;;

  *)
    echo_warning "unknown command: $(echo_yellow "$COMMAND")"
    echo
    echo
    show_help
    exit 1

esac
