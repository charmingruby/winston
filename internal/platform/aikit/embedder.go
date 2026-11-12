package aikit

import (
	"context"
	"errors"
	"math"

	"github.com/firebase/genkit/go/genkit"

	genai "github.com/firebase/genkit/go/ai"
	googleai "google.golang.org/genai"
)

type EmbeddingTask string

const (
	RetrievalDocumentEmbeddingTask = "RETRIEVAL_DOCUMENT"
	RetrievalQueryEmbeddingTask    = "RETRIEVAL_QUERY"
)

type Embedder struct {
	ak    *AIKit
	model string
	dim   int
}

func NewEmbedder(ak *AIKit, model string, dim int) *Embedder {
	return &Embedder{
		ak:    ak,
		model: model,
		dim:   dim,
	}
}

func (e *Embedder) Embed(
	ctx context.Context,
	text string,
	query bool,
) ([]float32, error) {
	task := RetrievalDocumentEmbeddingTask
	if query {
		task = RetrievalQueryEmbeddingTask
	}

	response, err := genkit.Embed(
		ctx,
		e.ak.Genkit,
		genai.WithEmbedderName(e.model),
		genai.WithDocs(genai.DocumentFromText(text, nil)),
		genai.WithConfig(
			&googleai.EmbedContentConfig{
				TaskType:             task,
				OutputDimensionality: new(int32(e.dim)),
			},
		),
	)
	if err != nil {
		return nil, err
	}

	if len(response.Embeddings) != 1 || response.Embeddings[0] == nil {
		return nil, errors.New("unexpected embedding response")
	}

	vector := response.Embeddings[0].Embedding
	if len(vector) != e.dim {
		return nil, errors.New("unexpected embedding dimensions")
	}

	var norm float64
	for _, value := range vector {
		if math.IsNaN(float64(value)) || math.IsInf(float64(value), 0) {
			return nil, errors.New("non-finite embedding")
		}

		norm += float64(value) * float64(value)
	}

	if norm == 0 {
		return nil, errors.New("zero embedding")
	}

	return vector, nil
}
