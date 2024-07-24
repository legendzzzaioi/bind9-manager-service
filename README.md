# bind9-manager-service
This backend project provides a solution for managing Bind9 using the Go-zero framework. The solution aims to simplify the configuration and management of the Bind9 DNS server by offering API interfaces for creating, updating, deleting, and querying domains and records. Additionally, the project supports the persistence of the /data/ directory to ensure the longevity of configuration data.


### Usage

Build Docker Image

```
docker build -t legendzzzaioi/bind9-manager-service:v1 .
```

Run Docker Container

```
docker run -d \
  --restart always \
  -p 53:53/tcp \
  -p 53:53/udp \
  -p 127.0.0.1:953:953/tcp \
  -p 8000:8000 \
  --name bind9-manager-service \
  -v /data:/data \
  -d legendzzzaioi/bind9-manager-service:v1
```

Usage with Kubernetes

```
# modify bind9-manager-service.yaml
kubectl -n xxx apply -f bind9-manager-service.yaml
```
