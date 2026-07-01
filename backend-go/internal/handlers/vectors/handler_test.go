package vectors

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/BenedictKing/ccx/internal/config"
	"github.com/BenedictKing/ccx/internal/metrics"
	"github.com/BenedictKing/ccx/internal/scheduler"
	"github.com/BenedictKing/ccx/internal/session"
	"github.com/gin-gonic/gin"
)

func newVectorsTestConfigManager(t *testing.T) *config.ConfigManager {
	t.Helper()
	cfgFile := t.TempDir() + "/config.json"
	if err := os.WriteFile(cfgFile, []byte(`{"upstream":[],"vectorsUpstream":[]}`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	cfgManager, err := config.NewConfigManager(cfgFile, "")
	if err != nil {
		t.Fatalf("config manager: %v", err)
	}
	return cfgManager
}

func newVectorsTestScheduler(cfgManager *config.ConfigManager, vectorsMetrics *metrics.MetricsManager) *scheduler.ChannelScheduler {
	if vectorsMetrics == nil {
		vectorsMetrics = metrics.NewMetricsManager()
	}
	return scheduler.NewChannelScheduler(
		cfgManager,
		metrics.NewMetricsManager(),
		metrics.NewMetricsManager(),
		metrics.NewMetricsManager(),
		metrics.NewMetricsManager(),
		metrics.NewMetricsManager(),
		session.NewTraceAffinityManager(),
		nil,
		vectorsMetrics,
	)
}

func newVectorsTestEnvConfig() *config.EnvConfig {
	envCfg := config.NewEnvConfig()
	envCfg.ProxyAccessKey = "test-proxy-key"
	return envCfg
}

func serveVectorsEmbeddingRequest(cfgManager *config.ConfigManager, sch *scheduler.ChannelScheduler, body string) *httptest.ResponseRecorder {
	r := gin.New()
	r.POST("/v1/embeddings", Handler(newVectorsTestEnvConfig(), cfgManager, sch))

	req := httptest.NewRequest(http.MethodPost, "/v1/embeddings", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer test-proxy-key")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestBuildEmbeddingsURL(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		want    string
	}{
		{name: "root", baseURL: "https://api.openai.com", want: "https://api.openai.com/v1/embeddings"},
		{name: "versioned", baseURL: "https://api.openai.com/v1", want: "https://api.openai.com/v1/embeddings"},
		{name: "hash", baseURL: "https://api.openai.com#", want: "https://api.openai.com/embeddings"},
		{name: "slash hash", baseURL: "https://api.openai.com/#", want: "https://api.openai.com/embeddings"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildEmbeddingsURL(tt.baseURL); got != tt.want {
				t.Fatalf("buildEmbeddingsURL(%q) = %q, want %q", tt.baseURL, got, tt.want)
			}
		})
	}
}

func TestBuildProviderRequestAppliesMappingAndHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/embeddings?encoding_format=float", strings.NewReader(""))
	c.Request.Header.Set("Authorization", "Bearer client-key")
	c.Request.Header.Set("X-Forwarded-For", "127.0.0.1")

	upstream := &config.UpstreamConfig{
		ServiceType:   "openai",
		AuthHeader:    "x-api-key",
		ModelMapping:  map[string]string{"embed-public": "text-embedding-3-small"},
		CustomHeaders: map[string]string{"X-Custom": "yes"},
	}
	bodyBytes := []byte(`{"model":"embed-public","input":"hello"}`)
	req, err := buildProviderRequest(c, upstream, "https://api.example.com/v1", "sk-test", bodyBytes, "embed-public")
	if err != nil {
		t.Fatalf("buildProviderRequest() error = %v", err)
	}
	if got := req.URL.String(); got != "https://api.example.com/v1/embeddings?encoding_format=float" {
		t.Fatalf("unexpected url: %s", got)
	}
	if got := req.Header.Get("x-api-key"); got != "sk-test" {
		t.Fatalf("x-api-key = %q, want sk-test", got)
	}
	if got := req.Header.Get("Authorization"); got != "" {
		t.Fatalf("Authorization should be removed, got %q", got)
	}
	if got := req.Header.Get("X-Custom"); got != "yes" {
		t.Fatalf("X-Custom = %q, want yes", got)
	}
	if got := req.Header.Get("X-Forwarded-For"); got != "" {
		t.Fatalf("X-Forwarded-For should be removed, got %q", got)
	}

	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(requestBody, &payload); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if got := payload["model"]; got != "text-embedding-3-small" {
		t.Fatalf("model = %v, want text-embedding-3-small", got)
	}
}

func TestParseEmbeddingsRequestValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name string
		body string
		ok   bool
	}{
		{name: "valid string", body: `{"model":"text-embedding-3-small","input":"hello"}`, ok: true},
		{name: "valid array", body: `{"model":"text-embedding-3-small","input":["hello"]}`, ok: true},
		{name: "missing model", body: `{"input":"hello"}`, ok: false},
		{name: "missing input", body: `{"model":"text-embedding-3-small"}`, ok: false},
		{name: "empty string input", body: `{"model":"text-embedding-3-small","input":""}`, ok: false},
		{name: "empty array input", body: `{"model":"text-embedding-3-small","input":[]}`, ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			_, _, ok := parseEmbeddingsRequest(c, []byte(tt.body))
			if ok != tt.ok {
				t.Fatalf("parseEmbeddingsRequest() ok = %v, want %v", ok, tt.ok)
			}
		})
	}
}

func TestExtractEmbeddingsUsage(t *testing.T) {
	usage := extractEmbeddingsUsage([]byte(`{"usage":{"total_tokens":17}}`))
	if usage == nil {
		t.Fatal("expected usage")
	}
	if usage.InputTokens != 17 || usage.OutputTokens != 0 || usage.PromptTokens != 17 || usage.PromptTokensTotal != 17 {
		t.Fatalf("unexpected usage: %+v", usage)
	}
}

func TestHandlerFailoverAndUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfgManager := newVectorsTestConfigManager(t)
	defer cfgManager.Close()

	var attempts int32
	upstreamServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		if r.URL.Path != "/v1/embeddings" {
			t.Errorf("unexpected upstream path: %s", r.URL.Path)
		}
		if r.Header.Get("X-Custom") != "yes" {
			t.Errorf("missing custom header")
		}
		if strings.Contains(r.Header.Get("Authorization"), "sk-bad") {
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":{"message":"rate limited"}}`))
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read upstream body: %v", err)
		}
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Errorf("decode upstream body: %v", err)
		}
		if got := payload["model"]; got != "text-embedding-3-small" {
			t.Errorf("upstream model = %v, want text-embedding-3-small", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"object":"list","data":[],"usage":{"prompt_tokens":12,"total_tokens":12}}`))
	}))
	defer upstreamServer.Close()

	if err := cfgManager.AddVectorsUpstream(config.UpstreamConfig{
		Name:          "vectors-test",
		ServiceType:   "openai",
		BaseURL:       upstreamServer.URL,
		APIKeys:       []string{"sk-bad", "sk-good"},
		ModelMapping:  map[string]string{"embed-public": "text-embedding-3-small"},
		CustomHeaders: map[string]string{"X-Custom": "yes"},
	}); err != nil {
		t.Fatalf("AddVectorsUpstream() error = %v", err)
	}

	vectorsMetrics := metrics.NewMetricsManager()
	sch := newVectorsTestScheduler(cfgManager, vectorsMetrics)
	r := gin.New()
	r.POST("/v1/embeddings", Handler(newVectorsTestEnvConfig(), cfgManager, sch))

	req := httptest.NewRequest(http.MethodPost, "/v1/embeddings", strings.NewReader(`{"model":"embed-public","input":"hello"}`))
	req.Header.Set("Authorization", "Bearer test-proxy-key")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	if got := atomic.LoadInt32(&attempts); got != 2 {
		t.Fatalf("attempts = %d, want 2", got)
	}
}

func TestHandlerAppliesVectorsModelMappingToUpstreamBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfgManager := newVectorsTestConfigManager(t)
	defer cfgManager.Close()

	var upstreamModel string
	upstreamServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read upstream body: %v", err)
			return
		}
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Errorf("decode upstream body: %v", err)
			return
		}
		upstreamModel, _ = payload["model"].(string)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"object":"list","data":[],"usage":{"prompt_tokens":1,"total_tokens":1}}`))
	}))
	defer upstreamServer.Close()

	if err := cfgManager.AddVectorsUpstream(config.UpstreamConfig{
		Name:         "jina-vectors",
		ServiceType:  "openai",
		BaseURL:      upstreamServer.URL,
		APIKeys:      []string{"sk-jina"},
		ModelMapping: map[string]string{"text-embedding-3-small": "jina-embeddings-v2-base-zh"},
	}); err != nil {
		t.Fatalf("AddVectorsUpstream() error = %v", err)
	}

	sch := newVectorsTestScheduler(cfgManager, metrics.NewMetricsManager())
	w := serveVectorsEmbeddingRequest(cfgManager, sch, `{"model":"text-embedding-3-small","input":"hello"}`)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	if upstreamModel != "jina-embeddings-v2-base-zh" {
		t.Fatalf("upstream model = %q, want jina-embeddings-v2-base-zh", upstreamModel)
	}
}

