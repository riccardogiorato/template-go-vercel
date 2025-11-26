package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/togethercomputer/together-go"
	"github.com/togethercomputer/together-go/option"
)

func TranscribeAudio(w http.ResponseWriter, r *http.Request) {
	togetherClient := together.NewClient(option.WithAPIKey(os.Getenv("TOGETHER_API_KEY")))

	// Download audio file from URL
	audioURL := "https://upload.wikimedia.org/wikipedia/commons/7/75/Marshall_Plan_Speech.wav"
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

	// Use the response body directly for transcription
	transcriptionResp, err := togetherClient.Audio.Transcriptions.New(context.Background(), together.AudioTranscriptionNewParams{
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

func TranslateAudio(w http.ResponseWriter, r *http.Request) {
	togetherClient := together.NewClient(option.WithAPIKey(os.Getenv("TOGETHER_API_KEY")))

	// Download audio file from URL
	audioURL := "https://upload.wikimedia.org/wikipedia/commons/7/75/Marshall_Plan_Speech.wav"
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

	// Use the response body directly for translation
	translationResp, err := togetherClient.Audio.Translations.New(context.Background(), together.AudioTranslationNewParams{
		File:  resp.Body,
		Model: "openai/whisper-large-v3",
	})
	if err != nil {
		http.Error(w, "Error translating", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Translation: %s", translationResp.Text)
}
