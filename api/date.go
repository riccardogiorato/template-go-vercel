package handler

import (
	"fmt"
	"net/http"
	"time"
)

func Date(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now().Format(time.RFC850)
	fmt.Fprint(w, currentTime)
}
