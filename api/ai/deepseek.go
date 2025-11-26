package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/togethercomputer/together-go"
	"github.com/togethercomputer/together-go/option"
)

func GenerateDeepSeekR1(w http.ResponseWriter, r *http.Request) {
	client := together.NewClient(option.WithAPIKey(os.Getenv("TOGETHER_API_KEY")))
	resp, err := client.Completions.New(context.Background(), together.CompletionNewParams{
		Model:     "deepseek-ai/deepseek-r1",
		Prompt:    "Which number is bigger: 9.9 or 9.11? Explain your reasoning step by step.",
		MaxTokens: together.Int(100),
	})
	if err != nil || len(resp.Choices) == 0 {
		http.Error(w, "Error generating response", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, resp.Choices[0].Text)
}
