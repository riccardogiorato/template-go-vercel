package handler

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

var client = redis.NewClient(&redis.Options{
	Addr:     "global-many-spaniel-31036.upstash.io:31036",
	Password: os.Getenv("UPSTASH_PASSWORD"),
	DB:       0,
	TLSConfig: &tls.Config{
		MinVersion: tls.VersionTLS12,
	},
})

func Redis(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	// set a foo key on upstash with this value
	client.Set(ctx, "foo", "value-from-upstash-redis", 0)

	// get the foo key from upstash
	foo := client.Get(ctx, "foo")

	resp := make(map[string]string)
	resp["foo"] = foo.Val()
	resp["github"] = "https://github.com/riccardogiorato/template-go-vercel/blob/main/api/redis.go"
	body, err := json.Marshal(resp)

	if err != nil {
		fmt.Printf("Error happened in JSON marshal. Err: %s", err)
	} else {
		w.Write(body)
	}
}
