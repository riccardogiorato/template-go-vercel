package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/togethercomputer/together-go"
	"github.com/togethercomputer/together-go/option"
)

func GenerateChat(w http.ResponseWriter, r *http.Request) {
	client := together.NewClient(option.WithAPIKey(os.Getenv("TOGETHER_API_KEY")))
	resp, err := client.Completions.New(context.Background(), together.CompletionNewParams{
		Model:     together.CompletionNewParamsModel(together.ChatCompletionNewParamsModelQwenQwen2_5_7BInstructTurbo),
		Prompt:    "User: debate the pros and cons of AI\nAssistant:",
		MaxTokens: together.Int(100),
	})
	if err != nil || len(resp.Choices) == 0 {
		http.Error(w, "Error generating chat", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"response": resp.Choices[0].Text})
}
