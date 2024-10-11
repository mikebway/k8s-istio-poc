# Install and Configure Docker, Colima, and Minikube

**IMPORTANT:** While it does not matter for running `brew` installations, etc., some of the following shell commands 
assume that your working directory is the `k8s-istio-poc` project root directory where they reference YAML file 
paths.

## Prerequisites

While this project can be used as a rough map for Windows and Linux installation, it's primary audience is
Mac users. As a Mac user, it is assumed that you have [Homebrew](https://brew.sh/) installed.

Again, as a Mac user you are almost certainly using Zsh as your command line shell. If so, and you have not done so
already, you might want to add this to your `.zshrc` file (or just run it on your command line as a one off):

```shell
# Ignore # comment lines in interactive command pastes
setopt interactivecomments
```

## Installation

1. Install Colima and Docker via Brew
   ```shell
   brew install colima docker docker-compose
   ```
   
2. Start Colima; adjust CPU and memory parameters for best use of your Mac’s resources
   ```shell
   colima start -c6 -m12 --edit
   ```

3. Configure Docker context to point at Colima
   ```shell
   docker context use colima
   ```
   
4. Install Minikube and configure to use Docker etc
   ```shell
   brew install minikube
   minikube config set driver docker
   
   # Adjust these values to be equal to or less than you used for 
   # Colima above
   minikube config set cpus 6
   minikube config set memory 12g
   ```
5. Start Minikube
   ```shell
   minikube start
   ```

## Configure Docker client to use Minikube image store

```shell
eval $(minikube docker-env)
```

**NOTE:** This only affects the current terminal shell behavior.

You can undo this environment variable configuration in the current shell with:

```shell
eval $(minikube docker-env -u) 
```

If in doubt, you can see if the Minikube environment is set with:

```shell
echo $MINIKUBE_ACTIVE_DOCKERD
```

If that comes up blank, Docker will not use the Minikube store. If it comes up `minikube`, it will.

## Confirm basic operation

1. Open the Kubernetes dashboard in a browser
   **IMPORTANT:** Run this in a separate shell window as it will b;ock while the dashboard is running.  
   ```shell
   minikube dashboard
   ```
   
2. Deploy a crude echo server
   ```shell
   kubectl create deployment hello-minikube --image=kicbase/echo-server:1.0
   kubectl expose deployment hello-minikube --type=NodePort --port=8080
   ```
   
3. View the echo server in a browser
   Either
   ```shell
   minikube service hello-minikube
   ```
   or
   ```shell
   kubectl port-forward service/hello-minikube 7080:8080
   ```

## Stopping and/or deleting the test deployment

We don't need the `hello-minikube` service anymore so lets get rid of it ...

* Stopping
  ```shell
  kubectl scale --replicas=0 deployment/hello-minikube
  ```
  
* Deleting
  ```shell
  kubectl delete deployment hello-minikube
  ```
  




