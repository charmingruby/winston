package data

import "github.com/charmingruby/winston/internal/model"

var Roster = []model.Player{
	{
		Name: "G",
		Roles: []model.PlayerRole{
			{
				Role: model.TankRole,
				HeroPool: []model.Hero{
					model.ZaryaHero,
					model.OrisaHero,
					model.DVaHero,
					model.WinstonHero,
				},
			},
			{
				Role: model.SupRole,
				HeroPool: []model.Hero{
					model.BaptisteHero,
					model.KirikoHero,
					model.IllariHero,
					model.JunoHero,
					model.MercyHero,
					model.LucioHero,
					model.ZenyattaHero,
				},
			},
		},
	},
	{
		Name: "S",
		Roles: []model.PlayerRole{
			{
				Role: model.DPSRole,
				HeroPool: []model.Hero{
					model.AsheHero,
					model.BastionHero,
					model.ReaperHero,
					model.HanzoHero,
				},
			},
			{
				Role: model.SupRole,
				HeroPool: []model.Hero{
					model.MoiraHero,
					model.KirikoHero,
				},
			},
		},
	},
}
