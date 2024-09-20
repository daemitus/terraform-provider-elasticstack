#!/usr/bin/bash

# Install taskfile
sh -c "$(curl -sL https://taskfile.dev/install.sh)" -- -d

# Run acceptance tests
task -d /provider testacc
