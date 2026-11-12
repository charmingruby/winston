package provider

import (
	genai "github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core/api"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	googleai "google.golang.org/genai"
)

type EmbeddingTask string

const (
	RetrievalDocumentEmbeddingTask = "RETRIEVAL_DOCUMENT"
	RetrievalQueryEmbeddingTask    = "RETRIEVAL_QUERY"
)

type Google struct {
	plugin *googlegenai.GoogleAI
}

func NewGoogle(apiKey string) *Google {
	return &Google{
		plugin: &googlegenai.GoogleAI{
			APIKey: apiKey,
		},
	}
}

func (p *Google) Plugin() api.Plugin {
	return p.plugin
}

func (p *Google) Model(name string) genai.ModelRef {
	return googlegenai.ModelRef(p.plugin.Name()+"/"+name, nil)
}

func (p *Google) Embedder(name string, dimensions int, query bool) (genai.EmbedderRef, error) {
	task := RetrievalDocumentEmbeddingTask
	if query {
		task = RetrievalQueryEmbeddingTask
	}

	return googlegenai.EmbedderRef(p.plugin.Name()+"/"+name, &googleai.EmbedContentConfig{
		TaskType:             task,
		OutputDimensionality: new(int32(dimensions)),
	}), nil
}
