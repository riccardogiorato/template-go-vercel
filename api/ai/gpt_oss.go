package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/togethercomputer/together-go"
	"github.com/togethercomputer/together-go/option"
)

func GenerateGPTOss(w http.ResponseWriter, r *http.Request) {
	client := together.NewClient(option.WithAPIKey(os.Getenv("TOGETHER_API_KEY")))
	resp, err := client.Completions.New(context.Background(), together.CompletionNewParams{
		Model:       "openai/gpt-oss-120b",
		Prompt:      "Solve this logic puzzle: If all roses are flowers and some flowers are red, can we conclude that some roses are red?",
		MaxTokens:   together.Int(200),
		Temperature: together.Float(1.0),
	})
	if err != nil || len(resp.Choices) == 0 {
		http.Error(w, "Error generating response", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, resp.Choices[0].Text)
}
