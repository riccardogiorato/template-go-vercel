package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/togethercomputer/together-go"
	"github.com/togethercomputer/together-go/option"
)

func GenerateImage(w http.ResponseWriter, r *http.Request) {
	client := together.NewClient(option.WithAPIKey(os.Getenv("TOGETHER_API_KEY")))
	resp, err := client.Images.New(context.Background(), together.ImageNewParams{
		Prompt: "A futuristic robot in a cityscape",
		Model:  together.ImageNewParamsModelBlackForestLabsFlux1SchnellFree,
		Width:  together.Int(1024),
		Height: together.Int(1024),
		Steps:  together.Int(4),
		N:      together.Int(1),
	})
	if err != nil || len(resp.Data) == 0 {
		http.Error(w, "Error generating image", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"image_url": resp.Data[0].URL})
}
