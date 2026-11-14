package winston

import (
	"github.com/charmingruby/erebor/aikit"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/ai/exp"
	middlewarex "github.com/firebase/genkit/go/plugins/middleware/exp"
)

const agentName = "winston"

type Winston struct {
	*aikit.Agent
}

func Define(
	ak *aikit.AIKit,
	provider aikit.Provider,
	modelName string,
	agents []exp.AgentRef,
	tools []aikit.Tool,
	opts ...exp.AgentOption[struct{}],
) (*Winston, error) {
	agent, err := aikit.DefineAgent(
		ak,
		agentName,
		provider.Model(modelName),
		tools,
		exp.InlinePrompt{
			ai.WithSystem(systemPrompt),
			ai.WithMaxTurns(8),
			ai.WithUse(&middlewarex.Agents{
				Agents:         agents,
				MaxDelegations: 3,
				HistoryLength:  0,
			}),
		},
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return &Winston{agent}, nil
}
