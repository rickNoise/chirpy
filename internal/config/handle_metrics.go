package config

import (
	"fmt"
	"net/http"
)

func (cfg *ApiConfig) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	template := `
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>	
`
	hits := fmt.Sprintf(template, cfg.fileserverHits.Load())
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(hits))
}
