package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func MyInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := make(map[string]string)
	resp["ip"] = r.RemoteAddr
	resp["user-agent"] = r.UserAgent()
	resp["accept-language"] = r.Header.Get("Accept-Language")
	resp["github"] = "https://github.com/riccardogiorato/template-go-vercel/blob/main/api/myinfo.go"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("Error happened in JSON marshal. Err: %s", err)
	} else {
		w.Write(jsonResp)
	}
}
