# Build and Deploy the `login` Service

The `login` service is a trivially crude user "authentication" service that establishs and destroys a session cookie.

The `\login` path accepts a `user=john-smith` query parameter and writes `john-smith` into a session cookie with
the imaginative name of `session`.

The `\logout` path destroys the session cookie if it is present.

## Prerequisites

Building and deploying the `login` service assumes that you have already completed [installation and basic configuration](Install.md)
of the Minikube cluster.

Minikube must have been started, which will be the case if you have just completed installation. If not, this should
do it, adjusting the Colima memory and CPU count as appropriate for your system:

```shell
colima start -c6 -m16 --edit
minkube start
```

## Build and manage with `make`

**IMPORTANT:** The following assumes that your current working directory is the [`k8s-istio-poc/login`](../login)
directory.

The easiest way to build, start and manage the `login` service is using `make`. Running `make` without any
parameters will display the following usage help:

```text
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

**IMPORTANT:** The following assumes that your current working directory is the [`k8s-istio-poc/login`](../login)
directory.

In a terminal shell, run the following:

```shell
# Configure environment variables to instruct docker to target the Minikube
# package store
eval $(minikube docker-env)

# Build the docker container image
docker build -t login:v1 .

# Create the Minikube application deployment
kubectl create -f deployment.yaml

# Expose the application as a NodeType service that can be accessed from
# outside the cluster
kubectl apply -f service.yaml
```

## Connect to the service from a browser

Either create a tunnel to the service and open a browser to it in one step with:

```shell
minikube service login
```

This will open a browser and show a 404 error. Add `/login` or `/logout` to the URL to get a 200 response.

Or, use `kubectl` to forward the service port and then open a browser on http://localhost:9080/login

```shell
kubectl port-forward service/login 9080:8080
```

On a Mac, both approaches will lock up your terminal shell until you Ctrl-C to shutdown the tunnel.

**NOTE:** The service will not be useful in this state beyond sending you to an unreachable `/dashboard` URL. You 
will need to complete the Istio ingress gateway configuration before you get to see any of the services fully function
as intended.