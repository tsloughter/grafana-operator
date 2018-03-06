FROM ubuntu:14.04
MAINTAINER Tristan Sloughter <t@crashfast.com>

RUN addgroup --system kube-operator && adduser --system --gecos kube-operator kube-operator

USER kube-operator

COPY bin/linux/grafana-operator .

ENTRYPOINT ["./grafana-operator"]
