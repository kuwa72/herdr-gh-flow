# RFC: テスタビリティ向上と堅牢化のためのシェルスクリプト脱却・実装言語移行検討

- **Issue**: [#22](https://github.com/kuwa72/herdr-gh-flow/issues/22)
- **親トラッキングIssue**: [#29](https://github.com/kuwa72/herdr-gh-flow/issues/29)
- **統合先Issue**: [#16 RFC: hgfの提供形態・基本UX方針の確定](https://github.com/kuwa72/herdr-gh-flow/issues/16)
- **ステータス**: Proposal / RFC
- **更新日**: 2026-09-06

---

## 1. 背景と目的

現在の `bin/hgf` は約 227 行の単一 bash スクリプトとして実装されている。Issue 選択・ブランチ作成・プロンプト生成・エージェント起動・CI 待機という一連のフローを `gh`, `fzf`, `jq`, `git`, `herdr` 等の外部コマンド呼び出しで実現している。

この方式の限界は以下の通り。

| 問題 | 影響 |
| :--- | :--- |
| **単体テスト不可** | 外部コマンドを直接 `$(...)` やパイプで呼んでいるため、戻り値・副作用・標準入出力を制御したモックテストが書きにくい。 |
| **TTY 競合** | `fzf` を子プロセスとして起動すると、標準入力の奪い合いやターミナル制御の不整合が発生しやすい。 |
| **外部コマンドの脆さ** | `gh` や `herdr` の出力形式・フラグ変更、環境依存（`jq` 未インストール等）で動作が変わる。 |
| **型安全性の欠如** | JSON パースや数値・文字列処理が bash 上でのみ行われ、入力検証・エラーハンドリングが弱い。 |
| **機能拡張の障壁** | マルチペインUX、エージェントローテーション、Issue 起票支援等を追加するたびにスクリプトが複雑化する。 |

本 RFC では、これらを解消するための実装言語・TUI ライブラリ・外部 CLI 抽象化・アーキテクチャの選定を行い、#16 の「提供形態・基本 UX 方針確定」への技術インプットとする。

---

## 2. 選定基準

| 基準 | 重視理由 |
| :--- | :--- |
| **テスタビリティ** | 外部コマンドをインターフェースで抽象化し、ユニットテスト時にモック・フェイクを注入できること。 |
| **TUI の堅牢性** | `fzf` 依存を減らし、組み込み TUI でプレビュー・キーバインド・入力制御を一貫して扱えること。 |
| **配布の容易さ** | GitHub Releases, Homebrew, `gh extension` 等への単一バイナリ配布が容易であること。 |
| **`gh` 連携** | `gh` CLI / GitHub API / `gh extension` 機構との親和性が高いこと。 |
| **`herdr` 連携** | Herdr のマルチペイン操作（`pane split`, `pane send-text` 等）を確実に呼び出せること。 |
| **保守性・習熟度** | エコシステムが充実し、開発者が継続的に保守・拡張しやすいこと。 |

---

## 3. 実装言語候補比較

### 3.1 Go

| 項目 | 評価 |
| :--- | :--- |
| **ビルド・配布** | `GOOS`/`GOARCH` で単一静的バイナリをクロスコンパイル。goreleaser 等との相性が非常に良い。 |
| **`gh` 連携** | `github.com/cli/go-gh` という公式 gh 拡張用 Go ライブラリがあり、API 呼び出し・リポジトリ情報取得・認証透過が楽。 |
| **TUI エコシステム** | `charmbracelet/bubbletea` (The ELM Architecture), `bubbles/list`, `bubbles/viewport`, `lipgloss` 等が成熟。 |
| **CLI フレームワーク** | `spf13/cobra`, `urfave/cli` 等が充実。標準 `flag` でも十分。 |
| **テスト** | 標準テスト + `stretchr/testify` で広く使われている。インターフェースベースの DI が自然。 |
| **並行処理** | エージェント起動・CI 待機等の並行タスクを goroutine + channel で表現しやすい。 |
| **欠点** | 実行時パニックの可能性、Rust に比べると型安全性・所有権は弱い。ただし CLI 用途では十分。 |

### 3.2 Rust

| 項目 | 評価 |
| :--- | :--- |
| **パフォーマンス・型安全性** | 高パフォーマンスで、所有権・借用により多くの実行時エラーをコンパイル時に排除できる。 |
| **TUI エコシステム** | `ratatui` が非常に強力で、複雑なレイアウト・キーイベント処理が得意。 |
| **CLI フレームワーク** | `clap` は業界標準で高機能。 |
| **`gh` 連携** | `gh` 公式の Rust SDK はない。`gh` サブコマンドを呼ぶか、GraphQL/REST クライアントを自前で組む必要がある。 |
| **配布** | クロスコンパイルは可能だが、リンカー・ターゲットツールチェインの整備が Go より重い。 |
| **テスト** | モックライブラリ（`mockall` 等）はあるが、DI パターンは Go ほど慣例化していない。 |
| **欠点** | ビルド時間、学習コスト、gh 拡張との親和性が Go に劣る。 |

### 3.3 比較まとめ

| 評価軸 | Go | Rust |
| :--- | :---: | :---: |
| 単体テスト・DI の容易さ | ◎ | ◯ |
| TUI ライブラリ成熟度 | ◎ | ◎ |
| `gh` / `gh extension` 親和性 | ◎ | △ |
| クロスコンパイル・配布 | ◎ | ◯ |
| ビルド速度 | ◎ | △ |
| 型安全性・所有権 | △ | ◎ |
| 保守性・人材プール | ◎ | ◯ |

**総合判断: Go を推奨。**
`gh` 公式 Go ライブラリの存在、TUI ライブラリ `bubbletea` の完成度、クロスコンパイル・配布の容易さが、`herdr-gh-flow` の性質（GitHub CLI 連携型オーケストレータ）に最も合致する。

---

## 4. TUI ライブラリ選定

### 4.1 Go: `charmbracelet/bubbletea` + `bubbles`

- **アーキテクチャ**: The ELM Architecture（Model-Update-View）を採用。状態遷移が明示的でテストしやすい。
- **既存 UX の再現**:
  - `bubbles/list` で Issue 一覧
  - `bubbles/viewport` でプレビュー（右ペイン相当）
  - `lipgloss` でスタイリング
- **キー割当**: `--expect` 相当のカスタムキーバインド（Enter, Ctrl-D, Ctrl-O 等）を `Update` 内で実装。
- **テスト**: `tea.Program` に `tea.WithInput` / `tea.WithOutput` を渡して入出力を制御可能。必要に応じ `charmbracelet/x/exp/teatest` 等のテスト支援ライブラリを検討。

### 4.2 Rust: `ratatui`

- **アーキテクチャ**: 即時モード（immediate mode）TUI。高度なカスタムレイアウトに強い。
- **テスト**: レンダリング出力の文字列比較は可能だが、イベント駆動の状態遷移テストは Go/bubbletea より煩雑。

### 4.3 結論

組み込み TUI には **`charmbracelet/bubbletea` + `bubbles` + `lipgloss`** を採用する。ELM アーキテクチャは `fzf` のプレビュー+選択 UX を型安全かつテスト可能に再実装するのに適している。

---

## 5. 外部 CLI / API の抽象化設計

テスタビリティと環境依存の切り離しのため、外部コマンドは全て Go インターフェースの背後に隠す。

```go
// internal/ports/ports.go
package ports

type IssueService interface {
	ListOpen(ctx context.Context) ([]IssueSummary, error)
	View(ctx context.Context, number int) (*Issue, error)
}

type PaneManager interface {
	Split(ctx context.Context, dir Direction, ratio float64) (PaneID, error)
	SendText(ctx context.Context, pane PaneID, text string) error
}

type AgentLauncher interface {
	Launch(ctx context.Context, agent string, prompt string) error
}

type PromptBuilder interface {
	Build(issue *Issue) (string, error)
}
```

### 5.1 `gh` 連携

- 推奨: `github.com/cli/go-gh` を優先し、`gh` CLI の認証・API 呼び出しを透過利用。
- フォールバック: `go-gh` がカバーしない操作（`gh pr merge` 等）は `exec.Command` で `gh` を呼び出すが、その Adapter も `IssueService` 等のインターフェースで抽象化。
- テスト: `IssueService` のフェイク実装で Issue 一覧・本文を固定値返却。

### 5.2 `herdr` 連携

- 現行スクリプトは `herdr pane split` / `herdr pane send-text` を `exec` 呼び出し。
- 移行後も `herdr` CLI を呼び出すが、`PaneManager` インターフェースを介して呼び出し、テストでは標準入出力・終了コードを記録するダミー `herdr` スクリプトで検証（AGENTS.md の「モックによる振る舞いテスト」に沿う）。
- Herdr 非検知環境では `InlinePaneManager` 等のフォールバック実装を提供し、単一端末でも動作可能にする。

### 5.3 エージェント連携

- `agy`, `devin`, `opencode`, `claude`, `codex` 等の起動を `AgentLauncher` インターフェースに集約。
- プロンプトは一時ファイルに書き出し、エージェントコマンドにパスを渡す方式は維持（エージェント側の文字列解析問題を回避）。
- テスト: ダミーエージェントスクリプトを `PATH` に配置し、受け取った引数とプロンプト内容をアサート。

### 5.4 プロンプト生成

- `PromptBuilder` インターフェースに切り出し、Issue タイトル・本文・受入条件を動的注入。
- プロンプトテンプレートは Go の `text/template` で管理し、テストでは golden file または文字列包含で検証。

---

## 6. 提案アーキテクチャ

```
.
├── bin/hgf              # 移行完了後は Go バイナリへの薄いラッパーまたは削除
├── cmd/hgf/main.go      # エントリポイント
├── internal/
│   ├── cli/             # コマンドパース (cobra or 標準 flag)
│   ├── config/          # 環境判定 (HERDR_ENV, エージェント優先順等)
│   ├── core/            # ドメインロジック（ブランチ名生成、ワークフロー制御）
│   ├── ports/           # 外部連携用インターフェース
│   ├── adapters/
│   │   ├── gh/          # github.com/cli/go-gh ベース
│   │   ├── herdr/       # herdr CLI 呼び出し
│   │   ├── agent/       # 各エージェント起動
│   │   └── prompter/    # プロンプトテンプレート
│   └── tui/             # bubbletea ベース Issue 選択 + プレビュー
├── go.mod
└── test/
    └── run-tests.sh     # go test ./... + 結合スクリプト
```

### 6.1 主要フロー

1. `cmd/hgf` がサブコマンドを解析。
2. `config` で `HERDR_ENV`・`gh` 認証状態・エージェント可用性を判定。
3. `core` が Issue 選択、ブランチ名生成、プロンプト生成を制御。
4. `adapters/gh` で Issue 一覧・本文を取得。
5. `tui` で `fzf` 相当のリスト+プレビューを表示。キー入力でエージェント決定。
6. `adapters/herdr` または `config` のフォールバックでエージェントを起動。
7. エージェント完了後、`bin/ci-wait` 相当の CI 監視を `adapters/gh` 経由で実行。

---

## 7. テスト戦略

1. **インターフェースモック**
   - `IssueService`, `PaneManager`, `AgentLauncher` 等に対し、テスト用フェイクを実装。
   - `core` の振る舞いはこれらのフェイクに対してユニットテスト。

2. **TUI テスト**
   - `bubbletea` の `tea.Program` に `WithInput`/`WithOutput` を指定し、キー入力シーケンスと最終出力をアサート。
   - 極力ビジネスロジックは `Update` 関数の純粋関数部分に分離し、状態遷移を直接テスト。

3. **外部コマンドの振る舞いテスト**
   - `herdr`, `gh`, `agy` 等を一時ディレクトリにダミースクリプトとして配置し、`PATH` を注入。
   - ダミーは受け取った引数・標準入力をログに書き出し、テスト後にそのログを検証（AGENTS.md ルール 3 に沿う）。

4. **Golden ファイル**
   - プロンプト生成結果は `test/fixtures/` 下の golden ファイルと比較。

---

## 8. 移行ロードマップ

| フェーズ | 内容 |
| :--- | :--- |
| **0** | RFC 承認（本ドキュメント） |
| **1** | `go.mod` 作成、`internal/ports`・`internal/adapters` 整備、`gh`/`herdr`/`agent` の Adapter をテスト込みで実装 |
| **2** | `internal/tui` を `bubbletea` で実装。`fzf` からの段階的移行（環境変数 or フラグで切り替え） |
| **3** | `cmd/hgf` 完成、`bin/hgf` を Go バイナリに置き換え、`test/run-tests.sh` を `go test` 中心に更新 |
| **4** | `gh extension` 用エントリポイント追加、配布・CI 整備（#26 連携） |

---

## 9. 決定事項

1. **実装言語**: Go
2. **TUI ライブラリ**: `charmbracelet/bubbletea` + `bubbles` + `lipgloss`
3. **外部 CLI 抽象化**: `internal/ports` に `IssueService`, `PaneManager`, `AgentLauncher`, `PromptBuilder` 等を定義し、Adapter パターンで実装
4. **テスト方針**: インターフェースモック + ダミーコマンド PATH 注入 + golden ファイル
5. **配布前提**: 単一静的バイナリ、`gh extension` 機構との親和性を重視

---

## 10. 親Issue #16 へのインプット

- **技術基盤**: シングルバイナリ（Go）による `hgf` 再実装が可能であること。
- **UX 基盤**: `fzf` 依存を脱却し、組み込み TUI（bubbletea）で左右マルチペイン・Human-in-the-loop 体験を再現可能。
- **提供形態**: `gh extension` との親和性が高く、GitHub Releases / Homebrew 等への配布も容易。
- **後続検討事項**: ブランチ命名規約・プロンプトテンプレート・エージェント選択キーの詳細は #25, #26 で詰める。
