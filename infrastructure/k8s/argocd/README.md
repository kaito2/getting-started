# argocd

[Getting Started - Argo CD - Declarative GitOps CD for Kubernetes](https://argoproj.github.io/argo-cd/getting_started/)

## Log

ローカルクラスタは [kind](https://kind.sigs.k8s.io/) でたてる。

```
$ kind create cluster
$ kubectx kind-kind
```

### 本編

#### クラスタに Argo CD の CRD, Service, アプリケーションリソースなどを導入

```
$ kubectl create namespace argocd
$ kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
customresourcedefinition.apiextensions.k8s.io/applications.argoproj.io created
customresourcedefinition.apiextensions.k8s.io/appprojects.argoproj.io created
serviceaccount/argocd-application-controller created
serviceaccount/argocd-dex-server created
serviceaccount/argocd-server created
role.rbac.authorization.k8s.io/argocd-application-controller created
role.rbac.authorization.k8s.io/argocd-dex-server created
role.rbac.authorization.k8s.io/argocd-server created
clusterrole.rbac.authorization.k8s.io/argocd-application-controller created
clusterrole.rbac.authorization.k8s.io/argocd-server created
rolebinding.rbac.authorization.k8s.io/argocd-application-controller created
rolebinding.rbac.authorization.k8s.io/argocd-dex-server created
rolebinding.rbac.authorization.k8s.io/argocd-server created
clusterrolebinding.rbac.authorization.k8s.io/argocd-application-controller created
clusterrolebinding.rbac.authorization.k8s.io/argocd-server created
configmap/argocd-cm created
configmap/argocd-rbac-cm created
configmap/argocd-ssh-known-hosts-cm created
configmap/argocd-tls-certs-cm created
secret/argocd-secret created
service/argocd-dex-server created
service/argocd-metrics created
service/argocd-redis created
service/argocd-repo-server created
service/argocd-server-metrics created
service/argocd-server created
deployment.apps/argocd-application-controller created
deployment.apps/argocd-dex-server created
deployment.apps/argocd-redis created
deployment.apps/argocd-repo-server created
deployment.apps/argocd-server created
```

色々作成されている

```
$ k get pods --namespace argocd
NAME                                             READY   STATUS              RESTARTS   AGE
argocd-application-controller-5cfb8d686c-cmmw9   1/1     Running             0          45s
argocd-dex-server-5cf8dd69f5-hgp2b               0/1     PodInitializing     0          45s
argocd-redis-6d7f9df848-2xj2c                    0/1     ContainerCreating   0          44s
argocd-repo-server-56b75988dc-nxsnh              0/1     ContainerCreating   0          44s
argocd-server-6766455855-2v4pz                   0/1     ContainerCreating   0          44s
```

#### Argo CLI を取得

```
brew tap argoproj/tap
brew install argoproj/tap/argocd
```

#### Argo CD API Server にアクセスする

方法は以下の3つで
* ServiceTypeを `LoadBalancer` に変更する
  * パッチをあてる
  * [Getting Started - Argo CD - Declarative GitOps CD for Kubernetes](https://argoproj.github.io/argo-cd/getting_started/#service-type-load-balancer)
* Ingress を置く
  * [Ingress Configuration - Argo CD - Declarative GitOps CD for Kubernetes](https://argoproj.github.io/argo-cd/operator-manual/ingress/)
* Port Forwarding する

今回は動かして見るだけなので Port Forwarding にしておく

別タブで以下を実行

```
kubectl port-forward svc/argocd-server -n argocd 8080:443
```

#### Login Using The CLI

初期パスワードはPod名になるらしいので以下のコマンドで取得（Podのパスを取得して`'/'`で切ってるだけ）

```
$ kubectl get pods -n argocd -l app.kubernetes.io/name=argocd-server -o name | cut -d'/' -f 2
argocd-server-7777777777-XXXXX
```

ログインする。（Port Forward しているので `<ARGOCD_SERVER>` は `localhost:8080`）
認証については今回は `y` で。

```
argocd login localhost:8080
WARNING: server certificate had error: x509: certificate signed by unknown authority. Proceed insecurely (y/n)? y
Username: admin
Password:
'admin' logged in successfully
Context 'localhost:8080' updated
```

パスワードの変更

```
$ argocd account update-password
*** Enter current password:
*** Enter new password:
*** Confirm new password:
Password updated
Context 'localhost:8080' updated
```

#### Argo CD からデプロイする対象クラスタを作成

kind でもう一つ作っておく

kind で実行すると謎のエラーになるので一旦中止。

```
$ kind create cluster --name deploy-target
...
$ argocd cluster add kind-deploy-target
INFO[0000] ServiceAccount "argocd-manager" created in namespace "kube-system"
INFO[0000] ClusterRole "argocd-manager-role" created
INFO[0000] ClusterRoleBinding "argocd-manager-role-binding" created
FATA[0000] rpc error: code = Unknown desc = REST config invalid: Get "https://127.0.0.1:50136/version?timeout=32s": dial tcp 127.0.0.1:50136: connect: connection refused
````

ふつうに `docker-for-desktop` を対象にする。

```
$ argocd cluster add docker-for-desktop
INFO[0000] ServiceAccount "argocd-manager" created in namespace "kube-system"
INFO[0000] ClusterRole "argocd-manager-role" created
INFO[0000] ClusterRoleBinding "argocd-manager-role-binding" created
Cluster 'https://kubernetes.docker.internal:6443' added
```

#### Creating Apps Via UI

以下の手順で Argo CD がチュートリアル用に用意してくれてる[GitHub - argoproj/argocd-example-apps: Example Apps to Demonstrate Argo CD](https://github.com/argoproj/argocd-example-apps)レポジトリから[guestbook](https://github.com/argoproj/argocd-example-apps/tree/master/guestbook)をデプロイする。

* ブラウザで `localhost:8080` にアクセスする。
* sign in
* `+ New App` をクリック
* 以下の項目を入力
  * `Application Name`: `guestbook`
  * `Project`: `default`
  * `SYNC POLICY`: `Maual`
  * `Repository URL`: `https://github.com/argoproj/argocd-example-apps.git`
  * `Revision`: `HEAD`
  * `Path`: `guestbook`
  * `Cluster`: `https://kubernetes.docker.internal:6443`
  * `Namespace`: `default`
* `CREATE` をクリック

#### Sync (Deploy) The Application

```
$ argocd app sync guestbook
TIMESTAMP                  GROUP        KIND   NAMESPACE                  NAME    STATUS    HEALTH        HOOK  MESSAGE
2020-06-14T11:54:11+09:00            Service     default          guestbook-ui  OutOfSync  Missing
2020-06-14T11:54:11+09:00   apps  Deployment     default          guestbook-ui  OutOfSync  Missing
2020-06-14T11:54:12+09:00            Service     default          guestbook-ui    Synced  Healthy
2020-06-14T11:54:12+09:00            Service     default          guestbook-ui    Synced   Healthy              service/guestbook-ui created
2020-06-14T11:54:12+09:00   apps  Deployment     default          guestbook-ui  OutOfSync  Missing              deployment.apps/guestbook-ui created
2020-06-14T11:54:12+09:00   apps  Deployment     default          guestbook-ui    Synced  Progressing              deployment.apps/guestbook-ui created

Name:               guestbook
Project:            default
Server:             https://kubernetes.docker.internal:6443
Namespace:          default
URL:                https://localhost:8080/applications/guestbook
Repo:               https://github.com/argoproj/argocd-example-apps.git
Target:             HEAD
Path:               guestbook
SyncWindow:         Sync Allowed
Sync Policy:        <none>
Sync Status:        Synced to HEAD (6bed858)
Health Status:      Progressing

Operation:          Sync
Sync Revision:      6bed858de32a0e876ec49dad1a2e3c5840d3fb07
Phase:              Succeeded
Start:              2020-06-14 11:54:11 +0900 JST
Finished:           2020-06-14 11:54:12 +0900 JST
Duration:           1s
Message:            successfully synced (all tasks run)

GROUP  KIND        NAMESPACE  NAME          STATUS  HEALTH       HOOK  MESSAGE
       Service     default    guestbook-ui  Synced  Healthy            service/guestbook-ui created
apps   Deployment  default    guestbook-ui  Synced  Progressing        deployment.apps/guestbook-ui created
```

guestbook がデプロイされているのが確認できる。
サービスを確認して Port Forward する。

```
$ kctx docker-for-desktop
Switched to context "docker-for-desktop".
$ k get pods
NAME                            READY   STATUS    RESTARTS   AGE
guestbook-ui-65b878495d-84c4j   1/1     Running   0          108s
$ k get services
NAME           TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)   AGE
guestbook-ui   ClusterIP   10.97.66.115   <none>        80/TCP    3m43s
kubernetes     ClusterIP   10.96.0.1      <none>        443/TCP   19d
$ kubectl port-forward svc/guestbook-ui 8888:80
Forwarding from 127.0.0.1:8888 -> 80
Forwarding from [::1]:8888 -> 80
...
```

ブラウザで `localhost:8888` にアクセスするとデプロイされているのが確認できる

[![Image from Gyazo](https://i.gyazo.com/fa23774afae5cc8bcde09374a0ff63c4.png)](https://gyazo.com/fa23774afae5cc8bcde09374a0ff63c4)

Argo CD の UI 上で `DELETE` を選択すると guestbook アプリケーションが削除されたことを確認できる。

```
$ k get pods
No resources found in default namespace.
```
