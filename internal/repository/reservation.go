package repository

import (
	"context"
	"fmt"
	"laundry-management-system/internal/model"

	"github.com/google/uuid"
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
		err = rows.Scan(
			&t.Id, &t.ApplianceId, &t.UserId, &t.StartTime, &t.EndTime,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reservations: %w", err)
		}
		reservations = append(reservations, t)
	}
	return reservations, nil
}

func (r *Repository) GetReservationsByApplianceId(ctx context.Context, applianceId uuid.UUID) ([]*model.Reservation, error) {
	query := `
		select * from reservations
		where appliance_id = $1
	`

	rows, err := r.pgPool.Query(ctx, query, applianceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []*model.Reservation
	for rows.Next() {
		t := &model.Reservation{}
		err = rows.Scan(
			&t.Id, &t.ApplianceId, &t.UserId, &t.StartTime, &t.EndTime,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reservations: %w", err)
		}
		reservations = append(reservations, t)
	}
	return reservations, nil
}

func (r *Repository) CreateReservationsByApplianceId(ctx context.Context, applianceId uuid.UUID, createReservation *model.CreateReservation) error {
	query := `insert into reservations(appliance_id, user_id, start_time, end_time) values ($1, $2, $3, $4)`

	_, err := r.pgPool.Exec(
		ctx,
		query,
		applianceId,
		createReservation.UserId,
		createReservation.StartTime,
		createReservation.EndTime,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteOldReservations(ctx context.Context) (int64, error) {
	query := `
	delete from reservations where reservation_id in (
		select reservation_id from reservations where end_time < now()
	)`

	cmdTag, err := r.pgPool.Exec(ctx, query)
	if err != nil {
		return 0, err
	}
	return cmdTag.RowsAffected(), nil
}
