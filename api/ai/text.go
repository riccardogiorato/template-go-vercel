package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/togethercomputer/together-go"
	"github.com/togethercomputer/together-go/option"
)

func GenerateText(w http.ResponseWriter, r *http.Request) {
	client := together.NewClient(option.WithAPIKey(os.Getenv("TOGETHER_API_KEY")))
	resp, err := client.Completions.New(context.Background(), together.CompletionNewParams{
		Model:     together.CompletionNewParamsModel(together.ChatCompletionNewParamsModelQwenQwen2_5_7BInstructTurbo),
		Prompt:    "Write a short story about a robot.",
		MaxTokens: together.Int(100),
	})
	if err != nil || len(resp.Choices) == 0 {
		http.Error(w, "Error generating text", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"text": resp.Choices[0].Text})
}
