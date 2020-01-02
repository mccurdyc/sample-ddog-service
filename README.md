https://app.datadoghq.com/account/settings#agent/kubernetes

has to use docker-compose because datadog recognizes all of the applications in
minikube as the same host.

1. docker-compose up --build

---

1. install and start minikube
https://matthewpalmer.net/kubernetes-app-developer/articles/guide-install-kubernetes-mac.html
  1. brew update
  2. brew install minikube
  3. brew install kubectl
  4. brew install hyperkit

`minikube start --vm-driver=hyperkit`

2. start datadog monitoring via `helm` chart
https://github.com/helm/charts/tree/master/stable/datadog#configuration

```bash
helm upgrade --install -f datadog-values.yaml datadog-monitoring stable/datadog --version 1.38.7
```

3. check agent status
`kubectl exec -it datadog-monitoring-xxxx agent status`

4. `kubectl create -f sample-deploy.yaml`
5. `kubectl create -f sample-service.yaml`
