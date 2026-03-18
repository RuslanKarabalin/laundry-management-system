package repository

import (
	"context"
	"fmt"
	"laundry-management-system/internal/model"
)

func (r *Repository) GetAppliances(ctx context.Context) ([]*model.Appliance, error) {
	query := `
		select * from appliances
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
	return appliances, nil
}
