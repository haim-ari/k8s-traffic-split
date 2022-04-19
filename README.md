# TrafficSplit

## Create local Cluster with k3d

### Install `k3d` cli

https://k3d.io/v5.3.0/#install-script

`brew install k3d`

### Create config file with traefik disabled

```yaml
kind: Simple
apiVersion: k3d.io/v1alpha2
name: playground
image: rancher/k3s:v1.20.8-k3s1
servers: 3
agents: 3
ports:
- port: 8080:80
  nodeFilters:
  - loadbalancer
options:
  k3s:
    extraServerArgs:
    - --no-deploy=traefik
```

### Create the cluster

`k3d cluster create -c k3d.yaml`


## Install Linkerd

https://linkerd.io/2.11/getting-started/

```bash
curl --proto '=https' --tlsv1.2 -sSfL https://run.linkerd.io/install | sh

linkerd check --pre

linkerd install --ha --set proxyInit.runAsRoot=true | kubectl apply -f -

linkerd check

linkerd viz install --ha | kubectl apply -f -
```

Or install with external prometheus

```
linkerd viz install --ha  --set prometheusUrl=http://prometheus-prometheus.prometheus:9090,prometheus.enabled=false | kubectl apply -f
```

```
linkerd check
```

Start the dashboard

```bash
linkerd viz dashboard
```

## Install SMI extension

https://linkerd.io/2.11/tasks/linkerd-smi/

```bash
linkerd smi install | kubectl apply -f -

linkerd smi check
```

## Install Ingress
https://kubernetes.github.io/ingress-nginx/deploy/

Create `ingress-nginx` ns and add annotation:


```bash
helm upgrade --install ingress-nginx ingress-nginx \
  --repo https://kubernetes.github.io/ingress-nginx \
  --namespace ingress-nginx --create-namespace
```

### Inject linkerd to ingress controller deployment

```bash
kubectl get deployment ingress-nginx-controller -n ingress-nginx -o yaml | linkerd inject  - | kubectl apply -f -
```

## Deploy the apps Charts

```bash
helm upgrade --install stg-demo . -f values.yaml -f app1-values.yaml
```

```bash
helm upgrade --install poc-demo . -f values.yaml -f app2-values.yaml
```

```bash
helm upgrade --install prd-demo . -f values.yaml -f app3-values.yaml
```

## Create the TrafficSplit

```bash
kubectl apply -f trafficsplit.yaml
```

## Load test the service internally

```bash
kubectl run -i --tty load-generator --rm --image=busybox --restart=Never -- /bin/sh -c "while sleep 0.01; do wget -q -O- http://prd-demo; done"
```

```bash
kubectl run -i --tty load-generator --rm --image=ghcr.io/six-ddc/plow --restart=Never -- "http://prd-demo" "-c" "100" "-d" "1m"
```

## Installing Linkerd with Helm

[Installing Linkerd with Helm](https://linkerd.io/2.11/tasks/install-helm/)

### Create SSL Secrets

```bash
step certificate create root.linkerd.cluster.local ca.crt ca.key \
--profile root-ca --no-password --insecure
```

```bash
step certificate create identity.linkerd.cluster.local issuer.crt issuer.key \
--profile intermediate-ca --not-after 8760h --no-password --insecure \
--ca ca.crt --ca-key ca.key
```

### Create Sealed Secrets

```bash
 kubectl -n linkerd create secret tls linkerd-tls-secret \
  --cert=issuer.crt \
  --key=issuer.key \
  --dry-run=client --output=yaml > linkerd-tls-secret.yaml
```

```bash
kubeseal -o yaml < linkerd-tls-secret.yaml >  linkerd-tls-sealed-secret.yaml
```

```bash
kubectl -n linkerd create secret generic linkerd-ca \
  --from-file=ca.crt \
  --dry-run=client \
  --output=yaml > linkerd-ca-secret.yaml
```

```bash
kubeseal -o yaml < linkerd-ca-secret.yaml >  linkerd-ca-sealed-secret.yaml
```
## Additional references 

https://linkerd.io/going-to-production/

https://linkerd.io/2.11/reference/proxy-configuration/

