# Step 2: 定期実行と設定ファイル

課題ページの原文: 「Use `setInterval` to check the URL every N seconds」「Log the results with timestamps to the console」「Add support for multiple URLs from a config file (JSON or YAML)」。
Go では `setInterval` を `time.Ticker` と goroutine に読み替える。

## 目標

JSON の設定ファイルに書いた複数の URL を、対象ごとの間隔でチェックし続ける常駐プロセスを作る。
結果は時刻付きで標準出力に流す。

```sh
go run ./cmd/uptime -config targets.json
# 2026-11-02T10:00:00Z UP   example   https://example.com   200  123ms
# 2026-11-02T10:00:00Z DOWN local     http://localhost:1    dial tcp ...: connection refused
# 2026-11-02T10:00:30Z UP   example   https://example.com   200  118ms
# ^C で終了
```

設定ファイルの例:

```json
{
  "targets": [
    { "name": "example", "url": "https://example.com", "interval": "30s" },
    { "name": "local", "url": "http://localhost:1", "interval": "10s" }
  ]
}
```

## 受け入れ条件

課題ページ由来:

- [ ] `-config` フラグで JSON ファイルを読み、複数の対象を取り込む。フラグ省略時は `targets.json`
- [ ] 対象ごとに `interval` の間隔で `check.URL` を呼び続ける。間隔は対象ごとに違ってよい
- [ ] 1 回のチェックごとに 1 行、時刻・UP/DOWN・名前・URL・ステータスコードまたはエラー・応答時間を出力する
- [ ] 設定ファイルが開けない、JSON として壊れている、`url` が空、`interval` が 0 以下または解釈できない、のいずれかなら起動せずに理由を表示して終了コード 1

Go 版として追加 (課題ページに無い。残すか外すかは着手前に決める):

- [ ] Ctrl-C (SIGINT) で全ての goroutine を止めてから終了する。実行中のチェックがあれば終わるのを待つ
- [ ] 設定の読み込みと検証は `internal/target`、定期実行は `internal/schedule` に置く。`cmd/uptime` はフラグ処理と起動・停止だけ
- [ ] `go vet ./...` と `go test ./...` が通る。テストは Step 1 と同じく別 PR でよい

## 学ぶ Go の概念

| 概念 | 読む場所 |
|---|---|
| goroutine | https://go.dev/tour/concurrency/1 。`go f()` で関数を並行に動かす |
| `time.Ticker` | https://pkg.go.dev/time#NewTicker 。`C` はチャネル。`Stop` を忘れない |
| チャネルと `select` | https://go.dev/tour/concurrency/2 、https://go.dev/tour/concurrency/5 。Ticker の `C` と `ctx.Done()` を同時に待つ |
| context によるキャンセル | https://pkg.go.dev/context#WithCancel 。Step 1 のタイムアウトと同じ型で「止めてよい」を伝える |
| シグナルで context を切る | https://pkg.go.dev/os/signal#NotifyContext 。Ctrl-C を `ctx.Done()` に変える |
| 全 goroutine の終了を待つ | https://pkg.go.dev/sync#WaitGroup 。`Add` / `Done` / `Wait` の対応 |
| JSON の読み込み | https://pkg.go.dev/encoding/json#Unmarshal 。構造体タグ `json:"name"`、未知のフィールドは無視される |
| 文字列から Duration | https://pkg.go.dev/time#ParseDuration 。`"30s"` を `time.Duration` にする。JSON の数値ではなく文字列で持つ理由を考える |
| ファイルを読む | https://pkg.go.dev/os#ReadFile |
| 構造化ログ | https://pkg.go.dev/log/slog 。`slog.Info("...", "key", value)` の形。時刻は自動で付く |
| エラーに文脈を足す | https://pkg.go.dev/fmt#Errorf の `%w` 。「設定ファイルが読めない: 元のエラー」の形にして呼び出し側に返す |

Step 1 では context を「時間切れ」に使った。Step 2 では同じ context を「Ctrl-C が来た」に使う。
どちらも `ctx.Done()` が閉じることで伝わる。この共通点を説明できるようにする。

## 自分で考える設計の問い

実装前に答えを docs/wip にメモしてから書き始める。

1. **スケジュールの単位**: 対象ごとに goroutine を 1 つ起こして各自 Ticker を持つか、Ticker を 1 つにして期限が来た対象をまとめてチェックするか。対象が 100 個のときに何が違うか
2. **設定の型**: `interval` を JSON でどう表すか。`"30s"` の文字列か、秒数の整数か。構造体のフィールドは `time.Duration` にするか `string` のままにして後で変換するか。`json.Unmarshal` は `time.Duration` に `"30s"` を直接入れられるか
3. **検証の場所**: `url` が空、`interval` が 0 以下、のような誤りはどこで見つけるか。読み込み関数の中か、呼び出し側か。見つけたらどう返すか (Step 1 の ADR 0005 では `Result` にエラーを含めたが、ここも同じでよいか)
4. **停止の伝え方**: Ctrl-C が来たとき、チェック中の goroutine にどうやって「止まれ」を伝えるか。Ticker を止めるだけで足りるか。`check.URL` に渡す context は `context.Background()` のままでよいか
5. **`check.URL` の変更**: Step 1 の `URL(url, timeout)` は内部で `context.Background()` から context を作っている。Step 2 で外から context を渡せるようにするか。変えるなら `cmd/check` も直す

## テストの観点 (別 PR で書くときの参考)

- `target`: 正しい JSON を読めるか。`url` 空、`interval` 不正、JSON 壊れ、ファイル無し、のそれぞれでエラーになるか
- `schedule`: 短い間隔 (10ms など) で動かし、一定時間に期待回数チェックが呼ばれるか。context を cancel したら止まるか
- `check.URL` を直接呼ばずに差し替えられる設計になっているか (問い 5 の答えがここで効く)

## 完了したら

1. `go vet ./...` の結果を確認。`-config` のファイル無し・JSON 壊れ・正常の 3 パターンと Ctrl-C で動作確認
2. `/go-review` でレビューを受ける。「これは何をしている?」に答えられるまで直す
3. 指摘を直してコミット。ブランチ `feature/step2-scheduled-checks` から PR を出す
4. 使った標準ライブラリを `docs/notes/stdlib.md` の「Step 2」に追記する (Step 1 分が未記入なら先に埋める)
5. 問い 1 と 5 は [ADR 0006](../adr/0006-goroutine-per-target.md) と [ADR 0007](../adr/0007-check-takes-context.md) に記録済み。他の問いで構造に関わる決定が出たら追記する
6. 次は Step 3 の課題カードを Claude が書く

## ブランチの切り方

PR 1 (Step 1) がまだマージされていないので、`feature/step1-single-check` から `feature/step2-scheduled-checks` を切る。
PR 1 がマージされたら `git rebase main` で載せ替える。
