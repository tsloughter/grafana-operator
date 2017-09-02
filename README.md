# Grafana Operator

**not meant for production use**

This Operator is based on the `grafana-watcher` sidecar from [Prometheus Operator](https://github.com/coreos/prometheus-operator) and the [Rolebinding Operator](https://github.com/treacher/namespace-rolebinding-operator) example.

Currently it simply watches for new `ConfigMaps` and if they define the annotation `grafana.net/dashboards` as `"true"` it will `POST` each dashboard from the `ConfigMap` to Grafana.

## Usage
```
--run-outside-cluster # Uses ~/.kube/config rather than in cluster configuration
--grafana-url # Sets the URL and authentication to use to access the Grafana API
```

## Development

### Build from source
1. `make install_deps`
2. `make build`
3. `./bin/grafana-operator --run-outside-cluster 1 --grafana-url <GRAFANA URL>`

Easiest way to install just Grafana to Kubernetes for playing with helm: `helm install stable/grafana` then add the dashboards, `kubectl apply -f examples/grafana-dashboards.yaml`