func TestHandlerPassesThroughVectorsModelWhenMappingMisses(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfgManager := newVectorsTestConfigManager(t)
	defer cfgManager.Close()

	var upstreamModel string
	upstreamServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read upstream body: %v", err)
			return
		}
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Errorf("decode upstream body: %v", err)
			return
		}
		upstreamModel, _ = payload["model"].(string)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"object":"list","data":[],"usage":{"prompt_tokens":1,"total_tokens":1}}`))
	}))
	defer upstreamServer.Close()

	if err := cfgManager.AddVectorsUpstream(config.UpstreamConfig{
		Name:         "openai-vectors",
		ServiceType:  "openai",
		BaseURL:      upstreamServer.URL,
		APIKeys:      []string{"sk-openai"},
		ModelMapping: map[string]string{"embed-public": "jina-embeddings-v2-base-zh"},
	}); err != nil {
		t.Fatalf("AddVectorsUpstream() error = %v", err)
	}

	sch := newVectorsTestScheduler(cfgManager, metrics.NewMetricsManager())
	w := serveVectorsEmbeddingRequest(cfgManager, sch, `{"model":"text-embedding-3-small","input":"hello"}`)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	if upstreamModel != "text-embedding-3-small" {
		t.Fatalf("upstream model = %q, want text-embedding-3-small", upstreamModel)
	}
}

func TestHandlerVectors422DoesNotFailoverOrAffectBreaker(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfgManager := newVectorsTestConfigManager(t)
	defer cfgManager.Close()

	var primaryAttempts int32
	primaryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&primaryAttempts, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"error":{"message":"Validation error: bad embedding model","type":"invalid_request_error","code":"invalid_request"}}`))
	}))
	defer primaryServer.Close()

	var secondaryAttempts int32
	secondaryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&secondaryAttempts, 1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"object":"list","data":[],"usage":{"prompt_tokens":1,"total_tokens":1}}`))
	}))
	defer secondaryServer.Close()

	if err := cfgManager.AddVectorsUpstream(config.UpstreamConfig{
		Name:        "secondary-vectors",
		ServiceType: "openai",
		BaseURL:     secondaryServer.URL,
		APIKeys:     []string{"sk-secondary"},
		Priority:    2,
	}); err != nil {
		t.Fatalf("AddVectorsUpstream(secondary) error = %v", err)
	}
	if err := cfgManager.AddVectorsUpstream(config.UpstreamConfig{
		Name:        "primary-vectors",
		ServiceType: "openai",
		BaseURL:     primaryServer.URL,
		APIKeys:     []string{"sk-primary"},
		Priority:    1,
	}); err != nil {
		t.Fatalf("AddVectorsUpstream(primary) error = %v", err)
	}

	vectorsMetrics := metrics.NewMetricsManager()
	sch := newVectorsTestScheduler(cfgManager, vectorsMetrics)
	w := serveVectorsEmbeddingRequest(cfgManager, sch, `{"model":"text-embedding-3-small","input":"hello"}`)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d, body=%s", w.Code, w.Body.String())
	}
	if got := atomic.LoadInt32(&primaryAttempts); got != 1 {
		t.Fatalf("primary attempts = %d, want 1", got)
	}
	if got := atomic.LoadInt32(&secondaryAttempts); got != 0 {
		t.Fatalf("secondary attempts = %d, want 0", got)
	}

	keyMetrics := vectorsMetrics.GetKeyMetrics(primaryServer.URL, "sk-primary", "openai")
	if keyMetrics == nil {
		t.Fatal("expected primary key metrics")
	}
	if keyMetrics.RequestCount != 1 || keyMetrics.FailureCount != 1 {
		t.Fatalf("metrics counts = requests:%d failures:%d, want 1/1", keyMetrics.RequestCount, keyMetrics.FailureCount)
	}
	if keyMetrics.ConsecutiveFailures != 0 {
		t.Fatalf("consecutive breaker failures = %d, want 0", keyMetrics.ConsecutiveFailures)
	}
	if got := vectorsMetrics.CalculateKeyFailureRate(primaryServer.URL, "sk-primary", "openai"); got != 0 {
		t.Fatalf("breaker failure rate = %v, want 0", got)
	}
	if state := vectorsMetrics.GetKeyCircuitState(primaryServer.URL, "sk-primary", "openai"); state != metrics.CircuitStateClosed {
		t.Fatalf("circuit state = %s, want closed", state)
	}
}

func TestAddUpstreamRejectsUnsupportedServiceType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfgManager := newVectorsTestConfigManager(t)
	defer cfgManager.Close()

	r := gin.New()
	r.POST("/api/vectors/channels", AddUpstream(cfgManager))

	req := httptest.NewRequest(http.MethodPost, "/api/vectors/channels", strings.NewReader(`{"name":"bad","serviceType":"gemini","baseUrl":"https://example.com","apiKeys":["sk-test"]}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "Vectors") {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}
