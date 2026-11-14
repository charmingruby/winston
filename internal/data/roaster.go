package data

import "github.com/charmingruby/winston/internal/model"

var Roaster = []model.Player{
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
}
