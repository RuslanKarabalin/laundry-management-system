package repository

import (
	"context"
	"fmt"
	"laundry-management-system/internal/model"
)

func (r *Repository) GetReservations(ctx context.Context) ([]*model.Reservation, error) {
	query := `
		select * from reservations
	`

	rows, err := r.pgPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []*model.Reservation
	for rows.Next() {
		t := &model.Reservation{}
		err := rows.Scan(
			&t.Id, &t.ApplianceId, &t.UserId, &t.StartTime, &t.EndTime,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reservations: %w", err)
		}
		reservations = append(reservations, t)
	}
	return reservations, nil
}
