# Kubebuilder

Git Repository: [kubernetes-sigs/kubebuilder](https://github.com/kubernetes-sigs/kubebuilder)
Getting Started: [Quick Start - The Kubebuilder Book](https://book.kubebuilder.io/quick-start.html)

その他参考リンク

-

## どんなもの?

CRD(Custom Resource Definition)を用いた Kubernetes API を作成するためのフレームワーク

## 特徴

## 周辺知識

## 動かす

brew で入るっぽい

```
$ brew install kubebuilder
```

## Quick Start

### プロジェクト作成

```
$ mkdir example
$ cd example
```

`GOPATH` の外の場合は `go mod init` が必要らしい

```
$ go mod init github.com/kaito2/gettingstarted
```

初期化

```
$ kubebuilder init --domain my.domain
```

TODO: `--domain` ってなに…調べる

色々生えた。

```
$ tree .
.
├── Dockerfile
├── Makefile
├── PROJECT
├── bin
│   └── manager
├── config
│   ├── certmanager
│   │   ├── certificate.yaml
│   │   ├── kustomization.yaml
│   │   └── kustomizeconfig.yaml
│   ├── default
│   │   ├── kustomization.yaml
│   │   ├── manager_auth_proxy_patch.yaml
│   │   ├── manager_webhook_patch.yaml
│   │   └── webhookcainjection_patch.yaml
│   ├── manager
│   │   ├── kustomization.yaml
│   │   └── manager.yaml
│   ├── prometheus
│   │   ├── kustomization.yaml
│   │   └── monitor.yaml
│   ├── rbac
│   │   ├── auth_proxy_client_clusterrole.yaml
│   │   ├── auth_proxy_role.yaml
│   │   ├── auth_proxy_role_binding.yaml
│   │   ├── auth_proxy_service.yaml
│   │   ├── kustomization.yaml
│   │   ├── leader_election_role.yaml
│   │   ├── leader_election_role_binding.yaml
│   │   └── role_binding.yaml
│   └── webhook
│       ├── kustomization.yaml
│       ├── kustomizeconfig.yaml
│       └── service.yaml
├── go.mod
├── go.sum
├── hack
│   └── boilerplate.go.txt
└── main.go

9 directories, 30 files
```

### API を作成

API (group/version) を `webapp/1` として作成し、Kind(CRD) 名を `Guestbook` として作成

```
$ kubebuilder create api --group webapp --version v1 --kind Guestbook
# すべての質問に y 
```

以下のファイルが作成される。

- `api/v1/guestbook_types.go`
- `controller/guestbook_controller.go`

### テスト実行

KIND(Kubernetes in Dockerの方)で用意したクラスタでテストを実行する。

```
$ kind create cluster
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.18.2) 🖼
 ✓ Preparing nodes 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Thanks for using kind! 😊
$ kubectx kind-kind
Switched to context "kind-kind".
```

Kustomize が必要なのでインストール

```
brew install kustomize
```

CRD をクラスタにインストール

```
$ make install
/Users/kaito2/.ghq/bin/controller-gen "crd:trivialVersions=true" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
kustomize build config/crd | kubectl apply -f -
customresourcedefinition.apiextensions.k8s.io/guestbooks.webapp.my.domain created
$ kubectl get crd
NAME                          CREATED AT
guestbooks.webapp.my.domain   2020-07-19T07:31:52Z
```

Controller を実行する。
以下のコマンドはフォアグラウンドで実行されるので新しいタブを開いてください。

```
$ make run
/Users/kaito2/.ghq/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
go fmt ./...
go vet ./...
/Users/kaito2/.ghq/bin/controller-gen "crd:trivialVersions=true" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
go run ./main.go
2020-07-19T16:33:12.809+0900	INFO	controller-runtime.metrics	metrics server is starting to listen	{"addr": ":8080"}
2020-07-19T16:33:12.809+0900	INFO	setup	starting manager
2020-07-19T16:33:12.809+0900	INFO	controller-runtime.manager	starting metrics server	{"path": "/metrics"}
2020-07-19T16:33:12.809+0900	INFO	controller-runtime.controller	Starting EventSource	{"controller": "guestbook", "source": "kind source: /, Kind="}
2020-07-19T16:33:12.910+0900	INFO	controller-runtime.controller	Starting Controller	{"controller": "guestbook"}
2020-07-19T16:33:12.910+0900	INFO	controller-runtime.controller	Starting workers	{"controller": "guestbook", "worker count": 1}
2020-07-19T16:35:39.699+0900	DEBUG	controller-runtime.controller	Successfully Reconciled	{"controller": "guestbook", "request": "default/guestbook-sample"}
```

## Custom Resource のインスタンスをデプロイ

```
$ kubectl apply -f config/samples/
guestbook.webapp.my.domain/guestbook-sample created
```

## Run It On the Cluster

**skip**

## Uninstall CRDs

```
make uninstall
```
