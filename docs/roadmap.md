# ロードマップ

課題ページ https://nst-sdc.github.io/project-ideas/web/url-monitoring.html の 10 ステップを Go に読み替えて進める。
順番は 1 → 2 → 3 → 4 → 5 → 6 → 8 → 9 → 10 → 7。マルチリージョンは構成を変えるので最後に回す。

各ステップは `docs/steps/NN-*.md` の課題カードに受け入れ条件を書き、それを満たしたらコミットして次へ進む。

| Step | 課題 | 成果物 | 学ぶ Go の概念 |
|---|---|---|---|
| 1 | 単発チェック CLI | `cmd/check`, `internal/check` | `flag`, `net/http.Client` の Timeout, `context.WithTimeout`, `errors.As` と `net.Error`, `time.Since`, 終了コード, `net/http/httptest` |
| 2 | 定期実行と設定ファイル | `cmd/uptime`, `internal/target`, `internal/schedule` | `time.Ticker` (setInterval 相当), goroutine, `sync.WaitGroup`, `context` と `signal.NotifyContext` による停止, `encoding/json`, `log/slog` |
| 3 | 結果の保存と 24h 稼働率 | `internal/store` | `database/sql`, `modernc.org/sqlite`, `go:embed` でスキーマ, UTC / RFC3339, SQL の集計 |
| 4 | ダッシュボード | `internal/httpserver`, `web/` | Go 1.22 以降の `ServeMux` メソッドパターン, JSON ハンドラ, `http.Flusher` で SSE, `embed.FS` で静的配信, goroutine 間の結果配信 (channel) |
| 5 | アラート | `internal/alert` | 状態遷移の検出, インタフェースで通知先を差し替え, `net/smtp` or webhook POST, クールダウン, 環境変数 |
| 6 | 公開ステータスページ | `internal/store` 拡張, `web/` 追加 | 日次ロールアップ, インシデントの開始/終了, `encoding/xml` で RSS |
| 8 | SSL 証明書監視 | `internal/check` 拡張 | `crypto/tls`, `resp.TLS.PeerCertificates`, `NotAfter` |
| 9 | 高度なチェック | `internal/check`, `internal/target` 拡張 | `io.LimitReader`, `strings.Contains`, カスタムメソッド/ヘッダ/ボディ, `net/http/httptrace` で DNS 時間 |
| 10 | コンテナ化 | `Dockerfile`, `compose.yaml` | マルチステージビルド, `CGO_ENABLED=0`, docker compose |
| 7 | マルチリージョン | `cmd/agent`, 中央 API | 役割の分離, HTTP クライアントとしての Go, リージョン別集計, 複数失敗時のみ通知 |

## 各ステップの受け入れ条件 (要約)

- **1**: URL を引数に取り、ステータスコードと応答時間を表示。失敗・タイムアウトは DOWN と理由を表示して終了コード 1
- **2**: JSON 設定の複数 URL を対象ごとの間隔でチェックし続け、時刻付きでログ出力。Ctrl-C で goroutine を止めて終了
- **3**: 全結果を SQLite に保存し、対象ごとの直近 24 時間の稼働率を出せる
- **4**: ブラウザで対象一覧と UP/DOWN を色分け表示。応答時間の推移をグラフ化。新しい結果が SSE で反映される
- **5**: UP→DOWN で通知が届く。連続した DOWN で通知が繰り返されない。DOWN→UP で復旧通知
- **6**: 認証なしページで各サービスの現在状態・90 日稼働率・インシデント履歴を表示。RSS を配信
- **8**: HTTPS 対象の証明書期限を記録し、14 日以内なら警告。ダッシュボードに期限を表示
- **9**: レスポンス本文のキーワード一致、メソッド/ヘッダ/ボディ指定、DNS 解決時間の個別計測
- **10**: `docker compose up` で全体が動く
- **7**: 2 つ以上の agent から結果を集め、リージョン別の応答時間を表示。1 リージョンだけの失敗では通知しない

## 未決の設計判断 (到達時に決める)

- Step 2: 対象ごとに goroutine + Ticker を持つか、1 つの Ticker で期限到来分をまとめて実行するか
- Step 4: 結果を SSE に流す経路 (channel の fan-out)
- Step 6: 90 日稼働率を生データから集計するか日次ロールアップ表を持つか
- Step 7: agent ↔ 中央の連携。Pull (HTTP) を第一候補とし、Pub/Sub は必要になったら再検討
- Step 10 以降: デプロイ先。未定。ローカル docker compose まで
