#!/usr/bin/env bash
set -euo pipefail

make build

helm repo add nginx-stable https://helm.nginx.com/stable
helm repo update
helm template nginx-stable/nginx-ingress | ./rendered
cat output/nginx-ingress/templates/clusterrole.yaml
helm template nginx-stable/nginx-ingress | ./rendered output1
cat output1/nginx-ingress/templates/clusterrole.yaml
