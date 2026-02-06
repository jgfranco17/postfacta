package core

import (
	"time"
)

type Severity string

const (
	LOW      Severity = "LOW"
	MEDIUM   Severity = "MEDIUM"
	HIGH     Severity = "HIGH"
	CRITICAL Severity = "CRITICAL"
)

type Status string

const (
	OPEN        Status = "OPEN"
	IN_PROGRESS Status = "IN_PROGRESS"
	RESOLVED    Status = "RESOLVED"
	CLOSED      Status = "CLOSED"
)

type Note struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
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
	i.Status = CLOSED
	i.EndTime = time.Now().UTC()
}

// ResolveIncident resolves the incident by updating its status and end time.
func (i *Incident) ResolveIncident() {
	i.Status = RESOLVED
	i.EndTime = time.Now().UTC()
}
