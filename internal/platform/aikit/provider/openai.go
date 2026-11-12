package provider

import (
	genai "github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core/api"
	"github.com/firebase/genkit/go/plugins/compat_oai/openai"
)

type OpenAI struct {
	plugin *openai.OpenAI
}

func NewOpenAI(apiKey string) *OpenAI {
	return &OpenAI{
		plugin: &openai.OpenAI{
			APIKey: apiKey,
		},
	}
}

func (p *OpenAI) Plugin() api.Plugin {
	return p.plugin
}

func (p *OpenAI) Model(name string) genai.ModelRef {
	return openai.ModelRef(name, nil)
}

func (p *OpenAI) Embedder(name string, dimensions int, _ bool) (genai.EmbedderRef, error) {
	return openai.NewEmbedderRef(
		name,
		&openai.TextEmbeddingConfig{
			Dimensions: dimensions,
		},
	), nil
}
