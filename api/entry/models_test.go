package entry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIncident_AddNote(t *testing.T) {
	t1 := time.Date(2026, 2, 7, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 2, 7, 11, 0, 0, 0, time.UTC)
	t3 := time.Date(2026, 2, 7, 12, 0, 0, 0, time.UTC)
	t0 := time.Date(2026, 2, 7, 9, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		incident      Incident
		notesToAdd    []Note
		expectedCount int
	}{
		{
			name: "add single note to empty incident",
			incident: Incident{
				ID:              "inc-1",
				AdditionalNotes: []Note{},
			},
			notesToAdd: []Note{
				{Timestamp: t1, Message: "First note"},
			},
			expectedCount: 1,
		},
		{
			name: "add multiple notes",
			incident: Incident{
				ID:              "inc-2",
				AdditionalNotes: []Note{},
			},
			notesToAdd: []Note{
				{Timestamp: t1, Message: "First note"},
				{Timestamp: t2, Message: "Second note"},
				{Timestamp: t3, Message: "Third note"},
			},
			expectedCount: 3,
		},
		{
			name: "add note to incident with existing notes",
			incident: Incident{
				ID: "inc-3",
				AdditionalNotes: []Note{
					{Timestamp: t0, Message: "Existing note"},
				},
			},
			notesToAdd: []Note{
				{Timestamp: t1, Message: "New note"},
			},
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, note := range tt.notesToAdd {
				tt.incident.AddNote(note)
			}

			assert.Len(t, tt.incident.AdditionalNotes, tt.expectedCount)

			// Verify the last note added matches
			if len(tt.notesToAdd) > 0 {
				lastAdded := tt.notesToAdd[len(tt.notesToAdd)-1]
				lastInIncident := tt.incident.AdditionalNotes[len(tt.incident.AdditionalNotes)-1]
				assert.Equal(t, lastAdded.Message, lastInIncident.Message)
				assert.Equal(t, lastAdded.Timestamp, lastInIncident.Timestamp)
			}
		})
	}
}

func TestIncident_GetNotes(t *testing.T) {
	t0 := time.Date(2026, 2, 7, 9, 0, 0, 0, time.UTC)
	t1 := time.Date(2026, 2, 7, 9, 30, 0, 0, time.UTC)
	t2 := time.Date(2026, 2, 7, 10, 0, 0, 0, time.UTC)
	t3 := time.Date(2026, 2, 7, 11, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		incident      Incident
		expectedCount int
		expectedOrder []string
	}{
		{
			name: "get notes from empty incident",
			incident: Incident{
				ID:              "inc-1",
				InitialNotes:    []Note{},
				AdditionalNotes: []Note{},
			},
			expectedCount: 0,
			expectedOrder: []string{},
		},
		{
			name: "get only initial notes",
			incident: Incident{
				ID: "inc-2",
				InitialNotes: []Note{
					{Timestamp: t0, Message: "Initial note 1"},
					{Timestamp: t1, Message: "Initial note 2"},
				},
				AdditionalNotes: []Note{},
			},
			expectedCount: 2,
			expectedOrder: []string{"Initial note 1", "Initial note 2"},
		},
		{
			name: "get only additional notes",
			incident: Incident{
				ID:           "inc-3",
				InitialNotes: []Note{},
				AdditionalNotes: []Note{
					{Timestamp: t2, Message: "Additional note 1"},
					{Timestamp: t3, Message: "Additional note 2"},
				},
			},
			expectedCount: 2,
			expectedOrder: []string{"Additional note 1", "Additional note 2"},
		},
		{
			name: "get combined initial and additional notes",
			incident: Incident{
				ID: "inc-4",
				InitialNotes: []Note{
					{Timestamp: t0, Message: "Initial note"},
				},
				AdditionalNotes: []Note{
					{Timestamp: t2, Message: "Additional note 1"},
					{Timestamp: t3, Message: "Additional note 2"},
				},
			},
			expectedCount: 3,
			expectedOrder: []string{"Initial note", "Additional note 1", "Additional note 2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notes := tt.incident.GetNotes()

			assert.Len(t, notes, tt.expectedCount)

			for i, expectedMsg := range tt.expectedOrder {
				assert.Equal(t, expectedMsg, notes[i].Message)
			}
		})
	}
}

func TestIncident_CloseIncident(t *testing.T) {
	tests := []struct {
		name           string
		incident       Incident
		expectedStatus Status
	}{
		{
			name: "close open incident",
			incident: Incident{
				ID:     "inc-1",
				Status: STATUS_OPEN,
			},
			expectedStatus: STATUS_CLOSED,
		},
		{
			name: "close in-progress incident",
			incident: Incident{
				ID:     "inc-2",
				Status: STATUS_IN_PROGRESS,
			},
			expectedStatus: STATUS_CLOSED,
		},
		{
			name: "close resolved incident",
			incident: Incident{
				ID:     "inc-3",
				Status: STATUS_RESOLVED,
			},
			expectedStatus: STATUS_CLOSED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeClose := time.Now().UTC()

			tt.incident.CloseIncident()

			afterClose := time.Now().UTC()

			assert.Equal(t, tt.expectedStatus, tt.incident.Status)
			require.False(t, tt.incident.EndTime.IsZero(), "EndTime should be set")
			assert.True(t, tt.incident.EndTime.After(beforeClose) || tt.incident.EndTime.Equal(beforeClose))
			assert.True(t, tt.incident.EndTime.Before(afterClose) || tt.incident.EndTime.Equal(afterClose))
		})
	}
}

