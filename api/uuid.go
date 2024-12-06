package handler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func TestUUID(w http.ResponseWriter, r *http.Request) {
	// Creating UUID Version 4
	id := uuid.New()
	// Print the UUID to the console
	fmt.Println(id.String())
	// Print the UUID to the response
	fmt.Fprint(w, id.String())
}
