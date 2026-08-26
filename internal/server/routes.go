package server

import "net/http"

func New(staticDir, examplePath string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/health", handleHealth)
	mux.HandleFunc("/api/example", makeExampleHandler(examplePath))
	mux.HandleFunc("/api/operate", handleOperate)
	mux.HandleFunc("/api/atmosphere", handleAtmosphere)
	mux.HandleFunc("/api/network", handleNetwork)
	mux.Handle("/", fileServerHandler(staticDir))
	return mux
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "GET required")
		return
	}
	writeJSON(w, http.StatusOK, HealthResponse{OK: true, Service: "fan-duct"})
}

func makeExampleHandler(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "GET required")
			return
		}
		data, err := readExample(path)
		if err != nil {
			writeError(w, http.StatusNotFound, "example file unavailable: "+err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}
}
