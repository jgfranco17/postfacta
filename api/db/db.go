package db

import (
	"context"

	"github.com/jgfranco17/postfacta/api/entry"
)

// DatabaseClient interface for database operations
type DatabaseClient interface {
	GetIncidentByID(ctx context.Context, incidentID string) (entry.Incident, error)
	GetAllIncidents(ctx context.Context) ([]entry.Incident, error)
	StoreIncident(ctx context.Context, incident entry.Incident) error
}

type DatabaseClientFactory func(ctx context.Context) (*DatabaseClient, error)

// NewClient creates a new Supabase database client
func NewClient(_ context.Context) (DatabaseClient, error) {
	return clientSingleton, nil
}
