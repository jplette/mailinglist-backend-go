package health

import (
	"fmt"
	"net/http"
)

// Ping godoc
// @Summary      Health check
// @Description  Returns OK if the service is running.
// @Tags         health
// @Produce      plain
// @Success      200  {string}  string  "OK"
// @Router       /health [get]
func Ping(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "OK")
	if err != nil {
		return
	}
}
