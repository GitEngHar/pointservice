# OpenAPI 運用手順書

本プロジェクトでは、OpenAPI 仕様書を分割管理し、メンテナビリティを向上させています。

## 1. ディレクトリ構成

```text
openapi/
├── openapi.yaml               # 【ルート定義】 全体の目次 ($ref のみ)
├── README.md                  # 本手順書
│
├── servers/                   # 【接続先】
│   └── local.yaml
│
├── paths/                     # 【API 定義】 URLパスごとの定義
│   ├── index.yaml             # パス一覧 (全パスをここで定義)
│   ├── point/                 # ポイント関連ドメイン
│   │   ├── add.yaml
│   │   ├── sub.yaml
│   │   └── confirm.yaml
│   └── reservation/           # 予約関連ドメイン
│       └── create.yaml
│
└── components/                # 【共通部品】
    └── schemas/               # 【データ型】 リクエスト/レスポンスの型定義
        ├── index.yaml         # スキーマ一覧
        ├── point.yaml
        ├── reservation.yaml
        └── common.yaml
```

## 2. 各ファイルの役割

- **`openapi.yaml` (ルート)**
  - 全体のエントリーポイントです。`paths` や `components` の実体は持たず、すべて `$ref` で参照します。

- **`paths/` (API 定義)**
  - ドメイン単位（例: `point`, `reservation`）でフォルダを分け、その中に YAML ファイルを配置します。
  - **追加手順**:
    1. `paths/<domain>/<action>.yaml` を作成
    2. `paths/index.yaml` にフルパス (`/point/add` 等) とファイル参照を追記

- **`components/schemas/` (データ型)**
  - リクエストボディやレスポンスの構造を定義します。
  - Go の構造体と対になるように管理します。
  - **追加手順**:
    1. `schemas/<domain>.yaml` に Schema 定義を追加
    2. `schemas/index.yaml` に参照を追加

## 3. API 追加フロー

新しい API を追加する際は、以下の順番で作業します。

1. **Schema 定義**: `components/schemas/` に必要なデータ型を追加・更新
2. **Path 定義**: `paths/<domain>/` に新しい YAML ファイルを作成
3. **Index 更新**: `paths/index.yaml` にパスと参照を追記
4. **実装**: 定義に合わせて Go ハンドラを実装

## 4. 運用ルール

- **`$ref` の活用**: ファイルが巨大になるのを防ぐため、定義は可能な限り分割し、`index.yaml` で集約します。
- **ドメイン駆動**: 機能単位ではなく、ビジネスドメイン単位（ポイント、予約など）でディレクトリを分けます。

## 5. ツール (Validation & Document Generation)

API 定義の検証やドキュメント生成には [Redocly CLI](https://redocly.com/docs/cli/) を推奨します。
特別なインストールは不要で、`npx` コマンドで実行可能です。

### 5.1. 構文チェック (Lint)

`$ref` のリンク切れや OpenAPI の構文エラーをチェックします。

```bash
# コマンド直接実行
npx --package @redocly/cli@latest redocly lint openapi/openapi.yaml

# または Makefile
make openapi-lint
```

### 5.2. ドキュメント生成 (HTML)

静的な HTML ファイルとしてドキュメントを生成します。

```bash
# コマンド直接実行 (ディレクトリは任意の場合)
npx --package @redocly/cli@latest redocly build-docs openapi/openapi.yaml --output html/openapi/index.html

# または Makefile
make openapi-docs
```

- 生成されたファイルはブラウザで開いて閲覧可能です。
- レビュー時や仕様共有時に活用してください。

