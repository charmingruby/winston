package reviewer

import (
	"github.com/charmingruby/erebor/aikit"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/ai/exp"
)

const agentName = "reviewer"

type Reviewer struct {
	*aikit.Agent
}

func Define(
	ak *aikit.AIKit,
	provider aikit.Provider,
	modelName string,
	opts ...exp.AgentOption[struct{}],
) (*Reviewer, error) {
	opts = append([]exp.AgentOption[struct{}]{
		exp.WithDescription[struct{}]("Reviews Overwatch recommendations for consistency, practicality, and unsupported claims."),
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

	return &Reviewer{agent}, nil
}
