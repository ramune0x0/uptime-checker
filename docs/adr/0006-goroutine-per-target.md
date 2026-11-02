# ADR 0006: 監視対象ごとに goroutine と Ticker を 1 つずつ持つ

日付: 2026-11-02
状態: 採用

## 文脈

Step 2 で複数の対象をそれぞれの間隔でチェックし続ける。課題ページの原文は「Use `setInterval` to check the URL every N seconds」。
対象ごとに独立した時計を持つか、時計を 1 つにして予定表で管理するかが論点。

## 決定

- 対象 1 つにつき goroutine を 1 つ起こし、その中で `time.Ticker` を持ってループする
- 各 goroutine は共通の context を受け取り、`ctx.Done()` が閉じたら Ticker を止めて抜ける
- 全 goroutine の終了は `sync.WaitGroup` で待つ

## 理由

- 課題ページの「対象ごとに `setInterval`」と 1 対 1 に対応し、要件の読み替えとして素直
- goroutine は OS スレッドより軽く、対象ごとに 1 つ起こすのは Go の常套手段。予定表を自分で管理するコードを持たずに済む
- goroutine / Ticker / `select` / context / WaitGroup という Go の並行処理の基本一式をこのステップで一度に練習できる

## 不利益

- 対象が数千になると goroutine も数千になる。この課題の規模では問題にならない

## 検討した代替案

- Ticker を 1 つにして毎秒「期限が来た対象」を選んで実行する: 対象数に依らず goroutine が増えないが、次回予定時刻の管理を自前で書く必要があり、要件に対して作り込みが大きい