func TestIncident_ResolveIncident(t *testing.T) {
	tests := []struct {
		name           string
		incident       Incident
		expectedStatus Status
	}{
		{
			name: "resolve open incident",
			incident: Incident{
				ID:     "inc-1",
				Status: STATUS_OPEN,
			},
			expectedStatus: STATUS_RESOLVED,
		},
		{
			name: "resolve in-progress incident",
			incident: Incident{
				ID:     "inc-2",
				Status: STATUS_IN_PROGRESS,
			},
			expectedStatus: STATUS_RESOLVED,
		},
		{
			name: "re-resolve already resolved incident",
			incident: Incident{
				ID:     "inc-3",
				Status: STATUS_RESOLVED,
			},
			expectedStatus: STATUS_RESOLVED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeResolve := time.Now().UTC()

			tt.incident.ResolveIncident()

			afterResolve := time.Now().UTC()

			assert.Equal(t, tt.expectedStatus, tt.incident.Status)
			require.False(t, tt.incident.EndTime.IsZero(), "EndTime should be set")
			assert.True(t, tt.incident.EndTime.After(beforeResolve) || tt.incident.EndTime.Equal(beforeResolve))
			assert.True(t, tt.incident.EndTime.Before(afterResolve) || tt.incident.EndTime.Equal(afterResolve))
		})
	}
}

func TestNote_Creation(t *testing.T) {
	t1 := time.Date(2026, 2, 7, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name              string
		timestamp         time.Time
		message           string
		expectedTimestamp time.Time
		expectedMessage   string
	}{
		{
			name:              "create note with standard timestamp",
			timestamp:         t1,
			message:           "Test message",
			expectedTimestamp: t1,
			expectedMessage:   "Test message",
		},
		{
			name:              "create note with empty message",
			timestamp:         t1,
			message:           "",
			expectedTimestamp: t1,
			expectedMessage:   "",
		},
		{
			name:              "create note with multiline message",
			timestamp:         t1,
			message:           "Line 1\nLine 2\nLine 3",
			expectedTimestamp: t1,
			expectedMessage:   "Line 1\nLine 2\nLine 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			note := Note{
				Timestamp: tt.timestamp,
				Message:   tt.message,
			}

			assert.Equal(t, tt.expectedTimestamp, note.Timestamp)
			assert.Equal(t, tt.expectedMessage, note.Message)
		})
	}
}

func TestNew(t *testing.T) {
	testCases := []struct {
		name           string
		request        IncidentRequest
		validateFields func(t *testing.T, incident Incident, request IncidentRequest)
	}{
		{
			name: "create incident with all fields",
			request: IncidentRequest{
				Title:       "Database outage",
				Description: "Primary database is not responding",
				Reporter:    "sre-team",
				Severity:    SEVERITY_CRITICAL,
				Owner:       "platform-team",
				Notes: []Note{
					{Timestamp: time.Date(2026, 3, 1, 10, 0, 0, 0, time.UTC), Message: "Initial detection"},
					{Timestamp: time.Date(2026, 3, 1, 10, 5, 0, 0, time.UTC), Message: "Escalated to oncall"},
				},
			},
			validateFields: func(t *testing.T, incident Incident, request IncidentRequest) {
				assert.NotEmpty(t, incident.ID)
				assert.Equal(t, request.Title, incident.Title)
				assert.Equal(t, request.Description, incident.Description)
				assert.Equal(t, request.Reporter, incident.Reporter)
				assert.Equal(t, request.Severity, incident.Severity)
				assert.Equal(t, request.Owner, incident.Owner)
				assert.Equal(t, STATUS_OPEN, incident.Status)
				assert.Len(t, incident.InitialNotes, 2)
				assert.False(t, incident.StartTime.IsZero())
			},
		},
		{
			name: "create incident with minimal fields",
			request: IncidentRequest{
				Title:       "API latency spike",
				Description: "Response times increased by 200%",
				Reporter:    "monitoring",
				Severity:    SEVERITY_MEDIUM,
			},
			validateFields: func(t *testing.T, incident Incident, request IncidentRequest) {
				assert.NotEmpty(t, incident.ID)
				assert.Equal(t, request.Title, incident.Title)
				assert.Equal(t, request.Description, incident.Description)
				assert.Equal(t, request.Reporter, incident.Reporter)
				assert.Equal(t, request.Severity, incident.Severity)
				assert.Empty(t, incident.Owner)
				assert.Equal(t, STATUS_OPEN, incident.Status)
				assert.Nil(t, incident.InitialNotes)
				assert.False(t, incident.StartTime.IsZero())
			},
		},
		{
			name: "verify unique IDs",
			request: IncidentRequest{
				Title:       "Test incident",
				Description: "Testing ID uniqueness",
				Reporter:    "test",
				Severity:    SEVERITY_LOW,
			},
			validateFields: func(t *testing.T, incident Incident, request IncidentRequest) {
				incident2 := New(request)
				assert.NotEqual(t, incident.ID, incident2.ID, "Each call to New should generate a unique ID")
			},
		},
		{
			name: "verify timestamp is recent",
			request: IncidentRequest{
				Title:       "Timestamp test",
				Description: "Testing that timestamp is set correctly",
				Reporter:    "test",
				Severity:    SEVERITY_LOW,
			},
			validateFields: func(t *testing.T, incident Incident, request IncidentRequest) {
				now := time.Now().UTC()
				timeDiff := now.Sub(incident.StartTime)
				assert.True(t, timeDiff < 1*time.Second, "StartTime should be very recent")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			incident := New(tc.request)
			tc.validateFields(t, incident, tc.request)
		})
	}
}
