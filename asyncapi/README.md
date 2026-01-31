# AsyncAPI 導入・運用手順書

本ドキュメントは、PointService アプリケーションにおける非同期メッセージング定義 (AsyncAPI) の導入背景、構成、および運用手順をまとめたものです。
新たに開発環境をセットアップする際や、定義を変更する際の手順として参照してください。

## 1. 導入の概要

RabbitMQ を介したメッセージングインターフェースの仕様を明確化するために、AsyncAPI を導入しました。
Go の構造体定義 (`internal/domain/point.go` 等) を正とし、それらと整合性の取れたインターフェース定義を YAML で管理します。

### 技術スタック・バージョン

- **AsyncAPI Spec**: `3.0.0` (推奨) / `2.6.0` (互換性維持)
  - **v3.0.0**: 最新仕様、より柔軟な定義が可能（`asyncapi-v3.yaml`）
  - **v2.6.0**: レガシー互換性重視（`asyncapi.yaml`）
- **Tools**:
  - `Node.js` (推奨 v18以上)
  - `npm` / `npx` (Node.js に同梱)
  - `@asyncapi/cli` (v5系以上)

## 2. 環境セットアップ

### 2.1. 前提条件

Node.js 18以上がインストールされていることを確認してください。

```bash
node --version  # v18.0.0 以上であることを確認
```

Node.jsがインストールされていない場合：

**macOS (Homebrew)**
```bash
brew install node
```

**Ubuntu/Debian**
```bash
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt-get install -y nodejs
```

