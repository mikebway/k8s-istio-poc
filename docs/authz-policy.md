# Configure Ingress to use the External Authorization Filter

The final step in configuring our Minikube cluster is to instruct the the Istio ingress gateway to have most if not
all requests assessed and modified by the [external authorization filter]() that we just deployed from the 
[`extauth`](../extauth) directory. 

This is a achieved by configuring the ingress gateway with an [Authorization Policy](https://istio.io/latest/docs/reference/config/security/authorization-policy/).
Specifically, this one: [../istio/authz-policy.yaml](../istio/authz-policy.yaml). The policy needs to reference the
[`extauth`](../extauth) service as an "extension provider." The extension provider must first be declared in the 
configuration map for the `istio-system` namespace.

The Istio documentation for the configuration map changes can be found here: [Define the external authorizer](https://istio.io/latest/docs/tasks/security/authorization/authz-custom/#define-the-external-authorizer).

The Istio documentation for configuring an authorization policy at the ingress gateway rather than within the service 
mesh can be found here: [Ingress Access Control](https://istio.io/latest/docs/tasks/security/authorization/authz-ingress/).

## Add the `extauth` service to the configuration map

Unfortunately, the `kubectl patch` command is a sledgehammer that is useless for making fine grained changes to 
Kubernetes configuration data maps so we have to do this the old fashioned way.

Open up an editor on the Istio configuration map with this command:

```shell
kubectl edit configmap/istio -n istio-system
```

You should see something that looks like this, with a number of "extension providers" already present:

```yaml
apiVersion: v1
data:
  mesh: |-
    accessLogFile: /dev/stdout
    defaultConfig:
      discoveryAddress: istiod.istio-system.svc:15012
    defaultProviders:
      metrics:
      - prometheus
    enablePrometheusMerge: true
    extensionProviders:
    - envoyOtelAls:
        port: 4317
        service: opentelemetry-collector.observability.svc.cluster.local
      name: otel
    - name: skywalking
      skywalking:
        port: 11800
        service: tracing.istio-system.svc.cluster.local
    - name: otel-tracing
      opentelemetry:
        port: 4317
        service: opentelemetry-collector.observability.svc.cluster.local
    rootNamespace: istio-system
    trustDomain: cluster.local
  meshNetworks: 'networks: {}'
kind: ConfigMap
metadata:
  creationTimestamp: "2024-09-16T14:11:43Z"
  labels:
    install.operator.istio.io/owning-resource: installed-state
    install.operator.istio.io/owning-resource-namespace: istio-system
    istio.io/rev: default
    operator.istio.io/component: Pilot
    operator.istio.io/managed: Reconcile
    operator.istio.io/version: 1.23.1
    release: istio
  name: istio
  namespace: istio-system
  resourceVersion: "820"
  uid: b8e604f4-4dc0-4596-ae74-2e435cc8c0b4
```

Copy and paste the following lines, inserting them at the end of the `extensionProviders:` section. For the sample 
shown above, that would be immediately above the `rootNamespace: istio-system` line.

```text
    - name: ext-authz-grpc
      envoyExtAuthzGrpc:
        service: extauth.default.svc.cluster.local
        port: 50051
```

Save the change and exit the editor. 

To verify that the patch applied, check the `extensionProviders:` section in the configuration map dump obtained with 
this command:

```shell
kubectl get configmap istio -n istio-system -o yaml
```

## Set `CUSTOM` authorization policy to reference the `extauth` extension 

To apply the [authz-policy.yaml](../istio/authz-policy.yaml), run the following from the project root directory:

```shell
kubectl apply -f istio/authz-policy.yaml
```

If you need to roll that back, you can do so with:

```shell
kubectl delete -f istio/authz-policy.yaml
```

## Viewing the access logs

This will be especially useful if the external authorization filter does not work out-of-the-box! You can dump
the Envoy logs from the ingress gateway with this command:

```shell
kubectl logs -l app=istio-ingressgateway -n istio-system --tail 100
```

To continuously follow the live log, replace `--tail 100` with `-f`:

```shell
kubectl logs -l app=istio-ingressgateway -n istio-system -f
```

It probably won't yield much useful information, but you can bump the logging detail by first finding the name of the 
currently active ingress pod and then setting the `debug` log level as follows ...

1. Find the ingress pod name with the `istioctl proxy-status` command:
   ```shell
   $ istioctl proxy-status
   NAME                                                   CLUSTER        CDS    ...
   authtest-68fb855b64-gtsjm.authtest                     Kubernetes     SYNCED ...
   istio-egressgateway-55d6d944d7-6txh5.istio-system      Kubernetes     SYNCED ...
   istio-ingressgateway-7968d6d777-26xll.istio-system     Kubernetes     SYNCED ...
   login-67b884b995-9c84t.authtest                        Kubernetes     SYNCED ...
   ```
   What you want to pick up is the `istio-ingressgateway-7968d6d777-26xll` value where the `7968d6d777-26xll` portion
   will be unique to your immediate context.
2. Substitute your value in the following to bump the logging level out of the ingress gateway up:
   ```shell
   istioctl proxy-config log istio-ingressgateway-7968d6d777-26xll -n istio-system --level debug
   ```
   
To revert, repeat the `istioctl proxy-config log` command but with the level set to `info`.

You can use `trace` in place of `debug` but, be warned, you will get a lot of noise in the log file from all the
background activity goin on inside Envoy and this can bury the few lines relating to requests that you submit.
