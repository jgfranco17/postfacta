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

var clientSingleton DatabaseClient = &internalClient{
	storage: make(map[string]core.Incident),
}

// NewClient creates a new Supabase database client
func NewClient() (DatabaseClient, error) {
	return clientSingleton, nil
}

type internalClient struct {
	storage map[string]core.Incident // In-memory storage for demonstration
}

func (ic *internalClient) GetIncidentByID(ctx context.Context, incidentID string) (core.Incident, error) {
	if incident, exists := ic.storage[incidentID]; exists {
		return incident, nil
	}
	return core.Incident{}, ErrNotFound
}

func (ic *internalClient) GetAllIncidents(ctx context.Context) ([]core.Incident, error) {
	allIncidents := []core.Incident{}
	for _, incident := range ic.storage {
		allIncidents = append(allIncidents, incident)
	}
	return allIncidents, nil
}

func (ic *internalClient) StoreIncident(ctx context.Context, incident core.Incident) error {
	if _, exists := ic.storage[incident.ID]; exists {
		return ErrConflict
	}
	ic.storage[incident.ID] = incident
	return nil
}
