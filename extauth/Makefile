.PHONY: help
help: ## List of available commands
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: build
build: ## Build the extauth Docker container
	eval $(minikube docker-env) &\
 	docker build -t extauth:v1 .

.PHONY: deploy
deploy: ## Deploy the container as a pod/service in Kubernetes
	eval $(minikube docker-env) &\
	kubectl create -f deployment.yaml --namespace istio-system &\
	kubectl apply -f service.yaml --namespace istio-system

.PHONY: undeploy
undeploy: ## Destroy the deployment/pod/service in Kubernetes
	kubectl delete -f service.yaml --namespace istio-system &\
	kubectl delete -f deployment.yaml --namespace istio-system

.PHONY: start
start: ## Start the service
	kubectl scale --replicas=1 -f deployment.yaml --namespace istio-system

.PHONY: stop
stop: ## Stop the service
	kubectl scale --replicas=0 -f deployment.yaml --namespace istio-system


