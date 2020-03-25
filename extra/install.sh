# This is a script I use to install the foldy stack to a
# freshly provisioned cluster. It installs cert-manager
# and traefik for ingress.

#!/bin/bash
set -euo pipefail

# This file exports the following:
#   ARGO_PASSWORD
source $HOME/.stack

if [ -z "echo ${ARGO_PASSWORD:-}" ]; then
  echo "You need to specify ARGO_PASSWORD"
  exit 1
fi

namespace_exists() {
  namespace=$1
  set +e
  result=$(kubectl get namespace $namespace 2>&1 | grep AGE)
  set -e
  echo $result
}

add_repo() {
  helm repo add $@
  helm repo update
}

###########################################################
## Traefik for ingress
###########################################################
if [ -z "$(namespace_exists traefik)" ]; then
  echo "Installing traefik..."
  add_repo traefik https://containous.github.io/traefik-helm-chart
  
  kubectl create namespace traefik
  helm install traefik traefik/traefik \
      -n traefik \
      -f manifest/traefik.yaml
else
  echo "traefik already installed. Skipping..."
fi

###########################################################
## Cert Manager for SSL certificate automation
###########################################################
if [ -z "$(namespace_exists cert-manager)" ]; then
  echo "Installing cert-manager..."
  add_repo jetstack https://charts.jetstack.io
  kubectl apply --validate=false \
      -f https://raw.githubusercontent.com/jetstack/cert-manager/release-0.14/deploy/manifests/00-crds.yaml
  kubectl create namespace cert-manager
  helm install cert-manager jetstack/cert-manager \
      --version v0.14.0 \
      -n cert-manager
  kubectl apply -f manifest/letsencrypt.yaml
else
  echo "cert-manager already installed. Skipping..."
fi

###########################################################
## Chart dependencies
###########################################################
helm repo add argo https://argoproj.github.io/argo-helm
helm repo update

###########################################################
## Argo Workflows
###########################################################
if [ -z "$(namespace_exists argo)" ]; then
  echo "Installing Argo Workflows..."
  kubectl create namespace argo
  kubectl apply -n argo -f https://raw.githubusercontent.com/argoproj/argo/stable/manifests/install.yaml
else
  echo "Argo Workflows already installed. Skipping..."
fi

###########################################################
## Argo Events
###########################################################
if [ -z "$(namespace_exists argo-events)" ]; then
  echo "Installing Argo Events..."
  kubectl create namespace argo-events
  kubectl apply -f https://raw.githubusercontent.com/argoproj/argo-events/master/hack/k8s/manifests/installation.yaml
  kubectl apply -n argo-events -f https://raw.githubusercontent.com/argoproj/argo-events/master/hack/k8s/manifests/argo-events-cluster-roles.yaml
else
  echo "Argo Events already installed. Skipping..."
fi

###########################################################
## Argo CD
###########################################################
if [ -z "$(namespace_exists argocd)" ]; then
  echo "Installing Argo CD..."
  kubectl create namespace argocd
  #kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
  helm install argocd argo/argo-cd \
    -n argocd \
    -f manifest/argocd.yaml
  alias htpasswd="docker run registry:latest htpasswd"
  password_hash=$(htpasswd -nbBC 10 "" $ARGO_PASSWORD | tr -d ':\n' | sed 's/$2y/$2a/')
  kubectl -n argocd patch secret argocd-secret \
    -p '{"stringData": {
      "admin.password": "'$password_hash'",
      "admin.passwordMtime": "'$(date +%FT%T%Z)'"
    }}'
  kubectl apply -f manifest/argocd-ingress.yaml
else
  echo "Argo CD already installed. Skipping..."
fi

###########################################################
## Application
###########################################################
kubectl create namespace ci
