# ADR 0003: Step 6 までは単一プロセス、Step 7 で agent を分離する

日付: 2026-10-27 / 状態: 採用

## 背景

チェッカー (定期実行) と Web サーバー (ダッシュボード・API) は役割が違う。
最初から別プロセスにすると、両者の連携 (DB ポーリング、Pub/Sub など) を先に作る必要があり、学習の順番が崩れる。
SQLite は書き込みプロセスを 1 つに限定したい ([ADR 0002](0002-sqlite-modernc.md))。

## 決定

- Step 2〜6: `cmd/uptime` 1 プロセスで、スケジューラと HTTP サーバーを goroutine として同時に動かす。両者は同じ `store` を共有する
- Step 7: `cmd/agent` を新設する。agent は中央 (`cmd/uptime`) から対象一覧を HTTP GET で取得し、結果を HTTP POST で送る Pull 方式。DB を書くのは中央だけ

## 比較 (Step 7 の連携方式)

| 候補 | 保有コード量 | AI の習熟度 | 機械検証 | 備考 |
|---|---|---|---|---|
| Pull (HTTP GET/POST) | 少ない | 高い | JSON スキーマ | 中央に API を足すだけ。agent は HTTP クライアント |
| DB ポーリング | 少ない | 高い | SQL | agent が DB に直接触るので SQLite の制約に反する |
| Pub/Sub (Redis, NATS 等) | 中 | 高い | 低い | 外部サービス追加。ローカル検証が重くなる |

## 結果

Step 4 では、スケジューラが出した結果を HTTP サーバー (SSE) に渡す経路が goroutine 間通信になる。channel の fan-out を学ぶ場になる。
Step 7 到達時に Pull で足りない要件 (即時反映など) があれば Pub/Sub を再検討する。
