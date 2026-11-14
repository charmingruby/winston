package tool

import (
	"context"

	"github.com/charmingruby/erebor/aikit"
	"github.com/charmingruby/winston/internal/model"
)

type GetTeamRoasterInput struct{}

type GetTeamRoasterOutput struct {
	Players []model.Player `json:"players" jsonschema:"description=Known team players. An empty list means no team information has been configured."`
}

func GetTeamRoaster() aikit.Tool {
	return aikit.DefineTool(
		"get_team_roster",
		"Returns the team's configured players, roles, main heroes, best-performing heroes, and playstyle notes. "+
			"Players may play multiple roles; hero preferences are grouped by role. Assign each player to only one role per composition. "+
			"Use before recommending team compositions or hero swaps. "+
			"These are team-reported preferences, not live profiles or verified statistics. "+
			"If the roster is empty or information is missing, ask the user instead of inventing it.",
		func(ctx context.Context, input GetTeamRoasterInput) (GetTeamRoasterOutput, error) {
			return GetTeamRoasterOutput{
				Players: []model.Player{
					{
						Name: "G",
						Roles: []model.PlayerRole{
							{
								Role: model.SupRole,
								HeroPool: []model.Hero{
									model.AnaHero,
									model.LifeweaverHero,
									model.MoiraHero,
								},
							},
						},
					},
				},
			}, nil
		},
	)
}
