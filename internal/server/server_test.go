package server_test

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"fan-duct/internal/server"
	"fan-duct/internal/solve"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	src := filepath.Join("..", "..", "example", "inline-fan.json")
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read example: %v", err)
	}
	dir := t.TempDir()
	ex := filepath.Join(dir, "inline-fan.json")
	if err := os.WriteFile(ex, data, 0o644); err != nil {
		t.Fatalf("write example: %v", err)
	}
	srv := httptest.NewServer(server.New(dir, ex))
	t.Cleanup(srv.Close)
	return srv
}

func TestHealthEndpoint(t *testing.T) {
	srv := newTestServer(t)
	resp, err := http.Get(srv.URL + "/health")
	if err != nil {
		t.Fatalf("GET /health: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	var out server.HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !out.OK || out.Service != "fan-duct" {
		t.Errorf("health = %+v", out)
	}
}

func TestOperateEndpointMatchesSolve(t *testing.T) {
	srv := newTestServer(t)
	raw, err := os.ReadFile(filepath.Join("..", "..", "example", "inline-fan.json"))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	in, err := solve.ParseInput(raw)
	if err != nil {
		t.Fatalf("ParseInput: %v", err)
	}
	b, err := in.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	op, err := b.OperatingPoint()
	if err != nil {
		t.Fatalf("OperatingPoint: %v", err)
	}
	resp, err := http.Post(srv.URL+"/api/operate", "application/json", strings.NewReader(string(raw)))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	var out server.OperateResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if math.Abs(out.Flow-op.Flow) > 1e-9 {
		t.Errorf("flow api=%v solve=%v", out.Flow, op.Flow)
	}
	if math.Abs(out.Pressure-op.Pressure) > 1e-6 {
		t.Errorf("pressure api=%v solve=%v", out.Pressure, op.Pressure)
	}
	if math.Abs(out.Residual) > 1e-6 {
		t.Errorf("residual %v", out.Residual)
	}
	if out.RespeedQ <= out.Flow {
		t.Errorf("respeed Q %v should exceed base %v", out.RespeedQ, out.Flow)
	}
}

func TestOperateEndpointRejectsBadInput(t *testing.T) {
	srv := newTestServer(t)
	cases := []string{
		`{not json`,
		`{"duct":{"length":50,"diameter":0},"fan":{"points":[{"q":0,"dp":1},{"q":0.3,"dp":0.5}]}}`,
		`{"duct":{"length":50,"diameter":0.15},"fan":{"points":[{"q":0,"dp":0},{"q":0.3,"dp":0}]}}`,
	}
	for _, body := range cases {
		resp, err := http.Post(srv.URL+"/api/operate", "application/json", strings.NewReader(body))
		if err != nil {
			t.Fatalf("POST: %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			resp.Body.Close()
			t.Errorf("status = %d for %s", resp.StatusCode, body)
			continue
		}
		var er server.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&er); err != nil {
			resp.Body.Close()
			t.Fatalf("decode: %v", err)
		}
		resp.Body.Close()
		if strings.TrimSpace(er.Error) == "" {
			t.Errorf("empty error for %s", body)
		}
	}
}

func TestAtmosphereEndpointSeaLevel(t *testing.T) {
	srv := newTestServer(t)
	body := `{"altitude_m":0,"rel_humidity":0.5}`
	resp, err := http.Post(srv.URL+"/api/atmosphere", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	var out server.AtmosphereResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if math.Abs(out.PressurePa-101325) > 1e-4 {
		t.Errorf("P = %v, want 101325", out.PressurePa)
	}
	if out.MoistDensity >= out.DensityKgM3 {
		t.Errorf("moist %v should be below dry %v", out.MoistDensity, out.DensityKgM3)
	}
}

func TestExampleEndpoint(t *testing.T) {
	srv := newTestServer(t)
	resp, err := http.Get(srv.URL + "/api/example")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	var in solve.Input
	if err := json.NewDecoder(resp.Body).Decode(&in); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if in.Duct.Diameter != 0.15 {
		t.Errorf("diameter = %v", in.Duct.Diameter)
	}
}

func TestNetworkSeriesDuctEndpoint(t *testing.T) {
	srv := newTestServer(t)
	raw, err := os.ReadFile(filepath.Join("..", "..", "example", "inline-fan.json"))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	payload := `{"mode":"series-duct","left":` + string(raw) + `,"right":` + string(raw) + `}`
	resp, err := http.Post(srv.URL+"/api/network", "application/json", strings.NewReader(payload))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	var out server.NetworkResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	in, _ := solve.ParseInput(raw)
	b, _ := in.Build()
	op, _ := b.OperatingPoint()
	if !(out.Flow < op.Flow) {
		t.Errorf("series-duct Q %v should be below single-duct %v", out.Flow, op.Flow)
	}
	if math.Abs(out.Residual) > 1e-5 {
		t.Errorf("residual %v", out.Residual)
	}
}
