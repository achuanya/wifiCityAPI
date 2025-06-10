package examples

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/wifiCityAPI/config"
	"github.com/gin-gonic/gin/wifiCityAPI/internal/models"
	"github.com/gin-gonic/gin/wifiCityAPI/internal/router"
	"github.com/gin-gonic/gin/wifiCityAPI/internal/service"
	"github.com/gin-gonic/gin/wifiCityAPI/pkg/database"
	"github.com/gin-gonic/gin/wifiCityAPI/pkg/security"
	"github.com/stretchr/testify/assert"
)

var (
	testServer *httptest.Server
	apiSecret  []byte
)

// TestMain sets up the test environment
func TestMain(m *testing.M) {
	// Setup
	gin.SetMode(gin.TestMode)
	os.Setenv("GIN_MODE", "test")

	// The config package's init() will be called, setting up test config.
	if err := config.LoadConfig("../config/config.yaml"); err != nil {
		log.Fatalf("Failed to load config for test: %v", err)
	}

	apiSecret = []byte(config.Cfg.Security.APISecret)

	database.Init()
	r := router.SetupRouter()
	testServer = httptest.NewServer(r)

	// Run tests
	exitCode := m.Run()

	// Teardown
	testServer.Close()
	// Clean up database if needed
	os.Exit(exitCode)
}

// APIClient is a helper to make signed & encrypted API requests
type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
	Secret     []byte
}

func NewAPIClient(baseURL string, secret []byte) *APIClient {
	return &APIClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		Secret:     secret,
	}
}

// sendRequest handles encryption, signing, and sending the request, then decrypts the response.
func (c *APIClient) sendRequest(t *testing.T, method, path string, bodyData any) (*models.EncryptedResponse, error) {
	var reqBody io.Reader
	var encryptedPayload string

	if bodyData != nil {
		jsonData, err := json.Marshal(bodyData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		encrypted, err := security.Encrypt(jsonData, c.Secret)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt request body: %w", err)
		}

		encryptedReqBody, err := json.Marshal(encrypted)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal encrypted request body: %w", err)
		}
		reqBody = bytes.NewBuffer(encryptedReqBody)
		encryptedPayload = encrypted.Data
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := "test-nonce-" + timestamp

	// Signature is timestamp + nonce + base64(encrypted_payload)
	// For GET/DELETE, the payload part is empty.
	sigPayload := timestamp + nonce + encryptedPayload
	mac := hmac.New(sha256.New, c.Secret)
	mac.Write([]byte(sigPayload))
	signature := hex.EncodeToString(mac.Sum(nil))

	req, err := http.NewRequest(method, c.BaseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Timestamp", timestamp)
	req.Header.Set("X-API-Nonce", nonce)
	req.Header.Set("X-API-Signature", signature)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBodyBytes, _ := io.ReadAll(resp.Body)
		t.Logf("Error response body: %s", string(respBodyBytes))

		// Try to decrypt error response
		var encryptedResp security.EncryptedData
		if err := json.Unmarshal(respBodyBytes, &encryptedResp); err == nil {
			decrypted, decryptErr := security.Decrypt(&encryptedResp, c.Secret)
			if decryptErr == nil {
				return nil, fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(decrypted))
			}
		}
		return nil, fmt.Errorf("API error: status=%d, raw_body=%s", resp.StatusCode, string(respBodyBytes))
	}

	var encryptedResp models.EncryptedResponse
	if err := json.NewDecoder(resp.Body).Decode(&encryptedResp); err != nil {
		return nil, fmt.Errorf("failed to decode encrypted response: %w", err)
	}

	return &encryptedResp, nil
}

