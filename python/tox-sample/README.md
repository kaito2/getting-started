## 概要

Pythonライブラリを複数バージョンでテストするツールらしい。（ライブラリ向け?）

厳密な目的はライブラリなどの互換性を気にするソフトウェアのテストを複数バージョンでできるようにする。ということだと思われるが、
実態的には、テストの実行やフォーマットのバリデーションを実行するための pipenv の scripts みたいな使い方をされているような気がする。

この記事とかも[Poetryとtoxを組み合わせてPythonのテストコードを動かす - Qiita](https://qiita.com/ninomiyt/items/01a3c3d93d99b3f551fd)

## メモ

[GitHub - tox-dev/tox: Command line driven CI frontend and development task automation tool](https://github.com/tox-dev/tox)を動かしてみた。

参考: 
* [Welcome to the tox automation project — tox 3.15.3.dev1 documentation](https://tox.readthedocs.io/en/latest/)
* [packaging — tox 3.15.3.dev1 documentation](https://tox.readthedocs.io/en/latest/example/package.html#poetry)
    * poetry 対応
* [Installation and Getting Started — pytest documentation](https://docs.pytest.org/en/stable/getting-started.html)
    * 使った

## 実行

```
$ poetry run tox
___________________________________________________________________________________________________________________________________ summary ____________________________________________________________________________________________________________________________________
  congratulations :)
```

## ダミーのテスト追加

