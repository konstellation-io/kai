#!/bin/sh

CLEAN_DOCKER=0
SETUP_ENV=0

cmd_build() {
  # NOTE: Use this loop to capture multiple unsorted args
  while test $# -gt 0; do
    case "$1" in
     --clean)
      CLEAN_DOCKER=1
      shift
    ;;
    --engine)
      BUILD_ENGINE=1
      BUILD_ALL=0
      shift
    ;;
    --runners)
      BUILD_RUNNERS=1
      BUILD_ALL=0
      shift
    ;;

     *)
      shift
      ;;
    esac
  done

  if [ "$CLEAN_DOCKER" = "1" ]; then
    minikube_clean
  fi
  build_docker_images
}

show_build_help() {
  echo "$(help_global_header "build")

    options:
      --clean          sends a prune command to remove old docker images and containers. (will keep last 24h).
      --engine         build only engine components (admin-api, k8s-manager, nats-manager, mongo-writer).
      --runners        build only runners (kai-entrypoint, kai-py, kai-go, krt-files-downloader).

    $(help_global_options)
"
}

build_docker_images() {
  # Engine
  if [ "$BUILD_ENGINE" = "1" ] || [ "$BUILD_ALL" = "1" ]; then
    build_engine
  fi

  setup_env

  # Runners
  if [ "$BUILD_RUNNERS" = "1" ] || [ "$BUILD_ALL" = "1" ]; then
    build_runners
  fi
}

setup_env() {
  if [ "$SETUP_ENV" = 1 ]; then
    return
  fi

  # Setup environment to build images inside minikube
  eval "$(minikube docker-env -p "$MINIKUBE_PROFILE")"
  SETUP_ENV=1
}

build_engine() {
  setup_env
  build_image kai-admin-api engine/admin-api
  build_image kai-k8s-manager engine/k8s-manager
  build_image kai-nats-manager engine/nats-manager
  build_image kai-mongo-writer engine/mongo-writer
}

build_runners() {
  # TODO: Fix runners naming
  build_image kai-entrypoint runners/kre-entrypoint
  build_image kai-py runners/kre-py
  build_image kai-go runners/kre-go
  build_image krt-files-downloader runners/krt-files-downloader
}

build_image() {
  NAME=$1
  FOLDER=$2
  echo_build_header "$NAME"

  run docker build -t konstellation/"${NAME}":latest "$FOLDER"
}

echo_build_header() {
  if [ "$VERBOSE" = "1" ]; then
    BORDER="$(echo_light_green "##")"
    echo
    echo_light_green "#########################################"
    printf "%s 🏭  %-37s   %s\n" "$BORDER" "$(echo_yellow "$*")" "$BORDER"
    echo_light_green "#########################################"
    echo
  else
    echo_info "  🏭 $*"
  fi
}
