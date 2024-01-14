package main

import (
	"context"
	"log"
	"os"

	_ "embed"

	"github.com/rhettg/openai-go"
	"github.com/rhettg/openai-go/chat"
)

//go:embed image.png
var imgData []byte

func main() {
	ctx := context.Background()
	s := openai.NewSession(os.Getenv("OPENAI_API_KEY"))

	client := chat.NewClient(s, "gpt-4-vision-preview")

	ic, err := chat.NewContentFromImage("image/png", imgData)
	if err != nil {
		log.Fatalf("Failed to create image content: %v", err)
	}
	tc := chat.NewContentFromText("please describe the image")

	resp, err := client.CreateMMCompletion(ctx, &chat.CreateMMCompletionParams{
		MaxTokens: 255,
		Messages: []*chat.MMMessage{
			{Role: "user", Content: []chat.Content{ic, tc}},
		},
	})
	if err != nil {
		log.Fatalf("Failed to complete: %v", err)
	}

	for _, choice := range resp.Choices {
		msg := choice.Message
		log.Printf("role=%q, content=%q", msg.Role, msg.Content)
	}
}
