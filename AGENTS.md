# herdr-gh-flow project instructions

## Issue 開発フロー

Issue に着手するときは、次の順序を必ず守る。

1. `issue/<number>-<slug>` ブランチを作成する。
2. TDD で実装する。まず `test/` に失敗するテストを追加し、Red を確認してから実装して Green を確認する。
3. ローカルテスト `./test/run-tests.sh` をすべてパスさせる。
4. base を `main` にして PR を作成する (`gh pr create --fill`)。
5. `bin/ci-wait <PR番号>` を実行して全 CI が pass するまで待つ。`gh pr checks` を手動で連打しない。
6. CI 成功後に PR をマージする (`gh pr merge --squash` または `--auto`)。
7. マージ後に関連 Issue が closed になったことと、マージコミットを確認する。

PR 作成や CI 起動だけでは Issue 対応を完了としない。CI 終了確認・マージ・Issue クローズ確認までを完了条件とする。
