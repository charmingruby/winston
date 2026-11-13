package analyst

import (
	"github.com/charmingruby/stan/internal/platform/aikit"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/ai/exp"
)

const agentName = "analyst"

type Analyst struct {
	*aikit.Agent
}

func Define(
	ak *aikit.AIKit,
	provider aikit.Provider,
	modelName string,
	opts ...exp.AgentOption[struct{}],
) (*Analyst, error) {
	opts = append([]exp.AgentOption[struct{}]{
		exp.WithDescription[struct{}]("Analyzes Overwatch hero compositions, synergies, counters, and role-specific swaps."),
	}, opts...)

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

	return &Analyst{agent}, nil
}
