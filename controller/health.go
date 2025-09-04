package controller

import (
	"fmt"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "OK")
	if err != nil {
		return
	}
}
