// 監視対象 1 件に HTTP GET を打ち、結果を Result にまとめる。
// 保存や表示は行わない。
package check

import (
	"net/http"
	"time"

	"example.com/checker/target"
)

type Result struct {
	Name       string
	URL        string
	StatusCode int
	LatencyMs  int64
	Error      string
	CheckedAt  time.Time
}

var client = &http.Client{Timeout: 5 * time.Second}

// 接続失敗やタイムアウトは error ではなく Result.Error に文字列で入れる。
// 「落ちていた」という事実も保存対象なので、呼び出し側を止めないため。
func Run(t target.Target) Result {
	start := time.Now()
	r := Result{Name: t.Name, URL: t.URL, CheckedAt: start}

	resp, err := client.Get(t.URL)
	r.LatencyMs = time.Since(start).Milliseconds()
	if err != nil {
		r.Error = err.Error()
		return r
	}
	defer resp.Body.Close()

	r.StatusCode = resp.StatusCode
	return r
}
