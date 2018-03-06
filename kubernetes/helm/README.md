Summary
=======
Support deploying `grafana-operator` to kubernetes via a [Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/).

Deploy
======
## Deploy Grafana
Deploy the current stable grafana chart
```
$ helm repo add stable https://kubernetes-charts.storage.googleapis.com
$ helm install stable/grafana --namespace=<NAMESPACE> --set server.ingress.enabled=true,server.adminPassword=<GRAFANA_ADMIN_PASSWORD>,server.ingress.hosts=[<INGRESS_URL>]
```

## Deploy grafana-operator
Now that a grafana instance has been deployed, the operator can be deployed with a reference to this instance-
``` bash
$ helm upgrade --install grafana-operator -f kubernetes/helm/values.yaml kubernetes/helm/.  --namespace=<NAMESPACE> --set grafana.url="<GRAFANA_URL>" --set grafana.auth.username=<GRAFANA_USERNAME> --set grafana.auth.username=<GRAFANA_PASSWORD>
```
## Create Dashboard
Create a dashboard, via a ConfigMap-
```
$ kubectl apply -f examples/grafana-dashboards.yaml
```
The operator will pick up the dashbaord defined in this ConfigMap and use the grafana api to create it in the Grafana instance.
