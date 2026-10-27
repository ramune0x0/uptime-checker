# ADR 0002: 永続化は SQLite、ドライバは modernc.org/sqlite

日付: 2026-10-27 / 状態: 採用

## 背景

Step 3 でチェック結果を保存する。課題ページは SQLite を想定。開発は 1 台の macOS で行い、デプロイ先は未定。

## 決定

`database/sql` + `modernc.org/sqlite` を使う。

## 比較

| 候補 | 保有コード量 | AI の習熟度 | 機械検証 | 備考 |
|---|---|---|---|---|
| modernc.org/sqlite | 少ない | 高い | SQL 文字列 | 純 Go。cgo 不要で `CGO_ENABLED=0` の静的ビルドが可能 |
| mattn/go-sqlite3 | 少ない | 高い | 同上 | cgo が必要。Docker ビルドで gcc が要る |
| PostgreSQL | 少ない | 高い | 同上 | 外部プロセスが要り、Step 3 時点では過剰 |

ORM (GORM, sqlc 等) は使わない。`database/sql` の `Query` / `Scan` を自分で書くことが学習目的に合う。

## 結果

SQLite は同時書き込みに弱いので、書き込みプロセスは常に 1 つにする ([ADR 0003](0003-single-process-until-multiregion.md))。
マルチリージョンやクラウドデプロイで書き込み元が増える場合は PostgreSQL への移行を再検討する。
