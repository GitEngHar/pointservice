package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"pointservice/domain"
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
	row := p.db.QueryRowContext(ctx, query, userID)
	var pointNum int
	if err := row.Scan(&pointNum); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Point{}, fmt.Errorf("user not found (user_id=%s) \n", userID)
		}
		return domain.Point{}, fmt.Errorf("failed to scan point row: %w (user_id:%s)", err, userID)
	}

	point, err := domain.NewPoint(userID, pointNum)
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
				INSERT INTO point_root (user_id, point_num) 
				VALUES(?, ?) 
				ON DUPLICATE KEY  UPDATE point_num = VALUES(point_num);
				`
	_, err := p.db.Exec(query, point.UserID, point.PointNum)
	if err != nil {
		return fmt.Errorf("failed to scan point row: %w (point_num:%d user_id:%s)", err, point.PointNum, point.UserID)
	}
	return nil
}
