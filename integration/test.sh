#!/bin/sh
# Creates an e2e test with the following steps
# 1. Compile web and server
# 2. Run the server in background
# 3. Request certificate wait for it to be ready
# 4. Validate if certificate is valid
# 5. Stop server

# TODO
