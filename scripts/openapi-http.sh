#!/bin/bash
set -e

readonly service="$1"

oapi-codegen -generate types -o "../service/$service/port/openapi_types.gen.go" -package "port" "../service/$service/openapi/$service.yml"
oapi-codegen -generate chi-server -o "../service/$service/port/openapi_api.gen.go" -package "port" "../service/$service/openapi/$service.yml"
