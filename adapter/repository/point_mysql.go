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
	// TODO: DB トランザクション管理をする
	var query = `SELECT point_num FROM point_root WHERE user_id=?`
	// TODO: 取得した値を処理する
	row := p.db.QueryRowContext(ctx, query, userID)
	var pointNum int
	if err := row.Scan(&pointNum); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Point{}, fmt.Errorf("no point record found for user_id=%s: %w", userID, err)
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
	// TODO: 取得した値を処理する
	_, err := p.db.Exec(query, point.PointNum, point.UserID)
	if err != nil {
		return fmt.Errorf("failed to scan point row: %w (point_num:%d user_id:%s)", err, point.PointNum, point.UserID)
	}
	return nil
}
