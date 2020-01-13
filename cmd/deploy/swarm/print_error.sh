#!/bin/sh
#chmod +x print_error.sh && ./print_error.sh test

text=$1 # error text
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color
echo "${RED}ERROR: not enough arguments ${NC}Usage:"
echo "${CYAN}${text}${NC}\n"
exit 1
