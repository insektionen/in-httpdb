#!/bin/bash
docker run -d --name in-httpdb --net=host --env-file ./env.list in-httpdb
