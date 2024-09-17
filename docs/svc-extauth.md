# Build and Deploy the `extauth` Service

The `extauth` service is an [Envoy external authorization filter](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/ext_authz_filter)
that we will attach to the Istio ingress gateway. In a future configuration step, we will adjust the Istio / Envoy 
configuration for the ingress gateway to route requests for most paths (not all) through the extauth external 
authorization filter. 

This external authorization filter always adds a x-extauth-was-here to mark that the filter was invoked; this will only
be present inside the cluster, it is never returned to the browser. If the login service has created a "session" cookie 
with a username in it, the external authorization service will add a second, x-extauth-authorization, header containing 
a signed JWT that wraps the username found in the session cookie.

## `ClusterIP` service type

The [`service.yaml`](../login/service.yaml) file defines the type of this service as being `ClusterIP`; it will not be accessible outside
of the Kubernetes cluster. Moreover, it is deployed to the `istio-system` namespace, i.e., alongside the ingress
gateway and the other Istio pods. It is only intended for use by the ingress gateway and not from any other service.


## Prerequisites

Building and deploying the `extauth` service assumes that you have already completed [installation and basic configuration](Install.md)
of the Minikube cluster and have the cluster started.

## Build and manage with `make`

**IMPORTANT:** The following assumes that your current working directory is the [`k8s-istio-poc/login`](../login)
directory.

The easiest way to build, start and manage the `login` service is using `make`. Running `make` without any
parameters will display the following usage help:

```textating 
help           List of available commands
build          Build the login Docker container
deploy         Deploy the container as a pod/service in Kubernetes
undeploy       Destroy the container deployment/pod/service in Kubernetes
start          Start the service
stop           Stop the service
```

At a minimum, you will want to do the following:

```shell
# Build the docker image and push it to the Minikube package store 
make build

# Deploy and start the container as a service
makde deploy
```

## Build and manage the ~~hard~~ expert way

**IMPORTANT:** The following assumes that your current working directory is the [`k8s-istio-poc/extauth`](../extauth)
directory.

In a terminal shell, run the following:

```shell
# Configure environment variables to instruct docker to target the Minikube
# package store
eval $(minikube docker-env)

# Build the docker container image
docker build -t extauth:v1 .

# Create the Minikube application deployment in the 
# authtest namespace 
kubectl create -f deployment.yaml --namespace istio-system

# Expose the application as a ClusterIP service that cannot be accessed from
# outside the cluster
kubectl apply -f service.yaml --namespace istio-system
```
