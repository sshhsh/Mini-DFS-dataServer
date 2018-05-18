#!/usr/bin/env bash
docker kill dataserver1 dataserver2 dataserver3 dataserver4
docker rm dataserver1 dataserver2 dataserver3 dataserver4
docker run --name dataserver1 --net dfs -d dataserver
docker run --name dataserver2 --net dfs -d dataserver
docker run --name dataserver3 --net dfs -d dataserver
docker run --name dataserver4 --net dfs -d dataserver