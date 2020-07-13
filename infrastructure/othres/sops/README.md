# SOPS

Git Repository: [mozilla/sops](https://github.com/mozilla/sops)
Getting Started: [mozilla/sops](https://github.com/mozilla/sops)

その他参考リンク

-

## どんなもの?
- JSONやYAMLなどのファイルの値のみを暗号化する
- GitOps などを行いたい場合に、secret ファイルをレポジトリに直接配置することを回避できる

## 特徴
- 値のみを暗号化するので、コンフィグファイルの可読性を損なわない
- AWS, Azure, GCP など複数プロバイダの鍵を暗号化に使用できる

## 周辺知識
競合
- [bitnami-labs/sealed-secrets](https://github.com/bitnami-labs/sealed-secrets)
  - SOPSとは異なり、秘密鍵・証明書をクラスタ内に Secret として保持する

## 動かす

### Install SOPS

```
brew install sops
```

```
$ sops -v
sops 3.5.0 (latest)
```

### 暗号化の対象になるファイルを作成

```
$ mkdir config
$ cd config
$ touch dev.yaml
```

```dev.yaml
user: dev-user
pass: dev-pass
```

### 暗号化のための鍵を作成

参考: [Google: google_kms_secret - Terraform by HashiCorp](https://www.terraform.io/docs/providers/google/d/kms_secret.html)

- 簡単化のために
  - 3環境用のモジュール化は見逃してください…
  - remote state は見逃してください…

```
$ cd ..
$ mkdir terraform
$ cd terraform
$ touch kms_dev.tf
```

```terraform:terraform/kms_dev.tf
resource "google_kms_key_ring" "key_ring" {
  name     = "sops-key-ring"
  location = "asia-northeast1-b"
}

resource "google_kms_crypto_key" "my_crypto_key" {
  name     = "sops-crypto-key"
  key_ring = google_kms_key_ring.my_key_ring.self_link
}
```

プロバイダ用のファイルも作成

```terraform:terraform/provider.tf
provider "google" {
  project = "kaito2"
  region  = "asia-northeast1"
}

// cloud kms を有効化
resource "google_project_service" "kms" {
  service = "cloudkms.googleapis.com"
}
```

```
$ terraform init
$ terraform plan
...
Plan: 2 to add, 0 to change, 0 to destroy.
...
$ terraform apply
...
Apply complete! Resources: 3 added, 0 changed, 0 destroyed.
```

```
$ gcloud kms keyrings list --location asia-northeast1
NAME
projects/kaito2/locations/asia-northeast1/keyRings/sops-key-ring

$ gcloud kms keys list --keyring sops-key-ring --location asia-northeast1
NAME                                                                                         PURPOSE          ALGORITHM                    PROTECTION_LEVEL  LABELS  PRIMARY_ID  PRIMARY_STATE
projects/kaito2/locations/asia-northeast1/keyRings/sops-key-ring/cryptoKeys/sops-crypto-key  ENCRYPT_DECRYPT  GOOGLE_SYMMETRIC_ENCRYPTION  SOFTWARE                  1           ENABLED
```

ちゃんと鍵が作成されていることが確認できました。

### 対象ファイルを暗号化

鍵を指定せずに実行すると怒られます。

```
$ cd ../config
$ sops -e dev.yaml
config file not found and no keys provided through command line options
```

GCP の場合は `--gcp-kms` オプションにリソースIDを渡す必要がある（or `SOPS_GCP_KMS_IDS` という名前の環境変数に設定する）。

リソースIDとは: [Object hierarchy &nbsp;|&nbsp; Cloud KMS Documentation &nbsp;|&nbsp; Google Cloud](https://cloud.google.com/kms/docs/object-hierarchy#key_resource_id)

コマンドを以下のように変更
（何も指定しないと標準出力に暗号化されたファイルの内容が出力される。）

```
$ sops -e --gcp-kms projects/kaito2/locations/asia-northeast1/keyRings/sops-key-ring/cryptoKeys/sops-crypto-key dev.yaml
user: ENC[AES256_GCM,data:8Hug2XT/PP8=,iv:X6/VCucFyTtqaspP20OOKbBm43AuFvVCN+mA7uReKTw=,tag:GTTIklpZw+kfQK+eamw1bQ==,type:str]
pass: ENC[AES256_GCM,data:Fw5s9NJVxY8=,iv:JTkKRLiSKa7i2o+Ek20Dq6+YUpxdUWPm2halpUtzi6o=,tag:4l57RVbGH1oaKdO+4zE80A==,type:str]
sops:
    kms: []
    gcp_kms:
    -   resource_id: projects/kaito2/locations/asia-northeast1/keyRings/sops-key-ring/cryptoKeys/sops-crypto-key
        created_at: '2020-07-13T11:42:45Z'
        enc: CiQAIyu+rA1vIggweC/+b2fWXzFVDlRYL9utDh9znmBKOoip8VcSSQCqguh30vH6TCl/mhT+AgucZ2sIRJwthwXk+M1GIzOIDvFoXCUHY1VvVyfzaSXCKhkw8boV4EW8QDqYEDyX1n6NVR5kIiSxSBQ=
    azure_kv: []
    lastmodified: '2020-07-13T11:42:46Z'
    mac: ENC[AES256_GCM,data:BXAavE6OhzKSfdmRE5cQmvyEKga0pnnSOAQLYmnjdadlRTV/TuDhSlriYNenEFXw4+RCHurL2xZfY3cOoqwuct6lZz2ATm81yzqjtymuH63xFG80DqjJRgKFQ+wBuggLPasKyWq8PHUQEr/0CEVZ/u9LgVO3CoJS5gWNISLzrow=,iv:KnaYCHB8r+clN1iVn0yCJnAKfI3r7h/H+sSCji++/Ic=,tag:/OzY47i++dfBO3VkdBiaUQ==,type:str]
    pgp: []
    unencrypted_suffix: _unencrypted
    version: 3.5.0
```

コマンドごとに `--gcp-kms` オプションを指定するのも面倒なので `.sops.yaml` という設定ファイルを設置することで sops がそれを読み込んでよしなにしてくれる。

See: [mozilla/sops](https://github.com/mozilla/sops#using-sops-yaml-conf-to-select-kms-pgp-for-new-files)

```yaml:config/.sops.yaml
creation_rules:
        - path_regex: \.dev\.yaml$
          gcp_kms: projects/kaito2/locations/asia-northeast1/keyRings/sops-key-ring/cryptoKeys/sops-crypto-key
```

コンフィグファイルを設定したので `--gcp-kms` オプションなしでも暗号化できるようになりました。

```
$ sops -e dev.yaml
user: ENC[AES256_GCM,data:cxkJBKyiQ4c=,iv:m00uQzik0OBw0tf5es4j1mZH+TM2JxynjifOhkW+ZHE=,tag:Jxtj4cfRe+TIdz3Vp9lw7w==,type:str]
pass: ENC[AES256_GCM,data:JUcITsSaKqc=,iv:cbNXopbCJPGtFbSkwn0TAXUKN0ojSjcMQrQJ5DpFOWo=,tag:i4/ey1aJh99PtZovpHU6Ww==,type:str]
sops:
    kms: []
    gcp_kms:
    -   resource_id: projects/kaito2/locations/asia-northeast1/keyRings/sops-key-ring/cryptoKeys/sops-crypto-key
        created_at: '2020-07-13T11:48:46Z'
        enc: CiQAIyu+rN+NkKTu3H/riwZ4NsNWh5SjmWiZQbp13qyaj8vPfdgSSQCqguh3JyE5RlKHVWKUkJ/cVRL4jNyKnIAMc6MUMZdeXURJ+I2tFuXRZRVbJpHKTf1CV4K0sKPdJyQwVvjdzbJAA2CrOZgdupw=
    azure_kv: []
    lastmodified: '2020-07-13T11:48:47Z'
    mac: ENC[AES256_GCM,data:/I+KswiR7VsJYDZJMwbUn1M/ygbLmbIdCuUVKR9uaZ+SFSudpU0zqX0mLg9R9jz5fB8t7JkcAPkXK5QiWJywvnE5YjtU3IsjrydNb7lEa630tqB0pEP8J4co+dGh9MzvyDQsldvwrREx9Ln8B3RfdyPaufYv7RfVqwBbWkgIYv0=,iv:g2jiRLD0wJYcj1OmhOzqDxOu5DWCcPTml9ROniyfXtQ=,tag:xXemqikO8oz8N8tlPUM+6A==,type:str]
    pgp: []
    unencrypted_suffix: _unencrypted
    version: 3.5.0
```

SOPS はデフォルトですべての値を暗号化しますが、例えば `user` の値は暗号化したくないなどの要件がある場合はオプションによって実現できます。

See: [mozilla/sops](https://github.com/mozilla/sops#encrypting-only-parts-of-a-file)

`.sops.yaml` を以下のように書き直します

```yaml:config/.sops.yaml
creation_rules:
  - path_regex: dev\.yaml
    encrypted_regex: pass
    gcp_kms: projects/kaito2/locations/asia-northeast1/keyRings/sops-key-ring/cryptoKeys/sops-crypto-key
```

再度暗号化を実行すると、 `user` の値は暗号化されておらず、`pass` の値は暗号化されていることが確認できます。

```
$ sops -e dev.yaml
user: dev-user
pass: ENC[AES256_GCM,data:nPihbY9Qe+o=,iv:naSvswd7v1Kdnq7VFfj2A3KT3SRV6HtLp2lTHFro3Ks=,tag:dtRf4pW0grYJqWBoNzuS7g==,type:str]
sops:
    kms: []
    gcp_kms:
    -   resource_id: projects/kaito2/locations/asia-northeast1/keyRings/sops-key-ring/cryptoKeys/sops-crypto-key
        created_at: '2020-07-13T11:53:47Z'
        enc: CiQAIyu+rFdLYVJKYgTiwH9L2sGwet7ugxWIlUPj39YhyyRR5ZISSQCqguh3cc4CB2tIC8KUhCa+5SXIw7HTpx1NP/FVBRIsXdscXHMnwn0FEgx2uH8MBrLGyv+Opb+WnomXWBZoBO+xQdDNOx5uvoI=
    azure_kv: []
    lastmodified: '2020-07-13T11:53:48Z'
    mac: ENC[AES256_GCM,data:liFOktDxSb/phDHyha8ollNe6fXlZkQYfhxiBPaKZJDuUmDBBw6d7DrQIF/lg1b4FBundJ5l9SW/m5PK58yeLNbzTqTwl/ACXH69M4cE+GQZqi4HRYdMa5hfe7VTwYn+TbKL3DHZFavPWl/Nih6by5PL8ooJds4Z7+Ey28Z5M0A=,iv:0HM9a2wrxJHV9ZALxQYoc5s6Kd2miADVwax8hAig29k=,tag:ABAxVKdOI/tJEt6JIvwIXA==,type:str]
    pgp: []
    encrypted_regex: pass
    version: 3.5.0
```

`-i` オプションを追加するとファイルを直接暗号化して置き換えます。

```
$ sops -e -i dev.yaml
```
