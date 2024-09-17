# Configure Ingress to use the External Authorization Filter

```diff
-  This page not yet complete
```

The final step in configuring our Minikube cluster is to instruct the the Istio ingress gateway to have most if not
all requests assessed and modified by the [external authorization filter]() that we just deployed from the 
[`extauth`](../extauth) directory. 

This is a achieved by configuring the ingress gateway with an [Authorization Policy](https://istio.io/latest/docs/reference/config/security/authorization-policy/).
Specifically, this one: [../extauth/authz-policy.yaml](../extauth/authz-policy.yaml).

The Istio documentation for configuring an authorization policy at the ingress gateway reather than within the service 
mesh can be found here: [Ingress Access Control](https://istio.io/latest/docs/tasks/security/authorization/authz-ingress/).

To apply the [authz-policy.yaml](../extauth/authz-policy.yaml), run the following from the project root directory:

```shell
kubectl apply -f extauth/authz-policy.yaml
```

If you need to roll that back, you can do so with:

```shell
kubectl delete -f extauth/authz-policy.yaml
```

