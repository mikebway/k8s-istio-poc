# Enable Metrics Collection

The `kubctl top` command is a useful tool for monitoring the resource usage of your Kubernetes cluster. It can be 
used to monitor the CPU and memory usage of your cluster's nodes and pods.

To install metrics collection services in your cluster, run the following command:

```shell
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

In theory, you should now be able to run the following command to see the CPU and memory usage of your cluster's node(s):

```shell
kubectl top nodes
```

which will display something like this:

```text
NAME       CPU(cores)   CPU%   MEMORY(bytes)   MEMORY%
minikube   435m         7%     2737Mi          22%
```

The Kubernetes web dashboard will also display CPU and memory graphics for both nodes and pods after metrics collection
has been installed.

## Troubleshooting

If `kubectl top nodes` displays `error: Metrics API not available`, follow the instructions in the 
[Fix “error: Metrics API not available” in Kubernetes](https://medium.com/@cloudspinx/fix-error-metrics-api-not-available-in-kubernetes-aa10766e1c2f)
article on [Medium](https://medium.com).