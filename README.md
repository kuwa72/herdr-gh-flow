# lead-cli

GitHub Issue駆動 × TDD × Herdr × マルチコーディングエージェント連携ワークフロー CLI。

## 特徴
- `gh issue list` + `fzf` による高速・インタラクティブな Issue 選択（本文・ラベルプレビュー付）。
- 複数エージェント（agy, claude, codex, devin 等）の選択・ローテーション対応。
- Herdr ペイン自動分割 & 最適化された TDD コンテキストの動的注入。
- 堅牢な CI 待機 & 自動マージ (`ci-wait`)。
