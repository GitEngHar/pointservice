package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"pointservice/internal/domain"
	"time"
)

type PointRepository struct {
	db *sql.DB
}

func NewPointRepository(db *sql.DB) PointRepository {
	return PointRepository{
		db: db,
	}
}

func (p PointRepository) GetPointByUserID(ctx context.Context, userID string) (domain.Point, error) {
	var (
		query     = `SELECT point_num,created_at,updated_at FROM point_root WHERE user_id=?`
		pointNum  int
		createdAt time.Time
		updatedAt time.Time
	)

	row := p.db.QueryRowContext(ctx, query, userID)
	if err := row.Scan(&pointNum, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Point{}, domain.ErrUserNotFound
		}
		return domain.Point{}, fmt.Errorf("%s: %s", domain.ErrSelectUserID, err.Error())
	}

	point, err := domain.NewPoint(userID, pointNum, createdAt, updatedAt)
	if err != nil {
		return domain.Point{}, err
	}
	return point, nil
}

func (p PointRepository) UpdatePointByUserID(ctx context.Context, point domain.Point) error {
	var query = `UPDATE point_root SET point_num=?, updated_at=? WHERE user_id=?`
	_, err := p.db.Exec(query, point.PointNum, point.UpdatedAt, point.UserID)
	if err != nil {
		return domain.ErrUpdatePoint
	}
	return nil
}

func (p PointRepository) UpdatePointOrCreateByUserID(ctx context.Context, point domain.Point) error {
	var query = `
				INSERT INTO point_root (user_id, point_num, created_at, updated_at) 
				VALUES(?, ?, ?, ?) 
				ON DUPLICATE KEY  UPDATE 
				    point_num = VALUES(point_num), 
					updated_at = VALUES(updated_at);
				`
	_, err := p.db.Exec(query, point.UserID, point.PointNum, point.CreatedAt, point.UpdatedAt)
	if err != nil {
		return domain.ErrCreateOrUpdatePoint

	}
	return nil
}

// AddPointIdempotent adds points to a user idempotently.
// Returns (true, nil) if points were added, (false, nil) if already processed, or (false, error) on failure.
func (p PointRepository) AddPointIdempotent(ctx context.Context, idempotencyKey string, userID string, pointAmount int) (bool, error) {
	now := time.Now()
	txID := idempotencyKey // Use idempotency key as transaction ID for simplicity

	// Try to insert the transaction record first (idempotency check)
	insertTxQuery := `
		INSERT IGNORE INTO point_transactions (id, idempotency_key, user_id, point_amount, created_at)
		VALUES (?, ?, ?, ?, ?)`
	result, err := p.db.ExecContext(ctx, insertTxQuery, txID, idempotencyKey, userID, pointAmount, now)
	if err != nil {
		return false, fmt.Errorf("failed to insert transaction record: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to get rows affected: %w", err)
	}

	// If no rows were inserted, the transaction was already processed
	if rowsAffected == 0 {
		return false, nil
	}

	// Update point_root with the new point amount
	upsertQuery := `
		INSERT INTO point_root (user_id, point_num, created_at, updated_at) 
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE 
			point_num = point_num + VALUES(point_num),
			updated_at = VALUES(updated_at)`
	_, err = p.db.ExecContext(ctx, upsertQuery, userID, pointAmount, now, now)
	if err != nil {
		return false, fmt.Errorf("failed to update point: %w", err)
	}

	return true, nil
}
