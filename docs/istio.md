# Configuring the Istio Service Mesh and Ingres

## Install the Minikube Istio addon

After installing and starting Colima and Minikube, and after building and deploying the [`authtest`](svc-authtest.md)
and [`login`](svc-login.md) services, run the commands below from the  [`k8s-istio-poc/istio`](../istio) directory.

**NOTE:** The services do not have to be installed to install Istio, but by doing that first you will be able to 
immediately test the Istio configuration at the end of this page.

The following is taken from the official Minikube [Using the Istio Addon](https://minikube.sigs.k8s.io/docs/handbook/addons/istio/)
documentation:

```shell
# Enable Istio
minikube addons enable istio-provisioner
minikube addons enable istio  
```

Confirm that installation was successful with:

```shell
kubectl get pods -n istio-system
```

You should see two `istio-...` pods running. If nothing shows, start debugging :-)

## Configure Istio

For some explanation of what is going on here, i.e. about gateways and virtual services, see 
[Setup Istio Ingress Traffic Management on Minikube](https://medium.com/codex/setup-istio-ingress-traffic-management-on-minikube-725c5e6d767a)
by [Rocky Chen](https://medium.com/@rocky-chen).

```shell
# Define the ingress gateway
kubectl apply -f istio/gateway.yaml

# Define the virtual service
kubectl apply -f istio/virtual-service.yaml
```


## To remove Istio

```shell
# Disable Istio
minikube addons disable istio-provisioner
minikube addons disable istio  
```

