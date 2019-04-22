build: main.go
	go build .

dockerize: build init.sh
	docker build . -t su225/istio-test:learning

k8s-deploy: k8s-deploy-a.yaml \
			k8s-deploy-b.yaml \
			k8s-deploy-c.yaml \
			dockerize
	which kubectl
	kubectl apply -f k8s-deploy-a.yaml
	kubectl apply -f k8s-deploy-b.yaml
	kubectl apply -f k8s-deploy-c.yaml

k8s-undeploy: istio-undeploy-traffic-settings
	which kubectl
	kubectl delete -f k8s-deploy-a.yaml
	kubectl delete -f k8s-deploy-b.yaml
	kubectl delete -f k8s-deploy-c.yaml

istio-deploy-traffic-settings: service-b-virtual-service.yaml k8s-deploy
	kubectl apply -f service-b-virtual-service.yaml

istio-undeploy-traffic-settings:
	kubectl delete -f service-b-virtual-service.yaml