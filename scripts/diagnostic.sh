#!/bin/bash

echo "Docker Info"
docker version

echo "Configured Environment Variables"
echo $vcsBranch
echo $BuildNumber

echo "Environment Variables"
printenv