package strategist

import (
	"github.com/charmingruby/winston/internal/platform/aikit"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/ai/exp"
)

const agentName = "strategist"

type Strategist struct {
	*aikit.Agent
}

func Define(
	ak *aikit.AIKit,
	provider aikit.Provider,
	modelName string,
	opts ...exp.AgentOption[struct{}],
) (*Strategist, error) {
	opts = append([]exp.AgentOption[struct{}]{
		exp.WithDescription[struct{}]("Plans Overwatch positioning, rotations, objectives, teamfights, and ultimate usage."),
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

	return &Strategist{agent}, nil
}
