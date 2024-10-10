# Configuring the Istio Service Mesh

## Make Istio responsible for the `authtest` namespace

Istio will only install Envoy sidecars in pods that belong to a namespace that has been labeled so, to have
Envoy proxy added to the [`authtest`](../authtest) and [`login`](../login) services we have to add a label to the
`authest` namespace:

```shell
kubectl label namespace authtest istio-injection=enabled
```

## Verify that Envoy has not already been deployed to the `authtest` pods

If you run this command:

```shell
 kubectl get pods  --namespace authtest
```

You should see something like this:

```text
NAME                        READY   STATUS    RESTARTS   AGE
authtest-68fb855b64-gtsjm   1/1     Running   0          5m46s
login-67b884b995-9c84t      1/1     Running   0          5m22s
```

Where the `1/1` indicates that each of the two pods is running 1 container of the 1 containers defined by the pod.

## Restart the services

Istio will not add the Envoy proxy to pods that are already running, so we must stop and restart the
[`authtest`](../authtest) and [`login`](../login) services:

```shell
# Shut the authtest and login pods down
kubectl scale --replicas=0 --namespace=authtest deployment/authtest
kubectl scale --replicas=0 --namespace=authtest deployment/login

# Delete the authtest and login services and deployments
kubectl delete -f authtest/service.yaml --namespace authtest
kubectl delete -f authtest/deployment.yaml --namespace authtest
kubectl delete -f login/service.yaml --namespace authtest
kubectl delete -f login/deployment.yaml --namespace authtest

# ... allow a little time for the pods to shut down ...

# Redeploy and start the authtest and login services 
kubectl create -f authtest/deployment.yaml --namespace authtest
kubectl apply -f authtest/service.yaml --namespace authtest
kubectl create -f login/deployment.yaml --namespace authtest
kubectl apply -f login/service.yaml --namespace authtest
```

## Verify that Envoy is now present in the `authtest` pods

Now, if you run this command:

```shell
 kubectl get pods  --namespace authtest
```

You should see something like this where each pod expects to have two running containers:

```text
NAME                        READY   STATUS    RESTARTS   AGE
authtest-5fd9dd7d5c-qzf4s   2/2     Running   0          75s
login-884d47dd-6j22x        2/2     Running   0          19s
```

Meanwhile, running this:

```shell
istioctl proxy-status
```

Should display something like this, with four proxies running for each of cluster ingress, cluster egress, the authtest
service pod, and the login service pod:

```text
NAME                                                   CLUSTER        CDS              LDS              EDS              RDS              ECDS        ISTIOD                      VERSION
authtest-68fb855b64-gtsjm.authtest                     Kubernetes     SYNCED (23m)     SYNCED (23m)     SYNCED (23m)     SYNCED (23m)     IGNORED     istiod-5cc65d99d5-s2rdf     1.23.1
istio-egressgateway-55d6d944d7-6txh5.istio-system      Kubernetes     SYNCED (23m)     SYNCED (23m)     SYNCED (23m)     IGNORED          IGNORED     istiod-5cc65d99d5-s2rdf     1.23.1
istio-ingressgateway-7968d6d777-26xll.istio-system     Kubernetes     SYNCED (23m)     SYNCED (23m)     SYNCED (23m)     SYNCED (23m)     IGNORED     istiod-5cc65d99d5-s2rdf     1.23.1
login-67b884b995-9c84t.authtest                        Kubernetes     SYNCED (23m)     SYNCED (23m)     SYNCED (23m)     SYNCED (23m)     IGNORED     istiod-5cc65d99d5-s2rdf     1.23.1
```

