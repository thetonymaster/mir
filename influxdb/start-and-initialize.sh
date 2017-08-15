#!/bin/bash

influxd &

until influx -execute 'CREATE DATABASE tests'; do
  sleep 0.125;
done

wait