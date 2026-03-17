package postgres

import (
	"context"
	"database/sql"
	"shipment/internal/domain"
	"shipment/internal/ports/outbound"
	"time"
)

type ShipmentModel struct {
	ID            string
	Reference     int32
	Origin        string
	Destination   string
	Status        string
	Cost          float64
	DriverRevenue float64
}

type EventModel struct {
	ShipmentID string
	Status     string
	Timestamp  int64
}

type Repository struct {
	db *sql.DB
}

var _ outbound.ShipmentRepository = (*Repository)(nil)

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, s *domain.Shipment) error {
	query := `INSERT INTO shipments (id, reference, origin, destination, status, cost, driver_revenue)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (id) DO UPDATE SET status = $5, cost = $6, driver_revenue = $7`

	_, err := r.db.ExecContext(ctx, query, s.ID, s.Reference, s.Origin, s.Destination, string(s.Status), s.Cost, s.DriverRevenue)
	if err != nil {
		return err
	}

	for _, event := range s.Events {
		_, err := r.db.ExecContext(ctx, `INSERT INTO shipment_events (shipment_id, status, timestamp) VALUES ($1, $2, $3)`,
			s.ID, string(event.Status), event.Timestamp.Unix())
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*domain.Shipment, error) {
	var shipment domain.Shipment
	var status string

	err := r.db.QueryRowContext(ctx, `SELECT id, reference, origin, destination, status, cost, driver_revenue FROM shipments WHERE id = $1`, id).
		Scan(&shipment.ID, &shipment.Reference, &shipment.Origin, &shipment.Destination, &status, &shipment.Cost, &shipment.DriverRevenue)
	if err != nil {
		return nil, err
	}

	shipment.Status = domain.Status(status)

	rows, err := r.db.QueryContext(ctx, `SELECT status, timestamp FROM shipment_events WHERE shipment_id = $1 ORDER BY timestamp ASC`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var eventStatus string
		var timestamp int64
		if err := rows.Scan(&eventStatus, &timestamp); err != nil {
			return nil, err
		}
		shipment.Events = append(shipment.Events, domain.Event{
			Status:    domain.Status(eventStatus),
			Timestamp: time.Unix(timestamp, 0),
		})
	}

	return &shipment, nil
}
