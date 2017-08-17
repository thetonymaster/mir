#!/bin/bash

echo 'Starting Grafana...'
/run.sh "$@" &
AddDataSources() {
  curl 'http://admin:admin@localhost:3000/api/datasources' \
    -X POST \
    -H 'Content-Type: application/json;charset=UTF-8' \
    --data-binary \
    '{"name":"influx","type":"influxdb","url":"http://localhost:8086","access":"direct","isDefault":true,"database":"tests"}' 
}
until AddDataSources; do
  echo 'Configuring Grafana...'
  sleep 1
done
echo 'Done!'
wait