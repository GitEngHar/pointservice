# AsyncAPI 導入・運用手順書

本ドキュメントは、PointService アプリケーションにおける非同期メッセージング定義 (AsyncAPI) の導入背景、構成、および運用手順をまとめたものです。
新たに開発環境をセットアップする際や、定義を変更する際の手順として参照してください。

## 1. 導入の概要
RabbitMQ を介したメッセージングインターフェースの仕様を明確化するために、AsyncAPI を導入しました。
Go の構造体定義 (`internal/domain/point.go` 等) を正とし、それらと整合性の取れたインターフェース定義を YAML で管理します。

### 技術スタック・バージョン
- **AsyncAPI Spec**: `2.6.0`
  - *理由*: 現在の HTML ドキュメント生成テンプレート (`@asyncapi/html-template`) が v3.0.0 に完全対応していないため、安定性を重視し v2.6.0 を採用しています。
- **Tools**:
  - `Node.js` (推奨 v18以上)
  - `npx` (Node.js に同梱)
  - `@asyncapi/cli` (v5系以上推奨)

## 2. ディレクトリ構成
```
pointservice/
├── asyncapi/
│   ├── asyncapi.yaml  # 定義ファイル本体 (v2.6.0)
│   └── README.md      # 本手順書
├── html/              # 生成されたドキュメント (HTML)
└── internal/          # 実装コード (Go)
```

## 3. 環境セットアップ
特別なツールのインストールは不要です。`npx` コマンドが利用可能な環境 (Node.js インストール済み) であれば、誰でも同じバージョンのツールを実行できます。

## 4. 開発フロー (Development Flow)

非同期メッセージングを含む開発を行う際の標準的なワークフローは以下の通りです。

1. **仕様検討 & YAML 更新**
   - 新しいメッセージやチャネルが必要になったら、まず `asyncapi/asyncapi.yaml` を更新します。
   - `schemas` 定義などを Go の構造体と合わせます。

2. **Validation (検証)**
   - 記述した YAML が正しいか確認します。
   - `npx --package @asyncapi/cli@latest asyncapi validate asyncapi/asyncapi.yaml`

3. **ドキュメント生成 & プレビュー**
   - 変更内容を HTML で生成し、ブラウザで確認してチームメンバーと共有・レビューします。
   - `npx @asyncapi/generator@2.6.0 asyncapi/asyncapi.yaml @asyncapi/html-template -o html --force-write`

4. **Go 実装**
   - 確定した仕様に基づいて、Go の Producer/Consumer コードを実装します。

5. **マージ**
   - `asyncapi.yaml`, `html/`, `internal/` のコードを含めてプルリクエストを作成します。

## 5. 運用手順

### 5.1. 定義の変更・更新
1. `asyncapi/asyncapi.yaml` を編集します。
2. Go の構造体が変更された場合は、`schemas` セクションも更新して同期をとってください。

### 5.2. 構文チェック (Validate)
編集した YAML ファイルが正しい構文か確認します。


```bash
npx --package @asyncapi/cli@latest asyncapi validate asyncapi/asyncapi.yaml
```

### 5.3. ドキュメント生成 (HTML)
YAML ファイルから可読性の高い HTML ドキュメントを生成します。

```bash
npx @asyncapi/generator@2.6.0 asyncapi/asyncapi.yaml @asyncapi/html-template -o html --force-write
```

- 生成されたファイル: `html/index.html`
- ブラウザで開くことで内容を確認できます。

## 6. トラブルシューティング

### 生成時に `Template is not compatible...` エラーが出る場合
`@asyncapi/html-template` と AsyncAPI Generator のバージョン互換性が原因です。
**解決策**: 上記コマンドのように `npx --package @asyncapi/cli@latest` を指定し、常に最新の CLI 経由で実行することで解決します。

### `asyncapi: 3.0.0` にしたらエラーになる
HTML テンプレートが v3.0.0 の一部機能に対応していない可能性があります。現状は **v2.6.0** を維持してください。
