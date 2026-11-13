package supervisor

import (
	"github.com/charmingruby/winstom/internal/platform/aikit"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/ai/exp"
	middlewarex "github.com/firebase/genkit/go/plugins/middleware/exp"
)

const agentName = "supervisor"

type Supervisor struct {
	Agent *aikit.Agent
}

func Define(
	ak *aikit.AIKit,
	provider aikit.Provider,
	modelName string,
	agents []exp.AgentRef,
	opts ...exp.AgentOption[struct{}],
) (*Supervisor, error) {
	agent, err := aikit.DefineAgent(
		ak,
		agentName,
		provider.Model(modelName),
		nil,
		exp.InlinePrompt{
			ai.WithSystem(systemPrompt),
			ai.WithMaxTurns(8),
			ai.WithUse(&middlewarex.Agents{
				Agents:         agents,
				MaxDelegations: 5,
				HistoryLength:  0,
			}),
		},
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return &Supervisor{
		Agent: agent,
	}, nil
}
