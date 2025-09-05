package health

import (
	"fmt"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "OK")
	if err != nil {
		return
	}
}
