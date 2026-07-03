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

const defaultAIModel = "Qwen/Qwen2.5-7B-Instruct-Turbo"

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
	content, err := generateChatContent(client, prompt, 100)
	if err != nil {
		http.Error(w, "Error generating chat", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"response": content})
}

func generateChatContent(client together.Client, prompt string, maxTokens int64) (string, error) {
	resp, err := client.Chat.Completions.New(context.Background(), together.ChatCompletionNewParams{
		Model: defaultAIModel,
		Messages: []together.ChatCompletionNewParamsMessageUnion{
			{
				OfChatCompletionNewsMessageChatCompletionUserMessageParam: &together.ChatCompletionNewParamsMessageChatCompletionUserMessageParam{
					Role: "user",
					Content: together.ChatCompletionNewParamsMessageChatCompletionUserMessageParamContentUnion{
						OfString: together.String(prompt),
					},
				},
			},
		},
		MaxTokens: together.Int(maxTokens),
	})
	if err != nil {
		return "", err
	}
	if resp == nil || len(resp.Choices) == 0 {
		return "", fmt.Errorf("empty chat response")
	}
	choice := resp.Choices[0]
	if choice.Message.Content != "" {
		return choice.Message.Content, nil
	}
	return choice.Text, nil
}

func handleText(w http.ResponseWriter, client together.Client, prompt string) {
	if prompt == "" {
		prompt = "Write a short story about a robot."
	}
	content, err := generateChatContent(client, prompt, 200)
	if err != nil {
		http.Error(w, "Error generating text", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"text": content})
}

func handleImage(w http.ResponseWriter, client together.Client, prompt string) {
	if prompt == "" {
		prompt = "A futuristic robot in a cityscape"
	}
	resp, err := client.Images.Generate(context.Background(), together.ImageGenerateParams{
		Prompt:         prompt,
		Model:          "black-forest-labs/FLUX.1-schnell",
		ResponseFormat: together.ImageGenerateParamsResponseFormatURL,
		N:              together.Int(1),
	})
	if err != nil || len(resp.Data) == 0 {
		http.Error(w, "Error generating image", http.StatusInternalServerError)
		return
	}
	// Response data is a union type; extract the URL variant.
	imageData := resp.Data[0].AsURL()
	json.NewEncoder(w).Encode(map[string]string{"image_url": imageData.URL})
}

func handleVision(w http.ResponseWriter, client together.Client, prompt string) {
	if prompt == "" {
		prompt = "You are a UX/UI designer. Describe this Trello board screenshot in detail. Pay close attention to background color, text color, font size, font family, padding, margin, border, etc. Use the exact text from the screenshot."
	}
	content, err := generateChatContent(client, prompt, 500)
	if err != nil {
		http.Error(w, "Error describing image", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, content)
}

func handleSpeech(w http.ResponseWriter, client together.Client, input string) {
	if input == "" {
		input = "Today is a wonderful day to build something people love!"
	}
	resp, err := client.Audio.Speech.New(context.Background(), together.AudioSpeechNewParams{
		Model:          together.AudioSpeechNewParamsModelCartesiaSonic,
		Input:          input,
		Voice:          "friendly sidekick",
		ResponseFormat: together.AudioSpeechNewParamsResponseFormatMP3,
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
			File: together.AudioTranslationNewParamsFileUnion{
				OfFile: resp.Body,
			},
			Model:          together.AudioTranslationNewParamsModelOpenAIWhisperLargeV3,
			ResponseFormat: together.AudioTranslationNewParamsResponseFormatJson,
		})
		if err != nil {
			http.Error(w, "Error translating", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Translation: %s", translationResp.Text)
	} else {
		transcriptionResp, err := client.Audio.Transcriptions.New(context.Background(), together.AudioTranscriptionNewParams{
			File: together.AudioTranscriptionNewParamsFileUnion{
				OfFile: resp.Body,
			},
			Model:          together.AudioTranscriptionNewParamsModelOpenAIWhisperLargeV3,
			ResponseFormat: together.AudioTranscriptionNewParamsResponseFormatJson,
			Language:       together.String("en"),
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
	content, err := generateChatContent(client, prompt, 100)
	if err != nil {
		http.Error(w, "Error generating JSON", http.StatusInternalServerError)
		return
	}
	var data interface{}
	if json.Unmarshal([]byte(content), &data) == nil {
		json.NewEncoder(w).Encode(data)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"json": content})
	}
}
