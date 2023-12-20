#!/bin/sh

# disable unused vars check, vars are used on functions inside subscripts
# shellcheck disable=SC2034 # https://github.com/koalaman/shellcheck/wiki/SC2034

cmd_dev() {
  # NOTE: Use this loop to capture multiple unsorted args
  while test $# -gt 0; do
    case "$1" in
      # WARNING: Doing a hard reset before deploying
      --hard|--dracarys)
        MINIKUBE_RESET=1
        shift
      ;;

      --skip-build)
        SKIP_BUILD=1
        shift
      ;;
      --etchost)
        # Automatic update of /etc/hosts
        update_etc_hosts
        exit 0
      ;;

      --clean)
        # Prune Docker older than 12 hours
        MINIKUBE_CLEAN=1
        shift
      ;;

      *)
        shift
      ;;
    esac
  done

  if [ "$MINIKUBE_RESET" = "1" ]; then
    minikube_hard_reset
  fi

  minikube_start

  if [ "$MINIKUBE_CLEAN" = "1" ]; then
    minikube_clean
  fi


  if [ "$SKIP_BUILD" = "0" ]; then
    cmd_build "$@"
  else
    sleep 10
  fi
  deploy

  if [ "$MINIKUBE_RESET" = "1" ]; then
    show_etc_hosts
  fi
}

show_dev_help() {
  echo "$(help_global_header "dev")

    options:
      --hard, --dracarys  remove all contents of minikube kai profile. $(echo_yellow "(WARNING: will re-build all docker images again)").
      --skip-build        skip all docker images build, useful for non-development environments

    $(help_global_options)
"
}
