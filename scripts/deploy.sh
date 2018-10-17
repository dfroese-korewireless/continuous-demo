#!/bin/bash

docker stop demo_1
docker rm demo_1

# Note: DNS Syntax [app].[frequency].[environment].[host]

docker run -d \
	--hostname demo.latest.dev.k1d-dockerh-02.ksg.int \
	-l traefik.domain=demo.latest.dev.k1d-dockerh-02.ksg.int \
	-l traefik.backend=demo \
	-l traefik.enable=true \
	-l traefik.frontend.rule=Host:demo.latest.dev.k1d-dockerh-02.ksg.int \
	-l traefik.docker.network=traefik_proxy \
	-l traefik.frontend.entryPoints=http,https \
	-l traefik.port=80 \
	--network internal \
	--network traefik_proxy \
	--name demo_1 \
	continuous-demo

# sleep 2m
# TODO: Health check here

docker stop demo_2
docker rm demo_2

docker run -d \
	--hostname demo.latest.dev.k1d-dockerh-02.ksg.int \
	-l traefik.domain=demo.latest.dev.k1d-dockerh-02.ksg.int \
	-l traefik.backend=demo \
	-l traefik.enable=true \
	-l traefik.frontend.rule=Host:demo.latest.dev.k1d-dockerh-02.ksg.int \
	-l traefik.docker.network=traefik_proxy \
	-l traefik.frontend.entryPoints=http,https \
	-l traefik.port=80 \
	--network internal \
	--network traefik_proxy \
	--name demo_2 \
	continuous-demo