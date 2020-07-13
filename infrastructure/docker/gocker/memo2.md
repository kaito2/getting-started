# bocker を参考に周辺知識を固める

[p8952/bocker](https://github.com/p8952/bocker/blob/master/bocker)

## ネットワーク編

[p8952/bocker](https://github.com/p8952/bocker/blob/master/bocker#L66)
=> `ip` コマンドってなに? ^^;

> 「ip」コマンドは、ネットワークデバイスやルーティング、ポリシーなどの表示と変更を行うコマンドです。従来は、ifconfig コマンドや netstat コマンド、route コマンドなど、net-tools パッケージに収録されているコマンド群を使用していました。現在は、いずれも ip コマンドへの移行が進んでいます。

[【 ip 】コマンド（基礎編）――ネットワークデバイスの IP アドレスを表示する](https://www.atmarkit.co.jp/ait/articles/1709/22/news019.html)

なので bocker で使われているオブジェクトのみ抜粋すると、

- `link`: ネットワークデバイス
- `netns`: ネットワークネームスペース

用語はわかったけど構成が図でみたいよ…

ということで`docker ネットワーク`あたりでググっていく

[Docker のネットワークの仕組み - sagantaf](http://sagantaf.hatenablog.com/entry/2019/12/18/234553)

veth についてわからないので調べる。

[veth(4) - Linux manual page](https://man7.org/linux/man-pages/man4/veth.4.html)

```
ip link add dev veth0_"$uuid" type veth peer name veth1_"$uuid"
```
