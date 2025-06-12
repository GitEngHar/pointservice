package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"pointservice/domain"
	"time"
)

type PointRepository struct {
	db *sql.DB
}

func NewPointSQL(db *sql.DB) PointRepository {
	return PointRepository{
		db: db,
	}
}

func (p PointRepository) GetPointByUserID(ctx context.Context, userID string) (domain.Point, error) {
	var query = `SELECT point_num FROM point_root WHERE user_id=?`
	var pointNum int
	var createdAt time.Time
	var updatedAt time.Time

	row := p.db.QueryRowContext(ctx, query, userID)
	if err := row.Scan(&pointNum, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Point{}, sql.ErrNoRows
		}
		return domain.Point{}, fmt.Errorf("failed to scan point row: %w (user_id:%s)", err, userID)
	}

	point, err := domain.NewPoint(userID, pointNum, createdAt, updatedAt)
	if err != nil {
		return domain.Point{}, fmt.Errorf("%w", err)
	}
	return point, nil
}

func (p PointRepository) UpdatePointByUserID(ctx context.Context, point domain.Point) error {
	var query = `UPDATE point_root SET point_num=? WHERE user_id=?`
	_, err := p.db.Exec(query, point.PointNum, point.UserID)
	if err != nil {
		return fmt.Errorf("failed to scan point row: %w (point_num:%d user_id:%s)", err, point.PointNum, point.UserID)
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
		return fmt.Errorf("failed to scan point row: %w (point_num:%d user_id:%s)", err, point.PointNum, point.UserID, point.CreatedAt, point.UpdatedAt)
	}
	return nil
}
