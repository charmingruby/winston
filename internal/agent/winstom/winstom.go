package winstom

import (
	"github.com/charmingruby/winstom/internal/platform/aikit"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/ai/exp"
)

const agentName = "Winstom"

type Winstom struct {
	*aikit.Agent
}

func Define(
	ak *aikit.AIKit,
	provider aikit.Provider,
	modelName string,
	opts ...exp.AgentOption[struct{}],
) (*Winstom, error) {
	agent, err := aikit.DefineAgent(
		ak,
		agentName,
		provider.Model(modelName),
		nil,
		exp.InlinePrompt{
			ai.WithSystem(systemPrompt),
		},
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return &Winstom{agent}, nil
}
