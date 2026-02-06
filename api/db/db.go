package db

import (
	"context"

	"github.com/jgfranco17/postfacta/api/core"
)

// DatabaseClient interface for database operations
type DatabaseClient interface {
	GetIncidentByID(ctx context.Context, incidentID string) (core.Incident, error)
	GetAllIncidents(ctx context.Context) ([]core.Incident, error)
	StoreIncident(ctx context.Context, incident core.Incident) error
}

// NewClient creates a new Supabase database client
func NewClient() (DatabaseClient, error) {
	client := internalClient{}
	return &client, nil
}

type internalClient struct {
}

func (ic *internalClient) GetIncidentByID(ctx context.Context, incidentID string) (core.Incident, error) {
	return core.Incident{}, nil
}

func (ic *internalClient) GetAllIncidents(ctx context.Context) ([]core.Incident, error) {
	return []core.Incident{}, nil
}

func (ic *internalClient) StoreIncident(ctx context.Context, incident core.Incident) error {
	return nil
}
