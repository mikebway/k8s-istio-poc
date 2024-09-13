.PHONY: help
help: ## List of available commands
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: build
build: ## Build the login Docker container
	eval $(minikube docker-env) &\
 	docker build -t login:v1 .

.PHONY: deploy
deploy: ## Deploy the container as a pod/service in Kubernetes
	eval $(minikube docker-env) &\
	kubectl apply -f ../authtest/namespace.yaml &\
	kubectl create -f deployment.yaml --namespace authtest &\
	kubectl apply -f service.yaml --namespace authtest

.PHONY: undeploy
undeploy: ## Destroy the deployment/pod/service in Kubernetes
	kubectl delete -f service.yaml --namespace authtest &\
	kubectl delete -f deployment.yaml --namespace authtest

.PHONY: start
start: ## Start the service
	kubectl scale --replicas=1 -f deployment.yaml --namespace authtest

.PHONY: stop
stop: ## Stop the service
	kubectl scale --replicas=0 -f deployment.yaml --namespace authtest
