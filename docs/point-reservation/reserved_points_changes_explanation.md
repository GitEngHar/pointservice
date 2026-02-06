# 予約ポイントAPI 変更内容解説

今回実装した `GET /users/{user_id}/reserved-points` に関連する変更箇所をファイルごとに解説します。

## 1. ドメイン層 (Domain Layer)

### `internal/domain/reservation.go`

`ReservationRepository` インターフェースに新しいメソッド定義を追加しました。

```go
type ReservationRepository interface {
    // ... (既存メソッド)
    
    // 【追加】 ユーザーIDに紐づく予約一覧を取得する
    // 定義: contextを受け取り、userID(string)を引数に、予約のリスト([]Reservation)とエラーを返す
    FindByUserID(ctx context.Context, userID string) ([]domain.Reservation, error)
}
```

## 2. インフラ層 (Infrastructure Layer)

### `internal/infra/repository/reservation_mysql.go`

MySQL用のリポジトリ実装に `FindByUserID` メソッドを追加しました。

```go
// 【追加】 ユーザーID指定で予約を検索する実装
func (r ReservationRepository) FindByUserID(ctx context.Context, userID string) ([]domain.Reservation, error) {
    // クエリ: 指定したuser_idのレコードを全カラム取得。実行日時(execute_at)の昇順で並べる。
    query := `
        SELECT id, user_id, point_amount, execute_at, status, idempotency_key, created_at, updated_at
        FROM point_reservations
        WHERE user_id = ?
        ORDER BY execute_at ASC`

    // クエリ実行
    rows, err := r.db.QueryContext(ctx, query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to query reservations: %w", err) // エラー時はラップして返す
    }
    defer rows.Close() // 処理終了後に必ずRowsを閉じる

    var reservations []domain.Reservation
    // 結果を1行ずつスキャン
    for rows.Next() {
        var res domain.Reservation
        if err := rows.Scan(
            &res.ID,
            &res.UserID,
            &res.PointAmount,
            &res.ExecuteAt,
            &res.Status,
            &res.IdempotencyKey,
            &res.CreatedAt,
            &res.UpdatedAt,
        ); err != nil {
            return nil, fmt.Errorf("failed to scan reservation: %w", err)
        }
        reservations = append(reservations, res) // スライスに追加
    }

    // イテレーション中のエラーチェック
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows iteration error: %w", err)
    }

    // 【重要】 0件の場合は nil ではなく空配列を返す
    // これにより JSON レスポンスが "null" ではなく "[]" になる
    if reservations == nil {
        return []domain.Reservation{}, nil
    }

    return reservations, nil
}
```

## 3. ユースケース層 (Usecase Layer)

### `internal/usecase/reservation_list.go` (新規作成)

「予約一覧取得」というビジネスロジックを担当する構造体を作成しました。

```go
// ... package & import ...

type (
    // ユースケース本体。リポジトリへの依存を持つ
    ReservationListUsecase struct {
        repo domain.ReservationRepository
    }

    // 入力データ構造
    ReservationListInput struct {
        UserID string
    }

    // 出力データ構造
    ReservationListOutput struct {
        ReservedPoints []ReservedPointDTO
    }

    // APIレスポンス用に整形されたデータ構造 (Data Transfer Object)
    ReservedPointDTO struct {
        Point   int
        Status  string
        AddDate string // RFC3339 format (例: "2026-01-01T10:00:00Z")
    }
)

// コンストラクタ
func NewReservationListUsecase(repo domain.ReservationRepository) *ReservationListUsecase {
    return &ReservationListUsecase{
        repo: repo,
    }
}

// 実行メソッド
func (u *ReservationListUsecase) Execute(ctx context.Context, input *ReservationListInput) (*ReservationListOutput, error) {
    // 1. リポジトリを使ってDBからデータを取得
    reservations, err := u.repo.FindByUserID(ctx, input.UserID)
    if err != nil {
        return nil, fmt.Errorf("failed to get reservations: %w", err)
    }

    // 2. ドメインモデル(Reservation) から レスポンス用DTO(ReservedPointDTO) へ変換
    dtos := make([]ReservedPointDTO, 0, len(reservations))
    for _, r := range reservations {
        dtos = append(dtos, ReservedPointDTO{
            Point:   r.PointAmount,
            Status:  string(r.Status),
            // 時間を文字列 (RFC3339形式) に変換
            AddDate: r.ExecuteAt.Format(time.RFC3339),
        })
    }

    return &ReservationListOutput{
        ReservedPoints: dtos,
    }, nil
}
```

## 4. プレゼンテーション層 (Presentation Layer)

### `internal/presentation/api_handler.go`

HTTPリクエストを受け取り、ユースケースを呼び出すハンドラを追加しました。

```go
// 【追加】 ユーザーの予約ポイント一覧を取得するハンドラ
func (p *PointHandler) GetReservedPoints(c echo.Context) error {
    ctx := c.Request().Context()
    userID := c.Param("user_id") // URLのパスパラメータ (:user_id) から値を取得

    // ユースケースへの入力を作成
    input := &usecase.ReservationListInput{
        UserID: userID,
    }

    // ユースケースを初期化して実行 (DI: p.reservationRepo を渡す)
    uc := usecase.NewReservationListUsecase(p.reservationRepo)
    result, err := uc.Execute(ctx, input)
    if err != nil {
        return handleErr(err) // エラーハンドリング (共通関数)
    }

    // 結果をJSONとして返す (HTTP 200 OK)
    return c.JSON(http.StatusOK, result)
}
```

### `internal/infra/router.go`

URLとハンドラを紐付けました。

```go
func (r *Router) Exec(h *presentation.PointHandler) {
    // ... (他のルート)
    
    // 【追加】 GETメソッドで /users/:user_id/reserved-points にアクセスが来たら
    // h.GetReservedPoints を実行するよう登録
    e.GET("/users/:user_id/reserved-points", h.GetReservedPoints)
    
    // ...
}
```