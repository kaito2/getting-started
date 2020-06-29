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
$ vagrant up$
```

`Kernel driver not installed` のエラーが出たので対応
[Solving VirtualBox “kernel driver not installed (rc=-1908)” Error on macOS](https://medium.com/@Aenon/mac-virtualbox-kernel-driver-error-df39e7e10cd8)
