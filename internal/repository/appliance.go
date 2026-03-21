package repository

import (
	"context"
	"fmt"
	"laundry-management-system/internal/model"
)

func (r *Repository) GetAppliances(ctx context.Context) ([]*model.Appliance, error) {
	query := `
		select appliance_id, name, type from appliances
	`

	rows, err := r.pgPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appliances []*model.Appliance
	for rows.Next() {
		t := &model.Appliance{}
		err := rows.Scan(
			&t.Id, &t.Name, &t.Type,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan appliances: %w", err)
		}
		appliances = append(appliances, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return appliances, nil
}

func (r *Repository) GetWashingMachines(ctx context.Context) ([]*model.Appliance, error) {
	query := `
		select appliance_id, name, type from appliances
		where type = 'washing_machine'
	`

	rows, err := r.pgPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appliances []*model.Appliance
	for rows.Next() {
		t := &model.Appliance{}
		err := rows.Scan(
			&t.Id, &t.Name, &t.Type,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan appliances: %w", err)
		}
		appliances = append(appliances, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return appliances, nil
}

func (r *Repository) GetTumbleDryers(ctx context.Context) ([]*model.Appliance, error) {
	query := `
		select appliance_id, name, type from appliances
		where type = 'tumble_dryer'
	`

	rows, err := r.pgPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appliances []*model.Appliance
	for rows.Next() {
		t := &model.Appliance{}
		err := rows.Scan(
			&t.Id, &t.Name, &t.Type,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan appliances: %w", err)
		}
		appliances = append(appliances, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return appliances, nil
}
