package types

type CreateServerConfig struct {
	Name         string `json:"name" validate:"required"`
	ServerType   string `json:"server_type" validate:"required"`
	PlayersCount int    `json:"players_count" validate:"required,min=1"`
	PlanType     string `json:"plan_type" validate:"required"`
	Difficulty   string `json:"difficulty" validate:"required"`
	OnlineMode   bool   `json:"online_mode" validate:"required,boolean"`
}
