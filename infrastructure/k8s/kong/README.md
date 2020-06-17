# Kong

Git Repository: [GitHub - Kong/kong: 🦍 The Cloud-Native API Gateway](https://github.com/kong/kong)
Getting Started:

- [kubernetes-ingress-controller/k4k8s.md at master · Kong/kubernetes-ingress-controller · GitHub](https://github.com/Kong/kubernetes-ingress-controller/blob/master/docs/deployment/k4k8s.md)
- [kubernetes-ingress-controller/getting-started.md at master · Kong/kubernetes-ingress-controller · GitHub](https://github.com/Kong/kubernetes-ingress-controller/blob/master/docs/guides/getting-started.md)

## どんなもの?

マイクロサービスに適した API Gateway

## 特徴

Gateway としてトラフィックをインターセプトし、様々な処理ができるそう。

e.g.

- Authentication
- Traffic Control
- Analytics
- Transoformations
- Logging
- Serverless

Nginx の Lua 拡張で書かれているっぽい

## 周辺知識

- kustomize
- kubectl の パッチ

## 動かす

kind でクラスタを作成

```
kind create cluster
```

kustomize を使って CRD をデプロイ

[kubernetes-ingress-controller/k4k8s.md at master · Kong/kubernetes-ingress-controller · GitHub](https://github.com/Kong/kubernetes-ingress-controller/blob/master/docs/deployment/k4k8s.md)

```
kubectl create namespace kong
kubectl apply -k github.com/kong/kubernetes-ingress-controller/deploy/manifests/base
```

```
$ kubens kong
$ kubectl port-forward svc/kong-proxy 8080:80
Forwarding from 127.0.0.1:8080 -> 8000
Forwarding from [::1]:8080 -> 8000
```

[kubernetes-ingress-controller/getting-started.md at master · Kong/kubernetes-ingress-controller · GitHub](https://github.com/Kong/kubernetes-ingress-controller/blob/master/docs/guides/getting-started.md)

```
$ curl -i localhost:8080
HTTP/1.1 404 Not Found
Date: Wed, 17 Jun 2020 14:27:36 GMT
Content-Type: application/json; charset=utf-8
Connection: keep-alive
Content-Length: 48
X-Kong-Response-Latency: 1
Server: kong/2.0.4

{"message":"no Route matched with those values"}
```

```
kubens default
kubectl apply -f https://bit.ly/echo-service
```

```
$ kubectl get po
NAME                   READY   STATUS    RESTARTS   AGE
echo-78b867555-sk4jz   1/1     Running   0          18s
$ kubectl get svc
NAME         TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)           AGE
echo         ClusterIP   10.105.229.4   <none>        8080/TCP,80/TCP   40s
kubernetes   ClusterIP   10.96.0.1      <none>        443/TCP           11m
```

Ingress をデプロイ

```
$ echo "
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: demo
spec:
  rules:
  - http:
      paths:
      - path: /foo
        backend:
          serviceName: echo
          servicePort: 80
" | kubectl apply -f -
```

正常にアクセスできるようになった。

```
$ curl -i localhost:8080/foo
HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Transfer-Encoding: chunked
Connection: keep-alive
Date: Wed, 17 Jun 2020 14:36:21 GMT
Server: echoserver
X-Kong-Upstream-Latency: 2
X-Kong-Proxy-Latency: 0
Via: kong/2.0.4



Hostname: echo-78b867555-sk4jz

Pod Information:
	node name:	kind-control-plane
	pod name:	echo-78b867555-sk4jz
	pod namespace:	default
	pod IP:	10.244.0.6

Server values:
	server_version=nginx: 1.12.2 - lua: 10010

Request Information:
	client_address=10.244.0.5
	method=GET
	real path=/foo
	query=
	request_version=1.1
	request_scheme=http
	request_uri=http://localhost:8080/foo

Request Headers:
	accept=*/*
	connection=keep-alive
	host=localhost:8080
	user-agent=curl/7.64.1
	x-forwarded-for=127.0.0.1
	x-forwarded-host=localhost
	x-forwarded-port=8000
	x-forwarded-proto=http
	x-real-ip=127.0.0.1

Request Body:
	-no body in request-
```

### プラグインを入れる

KongPlugin リソースをセットアップ

```
echo "
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: request-id
config:
  header_name: my-request-id
plugin: correlation-id
" | kubectl apply -f -
kongplugin.configuration.konghq.com/request-id created
```

プラグインを適用した新しい Ingress を作成

```
$ echo "
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: demo-example-com
  annotations:
    konghq.com/plugins: request-id
spec:
  rules:
  - host: example.com
    http:
      paths:
      - path: /bar
        backend:
          serviceName: echo
          servicePort: 80
" | kubectl apply -f -
ingress.extensions/demo-example-com created
```

新しく作成した Ingress 経由でアクセス

```
$ curl -i -H "Host: example.com" localhost:8080/bar/sample
HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Transfer-Encoding: chunked
Connection: keep-alive
Date: Wed, 17 Jun 2020 14:43:06 GMT
Server: echoserver
X-Kong-Upstream-Latency: 1
X-Kong-Proxy-Latency: 0
Via: kong/2.0.4



Hostname: echo-78b867555-sk4jz

Pod Information:
	node name:	kind-control-plane
	pod name:	echo-78b867555-sk4jz
	pod namespace:	default
	pod IP:	10.244.0.6

Server values:
	server_version=nginx: 1.12.2 - lua: 10010

Request Information:
	client_address=10.244.0.5
	method=GET
	real path=/bar/sample
	query=
	request_version=1.1
	request_scheme=http
	request_uri=http://example.com:8080/bar/sample

Request Headers:
	accept=*/*
	connection=keep-alive
	host=example.com
	my-request-id=1e8bbe86-4b3c-41ff-8035-591a3346e925#1
	user-agent=curl/7.64.1
	x-forwarded-for=127.0.0.1
	x-forwarded-host=example.com
	x-forwarded-port=8000
	x-forwarded-proto=http
	x-real-ip=127.0.0.1

Request Body:
	-no body in request-

```

プラグインの機能で `my-request-id` が付与されている。

### Service にプラグインを適用する

KongPlugin リソースを作成

```
$ echo "
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: rl-by-ip
config:
  minute: 5
  limit_by: ip
  policy: local
plugin: rate-limiting
" | kubectl apply -f -
kongplugin.configuration.konghq.com/rl-by-ip created
```

Service にアノテーションを追加

```
kubectl patch svc echo \
  -p '{"metadata":{"annotations":{"konghq.com/plugins": "rl-by-ip\n"}}}'
```

rate limit が適用されている

```
$ curl -I localhost:8080/foo
HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Connection: keep-alive
Date: Wed, 17 Jun 2020 14:48:04 GMT
Server: echoserver
X-RateLimit-Remaining-Minute: 4
X-RateLimit-Limit-Minute: 5
RateLimit-Remaining: 4
RateLimit-Limit: 5
RateLimit-Reset: 56
X-Kong-Upstream-Latency: 0
X-Kong-Proxy-Latency: 1
Via: kong/2.0.4

$ curl -I localhost:8080/foo
HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Connection: keep-alive
Date: Wed, 17 Jun 2020 14:48:13 GMT
Server: echoserver
X-RateLimit-Remaining-Minute: 3
X-RateLimit-Limit-Minute: 5
RateLimit-Remaining: 3
RateLimit-Limit: 5
RateLimit-Reset: 47
X-Kong-Upstream-Latency: 1
X-Kong-Proxy-Latency: 1
Via: kong/2.0.4

$ curl -I localhost:8080/foo
HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Connection: keep-alive
Date: Wed, 17 Jun 2020 14:48:15 GMT
Server: echoserver
X-RateLimit-Remaining-Minute: 2
X-RateLimit-Limit-Minute: 5
RateLimit-Remaining: 2
RateLimit-Limit: 5
RateLimit-Reset: 45
X-Kong-Upstream-Latency: 1
X-Kong-Proxy-Latency: 1
Via: kong/2.0.4

$ curl -I localhost:8080/foo
HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Connection: keep-alive
Date: Wed, 17 Jun 2020 14:48:16 GMT
Server: echoserver
X-RateLimit-Remaining-Minute: 1
X-RateLimit-Limit-Minute: 5
RateLimit-Remaining: 1
RateLimit-Limit: 5
RateLimit-Reset: 44
X-Kong-Upstream-Latency: 0
X-Kong-Proxy-Latency: 1
Via: kong/2.0.4

$ curl -I localhost:8080/foo
HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Connection: keep-alive
Date: Wed, 17 Jun 2020 14:48:17 GMT
Server: echoserver
X-RateLimit-Remaining-Minute: 0
X-RateLimit-Limit-Minute: 5
RateLimit-Remaining: 0
RateLimit-Limit: 5
RateLimit-Reset: 43
X-Kong-Upstream-Latency: 0
X-Kong-Proxy-Latency: 1
Via: kong/2.0.4

$ curl -I localhost:8080/foo
HTTP/1.1 429 Too Many Requests
Date: Wed, 17 Jun 2020 14:48:18 GMT
Content-Type: application/json; charset=utf-8
Connection: keep-alive
Retry-After: 42
Content-Length: 37
X-RateLimit-Remaining-Minute: 0
X-RateLimit-Limit-Minute: 5
RateLimit-Remaining: 0
RateLimit-Limit: 5
RateLimit-Reset: 42
X-Kong-Response-Latency: 1
Server: kong/2.0.4
```

最後のアクセスが `429` になっている。

### 最終的な構成

```
HTTP requests with /foo -> Kong enforces rate-limit -> echo server

HTTP requests with /bar -> Kong enforces rate-limit +   -> echo-server
   on example.com          injects my-request-id header
```
