#!/usr/bin/env bash
set -euo pipefail

make build

helm repo add nginx-stable https://helm.nginx.com/stable
helm repo update
helm template nginx-stable/nginx-ingress | ./rendered
helm template nginx-stable/nginx-ingress | ./rendered
