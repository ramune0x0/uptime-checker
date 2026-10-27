# Step 1: 単発チェック CLI

## 目標

URL をコマンドライン引数で受け取り、1 回だけ HTTP GET して結果を表示するコマンドを作る。

```sh
go run ./cmd/check https://example.com
# UP   200  123ms
go run ./cmd/check https://localhost:1
# DOWN  dial tcp [::1]:1: connect: connection refused
echo $?
# 1
```

## 受け入れ条件

- [ ] URL を 1 つ引数に取る。引数がなければ使い方を表示して終了コード 2
- [ ] 2xx〜3xx なら `UP`、ステータスコードと応答時間 (ms) を表示して終了コード 0
- [ ] 4xx 以上なら `DOWN`、ステータスコードを表示して終了コード 1
- [ ] 接続失敗・タイムアウトなら `DOWN` と理由を表示して終了コード 1
- [ ] タイムアウトは既定 5 秒。`-timeout 2s` のようなフラグで変更できる
- [ ] チェックのロジックは `internal/check` パッケージに置き、`cmd/check/main.go` は引数処理と表示だけ
- [ ] `internal/check` に `net/http/httptest` を使ったテストが 3 ケース以上 (200 / 500 / タイムアウト)
- [ ] `go vet ./...` と `go test ./...` が通る

## 学ぶ Go の概念

| 概念 | 読む場所 |
|---|---|
| コマンドライン引数とフラグ | `flag` パッケージ: https://pkg.go.dev/flag (`flag.Duration`, `flag.Args`) |
| HTTP クライアントとタイムアウト | https://pkg.go.dev/net/http#Client の `Timeout` フィールド |
| context によるタイムアウト | https://pkg.go.dev/context#WithTimeout と `http.NewRequestWithContext` |
| エラーの種類を見分ける | https://pkg.go.dev/errors#As, https://pkg.go.dev/net#Error の `Timeout()` |
| 経過時間の計測 | https://pkg.go.dev/time#Since, `Duration.Milliseconds` |
| 終了コード | https://pkg.go.dev/os#Exit。`defer` が実行されない点に注意 |
| HTTP のテスト | https://pkg.go.dev/net/http/httptest#NewServer |
| 内部パッケージ | https://go.dev/doc/go1.4#internalpackages |

`Client.Timeout` と `context.WithTimeout` はどちらでもタイムアウトを実現できる。両方読んで、違いを説明できるようにする。

## 自分で考える設計の問い

実装前に答えを docs/wip にメモしてから書き始めると、レビューで設計の意図を確認しやすい。

1. 「チェック結果」を表す型に何を持たせるか。ステータスコード、応答時間、エラー、時刻、UP/DOWN の判定。判定は保存する値か、他の値から導く関数か
2. 接続失敗やタイムアウトを `error` として返すか、結果の一部として返すか。「落ちていた」という事実は Step 3 で保存対象になる
3. `Client` をどこで作るか。関数内で毎回作る / パッケージ変数 / 構造体に持たせて外から渡す。テストでタイムアウトを短くしたいときに困らないのはどれか
4. 応答時間はどこからどこまでを測るか。リクエスト送信前〜ヘッダ受信まで / 本文読み切りまで
5. レスポンス本文を読まない場合でも `Body.Close()` は必要か。なぜか

## テストの観点

- `httptest.NewServer` で 200 / 500 を返すサーバーを立て、判定が正しいか
- ハンドラ内で `time.Sleep` してタイムアウトを起こし、DOWN かつ理由がタイムアウトになるか
- 存在しないポートへの接続拒否
- テスト用にタイムアウトを短くできる設計になっているか (問い 3 の答えがここで効く)

## 完了したら

1. `go vet ./...` と `go test ./...` の結果を確認
2. Claude にレビューを依頼
3. 指摘を直してコミット。次は Step 2 の課題カードを Claude が書く
