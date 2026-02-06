package db

import (
	"context"
)

// DatabaseClient interface for database operations
type DatabaseClient interface {
	StoreIncident(ctx context.Context, incidentID string) error
}
