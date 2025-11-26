package handler

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/togethercomputer/together-go"
	"github.com/togethercomputer/together-go/option"
)

func GenerateSpeech(w http.ResponseWriter, r *http.Request) {
	client := together.NewClient(option.WithAPIKey(os.Getenv("TOGETHER_API_KEY")))

	resp, err := client.Audio.New(context.Background(), together.AudioNewParams{
		Model:          together.AudioNewParamsModelCartesiaSonic,
		Input:          "Today is a wonderful day to build something people love!",
		Voice:          together.AudioNewParamsVoiceFriendlySidekick,
		ResponseFormat: together.AudioNewParamsResponseFormatMP3,
	})
	if err != nil {
		http.Error(w, "Error generating speech", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	w.Header().Set("Content-Type", "audio/mpeg")
	io.Copy(w, resp.Body)
}
