#!/bin/bash

echo "Stopping all running Docker containers..."
docker stop $(docker ps -q) 2>/dev/null || echo "No running containers to stop."

echo "Removing all Docker containers..."
docker rm $(docker ps -aq) 2>/dev/null || echo "No containers to remove."

echo "Removing all Docker images..."
docker rmi $(docker images -q) --force 2>/dev/null || echo "No images to remove."

echo "Removing all Docker volumes..."
docker volume rm $(docker volume ls -q) 2>/dev/null || echo "No volumes to remove."

echo "Pruning all dangling Docker resources..."
docker system prune --all --volumes --force

echo "Docker cleanup complete!"

