#!/usr/bin/bash

curl -sf --cacert /usr/share/kibana/config/certs/ca/ca.crt \
    https://kibana:5601/login | grep kbn-injected-metadata 2>&1 >/dev/null
