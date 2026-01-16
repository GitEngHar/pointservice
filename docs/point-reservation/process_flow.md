# 処理の流れ（ファイル・メソッド単位）

ポイント付与予約システムおよび集計バッチの各ステップにおける詳細な処理フローです。

## 1. 予約登録フェーズ (同期)

ユーザーが予約APIを呼び出してからDBに保存されるまでの流れ。

1.  **[point_handler.go](pointservice/internal/presentation/point_handler.go)**
    - `PointReserve(c echo.Context)`: リクエストを `ReservationCreateInput` にバインドし、ユースケースを呼び出す。
2.  **[reservation_create.go](pointservice/internal/usecase/reservation_create.go)**
    - `Execute(ctx, input)`: ドメイン層の `NewReservation` を呼び出してバリデーション済みのモデルを作成し、リポジトリに保存。
3.  **[reservation.go](pointservice/internal/domain/reservation.go)**
    - `NewReservation(...)`: ユーザーIDの形式、ポイント量、UUIDの生成などを行う。
4.  **[reservation_mysql.go](pointservice/internal/infra/repository/reservation_mysql.go)**
    - `Create(ctx, reservation)`: `point_reservations` テーブルにデータを `INSERT`。

---

## 2. スキャン＆配送フェーズ (非同期スケジューラ)

一定間隔で実行待ちの予約を見つけ出し、キューに積む流れ。

1.  **[scheduler/main.go](pointservice/cmd/scheduler/main.go)**
    - 10秒ごとのティッカーで `scanAndPublish` を実行。
2.  **[reservation_mysql.go](pointservice/internal/infra/repository/reservation_mysql.go)**
    - `GetPendingReservations(ctx, now)`: `execute_at <= now` かつ `status = 'PENDING'` のレコードを取得。
3.  **[reservation_mysql.go](pointservice/internal/infra/repository/reservation_mysql.go)**
    - `UpdateStatus(ctx, id, 'PROCESSING')`: 重複スキャンを防ぐために一時的にステータスを更新。
4.  **[rabbit_producer.go](pointservice/internal/infra/aync/mq/rabbit_producer.go)**
    - `PublishReservation(ctx, msg)`: `reservationQueue` に予約詳細情報を `Publish`。

---

## 3. ポイント付与フェーズ (非同期ワーカー)

キューからメッセージを受け取り、実際にポイントを加算する流れ。

1.  **[worker/main.go](pointservice/cmd/worker/main.go)**
    - `processMessage(ctx, msg, pointRepo, reservationRepo)`: メッセージをデコードして処理を開始。
2.  **[point_mysql.go](pointservice/internal/infra/repository/point_mysql.go)**
    - `AddPointIdempotent(ctx, idempotencyKey, userID, amount)`:
        - `point_transactions` に `INSERT IGNORE` を試行。
        - 成功（新規）なら `point_root` を `ON DUPLICATE KEY UPDATE` で加算。
        - 失敗（既に存在）なら何もせず成功として返す（冪等性）。
3.  **[reservation_mysql.go](pointservice/internal/infra/repository/reservation_mysql.go)**
    - `UpdateStatus(ctx, id, 'DONE')`: 予約完了を記録。
4.  **[worker/main.go](pointservice/cmd/worker/main.go)**
    - `msg.Ack(false)`: メッセージ処理完了をRabbitMQに通知。