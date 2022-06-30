## 7. (発展) multi-stage build

### Go で Web アプリを作ってコンテナで実行する
Go で Web サーバーを作成するパートは Docker とは関係がありませんが、がんばてみて下さい。

#### goenv, go インストール

goenv をインストールします。
```
$ git clone https://github.com/syndbg/goenv.git ~/.goenv
```

goenv にパスを通します。  
`.bashrc` などに以下の内容を追記します (お使いの環境に合わせて読み替えて下さい)。
```bash
export GOENV_ROOT="$HOME/.goenv"
export PATH="$GOENV_ROOT/bin:$PATH"
eval "$(goenv init -)"
```

fish の場合は以下のような設定内容になります。
```bash
set -gx GOENV_ROOT $HOME/.goenv
set -ga PATH $GOENV_ROOT/bin
source (goenv init - | psub)
```

設定内容を反映します。
```
$ exec $SHELL
```

#### go インストール
バージョン 1.18.3 をインストールします。
```
$ goenv install 1.18.3
```

デフォルトの Go バーションを先ほどインストールしたものに設定。
```
$ goenv global 1.18.3
```

ここまでで go コマンドが使えるようになっているはずです。  
`go version` でバーションを確認してみましょう。

### Web アプリ作成

プロジェクトのディレクトリを作成し、Go のプロジェクトとして初期化します。
```
$ mkdir goapp && cd goapp
$ go mod init goapp
```

プロジェクト直下に以下のファイルを作成します。  
`server.go`
```go
package main

import "net/http"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "hello"}`))
		default:
			http.Error(w, ``, http.StatusMethodNotAllowed)
		}
	})
	http.ListenAndServe(":8888", nil)
}
```

実行します。
```
$ go run server.go
```
ブラウザ等で http://localhost:8888 にアクセスして動作を確認しましょう。

### multi-stage building
先ほどはローカルで実行しましたが、コンテナとして実行できるよう image を作成します。

```Dockerfile
FROM golang:1.18.3-bullseye AS build
WORKDIR /app
COPY . /app
RUN go mod download
RUN go build -o /server

FROM gcr.io/distroless/base-debian11
WORKDIR /
COPY --from=build /server /server
EXPOSE 8888
USER nonroot:nonroot
ENTRYPOINT ["/server"]
```
そもそも説明していない命令コマンドがいくつかできてきました。それらについては後に回すとして、全体の説明をするとこの Dockerfile は 2 つのパートに分かれています。`FROM` が二箇所出現しているところがそれで、image が 2 つ作成されます。一つ目の image で go のプログラムをコンパイルし、できあがった実行ファイルを二つ目の image にコピー (`COPY --from=build` の箇所) します。コンテナとして使用するのは 2 つ目の image です。

このように、多段階に分けて image をビルドすること、あるいはその機能のことを multi-stage build といいます。  
このようなことをして何が嬉しいのかというと、image サイズを小さくできるという点があります。コンパイルを行うためにはそのための開発ツールが必要ですが、コンテナ実行時にそれらは必要ありません。そのためコンパイルする image と実行する image を分けると都合がいいのです。また実行コンテナに不要なものがあると、それがセキュリティホールになる可能性があるため、できる限り不要なものは削ぐべきです。

初出の命令コマンドなどを項目ごとに説明します。
- `FROM golang:1.18.3-bullseye AS build`  
  golang version 1.18.3 のビルド&実行環境をもつ image を使います。`AS` によってここで作成される image に名前を付けています。これは 2 つ目の image ビルド時の `COPY` 命令コマンドで参照するためです。
- `FROM gcr.io/distroless/base-debian11`  
  Google がホストしている distroless と呼ばれる一連の image の一つです。この image はとても軽量であるため実行用コンテナとして優れています。
- `WORKDIR`  
  `cd` コマンドに対応するものです。Dockefile 中でワーキングディレクトリを変更します。
- `EXPOSE`
  コンテナのポート開放を指定する命令ですが、実際には何もしません。実際のポート開放は実行時のコマンドでしていします。
- `USER nonroot:nonroot`
  `USER` は `CMD` や `ENTRYPOINT` の実行ユーザーを指定します。distroless image では nonroot(user):nonroot(group) を指定することができ、これによってプロセスを非ルートユーザーで実行できます。
- `ENTRYPOINT` コンテナ立ち上げ時の起動コマンドを指定します。`CMD` と似ていますが細かな違いがあります。ここでは詳細には立ち入らないので、詳しく知りたい場合は調べてみて下さい。

ビルド & 実行
```
$ docker build . -t goapp
$ docker run -it --rm -p 8888:8888 goapp
```

これにて終わりです。

### 参考文献
- 渋川よしき・辻大志郎・真野隼記著 (2022)『実用 GO言語』 オライリー・ジャパン
- https://zenn.dev/suiudou/articles/5e1dfd1008bf29
- https://qiita.com/tanan/items/e79a5dc1b54ca830ac21#user