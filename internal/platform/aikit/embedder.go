package aikit

import (
	"context"
	"errors"
	"math"

	"github.com/charmingruby/stan/internal/platform/aikit/provider"
	"github.com/firebase/genkit/go/genkit"

	genai "github.com/firebase/genkit/go/ai"
)

type Embedder struct {
	provider provider.Provider
	genkit   *genkit.Genkit
	model    string
	dim      int
}

func NewEmbedder(ak *AIKit, provider provider.Provider, model string, dim int) (*Embedder, error) {
	if _, err := provider.Embedder(model, dim, false); err != nil {
		return nil, err
	}

	return &Embedder{
		genkit:   ak.Genkit,
		model:    model,
		provider: provider,
		dim:      dim,
	}, nil
}

func (e *Embedder) Embed(
	ctx context.Context,
	text string,
	query bool,
) ([]float32, error) {
	ref, err := e.provider.Embedder(e.model, e.dim, query)
	if err != nil {
		return nil, err
	}

	response, err := genkit.Embed(
		ctx,
		e.genkit,
		genai.WithEmbedder(ref),
		genai.WithDocs(genai.DocumentFromText(text, nil)),
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
