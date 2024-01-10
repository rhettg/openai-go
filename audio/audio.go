// Package audio implements a client for OpenAI's Whisper
// audio transcriber.
package audio

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/rakyll/openai-go"
)

const (
	defaultCreateTranscriptionEndpoint = "https://api.openai.com/v1/audio/transcriptions"
	defaultCreateSpeechEndpoint        = "https://api.openai.com/v1/audio/speech"
)

// Client is a client to communicate with Open AI's ChatGPT APIs.
type Client struct {
	s     *openai.Session
	model string

	// CreateTranscriptionEndpoint allows overriding the default API endpoint.
	// Set this field before using the client.
	CreateTranscriptionEndpoint string
	CreateSpeechEndpoint        string
}

// NewClient creates a new default client that uses the given session
// and defaults to the given model.
func NewClient(session *openai.Session, model string) *Client {
	if model == "" {
		model = "whisper-1"
	}
	return &Client{
		s:                           session,
		model:                       model,
		CreateTranscriptionEndpoint: defaultCreateTranscriptionEndpoint,
		CreateSpeechEndpoint:        defaultCreateSpeechEndpoint,
	}
}

type CreateTranscriptionParams struct {
	Model    string
	Language string

	Audio       io.Reader
	AudioFormat string // such as "mp3" or "wav", etc.

	Prompt string // optional
	// TODO: Add temperature.
}

type CreateTranscriptionResponse struct {
	Text string `json:"text,omitempty"`
}

func (c *Client) CreateTranscription(ctx context.Context, p *CreateTranscriptionParams) (*CreateTranscriptionResponse, error) {
	if p.AudioFormat == "" {
		return nil, fmt.Errorf("audio format is required")
	}
	if p.Model == "" {
		p.Model = c.model
	}
	params := url.Values{}
	params.Set("model", p.Model)
	if p.Language != "" {
		params.Set("language", p.Language)
	}
	if p.Prompt != "" {
		params.Set("prompt", p.Prompt)
	}
	var r CreateTranscriptionResponse
	return &r, c.s.Upload(ctx, c.CreateTranscriptionEndpoint, p.Audio, p.AudioFormat, params, &r)
}

type CreateSpeechParams struct {
	Model          string `json:"model"`
	Input          string `json:"input"`
	Voice          string `json:"voice,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	Speed          int    `json:"speed,omitempty"`
}

func (c *Client) CreateSpeech(ctx context.Context, p *CreateSpeechParams, w io.Writer) error {
	if p.Voice == "" {
		return fmt.Errorf("voice is required")
	}
	if p.Model == "" {
		p.Model = c.model
	}

	err := c.s.Download(ctx, c.CreateSpeechEndpoint, p, w)
	if err != nil {
		return fmt.Errorf("failed to download speech: %v", err)
	}

	return nil
}
