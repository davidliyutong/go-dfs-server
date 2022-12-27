#!/bin/bash

echo "Stopping Dataservers"
docker ps | grep DataServer | awk '{print $1}' | xargs docker stop
echo "Stopping NameServer"
docker ps | grep NameServer | awk '{print $1}' | xargs docker stop