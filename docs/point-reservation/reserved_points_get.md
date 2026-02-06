# 予約ポイントAPI 実装ドキュメント

このドキュメントは、ユーザーの予約ポイントを取得するAPIエンドポイント (`GET /users/{user_id}/reserved-points`) の実装詳細と設計判断をまとめたものです。

## 1. 設計判断 (Design Decisions)

> [!IMPORTANT]
> **レスポンスの `status` フィールドについて**
> 要件定義において矛盾がありましたが、バックエンドが情報を持っているなら提供する方が安全であり、クライアント側で不要であれば無視できるため、レスポンスに `status` フィールドを**含める**こととしました。
>
> **URL設計**
> ファイル構成の `reservation` ドメインにマッピングしつつ、よりRESTfulな `users/{user_id}/reserved-points` を採用しました。

## 2. 実装詳細 (Implementation Details)

OpenAPI仕様書に以下の定義を追加しました。

### スキーマ定義

[reservation.yaml](file:///Users/pinoko/workspace/dev/個人開発/pointservice/openapi/components/schemas/reservation.yaml) に `ReservedPoint` と `ReservedPointList` を追加しました。

```yaml
ReservedPoint:
  type: object
  properties:
    point:
      type: integer
      description: 付与予定ポイント数
      example: 100
    status:
      type: string
      description: 予約ステータス
      example: reserved
    add_date:
      type: string
      format: date-time
      description: 付与予定日
      example: "2026-01-10T13:10:00+09:00"
  required:
    - point
    - status
    - add_date
```

### パス定義

`GET` オペレーションを定義した [list.yaml](file:///Users/pinoko/workspace/dev/個人開発/pointservice/openapi/paths/reservation/list.yaml) を作成しました。

### パス登録

[index.yaml](file:///Users/pinoko/workspace/dev/個人開発/pointservice/openapi/paths/index.yaml) に新しいパスを登録しました。

```yaml
/users/{user_id}/reserved-points:
  $ref: ./reservation/list.yaml
```

## 3. 検証結果

### Lint (構文チェック)
`make openapi-lint` が通過しました（localhostやoperationIdに関する軽微な警告のみ）。

### ドキュメント生成
`make openapi-docs` により、`html/openapi/index.html` にHTMLドキュメントが正常に生成されました。

## 4. API呼び出し例 (Usage Example)

### Request

```bash
curl -X GET "http://localhost:8080/users/User01/reserved-points"
```

### Response

```json
{
  "reserved_points": [
    {
      "point": 100,
      "status": "reserved",
      "add_date": "2026-01-10T13:10:00+09:00"
    }
  ]
}
```
