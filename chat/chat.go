// Package chat contains a client for Open AI's ChatGPT APIs.
package chat

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/rakyll/openai-go"
)

const defaultModel = "gpt-3.5-turbo"

const defaultCreateCompletionsEndpoint = "https://api.openai.com/v1/chat/completions"

// Client is a client to communicate with Open AI's ChatGPT APIs.
type Client struct {
	s     *openai.Session
	model string

	// CreateCompletionsEndpoint allows overriding the default API endpoint.
	// Set this field before using the client.
	CreateCompletionEndpoint string
}

// NewClient creates a new default client that uses the given session
// and defaults to the given model.
func NewClient(session *openai.Session, model string) *Client {
	if model == "" {
		model = defaultModel
	}
	return &Client{
		s:                        session,
		model:                    model,
		CreateCompletionEndpoint: defaultCreateCompletionsEndpoint,
	}
}

type CreateCompletionParams struct {
	Model string `json:"model,omitempty"`

	Messages []*Message `json:"messages,omitempty"`
	Stop     []string   `json:"stop,omitempty"`
	Stream   bool       `json:"stream,omitempty"`

	Functions    []Function `json:"functions,omitempty"`
	FunctionCall string     `json:"function_call,omitempty"`

	N           int     `json:"n,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	MaxTokens   int     `json:"max_tokens,omitempty"`

	PresencePenalty  float64 `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`

	User string `json:"user,omitempty"`
}

type CreateMMCompletionParams struct {
	Model string `json:"model,omitempty"`

	Messages []*MMMessage `json:"messages,omitempty"`
	Stop     []string     `json:"stop,omitempty"`
	Stream   bool         `json:"stream,omitempty"`

	N           int     `json:"n,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	MaxTokens   int     `json:"max_tokens,omitempty"`

	PresencePenalty  float64 `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`

	User string `json:"user,omitempty"`
}

type CreateCompletionResponse struct {
	ID        string    `json:"id,omitempty"`
	Object    string    `json:"object,omitempty"`
	CreatedAt int64     `json:"created_at,omitempty"`
	Choices   []*Choice `json:"choices,omitempty"`

	Usage *openai.Usage `json:"usage,omitempty"`
}

type Choice struct {
	Message      *Message `json:"message,omitempty"`
	Index        int      `json:"index,omitempty"`
	LogProbs     int      `json:"logprobs,omitempty"`
	FinishReason string   `json:"finish_reason,omitempty"`
}

type Message struct {
	Role         string        `json:"role,omitempty"`
	FunctionCall *FunctionCall `json:"function_call,omitempty"`
	Content      string        `json:"content,omitempty"`
	Name         string        `json:"name,omitempty"`
}

type MMMessage struct {
	Role    string    `json:"role,omitempty"`
	Content []Content `json:"content,omitempty"`
}

type ImageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}

type Content struct {
	Type     string   `json:"type"`
	Text     string   `json:"text,omitempty"`
	ImageURL ImageURL `json:"image_url,omitempty"`
}

// NewContentFromImage creates a new Content from image data.
func NewContentFromImage(mime_type string, d []byte) (Content, error) {
	if !strings.HasPrefix(mime_type, "image/") {
		return Content{}, errors.New("mime_type must be image/*")
	}

	// Based on the python reference code in
	// https://platform.openai.com/docs/guides/vision/uploading-base-64-encoded-images
	// this should be the parallel of:
	//     base64.b64encode(image_file.read()).decode('utf-8')
	// which defaults to the standard base64 encoding.  I would have guessed
	// it would be using the URL-safe encoding but that isn't what the code is
	// saying.
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(d)))
	base64.StdEncoding.Encode(dst, d)

	image_url := strings.Builder{}
	image_url.WriteString("data:")
	image_url.WriteString(mime_type)
	image_url.WriteString(";base64,")
	image_url.Write(dst)

	url := ImageURL{
		URL: image_url.String(),
	}

	return Content{Type: "image_url", ImageURL: url}, nil
}

func NewContentFromImageURL(url string) Content {
	return Content{
		Type: "image_url",
		ImageURL: ImageURL{
			URL: url,
		},
	}
}

func NewContentFromText(text string) Content {
	return Content{
		Type: "text",
		Text: text,
	}
}

func (c *Client) CreateCompletion(ctx context.Context, p *CreateCompletionParams) (*CreateCompletionResponse, error) {
	if p.Model == "" {
		p.Model = c.model
	}
	if p.Stream {
		return nil, errors.New("use StreamingClient instead")
	}

	var r CreateCompletionResponse
	if err := c.s.MakeRequest(ctx, c.CreateCompletionEndpoint, p, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// CreateMMCompletion is the multi-modal version of CreateCompletion.
func (c *Client) CreateMMCompletion(ctx context.Context, p *CreateMMCompletionParams) (*CreateCompletionResponse, error) {
	if p.Model == "" {
		p.Model = c.model
	}
	if p.Stream {
		return nil, errors.New("use StreamingClient instead")
	}

	var r CreateCompletionResponse
	if err := c.s.MakeRequest(ctx, c.CreateCompletionEndpoint, p, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
