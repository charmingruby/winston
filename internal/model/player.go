package model

type Player struct {
	Name  string       `json:"name" jsonschema:"description=Player name or nickname"`
	Roles []PlayerRole `json:"roles" jsonschema:"description=Roles this player can play, with hero preferences for each role. These are options, not simultaneous assignments."`
}

type PlayerRole struct {
	Role     Role   `json:"role" jsonschema:"description=Playable role: tank, damage, or support"`
	HeroPool []Hero `json:"main_heroes" jsonschema:"description=Heroes the player can play in this role"`
}
