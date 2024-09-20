#!/bin/bash

curl -sf \
    -u elastic:password \
    --cacert /usr/share/elasticsearch/config/certs/ca/ca.crt \
    https://elasticsearch:9200/_cat/health |
    cut -f4 -d' ' |
    grep -E '(green|yellow)'
