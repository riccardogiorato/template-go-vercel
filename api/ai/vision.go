package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/togethercomputer/together-go"
	"github.com/togethercomputer/together-go/option"
)

func DescribeImage(w http.ResponseWriter, r *http.Request) {
	client := together.NewClient(option.WithAPIKey(os.Getenv("TOGETHER_API_KEY")))

	prompt := "You are a UX/UI designer. Describe this Trello board screenshot in detail. Pay close attention to background color, text color, font size, font family, padding, margin, border, etc. Use the exact text from the screenshot."

	resp, err := client.Completions.New(context.Background(), together.CompletionNewParams{
		Model:       together.CompletionNewParamsModel(together.ChatCompletionNewParamsModelQwenQwen2_5_7BInstructTurbo),
		Prompt:      prompt,
		MaxTokens:   together.Int(500),
		Temperature: together.Float(0.2),
	})
	if err != nil || len(resp.Choices) == 0 {
		http.Error(w, "Error describing image", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, resp.Choices[0].Text)
}
