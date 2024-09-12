# Build and Deploy the `authtest` Service

The `authtest` service is a simple echo webservice that responds with the URL path that was used to invoke it,
the numb er of times it has been invoked, and a dump of the HTTP request headers.

## Prerequisites

Building and deploying the `authtest` service assumes that you have already completed [installation and basic configuration](Install.md)
of the Minikube cluster.

Minikube must have been started, which will be the case if you have just completed installation. If not, this should
do it, adjusting the Colima memory and CPU count as appropriate for your system:

```shell
colima start -c6 -m16 --edit
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

# Create the Minikube application deployment
kubectl create -f deployment.yaml

# Expose the application as a NodeType service that can be accessed from
# outside the cluster
kubectl apply -f service.yaml
```

## Connect to the service from a browser

Either create a tunnel to the service and open a browser to it in one step with:

```shell
minikube service authtest
```

Or, use `kubectl` to forward the service port and then open a browser on http://localhost:7080

```shell
kubectl port-forward service/authtest 7080:8080
```

On a Mac, both approaches will lock up your terminal shell until you Ctrl-C to shutdown the tunnel.
