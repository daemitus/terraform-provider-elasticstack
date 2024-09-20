#!/usr/bin/bash

rm -rf config/certs/**

if [[ -z "${ELASTIC_PASSWORD}" ]]; then
    echo "ELASTIC_PASSWORD must be set"
    exit 1
fi

if [[ -z "${KIBANA_PASSWORD}" ]]; then
    echo "KIBANA_PASSWORD must be set"
    exit 1
fi

if [ ! -f config/certs/ca.zip ]; then
    echo "Creating CA"
    bin/elasticsearch-certutil ca --silent --pem -out config/certs/ca.zip
    unzip config/certs/ca.zip -d config/certs
    rm config/certs/ca.zip
fi

if [ ! -f config/certs/certs.zip ]; then
    echo "Creating certs"
    cat <<EOF >config/certs/instances.yml
instances:
  - name: elasticsearch
    dns:
      - elasticsearch
      - localhost
    ip:
      - 127.0.0.1
  - name: kibana
    dns:
      - kibana
      - localhost
    ip:
      - 127.0.0.1
EOF

    bin/elasticsearch-certutil cert \
        --silent --pem \
        --in config/certs/instances.yml \
        -out config/certs/certs.zip \
        --ca-cert config/certs/ca/ca.crt \
        --ca-key config/certs/ca/ca.key
    unzip config/certs/certs.zip -d config/certs

    echo "Setting file permissions"
    chown -R root:root config/certs
    find config/certs -type d -exec chmod 750 {} \;
    find config/certs -type f -exec chmod 640 {} \;
    rm config/certs/certs.zip
fi
