## 3. docker の使い方

### コンテナの作成・実行・削除
Linux ディストリビューションの 1 つである debian の image を使って、実際にコンテナに触れてみましょう。

#### コンテナの作成
以下のコマンドを実行します。
```
$ docker create -it --init --name testvm debian:bullseye
```

これは、
- debian:bullseye イメージを元に (`debian:bullseye`)
- testvm という名前の (`--name testvm`)
- コンテナを作成しなさい

という意味のコマンドです。  
(`-it`, `--init` オプションの意味は今は気にしないで下さい。)

`debian:bullseye` は debian11 の docker image で (bullseye というのは version 11 のコードネームです)、この image は Docker 公式の [DockerHub](https://hub.docker.com/_/debian?tab=tags&page=1&name=bullseye) で公開されています。
docker はローカルに image が存在しなければ、基本的にこの DockerHub から取得してきます。

#### コンテナの作成確認
さて、コンテナが作成されていることを確認しましょう。
```
$ docker ps
CONTAINER ID   IMAGE     COMMAND   CREATED   STATUS    PORTS     NAMES

$ docker ps -a
CONTAINER ID   IMAGE             COMMAND   CREATED          STATUS    PORTS     NAMES
d12b58f32f35   debian:bullseye   "bash"    43 seconds ago   Created             testvm
```
`docker ps` コマンドは実行中のコンテナの一覧を表示するコマンドです。  
`-a` オプションをつけることで停止中のコンテナも含めて表示します。  
今回、コンテナ testvm はまだ起動していないので `-a` オプションを付けないと表示されません。

#### コンテナの実行
コンテナを実行し、ステータスを確認します。
```
$ docker start testvm
$ docker ps
CONTAINER ID   IMAGE             COMMAND   CREATED         STATUS         PORTS     NAMES
f80a2cfcb3ae   debian:bullseye   "bash"    6 seconds ago   Up 2 seconds             testvm
```
コンテナには「起動コマンド」というものが設定されていて、`docker start` コマンドでそれが実行されます。testvm コンテナの起動コマンドは `bash` であることがわかります。  

#### コンテナにアタッチ
コンテナの中に入ってみましょう。ヴァーチャルマシーン (VM) にログインするようなイメージです。
```
$ docker attach testvm
```

コンテナ内で、debian OS を自由に操作できます。
Linux コマンドを実行して遊んでみましょう。  

コンテナから出るには `C-p C-q` でデタッチします。  
(`exit` コマンドでも抜けられますが、この場合 `bash` プロセスが終了しコンテナが停止してしまいます。)

#### コマンドの実行
アタッチしたコンテナ内で Linux コマンドを実行できましたが、コンテナ外からも実行することができます。  
`pwd` コマンドは以下のように実行します。
```
$ docker exec testvm pwd
```

#### コンテナの停止・破棄
コンテナを停止します。
```
$ docker stop testvm
```

停止しているか確認しましょう。
```
$ docker ps
$ docker ps -a
```

コンテナを削除します。
```
docker rm testvm
```
これによってコンテナが完全に削除され、`docker ps -a` をしても表示されなくなります。

#### 補足 (docker run コマンド)

コンテナの生成・実行は 1 つのコマンドで済ますことができます。
```
$ docker run -d -it --init --name testvm debian:bullseye
```
これは以下の 2 コマンドを実行したのと同じです
- `docker create -it --init --name testvm debian:bullseye`
- `docker start testvm`

チュートリアルでは説明のために `docker create`, `docker start` コマンドを用いましたが、大抵の場合では `docker run` コマンドを使う機会のほうが多いでしょう。

#### ここまでのまとめ
ここまでに出てきたコマンドのうち覚えておきたいものをまとめました。
- コンテナ作成・実行 (`docker run`)
  ```
  $ docker run [-d] -it [--rm] [--init] [--name <container_name>] <image_name[:tag]> [command]
  ```
- コンテナ停止 (`docker stop`)
  ```
  $ docker stop <container_name>
  ```
- コンテナ削除 (`docker rm`)
  ```
  $ docker rm <container_name>
  ```
- コマンド実行 (`docker exec`)
  ```
  $ docker exec <container_name> <command>
  ```
- コンテナの確認 (`docker ps`)
  ```
  $ docker ps [-a]
  ```

### コンテナで色々なコマンドを実行する
前の節ではコンテナを立ち上げただけでしたが、ここでは立ち上げたコンテナにいくつかコマンドを実行させてみましょう。  
コマンド実行を通じて、実行コマンドのオプションの意味やコンテナの揮発性を学びます。

再び前節でつくったコンテナを用意します。
```
$ docker run -d -it --rm --init --name testvm debian:bullseye
```
`--rm` オプションをつけると、コンテナの停止時にコンテナの破棄も同時に行われます。

#### 普通にコマンド実行
`ls -la` コマンドを実行してみましょう。
```
$ docker exec testvm ls -l
```

#### `-i` オプションについて
debian は apt というパッケージマネージャを採用しています。  
`apt update` でパッケージ一覧を更新します。
```
$ docker exec testvm apt update
```

`apt upgrade` はインストール済みのパッケージのバージョンを最新にするコマンドです。
```
$ docker exec testvm apt upgrade
```
このコマンドは途中で失敗します。
本当にアップグレードしますか？という確認メッセージが表示されたあとコマンドはユーザーの入力を待つのですが、標準入力が存在しないためにコマンドが終了してしまいます。

標準入力を有効にするには `-i` オプションを有効にします。
```
$ docker exec -i testvm apt upgrade
```
これで入力を受け付けるようになったはずです。

#### `-t` オプションについて
`bash` コマンドを実行してみましょう。
```
$ docker exec testvm bash
```
何も起きずにコマンドが終了したかと思います。  
`-t` オプションでプロセスに擬似端末を紐付け、フォアグラウンドで実行できます。やってみましょう。
```
$ docker exec -t testvm bash
```
bash が起動したのはいいものの、今度はコマンドを入力しても何も起こりません。  
`-i` オプションがないので標準入力を受け付けてくれないからです。  
一度 `C-c` で抜けて、`-it` オプションを付けて再実行してみましょう。
```
$ docker exec -it testvm bash
```

`C-p C-q` または `exit` コマンドで終了します。  
(ここで実行している bash は起動コマンドとは別のプロセスなので `exit` で抜けてもコンテナは停止しません。)

**NOTE**  
以降、特に必要ない場合であっても `-it` オプションを付けて実行することがよくあります。
必要ない場合に有効にしていても基本的に問題が発生することはありません。  
特に初心者のうちは何も考えずに `-it` オプションを付けてしまってもよいでしょう。

#### コンテナの揮発性
最後にちょっとしたコマンドをインストールして遊びます
```
$ docker exec -it testvm apt install -y bpytop
```

インストールが完了したらコマンドを実行してみます。
```
$ docker exec -it testvm bpytop
```
`Esc` で終了できます。

`bpytop` をインストールしてみましたが、コンテナを破棄するとこの変更も消えてしまいます (`--rm` オプションをつけてコンテナ生成をしていたことに注意)。
```
$ docker stop testvm
$ docker run -d -it --rm --init --name testvm debian:bullseye
$ docker exec -it testvm bpytop
OCI runtime exec failed: exec failed: unable to start container process: exec: "bpytop": executable file not found in $PATH: unknown
```
コンテナは image を元に作成されるため、コンテナに対する変更はコンテナの再作成時に引き継がれません。

基本的にコンテナに変更を加えていくような使い方はせず、`bpytop` のインストールされた debian コンテナが欲しいという様な場合は image をカスタマイズして使います。

次章では image をカスタマイズする方法を学んでいきます。