**Windows**
- [Node.js公式サイト](https://nodejs.org/)からインストーラーをダウンロード

### 2.2. プロジェクトセットアップ

プロジェクトルートディレクトリで以下を実行：

```bash
# 依存関係のインストール
npm install

# または Yarn を使用する場合
yarn install
```

これにより `@asyncapi/cli` がローカルにインストールされ、`npm run` でコマンドが実行可能になります。

### 2.3. 動作確認

```bash
# バリデーション実行
npm run asyncapi:validate

# HTML生成
npm run asyncapi:generate

# ブラウザでプレビュー（自動でブラウザが開きます）
npm run asyncapi:preview
```

## 3. ディレクトリ構成とファイルの役割

```text
pointservice/
├── asyncapi/
│   ├── asyncapi.yaml              # 【AsyncAPI 2.6.0】ルート定義
│   ├── asyncapi-v3.yaml           # 【AsyncAPI 3.0.0】ルート定義（推奨）
│   ├── README.md                  # 本手順書
│   │
│   ├── servers/                   # 【接続先】（v2のみ）
│   │   └── rabbitmq.yaml
│   │
│   ├── channels/                  # 【通信経路】
│   │   ├── index.yaml             # チャネル一覧（v2のみ）
│   │   ├── point/                 # ポイントドメイン
│   │   │   ├── updated.yaml       # ポイント更新イベント
│   │   │   └── history_archived.yaml
│   │   └── reservation/           # 予約ドメイン
│   │       └── created.yaml       # 予約作成イベント
│   │
│   └── components/                # 【部品】
│       ├── index.yaml             # Components一覧（v2のみ）
│       ├── securitySchemes.yaml
│       ├── messages/              # 【メッセージ定義】
│       │   ├── index.yaml
│       │   ├── point.yaml
│       │   ├── reservation.yaml
│       │   └── point-history-archive.yaml
│       └── schemas/               # 【データ型定義】
│           ├── index.yaml
│           ├── point.yaml
│           ├── reservation.yaml
│           └── point-history-archive.yaml
│
├── html/                          # 生成ドキュメント (v2)
├── html-v3/                       # 生成ドキュメント (v3)
├── docs/asyncapi/                 # 生成Markdownドキュメント
├── package.json                   # Node.js設定・スクリプト定義
└── internal/                      # 実装コード (Go)
```

### ファイルの使い分け

#### AsyncAPI 3.0 (`asyncapi-v3.yaml`) - 推奨
- **特徴**:
  - 最新仕様、より直感的な構造
  - **既存のスキーマファイルを参照** - `components/schemas/` 配下のYAMLを再利用
  - channels/operationsはインライン定義、schemasは外部ファイル参照のハイブリッド構造
  - `operations` による明示的なアクション定義
- **用途**: 新規開発、将来的な標準
- **ビルド**: `bundle`コマンドで全参照を解決してから生成

#### AsyncAPI 2.6 (`asyncapi.yaml`) - レガシー互換
- **特徴**:
  - 既存ツールとの互換性重視
  - ファイル分割による管理（channels, servers, components全て）
  - `publish/subscribe` による定義
- **用途**: 既存システムとの互換性が必要な場合
- **ビルド**: `bundle`コマンドで全参照を解決してから生成

### 外部ファイル参照の仕組み

両バージョンとも、**既存の `components/schemas/` ファイルを参照**する構造になっています：

- **v3**: `asyncapi-v3.yaml` が `components/schemas/*.yaml` を `$ref` で参照
- **v2**: `asyncapi.yaml` が `components/index.yaml` 経由で各ファイルを参照

スキーマ定義の変更は `components/schemas/` 配下のファイルを編集するだけで、v2/v3両方に反映されます。

**重要**: ドキュメント生成時は、`bundle` コマンドで全ての外部参照を1つのファイルに統合してから生成します（`npm run asyncapi:generate` で自動実行）。

## 4. 利用可能なコマンド

### 4.1. AsyncAPI 3.0 用コマンド（推奨）

```bash
# バリデーション（外部参照を含めて検証）
npm run asyncapi:validate

# 外部参照をバンドル（1ファイルに統合）
npm run asyncapi:bundle

# HTMLドキュメント生成（自動でbundle → 生成を実行）
npm run asyncapi:generate

# ブラウザでプレビュー
npm run asyncapi:preview

# Markdownドキュメント生成（自動でbundle → 生成を実行）
npm run asyncapi:docs
```

### 4.2. AsyncAPI 2.6 用コマンド（互換性）

```bash
# バリデーション
npm run asyncapi:validate:v2

# 外部参照をバンドル
npm run asyncapi:bundle:v2

# HTMLドキュメント生成（自動でbundle → 生成を実行）
npm run asyncapi:generate:v2

# ブラウザでプレビュー
npm run asyncapi:preview:v2
```

### 4.3. バンドルコマンドについて

`npm run asyncapi:generate` コマンドは内部で以下を実行します：

1. `asyncapi:bundle` - 外部参照 (`$ref`) を全て解決して1つのファイルに統合
2. 統合されたファイルからHTMLを生成

**生成される中間ファイル**:
- `asyncapi/asyncapi-v3-bundled.yaml` (v3用)
- `asyncapi/asyncapi-v2-bundled.yaml` (v2用)

これらは自動生成されるため、`.gitignore` に追加済みです。

### 4.3. npx による直接実行（グローバルインストール不要）

```bash
# AsyncAPI 3.0
npx @asyncapi/cli validate asyncapi/asyncapi-v3.yaml
npx @asyncapi/cli generate fromTemplate asyncapi/asyncapi-v3.yaml @asyncapi/html-template -o html-v3 --force-write

# AsyncAPI 2.6
npx @asyncapi/cli validate asyncapi/asyncapi.yaml
npx @asyncapi/cli generate fromTemplate asyncapi/asyncapi.yaml @asyncapi/html-template -o html --force-write
```

## 5. 開発フロー (Development Flow)

非同期メッセージングを含む開発を行う際の標準的なワークフローは以下の通りです。

### 5.1. 新規メッセージ・チャネルの追加

**共通: スキーマファイルの作成・更新**

まず、`components/schemas/` に新しいスキーマファイルを作成または既存ファイルを更新：

```bash
# 例: 新しいスキーマファイルを作成
vi asyncapi/components/schemas/new-event.yaml
```

```yaml
type: object
description: 新しいイベントのスキーマ
required:
  - event_id
properties:
  event_id:
    type: string
    description: イベントID
```

**AsyncAPI 3.0の場合（推奨）:**

1. `asyncapi-v3.yaml` を編集
   - `channels` セクションに新しいチャネルを追加
   - `components/schemas` に外部スキーマファイルへの参照を追加
   - `components/messages` にメッセージを定義（スキーマを参照）
   - `operations` にアクションを定義

```yaml
components:
  schemas:
    NewEvent:
      $ref: ./components/schemas/new-event.yaml  # 外部ファイル参照
  messages:
    NewEvent:
      payload:
        $ref: '#/components/schemas/NewEvent'
channels:
  newEvent:
    address: new.event
    messages:
      NewEvent:
        $ref: '#/components/messages/NewEvent'
```

2. バリデーション実行
   ```bash
   npm run asyncapi:validate
   ```

3. ドキュメント生成で確認
   ```bash
   npm run asyncapi:generate
   ```

**AsyncAPI 2.6の場合:**

1. 適切なディレクトリに新しいYAMLファイルを作成
   - スキーマ: `components/schemas/new-event.yaml` （上記で作成済み）
   - メッセージ: `components/messages/new-event.yaml`
   - チャネル: `channels/[domain]/new-event.yaml`

2. 各 `index.yaml` に参照を追加

3. バリデーション実行
   ```bash
   npm run asyncapi:validate:v2
   ```

### 5.2. Goコードとの同期

**重要**: スキーマは `components/schemas/` の個別ファイルで管理されており、v2/v3両方から参照されています。

1. Goの構造体を変更したら、対応する **`components/schemas/[name].yaml`** を更新
   ```bash
   # 例: Point構造体を変更した場合
   vi asyncapi/components/schemas/point.yaml
   ```

2. フィールド名、型、必須項目などを一致させる
   - Go `string` → YAML `type: string`
   - Go `int` → YAML `type: integer`
   - Go `time.Time` → YAML `type: string, format: date-time`

3. バリデーションとドキュメント生成で確認
   ```bash
   npm run asyncapi:validate      # v3検証
   npm run asyncapi:validate:v2   # v2検証（両方通ることを確認）
   npm run asyncapi:generate      # ドキュメント生成
   ```

**メリット**: スキーマファイルを1度更新するだけで、v2/v3両方のドキュメントに反映されます。

### 5.3. ドキュメント公開

```bash
# HTML生成
npm run asyncapi:generate

# 生成されたファイルをブラウザで開く
open html-v3/index.html  # macOS
xdg-open html-v3/index.html  # Linux
```

## 6. トラブルシューティング

### 6.1. `npm: command not found`

Node.jsがインストールされていません。「2.1. 前提条件」を参照してセットアップしてください。

### 6.2. バリデーションエラーが発生する

```bash
# 詳細なエラー情報を表示
npx @asyncapi/cli validate asyncapi/asyncapi-v3.yaml --diagnostics
```

**よくあるエラー:**
- YAML構文エラー: インデントを確認
- `$ref` の参照先が見つからない: パスを確認
- 必須フィールドの欠落: AsyncAPI仕様を確認

### 6.3. ドキュメント生成が失敗する

```bash
# キャッシュをクリア
rm -rf html-v3 html

# 再生成
npm run asyncapi:generate
```

### 6.4. v2とv3のどちらを使うべきか？

- **新規開発**: AsyncAPI 3.0 (`asyncapi-v3.yaml`) を推奨
- **既存システムとの互換性が必要**: AsyncAPI 2.6 (`asyncapi.yaml`)
- **移行**: v2からv3への移行は段階的に実施可能（両方を併存させて運用可能）

## 7. CI/CDへの組み込み

### GitHub Actions の例

```yaml
name: AsyncAPI Validation

on: [push, pull_request]

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '20'
      - run: npm install
      - run: npm run asyncapi:validate
      - run: npm run asyncapi:generate
      - uses: actions/upload-artifact@v3
        with:
          name: asyncapi-docs
          path: html-v3/
```

## 8. 参考リンク

- [AsyncAPI公式ドキュメント](https://www.asyncapi.com/docs)
- [AsyncAPI 3.0 仕様](https://www.asyncapi.com/docs/reference/specification/v3.0.0)
- [AsyncAPI CLI ドキュメント](https://github.com/asyncapi/cli)
- [HTMLテンプレート](https://github.com/asyncapi/html-template)
