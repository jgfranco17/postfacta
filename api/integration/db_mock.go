package integration_test

import (
	"context"

	"github.com/jgfranco17/postfacta/api/entry"
	"github.com/stretchr/testify/mock"
)

type DatabaseClientMock struct {
	mock.Mock
}

func (m *DatabaseClientMock) GetIncidentByID(ctx context.Context, incidentID string) (entry.Incident, error) {
	args := m.Called(ctx, incidentID)
	incident, _ := args.Get(0).(entry.Incident)
	return incident, args.Error(1)
}

func (m *DatabaseClientMock) GetAllIncidents(ctx context.Context) ([]entry.Incident, error) {
	args := m.Called(ctx)
	incidents, _ := args.Get(0).([]entry.Incident)
	return incidents, args.Error(1)
}

func (m *DatabaseClientMock) StoreIncident(ctx context.Context, incident entry.Incident) error {
	args := m.Called(ctx, incident)
	return args.Error(0)
}
