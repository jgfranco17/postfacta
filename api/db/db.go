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

var clientSingleton DatabaseClient = &internalClient{
	storage: make(map[string]entry.Incident),
}

// NewClient creates a new Supabase database client
func NewClient() (DatabaseClient, error) {
	return clientSingleton, nil
}

type internalClient struct {
	storage map[string]entry.Incident // In-memory storage for demonstration
}

func (ic *internalClient) GetIncidentByID(ctx context.Context, incidentID string) (entry.Incident, error) {
	if incident, exists := ic.storage[incidentID]; exists {
		return incident, nil
	}
	return entry.Incident{}, ErrNotFound
}

func (ic *internalClient) GetAllIncidents(ctx context.Context) ([]entry.Incident, error) {
	allIncidents := []entry.Incident{}
	for _, incident := range ic.storage {
		allIncidents = append(allIncidents, incident)
	}
	return allIncidents, nil
}

func (ic *internalClient) StoreIncident(ctx context.Context, incident entry.Incident) error {
	if _, exists := ic.storage[incident.ID]; exists {
		return ErrConflict
	}
	ic.storage[incident.ID] = incident
	return nil
}
