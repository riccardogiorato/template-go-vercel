package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/togethercomputer/together-go"
	"github.com/togethercomputer/together-go/option"
)

func GenerateJSON(w http.ResponseWriter, r *http.Request) {
	client := together.NewClient(option.WithAPIKey(os.Getenv("TOGETHER_API_KEY")))
	resp, err := client.Completions.New(context.Background(), together.CompletionNewParams{
		Model:     together.CompletionNewParamsModel(together.ChatCompletionNewParamsModelQwenQwen2_5_7BInstructTurbo),
		Prompt:    "Return only a JSON object representing a person with name, age, and city. No additional text.",
		MaxTokens: together.Int(100),
	})
	if err != nil || len(resp.Choices) == 0 {
		http.Error(w, "Error generating JSON", http.StatusInternalServerError)
		return
	}
	var data interface{}
	if json.Unmarshal([]byte(resp.Choices[0].Text), &data) == nil {
		json.NewEncoder(w).Encode(data)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"json": resp.Choices[0].Text})
	}
}
