# ADR 0004: フロントは TypeScript、リアルタイム更新は SSE

日付: 2026-10-27 / 状態: 採用

## 背景

課題ページは Chart.js と Socket.io を想定。学習対象は Go であり、フロントに時間をかけたくない。ユーザーは TypeScript に馴染みがある。

## 決定

- `web/` に Vite + TypeScript + Chart.js の最小構成を置く。Go は JSON API と静的ファイル配信だけを担当し、ビルド成果物 `web/dist` を `embed.FS` で配信する
- リアルタイム更新は SSE (Server-Sent Events)。Go 側は `http.Flusher` だけで書け、フロントは `EventSource` で受ける

## 比較 (リアルタイム更新)

| 候補 | 保有コード量 | AI の習熟度 | 機械検証 | 備考 |
|---|---|---|---|---|
| SSE | 少ない | 高い | 低い | 標準ライブラリのみ。サーバー→クライアントの片方向で要件を満たす |
| WebSocket | 中 | 高い | 低い | 外部ライブラリ (nhooyr / gorilla) が要る。双方向は不要 |
| ポーリング | 少ない | 高い | 低い | 最も単純だが無駄なリクエストが多い |

## 結果

`html/template` でのサーバーサイド描画は使わない。公開ステータスページ (Step 6) も同じ TS フロントのページとして作る。
フロントのコードはユーザーが望めば Claude が書いてよい (CLAUDE.md)。
