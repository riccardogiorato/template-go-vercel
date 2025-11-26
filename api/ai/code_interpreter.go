package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/togethercomputer/together-go"
	"github.com/togethercomputer/together-go/option"
)

func RunCodeInterpreter(w http.ResponseWriter, r *http.Request) {
	client := together.NewClient(option.WithAPIKey(os.Getenv("TOGETHER_API_KEY")))
	resp, err := client.CodeInterpreter.Execute(context.Background(), together.CodeInterpreterExecuteParams{
		Code:     "print(\"Welcome to Together Code Interpreter!\")",
		Language: together.CodeInterpreterExecuteParamsLanguagePython,
	})
	if err != nil {
		http.Error(w, "Error running code interpreter", http.StatusInternalServerError)
		return
	}

	// Handle the union response
	successResp := resp.AsSuccessfulExecution()
	json.NewEncoder(w).Encode(successResp.Data)
}
