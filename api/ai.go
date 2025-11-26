package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/togethercomputer/together-go"
	"github.com/togethercomputer/together-go/option"
)

type AIRequest struct {
	Type     string `json:"type"`
	Prompt   string `json:"prompt,omitempty"`
	Input    string `json:"input,omitempty"`
	AudioURL string `json:"audio_url,omitempty"`
	Code     string `json:"code,omitempty"`
	Language string `json:"language,omitempty"`
	Action   string `json:"action,omitempty"` // for transcription: transcribe or translate
}

func Handler(w http.ResponseWriter, r *http.Request) {
	var req AIRequest

	if r.Method == "GET" {
		// Parse from query params
		req.Type = r.URL.Query().Get("type")
		req.Prompt = r.URL.Query().Get("prompt")
		req.Input = r.URL.Query().Get("input")
		req.AudioURL = r.URL.Query().Get("audio_url")
		req.Code = r.URL.Query().Get("code")
		req.Language = r.URL.Query().Get("language")
		req.Action = r.URL.Query().Get("action")
	} else {
		// Parse from JSON body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
	}

	client := together.NewClient(option.WithAPIKey(os.Getenv("TOGETHER_API_KEY")))

	switch req.Type {
	case "chat":
		handleChat(w, client, req.Prompt)
	case "text":
		handleText(w, client, req.Prompt)
	case "image":
		handleImage(w, client, req.Prompt)
	case "vision":
		handleVision(w, client, req.Prompt)
	case "speech":
		handleSpeech(w, client, req.Input)
	case "transcription":
		handleTranscription(w, client, req.AudioURL, req.Action)
	case "code_interpreter":
		handleCodeInterpreter(w, client, req.Code, req.Language)
	case "json":
		handleJSON(w, client, req.Prompt)
	default:
		http.Error(w, "Unknown type", http.StatusBadRequest)
	}
}

func handleChat(w http.ResponseWriter, client together.Client, prompt string) {
	if prompt == "" {
		prompt = "User: debate the pros and cons of AI\nAssistant:"
	}
	resp, err := client.Completions.New(context.Background(), together.CompletionNewParams{
		Model:     together.CompletionNewParamsModel(together.ChatCompletionNewParamsModelQwenQwen2_5_7BInstructTurbo),
		Prompt:    prompt,
		MaxTokens: together.Int(100),
	})
	if err != nil || len(resp.Choices) == 0 {
		http.Error(w, "Error generating chat", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"response": resp.Choices[0].Text})
}

func handleText(w http.ResponseWriter, client together.Client, prompt string) {
	if prompt == "" {
		prompt = "Write a short story about a robot."
	}
	resp, err := client.Completions.New(context.Background(), together.CompletionNewParams{
		Model:     together.CompletionNewParamsModel(together.ChatCompletionNewParamsModelQwenQwen2_5_7BInstructTurbo),
		Prompt:    prompt,
		MaxTokens: together.Int(100),
	})
	if err != nil || len(resp.Choices) == 0 {
		http.Error(w, "Error generating text", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"text": resp.Choices[0].Text})
}

func handleImage(w http.ResponseWriter, client together.Client, prompt string) {
	if prompt == "" {
		prompt = "A futuristic robot in a cityscape"
	}
	resp, err := client.Images.New(context.Background(), together.ImageNewParams{
		Prompt: prompt,
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

func handleVision(w http.ResponseWriter, client together.Client, prompt string) {
	if prompt == "" {
		prompt = "You are a UX/UI designer. Describe this Trello board screenshot in detail. Pay close attention to background color, text color, font size, font family, padding, margin, border, etc. Use the exact text from the screenshot."
	}
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

func handleSpeech(w http.ResponseWriter, client together.Client, input string) {
	if input == "" {
		input = "Today is a wonderful day to build something people love!"
	}
	resp, err := client.Audio.New(context.Background(), together.AudioNewParams{
		Model:          "cartesia/sonic-2",
		Input:          input,
		Voice:          "friendly sidekick",
		ResponseFormat: "mp3",
	})
	if err != nil {
		http.Error(w, "Error generating speech", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	w.Header().Set("Content-Type", "audio/mpeg")
	io.Copy(w, resp.Body)
}

func handleTranscription(w http.ResponseWriter, client together.Client, audioURL, action string) {
	if audioURL == "" {
		audioURL = "https://upload.wikimedia.org/wikipedia/commons/7/75/Marshall_Plan_Speech.wav"
	}
	req, err := http.NewRequest("GET", audioURL, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating request: %v", err), http.StatusInternalServerError)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Go HTTP Client)")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error downloading audio file: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Failed to download audio file: HTTP %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}

	if action == "translate" {
		translationResp, err := client.Audio.Translations.New(context.Background(), together.AudioTranslationNewParams{
			File:  resp.Body,
			Model: "openai/whisper-large-v3",
		})
		if err != nil {
			http.Error(w, "Error translating", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Translation: %s", translationResp.Text)
	} else {
		transcriptionResp, err := client.Audio.Transcriptions.New(context.Background(), together.AudioTranscriptionNewParams{
			File:     resp.Body,
			Model:    "openai/whisper-large-v3",
			Language: together.String("en"),
		})
		if err != nil {
			http.Error(w, "Error transcribing", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Transcription: %s", transcriptionResp.Text)
	}
}

func handleCodeInterpreter(w http.ResponseWriter, client together.Client, code, language string) {
	if code == "" {
		code = "print(\"Welcome to Together Code Interpreter!\")"
	}
	if language == "" {
		language = "python"
	}
	var lang together.CodeInterpreterExecuteParamsLanguage
	if language == "python" {
		lang = together.CodeInterpreterExecuteParamsLanguagePython
	} else {
		http.Error(w, "Unsupported language", http.StatusBadRequest)
		return
	}
	resp, err := client.CodeInterpreter.Execute(context.Background(), together.CodeInterpreterExecuteParams{
		Code:     code,
		Language: lang,
	})
	if err != nil {
		http.Error(w, "Error running code interpreter", http.StatusInternalServerError)
		return
	}

	successResp := resp.AsSuccessfulExecution()
	json.NewEncoder(w).Encode(successResp.Data)
}

func handleJSON(w http.ResponseWriter, client together.Client, prompt string) {
	if prompt == "" {
		prompt = "Return only a JSON object representing a person with name, age, and city. No additional text."
	}
	resp, err := client.Completions.New(context.Background(), together.CompletionNewParams{
		Model:     together.CompletionNewParamsModel(together.ChatCompletionNewParamsModelQwenQwen2_5_7BInstructTurbo),
		Prompt:    prompt,
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
