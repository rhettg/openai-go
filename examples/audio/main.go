package main

import (
	"context"
	"log"
	"os"

	"github.com/rhettg/openai-go"
	"github.com/rhettg/openai-go/audio"
)

func main() {
	ctx := context.Background()

	s := openai.NewSession(os.Getenv("OPENAI_API_KEY"))
	client := audio.NewClient(s, "")
	filePath := os.Getenv("AUDIO_FILE_PATH")
	if filePath == "" {
		log.Fatal("must provide an AUDIO_FILE_PATH env var")
	}
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error opening audio file: %v", err)
	}
	defer f.Close()
	resp, err := client.CreateTranscription(ctx, &audio.CreateTranscriptionParams{
		Language:    "en",
		Audio:       f,
		AudioFormat: "mp3",
	})
	if err != nil {
		log.Fatalf("error transcribing file: %v", err)
	}

	log.Println(resp.Text)

	of, err := os.OpenFile("speech.aac", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error opening speech file: %v", err)
	}

	err = client.CreateSpeech(ctx, &audio.CreateSpeechParams{
		Model:          "tts-1",
		Voice:          "nova",
		ResponseFormat: "aac",
		Input:          resp.Text,
	}, of)
	if err != nil {
		log.Fatalf("error creating speech: %v", err)
	}

	err = of.Close()
	if err != nil {
		log.Fatalf("error saving output file: %v", err)
	}

	log.Println("saved speech file to speech.aac")
}
