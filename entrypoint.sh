#!/bin/sh
set -e

echo "Running migrations..."
./migrate

echo "Starting server..."
./bridge
