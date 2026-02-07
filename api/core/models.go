package core

import (
	"time"
)

type Severity string

const (
	SEVERITY_LOW      Severity = "LOW"
	SEVERITY_MEDIUM   Severity = "MEDIUM"
	SEVERITY_HIGH     Severity = "HIGH"
	SEVERITY_CRITICAL Severity = "CRITICAL"
)

type Status string

const (
	STATUS_OPEN        Status = "OPEN"
	STATUS_IN_PROGRESS Status = "IN_PROGRESS"
	STATUS_RESOLVED    Status = "RESOLVED"
	STATUS_CLOSED      Status = "CLOSED"
)

type Note struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

type Incident struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Severity        Severity  `json:"severity"`
	Status          Status    `json:"status"`
	Reporter        string    `json:"reporter"`
	StartTime       time.Time `json:"start_time"`
	InitialNotes    []Note    `json:"initial_notes"`
	AdditionalNotes []Note    `json:"additional_notes"`
	Owner           string    `json:"owner,omitempty"`
	EndTime         time.Time `json:"end_time,omitempty"`
}

// AddNote adds an additional note to the incident.
func (i *Incident) AddNote(note Note) {
	i.AdditionalNotes = append(i.AdditionalNotes, note)
}

// GetNotes generates a summary report of the incident.
func (i *Incident) GetNotes() []Note {
	return append(i.InitialNotes, i.AdditionalNotes...)
}

// CloseIncident closes the incident by updating its status and end time.
func (i *Incident) CloseIncident() {
	i.Status = STATUS_CLOSED
	i.EndTime = time.Now().UTC()
}

// ResolveIncident resolves the incident by updating its status and end time.
func (i *Incident) ResolveIncident() {
	i.Status = STATUS_RESOLVED
	i.EndTime = time.Now().UTC()
}
