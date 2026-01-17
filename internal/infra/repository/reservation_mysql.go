package repository

import (
	"context"
	"database/sql"
	"fmt"
	"pointservice/internal/domain"
	"time"
)

type ReservationRepository struct {
	db *sql.DB
}

func NewReservationSQL(db *sql.DB) ReservationRepository {
	return ReservationRepository{
		db: db,
	}
}

// データベースに新しい予約を保存（INSERT）する。
func (r ReservationRepository) Create(ctx context.Context, reservation domain.Reservation) error {
	query := `
		INSERT INTO point_reservations 
		(id, user_id, point_amount, execute_at, status, idempotency_key, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		reservation.ID,
		reservation.UserID,
		reservation.PointAmount,
		reservation.ExecuteAt,
		reservation.Status,
		reservation.IdempotencyKey,
		reservation.CreatedAt,
		reservation.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("%w: %s", domain.ErrReservationCreateFailed, err.Error())
	}
	return nil
}

// 「待ち状態」で「そろそろ実行すべき」予約をリストアップしてくる。
func (r ReservationRepository) GetPendingReservations(ctx context.Context, executeAt time.Time) ([]domain.Reservation, error) {

	// 「実行時間が指定時間より前（<= ?）」で、かつ「まだ未処理（PENDING）」なやつを探してくる。
	query := `
		SELECT id, user_id, point_amount, execute_at, status, idempotency_key, created_at, updated_at
		FROM point_reservations
		WHERE execute_at <= ? AND status = ?`

	rows, err := r.db.QueryContext(ctx, query, executeAt, domain.StatusPending)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", domain.ErrReservationGetFailed, err.Error())
	}
	defer rows.Close()

	var reservations []domain.Reservation

	// データベースから取り出した1行のデータを、Goの変数（res）にコピー（Scan）する。
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
			return nil, fmt.Errorf("%w: %s", domain.ErrReservationGetFailed, err.Error())
		}
		reservations = append(reservations, res) // コピーできたデータをリストに追加。
	}

	// ループ中に何かエラーが起きてなかったか確認。
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %s", domain.ErrReservationGetFailed, err.Error())
	}

	return reservations, nil
}

// 予約の状態を更新する。
func (r ReservationRepository) UpdateStatus(ctx context.Context, id string, status domain.ReservationStatus) error {

	// 「IDが一致するやつのステータスと更新日時を書き換える」。
	query := `UPDATE point_reservations SET status = ?, updated_at = ? WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("%w: %s", domain.ErrReservationUpdateFailed, err.Error())
	}

	rowsAffected, err := result.RowsAffected() // 更新された行数を確認。
	if err != nil {
		return fmt.Errorf("%w: %s", domain.ErrReservationUpdateFailed, err.Error())
	}
	if rowsAffected == 0 {
		return domain.ErrReservationNotFound
	}
	return nil
}
