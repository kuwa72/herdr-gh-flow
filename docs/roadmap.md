# 実装フェーズの着手順ロードマップ

- **親トラッキングIssue**: [#50](https://github.com/kuwa72/lead-cli/issues/50)
- **設計トラッキング**: [#29](https://github.com/kuwa72/lead-cli/issues/29)（CLOSED）
- **正本**: 着手順の正本は #50 の sub-issues と blocked-by 関係とする。本書は可読用のミラーであり、変更時は両方を更新する。

## Phase 1 — 土台（順次）

| 順 | Issue | 内容 | 先行 |
|----|-------|------|------|
| 1 | #35 | 足場（go.mod・cmd/lead・version・CI配線） | なし |
| 2 | #38 | テスト基盤（モック・ダミーPATH・golden） | #35 |
| 3 | #36 | ports＋Adapter（gh・herdr・エージェント） | #35 |

## Phase 2 — 操作面（#35・#36 の後）

| 順 | Issue | 内容 | 先行 |
|----|-------|------|------|
| 4 | #39 | CLI骨格（cobra採否の確定含む） | #35・#36 |
| 5 | #37 | 組み込みTUI picker | #36・#39 |
| 6 | #40 | ワークフロー状態（workflows.json・worktree） | #36・#39 |

## Phase 3 — 仕上げ（相互にほぼ独立・並行可）

| 順 | Issue | 内容 | 先行 |
|----|-------|------|------|
| 7 | #41 | マージ・クローズ自動化 | #40 |
| 8 | #44 | setup・completion・doctor・update | #39 |
| 9 | #42 | リリースパイプライン（goreleaser） | #39 |
| 10 | #43 | install.sh | #42 |
| 11 | #48 | 自己外部操作機構（server・socket API） | #40 |
| — | #45 | エージェント向けプロンプト文書 | なし（随時可） |

クリティカルパス: #35 → #36 → #39 → #40 → #41。
