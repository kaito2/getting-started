# oauth2_proxy を動かす

[Home - OAuth2 Proxy](https://oauth2-proxy.github.io/oauth2-proxy/)

参考:
* [oauth2_proxy と Auth0 を用いた Nginx のお手軽 OAuth 化 · Yutaka 🍊 Kato](https://mikan.github.io/2018/05/23/enable-oauth-to-your-nginx-by-oauth2-proxy-and-auth0/)
* [セキュアなリバースプロキシの作り方 - BBSakura Networks Blog](https://bbsakura.github.io/posts/secured-reversproxy/)
* [手っ取り早くウェブアプリケーションにOAuth2認証を導入する - その手の平は尻もつかめるさ](https://moznion.hatenadiary.com/entry/2017/12/14/230945)
	* Docker Compose 化の参考に

## 実行

`.env` に設定する値は [oauth2_proxy と Auth0 を用いた Nginx のお手軽 OAuth 化 · Yutaka 🍊 Kato](https://mikan.github.io/2018/05/23/enable-oauth-to-your-nginx-by-oauth2-proxy-and-auth0/) を参照されたし。

```
cp .env.sample .env
# edit .env
docker-compose up
# access to http://localhost/
```

## Log

`https://your-site.com/oauth2/callback` を適当に `https://kaito2.example.com/oauth2/callback` とかにしておいた。

## TODO

― [ ] redirect 後にぶっ壊れるので修正