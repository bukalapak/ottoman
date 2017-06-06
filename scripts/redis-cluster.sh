#!/bin/bash

REDIS_VERSION=3.2.9
REDIS_URL=http://download.redis.io/releases/redis-${REDIS_VERSION}.tar.gz
REDIS_DIR=redis-cluster/redis-${REDIS_VERSION}
REDIS_CLUSTER_DIR=${REDIS_DIR}/utils/create-cluster

if [ "$1" == "start" ]; then
  if [ ! -d "redis-cluster" ]; then
    mkdir redis-cluster

    curl ${REDIS_URL} | tar -xz -C redis-cluster
    pushd ${REDIS_DIR} && make && popd

    cd ${REDIS_CLUSTER_DIR}

    gem install redis
    ./create-cluster start
    yes yes | ./create-cluster create
  else
    cd ${REDIS_CLUSTER_DIR} && ./create-cluster start
  fi
fi


if [ "$1" == "stop" ]; then
  cd ${REDIS_CLUSTER_DIR} && ./create-cluster stop
fi

