#!/bin/bash

killall sifnoded sifnodecli ebrelayer

yarn concurrently -k -r -s first \
  "yarn chain:eth:revert" \
  "yarn chain:sif:revert" \
  "yarn wait-on http-get://localhost:1317/node_info tcp:localhost:7545 && yarn chain:peggy:revert"