# Containers the hard way: Gocker: A mini Docker written in Go

[Containers the hard way: Gocker: A mini Docker written in Go - Unixism](https://unixism.net/2020/06/containers-the-hard-way-gocker-a-mini-docker-written-in-go/)をやってみたメモ

## 概要

- コンテナ技術はアプリケーションをパッケージングする手法として有名になっているがその中で、Docker の立ち位置が誤解されている。
  - Docker はコンテナの作成・実行・削除・ネットワークなどの管理機能を提供するが、それ自体を実現しているのは OS のプリミティブな機能である。
  - とはいえ Linux のシステムコール一発でコンテナができるわけではなく、Linux の namespac や cgroups などの機能の集まりによって実現される。

## What is Goker

- Go で Docker の最低限の機能を実装したもの
  - create container
  - manage container image
  - excute process in existing container
  - etc.

![fed54ee9.png](https://unixism.net/wp-content/uploads/2020/06/Gocker.png)

## Gocker capabilities

- `gocker run <--cpus=cpus-max> <--mem=mem-max> <--pids=pids-max>`
- `gocker ps`
- `gocker exec </path/to/command>`
- `gocker images`
- `gocker rmi`

## Other capabilities

- Overlay file system
  - 複数コンテナ間でファイルししテムを共有できる
- コンテナは独立したネットワークネームスペースを得る
- リソース制限ができる

## Gocker container isolation

- File system (via `chroot`)
- PID
- IPC
  - [Instructions per cycle - Wikipedia](https://en.wikipedia.org/wiki/Instructions_per_cycle)?
- UTS (hostname)
- Mount
- Network

以下は cgroups で制限できるので実現可能

- Number of CPU cores
- RAM
- Nubmer of PIDs (to limit processes)

## Memo

こっからは初回にまとめるのはあれなので断片的にメモを取っていく

ns は `lsns` コマンドで取得できるらしいので Docker を立ち上げてみていく。

```
$ docker run --rm -it ubuntu /bin/bash
# lsns
        NS TYPE   NPROCS PID USER COMMAND
4026531835 cgroup      2   1 root /bin/bash
4026531837 user        2   1 root /bin/bash
4026533011 mnt         2   1 root /bin/bash
4026533012 uts         2   1 root /bin/bash
4026533013 ipc         2   1 root /bin/bash
4026533014 pid         2   1 root /bin/bash
4026533016 net         2   1 root /bin/bash
```

shell プロセスが所属しているネームスペースは以下のコマンドで確認できる

```
# ls -l /proc/self/ns
total 0
lrwxrwxrwx 1 root root 0 Jun 20 05:39 cgroup -> 'cgroup:[4026531835]'
lrwxrwxrwx 1 root root 0 Jun 20 05:39 ipc -> 'ipc:[4026533013]'
lrwxrwxrwx 1 root root 0 Jun 20 05:39 mnt -> 'mnt:[4026533011]'
lrwxrwxrwx 1 root root 0 Jun 20 05:39 net -> 'net:[4026533016]'
lrwxrwxrwx 1 root root 0 Jun 20 05:39 pid -> 'pid:[4026533014]'
lrwxrwxrwx 1 root root 0 Jun 20 05:39 pid_for_children -> 'pid:[4026533014]'
lrwxrwxrwx 1 root root 0 Jun 20 05:39 user -> 'user:[4026531837]'
lrwxrwxrwx 1 root root 0 Jun 20 05:39 uts -> 'uts:[4026533012]'
```

独自の名前空間に所属する shell プロセスを作る。

```
# unshare
unshare: unshare failed: Operation not permitted
```

ですよね。

おとなしく vagrant で環境を用意する。

vagrant インストール

```
brew cask install virtualbox
brew cask install vagrant
```

```
$ vagrant init bento/ubuntu-18.04
$ vagrant up
```

`Kernel driver not installed` のエラーが出たので対応
[Solving VirtualBox “kernel driver not installed (rc=-1908)” Error on macOS](https://medium.com/@Aenon/mac-virtualbox-kernel-driver-error-df39e7e10cd8)

```
$ vagrant ssh
#
```

気を取り直して、独自の名前空間に所属する shell プロセスを作る。

```
vagrant@vagrant:~$ sudo unshare --fork --pid --mount-proc /bin/bash
root@vagrant:~#
```

以下のように独立した PID 空間を持っていることがわかる。

```
# ps
  PID TTY          TIME CMD
    1 pts/0    00:00:00 bash
    9 pts/0    00:00:00 ps
```

演習としてコンテナがホスト上で動作する際の shell の PID を確認するとのことだったので確認してみる。

```
$ docker run --rm -it ubuntu /bin/bash
# ps
  PID TTY          TIME CMD
    1 pts/0    00:00:00 bash
    8 pts/0    00:00:00 ps
```

やはこコンテナでも unshare を実行した場合と同じ挙動になっている。

(Types of namespaces)

PID ネームスペースの仕組みを理解したところで、他のネームスペースの役割を理解していく。

[namespaces(7) - Linux manual page](https://man7.org/linux/man-pages/man7/namespaces.7.html)では 8 つのネームスペースが説明されている。

namespace を分割することで同じホストカーネル上で動作しているにも関わらず異なる VM で動作しているかのように独立させることができる。

(Creating new namespaces or joining existing ones)

デフォルトの挙動では `fork()` を呼び出して作成した子プロセスは親プロセスの PID ネームスペースを継承する。
しかも、`fork()` の引数は 0 こであり、作成前にプロパティの制御ができない。

そこで `clone()` を使いより細かな制御を行う。

(A side note on clone()) -> `clone()` に関する補足
内部的には `fork()`, `vfork()` は `clone()` を読んでいるだけ。
(ソースコードは省略)

(Using clone() to create processes with new namespaces)
このタイミングではネットワークが分離されていないが、Gocker ではこのあとに他のシステムコールによって仮想イーサネットインタフェースをセットアップする。

(Using unshare() to create and join new namespaces)
新しい子プロセスを作らずに既存のプロセスに新しい PID ネームスペースを割り当てたい場合は `unshare()`を用いる。

(Joining namespaces other processes belong to)
他のプロセスが所属してるネームスペースに参加するためには `setun()` システムコールを用いる。

**(How Gocker creates containers)**

この関数を追う [shuveb/containers-the-hard-way](https://github.com/shuveb/containers-the-hard-way/blob/40f4648feb9670709bd0c6f8a3c2a812d1103fd6/run.go#L192)

- コンテナIDを発行
- コンテナID用のルードティレク取りを階層的に作成
  - ここで階層化されていることによりいわゆるベースイメージを共有することができる
- 仮想イーサネットデバイス（ペア）を作成
  - パイプみたいなもんらしい
  - `veth0_<container-id>`, `veth1_<container-id>` のような形を取るらしい
  - 独自のネットワーク空間をとり、通信内容は秘密になっている
- 最後に `prepareAndExecuteContainer` によって実際にプロセスをコンテナ内で実行
  - この関数から抜けるタイミングにはコマンドは実行完了している
- クリーンアップ処理

`prepareAndExecuteContainer` についてもう少し詳しく見ていく

`prepareAndExecuteContainer` は3つのプロセスを作成する。
`gocker run` すると表示されるこれらである。

- `setup-netns`
- `setup-veth`
- `child-mode`

(Setting up networking that works inside containers)

新しいネットワークネームスペースを設定することはとても簡単で、`clone()` を呼ぶ際に `CLONE_NEWNET` をビットマスクにセットすれば良い。
トリッキーなのはコンテナ内から外部と通信できるようにすること。

gocker で作るネームスペースの中で一番最初に作るネームスペースがネットワークネームスペースである。

`setup-ns` と `setup-veth` を呼ぶと、

> nsの実態はファイルなの?

- 新しいネットワークネームスペースをセットアップ
- `setns()` システムコールはプロセスが属するすべてのネームスペースを列挙する `/proc/<pid>/ns` 内のファイルディスクリプタによって参照されるネームスペースに呼び出しプロセスのネットワークネームスペースを設定
- 
