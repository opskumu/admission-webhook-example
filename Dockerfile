FROM centos:7

COPY pod-admission-webhook /pod-admission-webhook 
COPY config.yaml /config/config.yaml

WORKDIR /

ENTRYPOINT ["/pod-admission-webhook"]
