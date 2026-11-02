# ADR 0007: check.URL は第 1 引数に context を取る

日付: 2026-11-02
状態: 採用

## 文脈

Step 1 の `check.URL(url, timeout)` は内部で `context.Background()` からタイムアウト付き context を作っている。
Step 2 で Ctrl-C による停止を入れるにあたり、進行中のチェックを外から中断できるようにするかが論点。

## 決定

- シグネチャを `URL(ctx context.Context, url string, timeout time.Duration) Result` に変える
- 関数内では渡された `ctx` を親にして `context.WithTimeout` を作る。`timeout` の意味は変わらない
- `cmd/check` は `context.Background()` を渡す。`cmd/uptime` は Ctrl-C で切れる context を渡す

## 理由

- Go の標準ライブラリと同じ形。「context を取る関数は第 1 引数 `ctx`」は pkg.go.dev/context の冒頭に書かれた慣習
- Ctrl-C を押した瞬間に進行中の HTTP 通信も中断できる。`Background` のままだと最大 `timeout` 秒待つことになる
- Step 7 の agent でも同じ関数を使うので、停止の口を今のうちに開けておく

## 不利益

- `cmd/check` の呼び出しを直す必要がある。変更は 1 行
- Ctrl-C の瞬間に進行中だったチェックは `context.Canceled` のエラーで返り、本物の DOWN と区別なく記録される。Step 2 の範囲 (ログ出力のみ) では許容し、Step 3 で DB に保存する前に `errors.Is(r.Err, context.Canceled)` で除外するか判断する (2026-11-02 レビューで持ち越し決定)

## 検討した代替案

- 変えない。Ctrl-C で「次のチェックを始めない」だけにする: 要件の範囲では足りるが、進行中の 1 回を待つ分だけ終了が遅れる
- `timeout` を引数から外し、呼び出し側が `WithTimeout` を作って渡す: より純粋な形だが、呼び出し側ごとに同じ数行を書くことになる。Step 8 以降でチェックの設定が増えたら `Checker` 構造体に `Timeout` を持たせる形 (`http.Client` と同じ) に移す
