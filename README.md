# connectivity-listener

This service implements two ping listeners, one on HTTP and one on GRPC.

The ping endpoints have been generated with `protoc`.

```
sudo apt install -y protobuf-compiler
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/ping/ping.proto
```

(this is for historical accuracy, does not need to be executed again)


## Prerequisites

* `go`
* `docker`

I have used `minikube` as my local kubernetes dev environment. I assume the listener works in other kubernetes implementations
as well, but the Install and run commands need to be adapted.

```
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube_latest_amd64.deb
sudo dpkg -i minikube_latest_amd64.deb
alias kubectl='minikube kubectl --'
```


## Install and run the service


```
alias kubectl='minikube kubectl --'
docker build -t connectivity-listener:latest .
minikube image load connectivity-listener:latest
kubectl apply -f k8s/deployment.yaml
kubectl get pods -A
```
