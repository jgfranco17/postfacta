package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type serviceInfoResponse struct {
	Name        string `json:"name"`
	Environment string `json:"environment"`
	Uptime      string `json:"uptime"`
	Codebase    struct {
		Author     string   `json:"author"`
		Repository string   `json:"repository"`
		Version    string   `json:"version"`
		License    string   `json:"license"`
		Languages  []string `json:"languages"`
		Active     bool     `json:"active"`
	} `json:"codebase"`
}

type healthStatusResponse struct {
	Status string `json:"status"`
}

var _ = Describe("System endpoints", func() {
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

	It("returns ok for root", func() {
		response, _, err := runner.Do(http.MethodGet, "/", nil, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(response.StatusCode).To(Equal(http.StatusOK))
	})

	It("returns ok for /home", func() {
		response, _, err := runner.Do(http.MethodGet, "/home", nil, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(response.StatusCode).To(Equal(http.StatusOK))
	})

	It("returns health status", func() {
		response, body, err := runner.Do(http.MethodGet, "/healthz", nil, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(response.StatusCode).To(Equal(http.StatusOK))

		var payload healthStatusResponse
		Expect(json.Unmarshal(body, &payload)).To(Succeed())
		Expect(payload.Status).To(Equal("healthy"))
	})

	It("returns service info", func() {
		response, body, err := runner.Do(http.MethodGet, "/service-info", nil, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(response.StatusCode).To(Equal(http.StatusOK))

		var payload serviceInfoResponse
		Expect(json.Unmarshal(body, &payload)).To(Succeed())
		Expect(payload.Name).To(Equal("PostFacta API"))
		Expect(payload.Codebase.Author).To(Equal("test"))
		Expect(payload.Codebase.Version).To(Equal("0.0.0"))
		Expect(payload.Codebase.Active).To(BeTrue())
		Expect(payload.Uptime).NotTo(BeEmpty())
	})

	It("returns not found for unknown routes", func() {
		response, body, err := runner.Do(http.MethodGet, "/does-not-exist", nil, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(response.StatusCode).To(Equal(http.StatusNotFound))

		var payload errorResponse
		Expect(json.Unmarshal(body, &payload)).To(Succeed())
		Expect(payload.Message).To(ContainSubstring("does not exist"))
		Expect(payload.Message).To(ContainSubstring("/does-not-exist"))
	})
})
