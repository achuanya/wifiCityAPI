package examples

import (
	"app/config"
	"app/internal/models"
	"app/internal/router"
	"app/internal/service"
	"app/pkg/database"
	"app/pkg/security"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

var (
	testServer *httptest.Server
	apiSecret  []byte
)

// EncryptedResponse is a standard wrapper for encrypted API responses.
type EncryptedResponse struct {
	Data security.EncryptedData `json:"data"`
}

// TestMain sets up the test environment
func TestMain(m *testing.M) {
	// ... existing code ...
	HTTPClient *http.Client
	Secret     []byte
}

func NewAPIClient(baseURL string, secret []byte) *APIClient {
	// ... existing code ...
}

// sendRequest handles encryption, signing, and sending the request, then decrypts the response.
func (c *APIClient) sendRequest(t *testing.T, method, path string, bodyData any) (*EncryptedResponse, error) {
	var reqBody io.Reader
	var encryptedPayload string

	// ... existing code ...
		return nil, fmt.Errorf("API error: status=%d, raw_body=%s", resp.StatusCode, string(respBodyBytes))
	}

	var encryptedResp EncryptedResponse
	if err := json.NewDecoder(resp.Body).Decode(&encryptedResp); err != nil {
		return nil, fmt.Errorf("failed to decode encrypted response: %w", err)
	}

	return &encryptedResp, nil
}

func TestStoreAndWifiAPI(t *testing.T) {
	client := NewAPIClient(testServer.URL, apiSecret)
	// ... existing code ...
}
