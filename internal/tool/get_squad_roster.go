package tool

import (
	"context"

	"github.com/charmingruby/erebor/aikit"
	"github.com/charmingruby/winston/internal/data"
	"github.com/charmingruby/winston/internal/model"
)

type GetSquadRosterInput struct{}

type GetSquadRosterOutput struct {
	Players []model.Player `json:"players" jsonschema:"description=Known squad players. An empty list means no squad information has been configured."`
}

func GetSquadRoster() aikit.Tool {
	return aikit.DefineTool(
		"get_squad_roster",
		"Returns the squad's configured players, roles, main heroes, best-performing heroes, and playstyle notes. "+
			"Players may play multiple roles; hero preferences are grouped by role. Assign each player to only one role per composition. "+
			"Use before recommending team compositions or hero swaps. "+
			"These are squad-reported preferences, not live profiles or verified statistics. "+
			"If the roster is empty or information is missing, ask the user instead of inventing it.",
		func(ctx context.Context, input GetSquadRosterInput) (GetSquadRosterOutput, error) {

			return GetSquadRosterOutput{
				Players: data.Roster,
			}, nil
		},
	)
}
