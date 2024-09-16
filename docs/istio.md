# Configuring Istio ingress

## Install `istioctl`

`'istioctl` is a command line interface for Istion that, among amny otehr things, can install Isto into
a Kubernetes cluster.

The easiest way to install `istioctl` on a Mac is to use Homebrew:

```shell
brew install istioctl
```

But you can do it the hard way if you prefer by following the official instructions at
[Download the Istio release](https://istio.io/latest/docs/setup/additional-setup/download-istio-release/).

### DO NOT USE `minikube addons enable`
Do not use the `minikube addons enable` approach to install Istio that you may find if you Google "Minikube Istio". 
This will work up to a point but you wont' have the `istioctl` command that you need to manage and monitor how
Istio functions.

## Install Istio into your Minikube cluster

Run the following to confirm whether or not you are missing any prerequisites:

```shell
istioctl analyze
```

Most likely, this will tell you that you need to run the following:

```shell
kubectl label namespace default istio-injection=enabled
```

Now you can deploy Istio to the cluster with this command:

```shell
istioctl install --set profile=demo -y
```

To understand what the canned `demo` profile gives you, see 
[Installation Configuration Profiles](https://istio.io/latest/docs/setup/additional-setup/config-profiles/).

Now confirm that installation was successful with:

```shell
kubectl get pods -n istio-system
```

You should see three `istio-...` pods running for: ingress, egress, and the Istio daemon. If nothing shows, start debugging :-)

## Configure Istio ingress routing

For some explanation of what is going on here, i.e. about gateways and virtual services, see 
[Setup Istio Ingress Traffic Management on Minikube](https://medium.com/codex/setup-istio-ingress-traffic-management-on-minikube-725c5e6d767a)
by [Rocky Chen](https://medium.com/@rocky-chen).

```shell
# Define the ingress gateway
kubectl apply -f istio/gateway.yaml

# Define the virtual service
kubectl apply -f istio/virtual-service.yaml
```

## Open a tunnel to the cluster

The ingress gateway on the cluster is now listening, but we still need to be able to connect from localhost to it.
We achieve that by having Minikube run a TCP tunnel (you will be prompted for your sudo password):

```shell
# Open a tunnel to the cluster network
minikube tunnel
```

The tunnel will remain loaded until you Ctrl-C to kill the app. You will probably want to open a second shell window
for this.

