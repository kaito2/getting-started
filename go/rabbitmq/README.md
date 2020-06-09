## 概要

カテゴリとしてはメッセージブローカーらしい

See: [RabbitMQ tutorial - “Hello World!”  — RabbitMQ](https://www.rabbitmq.com/tutorials/tutorial-one-go.html)

## 実行

RabbitMQ のサーバーをDockerで建てる必要がある
ref: [Downloading and Installing RabbitMQ — RabbitMQ](https://www.rabbitmq.com/download.html)

```
docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management
```

```
$ go run send/main.go
$ go run receive/main.
2020/06/09 21:07:01  [*] Waiting for messages. To exit press CTRL + C
2020/06/09 21:07:01 Received a message: Hello World!
```

## 新しい概念

RabbitMQ は複数プロトコルをサポートしているのが売りっぽい
参考: [Which protocols does RabbitMQ support? — RabbitMQ](https://www.rabbitmq.com/protocols.html)

チュートリアル出でてきた AMQP（Advanced Message Queuing Protocol）は金融がルーツのメッセージプロトコルらしい。
他にもMQTTやHTTPなどもサポートされていた。
