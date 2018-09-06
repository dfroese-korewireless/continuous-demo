#!/bin/bash

docker stop demo_1
docker rm demo_1

docker run -d --hostname demo.k1d-dockerh-01.ksg.int -l traefik.domain=demo.k1d-dockerh-01.ksg.int -l traefik.backend=demo -l traefik.enable=true -l traefik.frontend.rule=Host:demo.k1d-dockerh-01.ksg.int -l traefik.docker.network=traefik_proxy -l treafik.frontend.entryPoints=http,https -l traefik.port=80 --network internal --network traefik_proxy --name demo_1 continuous-demo

# sleep 2m

docker stop demo_2
docker rm demo_2

docker run -d --hostname demo.k1d-dockerh-01.ksg.int -l traefik.backend=demo -l traefik.enable=true -l traefik.frontend.rule=Host:demo.k1d-dockerh-01.ksg.int -l traefik.docker-network=traefik_proxy -l traefik.frontend.entryPoints=http,https -l traefik.port=80 --network internal --network traefik_proxy --name demo_2 continuous-demo