package integration_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/jgfranco17/postfacta/api/db"
	"github.com/jgfranco17/postfacta/api/entry"
	"github.com/stretchr/testify/mock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type startIncidentResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

var _ = Describe("Incident endpoints", func() {
	var server *httptest.Server
	var dbMock *DatabaseClientMock
	var runner *HttpTestRunner

	BeforeEach(func() {
		dbMock = &DatabaseClientMock{}
		var err error
		server, err = NewMockServer(dbMock)
		Expect(err).NotTo(HaveOccurred())
		runner = &HttpTestRunner{BaseURL: server.URL}
	})

	AfterEach(func() {
		server.Close()
		dbMock.AssertExpectations(GinkgoT())
	})

	It("returns empty incident list", func() {
		dbMock.On("GetAllIncidents", mock.Anything).Return([]entry.Incident{}, nil).Once()

		response, body, err := runner.Do(http.MethodGet, "/api/v0/incidents", nil, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(response.StatusCode).To(Equal(http.StatusOK))

		var payload []entry.Incident
		Expect(json.Unmarshal(body, &payload)).To(Succeed())
		Expect(payload).To(BeEmpty())
	})

	It("returns incidents", func() {
		incidents := []entry.Incident{
			{
				ID:          "incident-1",
				Title:       "DB outage",
				Description: "Primary database unavailable",
				Severity:    entry.Severity("HIGH"),
				Status:      entry.Status("OPEN"),
				Reporter:    "sre",
				StartTime:   time.Date(2025, time.November, 1, 10, 0, 0, 0, time.UTC),
			},
			{
				ID:          "incident-2",
				Title:       "Cache saturation",
				Description: "Cache hit rate dropped",
				Severity:    entry.Severity("MEDIUM"),
				Status:      entry.Status("IN_PROGRESS"),
				Reporter:    "infra",
				StartTime:   time.Date(2025, time.November, 2, 11, 30, 0, 0, time.UTC),
			},
		}
		dbMock.On("GetAllIncidents", mock.Anything).Return(incidents, nil).Once()

		response, body, err := runner.Do(http.MethodGet, "/api/v0/incidents", nil, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(response.StatusCode).To(Equal(http.StatusOK))

		var payload []entry.Incident
		Expect(json.Unmarshal(body, &payload)).To(Succeed())
		Expect(payload).To(HaveLen(2))
		Expect(payload[0].ID).To(Equal("incident-1"))
		Expect(payload[1].ID).To(Equal("incident-2"))
	})

	It("returns an error for incident list failures", func() {
		dbMock.On("GetAllIncidents", mock.Anything).Return(nil, errors.New("db unavailable")).Once()

		response, _, err := runner.Do(http.MethodGet, "/api/v0/incidents", nil, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
	})

	It("starts a new incident", func() {
		requestBody := map[string]any{
			"title":       "API outage",
			"description": "External API is failing",
			"reporter":    "oncall",
			"severity":    "HIGH",
			"owner":       "platform",
		}
		reader, err := toJSONReader(requestBody)
		Expect(err).NotTo(HaveOccurred())

		dbMock.On("StoreIncident", mock.Anything, mock.MatchedBy(func(incident entry.Incident) bool {
			return incident.Title == "API outage" &&
				incident.Description == "External API is failing" &&
				incident.Reporter == "oncall" &&
				incident.Severity == entry.Severity("HIGH") &&
				incident.Owner == "platform" &&
				incident.Status == entry.Status("OPEN") &&
				!incident.StartTime.IsZero() &&
				incident.ID != ""
		})).Return(nil).Once()

		headers := map[string]string{"Content-Type": "application/json"}
		response, body, err := runner.Do(http.MethodPost, "/api/v0/incidents/start", reader, headers)
		Expect(err).NotTo(HaveOccurred())
		Expect(response.StatusCode).To(Equal(http.StatusCreated))

		var payload startIncidentResponse
		Expect(json.Unmarshal(body, &payload)).To(Succeed())
		Expect(payload.ID).NotTo(BeEmpty())
		Expect(payload.Message).To(Equal("Incident started successfully"))
	})

	It("rejects malformed incident creation requests", func() {
		response, body, err := runner.Do(http.MethodPost, "/api/v0/incidents/start", nil, map[string]string{
			"Content-Type": "application/json",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(response.StatusCode).To(Equal(http.StatusBadRequest))

		var payload errorResponse
		Expect(json.Unmarshal(body, &payload)).To(Succeed())
		Expect(payload.Message).To(ContainSubstring("Invalid request body"))
	})

	It("handles incident conflicts", func() {
		requestBody := map[string]any{
			"title":       "Duplicate incident",
			"description": "Testing conflict response",
			"reporter":    "qa",
			"severity":    "LOW",
		}
		reader, err := toJSONReader(requestBody)
		Expect(err).NotTo(HaveOccurred())

		dbMock.On("StoreIncident", mock.Anything, mock.Anything).Return(db.ErrConflict).Once()

		response, body, err := runner.Do(http.MethodPost, "/api/v0/incidents/start", reader, map[string]string{
			"Content-Type": "application/json",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(response.StatusCode).To(Equal(http.StatusConflict))

		var payload errorResponse
		Expect(json.Unmarshal(body, &payload)).To(Succeed())
		Expect(payload.Message).To(ContainSubstring("same ID"))
	})

	It("handles incident storage failures", func() {
		requestBody := map[string]any{
			"title":       "Storage failure",
			"description": "Testing error response",
			"reporter":    "qa",
			"severity":    "MEDIUM",
		}
		reader, err := toJSONReader(requestBody)
		Expect(err).NotTo(HaveOccurred())

		dbMock.On("StoreIncident", mock.Anything, mock.Anything).Return(errors.New("write failed")).Once()

		response, _, err := runner.Do(http.MethodPost, "/api/v0/incidents/start", reader, map[string]string{
			"Content-Type": "application/json",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
	})

	It("rejects incident with empty title", func() {
		requestBody := map[string]any{
			"title":       "",
			"description": "Valid description text",
			"reporter":    "sre",
			"severity":    "HIGH",
		}
		reader, err := toJSONReader(requestBody)
		Expect(err).NotTo(HaveOccurred())

		response, body, err := runner.Do(http.MethodPost, "/api/v0/incidents/start", reader, map[string]string{
			"Content-Type": "application/json",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(response.StatusCode).To(Equal(http.StatusBadRequest))

		var payload errorResponse
		Expect(json.Unmarshal(body, &payload)).To(Succeed())
		Expect(payload.Message).To(ContainSubstring("Validation failed"))
		Expect(payload.Message).To(ContainSubstring("title"))
	})

	It("rejects incident with title too short", func() {
		requestBody := map[string]any{
			"title":       "AB",
			"description": "Valid description text",
			"reporter":    "sre",
			"severity":    "HIGH",
		}
		reader, err := toJSONReader(requestBody)
		Expect(err).NotTo(HaveOccurred())

		response, body, err := runner.Do(http.MethodPost, "/api/v0/incidents/start", reader, map[string]string{
			"Content-Type": "application/json",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(response.StatusCode).To(Equal(http.StatusBadRequest))

		var payload errorResponse
		Expect(json.Unmarshal(body, &payload)).To(Succeed())
		Expect(payload.Message).To(ContainSubstring("Validation failed"))
		Expect(payload.Message).To(ContainSubstring("title"))
		Expect(payload.Message).To(ContainSubstring("at least 3"))
	})

	It("rejects incident with description too short", func() {
		requestBody := map[string]any{
			"title":       "Valid title",
			"description": "Short",
			"reporter":    "sre",
			"severity":    "HIGH",
		}
		reader, err := toJSONReader(requestBody)
		Expect(err).NotTo(HaveOccurred())

		response, body, err := runner.Do(http.MethodPost, "/api/v0/incidents/start", reader, map[string]string{
			"Content-Type": "application/json",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(response.StatusCode).To(Equal(http.StatusBadRequest))

		var payload errorResponse
		Expect(json.Unmarshal(body, &payload)).To(Succeed())
		Expect(payload.Message).To(ContainSubstring("Validation failed"))
		Expect(payload.Message).To(ContainSubstring("description"))
		Expect(payload.Message).To(ContainSubstring("at least 10"))
	})

	It("rejects incident with invalid severity", func() {
		requestBody := map[string]any{
			"title":       "Valid title",
			"description": "Valid description text here",
			"reporter":    "sre",
			"severity":    "INVALID",
		}
		reader, err := toJSONReader(requestBody)
		Expect(err).NotTo(HaveOccurred())

		response, body, err := runner.Do(http.MethodPost, "/api/v0/incidents/start", reader, map[string]string{
			"Content-Type": "application/json",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(response.StatusCode).To(Equal(http.StatusBadRequest))

		var payload errorResponse
		Expect(json.Unmarshal(body, &payload)).To(Succeed())
		Expect(payload.Message).To(ContainSubstring("Validation failed"))
		Expect(payload.Message).To(ContainSubstring("severity"))
		Expect(payload.Message).To(ContainSubstring("must be one of"))
	})

	It("rejects incident with missing required fields", func() {
		requestBody := map[string]any{
			"title": "Valid title",
		}
		reader, err := toJSONReader(requestBody)
		Expect(err).NotTo(HaveOccurred())

		response, body, err := runner.Do(http.MethodPost, "/api/v0/incidents/start", reader, map[string]string{
			"Content-Type": "application/json",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(response.StatusCode).To(Equal(http.StatusBadRequest))

		var payload errorResponse
		Expect(json.Unmarshal(body, &payload)).To(Succeed())
		Expect(payload.Message).To(ContainSubstring("Validation failed"))
	})
})
