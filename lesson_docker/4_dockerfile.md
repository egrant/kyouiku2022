## 4. Dockerfile の書き方

### debian:bullseye image のカスタマイズ
前章の bullseye に bpytop やの git やのを色々インストールして見るだけのセクション
Apache もインストールして CMD [httpd], EXPOSE 80 をするか〜〜〜

まず前章を踏まえて、`bpytop` がインストールされた debian image を作ってみましょう。

適当なディレクトリ (ここでは mydebian/ とします) を用意し、その中に `Dockerfile` というファイル名で以下の内容のファイルを作成して下さい。

```Dockerfile
FROM debian:bullseye       # (1)

RUN apt update             # (2)
RUN apt install -y bpytop  # (2)
```

Dockerfile の内容は以下の様になっています。

1. まず `FROM` でベースとなる image を指定します
2. `RUN` でコマンドを実行できます
  - apt update コマンドを使ってパッケージ一覧を更新します
  - apt install コマンドで `bpytop` をインストールします

Dockerfile は上のように `FROM` 命令から始まるのが基本です。
`FROM` 命令で指定したベースの image を `RUN` を始めとする Dockerfile の命令コマンドを使ってカスタマイズしていきます。

Dokcerfile をもとに image を作成しましょう。
次のコマンドを実行して下さい。(// の部分はコマンドの説明コメントです)
```
$ docker build ./mydebian -t mydebian
// docker build <Dockerfile のあるディレクトリのパス> -t <image_tag>
```

`mydebian:latest` という名前の image が作成されます。
以下のコマンドで image の一覧を取得でるので確認してみましょう。
```
$ docker image ls
```

kdjfallfkasjf

### 実践的なやつ
nginx image を使って静的ページを表示できるようにしていく