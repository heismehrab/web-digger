package http

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"web-digger/internal/core/domain/models"
	"web-digger/pkg/mocks"
)

const analyzerURI = "/api/v1/analyze-html"

func TestAnalyzerOnSuccessStatusCode(t *testing.T) {
	mockedAnalyzerService := new(mocks.MockedAnalyzerService)

	mockedServiceResult := models.AnalyzerResult{
		Title:             "mock title",
		Version:           "10.0",
		InaccessibleLinks: 15,
		HasLoginForm:      true,
		Hs: models.HTags{
			1: 0,
			2: 0,
			3: 2,
			4: 5,
			5: 9,
			6: 1,
		},
		InternalLinks: models.Links{
			"in-link-1",
			"in-link-2",
		},
		ExternalLinks: models.Links{
			"ex-link-1",
			"ex-link-2",
		},
	}

	mockedAnalyzerService.On("Analyze", mock.Anything, mock.Anything).Return(
		&mockedServiceResult,
		nil,
	)

	// Configure handler.
	h := setupHandler(t, mockedAnalyzerService)
	router := h.setupRouter()

	reqBody := urlToAnalyzeRequest{
		Url: "https://www.google.com",
	}

	req := createRequest(t, http.MethodPost, analyzerURI, reqBody, map[string]string{
		"content-type": "application/json",
	})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response StandardResponse

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("json unmarshal of target endpoint failed with error: %s", err)
	}

	data := response.Data.(map[string]interface{})

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, mockedServiceResult.Title, data["title"])
	assert.Equal(t, mockedServiceResult.Version, data["version"])
	assert.Equal(t, mockedServiceResult.HasLoginForm, data["hasLoginForm"])
	assert.Equal(t, float64(mockedServiceResult.InaccessibleLinks), data["inaccessibleLinks"])

	headings := data["headings"].(map[string]interface{})

	for h, count := range headings {
		atoi, _ := strconv.Atoi(h)
		assert.Equal(t, mockedServiceResult.Hs[atoi], int(count.(float64)))
	}

	links := data["internalLinks"].([]interface{})
	assert.Equal(t, mockedServiceResult.InternalLinks[0], links[0].(string))
	assert.Equal(t, mockedServiceResult.InternalLinks[1], links[1].(string))

	links = data["externalLinks"].([]interface{})
	assert.Equal(t, mockedServiceResult.ExternalLinks[0], links[0].(string))
	assert.Equal(t, mockedServiceResult.ExternalLinks[1], links[1].(string))
}

func TestAnalyzerOnBadRequestStatusCode(t *testing.T) {
	mockedAnalyzerService := new(mocks.MockedAnalyzerService)
	mockedAnalyzerService.On("Analyze", mock.Anything, mock.Anything).Return(
		nil,
		nil,
	)

	// Configure handler.
	h := setupHandler(t, mockedAnalyzerService)
	router := h.setupRouter()

	// Send an invalid URL (Without HTTPS scheme).
	reqBody := urlToAnalyzeRequest{
		Url: "www.google.com",
	}

	req := createRequest(t, http.MethodPost, analyzerURI, reqBody, map[string]string{
		"content-type": "application/json",
	})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
