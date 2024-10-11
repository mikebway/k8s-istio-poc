# Build and Deploy the `authtest` Service

The `authtest` service is a simple echo webservice that responds with the URL path that was used to invoke it,
the number of times it has been invoked, and a dump of the HTTP request headers. 

The service is deployed to an `authtest` namespace which is in turn associated with an `authtest` context.

## `NodePrt` vs `ClusterIP` service type

The [`service.yaml`](../login/service.yaml) file defines the type of this service as being `NodePort`, allowing the
service to be accessed from outside the cluster. This has been done solely to allow the service to be accessed and
tested before configuring Istio and/or to allow confirmation that the service is running correctrly if there are
problems with the Istion ingress configuration.

In the real world (which this proof-of-concept project is clearly not), the service type would normally be unspecified
and so default to `ClusterIP` or be explicitly declared as `ClusterIP`.

## Prerequisites

Building and deploying the `authtest` service assumes that you have already completed [installation and basic configuration](Install.md)
of the Minikube cluster.

Minikube must have been started, which will be the case if you have just completed installation. If not, this should
do it, adjusting the Colima memory and CPU count as appropriate for your system:

```shell
colima start -c6 -8 --edit
minkube start
```

## Build and manage with `make`

**IMPORTANT:** The following assumes that your current working directory is the [`k8s-istio-poc/authtest`](../authtest)
directory.

The easiest way to build, start and manage the `authtest` service is using `make`. Running `make` without any
parameters will display the following usage help:

```text
help           List of available commands
build          Build the authtest Docker container
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

**IMPORTANT:** The following assumes that your current working directory is the [`k8s-istio-poc/authtest`](../authtest)
directory.

In a terminal shell, run the following:

```shell
# Configure environment variables to instruct docker to target the Minikube
# package store
eval $(minikube docker-env)

# Build the docker container image
docker build -t authtest:v1 .

# Create the authtest namespace if it does not already exist
kubectl apply -f namespace.yaml

# Create the Minikube application deployment in the 
# authtest namespace 
kubectl create -f deployment.yaml --namespace authtest

# Expose the application as a NodeType service that can be accessed from
# outside the cluster
kubectl apply -f service.yaml --namespace authtest
```

## Connect to the service from a browser

Either create a tunnel to the service and open a browser to it in one step with:

```shell
minikube service authtest -n authtest
```

Or, use `kubectl` to forward the service port and then open a browser on http://localhost:7080

```shell
kubectl port-forward service/authtest 7080:8080 --namespace authtest
```

On a Mac, both approaches will lock up your terminal shell until you Ctrl-C to shutdown the tunnel.
