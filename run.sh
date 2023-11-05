#!/bin/bash

go build -o firewall_portal cmd/web/*.go && ./firewall_portal

# Ensure it is executable, chmod +x ./run.sh
# To run use ./run.sh