func TestStoreAndWifiAPI(t *testing.T) {
	client := NewAPIClient(testServer.URL, apiSecret)
	var createdStore models.Store
	var createdWifi models.WifiConfig

	// --- Test Create Store ---
	t.Run("CreateStore", func(t *testing.T) {
		createStoreInput := service.CreateStoreInput{
			Name:     "Test Store",
			Province: "Test Province",
			City:     "Test City",
			Address:  "123 Test Street",
			Phone:    "1234567890",
		}

		encryptedResp, err := client.sendRequest(t, "POST", "/api/v1/stores", createStoreInput)
		assert.NoError(t, err)
		assert.NotNil(t, encryptedResp)

		decryptedData, err := security.Decrypt(&encryptedResp.Data, client.Secret)
		assert.NoError(t, err)

		err = json.Unmarshal(decryptedData, &createdStore)
		assert.NoError(t, err)
		assert.Equal(t, createStoreInput.Name, createdStore.Name)
		assert.True(t, createdStore.StoreID > 0)
		t.Logf("Created store with ID: %d", createdStore.StoreID)
	})

	if createdStore.StoreID == 0 {
		t.Fatal("Store creation failed, cannot proceed with other tests.")
	}

	// --- Test Get Store ---
	t.Run("GetStore", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/stores/%d", createdStore.StoreID)
		encryptedResp, err := client.sendRequest(t, "GET", path, nil)
		assert.NoError(t, err)

		var fetchedStore models.Store
		decryptedData, err := security.Decrypt(&encryptedResp.Data, client.Secret)
		assert.NoError(t, err)
		err = json.Unmarshal(decryptedData, &fetchedStore)
		assert.NoError(t, err)

		assert.Equal(t, createdStore.StoreID, fetchedStore.StoreID)
		assert.Equal(t, createdStore.Name, fetchedStore.Name)
	})
	// --- Test Create Wifi Config ---
	t.Run("CreateWifiConfig", func(t *testing.T) {
		// Example password, in a real scenario this would be pre-encrypted by the client
		dummyEncryptedPassword := "dummy-encrypted-password"

		createWifiInput := service.CreateWifiConfigInput{
			StoreID:           uint(createdStore.StoreID),
			SSID:              "Test-WIFI",
			PasswordEncrypted: dummyEncryptedPassword,
			EncryptionType:    "WPA2",
			WifiType:          "CUSTOMER",
		}

		encryptedResp, err := client.sendRequest(t, "POST", fmt.Sprintf("/api/v1/stores/%d/wifis", createdStore.StoreID), createWifiInput)
		assert.NoError(t, err)

		decryptedData, err := security.Decrypt(&encryptedResp.Data, client.Secret)
		assert.NoError(t, err)

		err = json.Unmarshal(decryptedData, &createdWifi)
		assert.NoError(t, err)
		assert.Equal(t, createWifiInput.SSID, createdWifi.SSID)
		assert.True(t, createdWifi.WifiID > 0)
		t.Logf("Created WiFi config with ID: %d for store %d", createdWifi.WifiID, createdStore.StoreID)
	})

	if createdWifi.WifiID == 0 {
		t.Fatalf("WiFi creation failed for store %d.", createdStore.StoreID)
	}

	// --- Test Delete Wifi Config ---
	t.Run("DeleteWifiConfig", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/stores/%d/wifis/%d", createdStore.StoreID, createdWifi.WifiID)
		_, err := client.sendRequest(t, "DELETE", path, nil)
		// Expecting no error for a successful 204 No Content response
		// The helper might return an error if status is not 2xx. Our helper handles this.
		// A successful DELETE should not return a body to decrypt. The helper needs adjustment or the test should check for this.
		// For now, let's assume the helper is fine with no content on success.
		// Actually sendRequest expects encrypted body, which is not returned on 204.
		// We need to adjust the test or the helper. Let's make a direct request for DELETE.

		// Direct request for DELETE to handle 204 No Content
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		nonce := "test-nonce-delete-" + timestamp
		sigPayload := timestamp + nonce
		mac := hmac.New(sha256.New, client.Secret)
		mac.Write([]byte(sigPayload))
		signature := hex.EncodeToString(mac.Sum(nil))

		req, _ := http.NewRequest("DELETE", client.BaseURL+path, nil)
		req.Header.Set("X-API-Timestamp", timestamp)
		req.Header.Set("X-API-Nonce", nonce)
		req.Header.Set("X-API-Signature", signature)

		resp, err := client.HTTPClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	// --- Test Delete Store ---
	t.Run("DeleteStore", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/stores/%d", createdStore.StoreID)

		// Direct request for DELETE to handle 204 No Content
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		nonce := "test-nonce-delete-" + timestamp
		sigPayload := timestamp + nonce
		mac := hmac.New(sha256.New, client.Secret)
		mac.Write([]byte(sigPayload))
		signature := hex.EncodeToString(mac.Sum(nil))

		req, _ := http.NewRequest("DELETE", client.BaseURL+path, nil)
		req.Header.Set("X-API-Timestamp", timestamp)
		req.Header.Set("X-API-Nonce", nonce)
		req.Header.Set("X-API-Signature", signature)

		resp, err := client.HTTPClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify deletion
		encryptedResp, err := client.sendRequest(t, "GET", path, nil)
		assert.Error(t, err, "Expected an error when getting a deleted store")
		assert.Nil(t, encryptedResp, "Expected no response object on error")
		// The error message from sendRequest should indicate a 404
		assert.Contains(t, err.Error(), "status=404")
	})
}
