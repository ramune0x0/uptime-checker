# ADR 0008: 設定はファイル形の型で受け、検証して別の型に詰め替える

日付: 2026-11-02
状態: 採用

## 文脈

Step 2 で監視対象を JSON 設定ファイルから読む。`interval` は `"30s"` のような文字列で書きたいが、`json.Unmarshal` は文字列を `time.Duration` に入れられない (`time.Duration` の実体は `int64` で、JSON の数値しか受け付けない)。
`url` が空、`interval` が 0 以下、といった誤りをどこで見つけ、どう返すかも論点。

## 決定

- `internal/target` に 2 種類の型を置く
  - ファイルの形そのままの非公開の型。`interval` は `string`
  - 検証済みの公開の型 `Target`。`Interval` は `time.Duration`
- `Load(path string) ([]Target, error)` がファイルを読み、非公開の型に `json.Unmarshal` し、1 件ずつ検証しながら `Target` に詰め替える
- 検証で 1 件でも誤りがあれば `nil` と `error` を返す。エラーにはファイル名、何番目の対象か、何が悪いかを `fmt.Errorf` の `%w` で足す
- 検証する項目: `name` と `url` が空でない、`interval` が `time.ParseDuration` で解釈でき 0 より大きい
- 設定ファイルの実体はリポジトリのルートに `targets.json` として置く。秘密情報を含まないのでコミットする

## 理由

- `Target` を持っている時点で値が正しいことが保証される。使う側で検証を繰り返さなくてよい
- ファイルの形と内部の形を分けておくと、Step 9 で `method` / `headers` / `keyword` が増えたときに、JSON の都合と内部の都合を別々に変えられる
- 設定の誤りは「起動できない」ので、チェック結果 (ADR 0005) と違い通常の `error` で返して呼び出し側を止める

## 不利益

- 型が 2 つになり、フィールドの詰め替えを書く手間がある。対象のフィールドは 3 つなので許容する

## 検討した代替案

- `Interval string` の 1 つの型だけにして使う側で毎回 `ParseDuration` する: 使う側が増えるたびに検証が散らばる
- `time.Duration` を包む独自型に `UnmarshalText` を実装する: 文字列を直接 `time.Duration` 相当に入れられる Go の定石だが、使う側で毎回 `time.Duration(x)` の変換が要る。フィールドが増えてきたら再検討する
- `interval` を秒数の整数にする: `json.Unmarshal` だけで済むが、`"30s"` の方が読み書きしやすく、ナノ秒との混同も防げる
