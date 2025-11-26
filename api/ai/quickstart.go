package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/togethercomputer/together-go"
	"github.com/togethercomputer/together-go/option"
)

func GenerateQuickstart(w http.ResponseWriter, r *http.Request) {
	client := together.NewClient(option.WithAPIKey(os.Getenv("TOGETHER_API_KEY")))
	resp, err := client.Completions.New(context.Background(), together.CompletionNewParams{
		Model:     together.CompletionNewParamsModel(together.ChatCompletionNewParamsModelQwenQwen2_5_7BInstructTurbo),
		Prompt:    "What are the top 3 things to do in New York?",
		MaxTokens: together.Int(200),
	})
	if err != nil || len(resp.Choices) == 0 {
		http.Error(w, "Error generating quickstart", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, resp.Choices[0].Text)
}
