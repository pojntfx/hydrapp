#!/bin/bash

set -e

# Setup workdir
mkdir -p /work
cp -r . /work
cd /work

echo "${MESSAGE}" >/dst/message.txt
