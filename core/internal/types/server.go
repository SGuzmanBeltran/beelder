package types

type CreateServerData struct {
	ServerID     string
	ServerConfig *CreateServerConfig
	ImageName 	 string
}
type CreateServerConfig struct {
	Name         string `json:"name" validate:"required,min=3,max=64"`
	ServerType   string `json:"server_type" validate:"required"`
	PlayersCount int    `json:"players_count" validate:"required,min=1,max=100"`
	PlanType     string `json:"plan_type" validate:"required,oneof=free budget premium"`
	Difficulty   string `json:"difficulty" validate:"required,oneof=peaceful easy normal hard"`
	OnlineMode   bool   `json:"online_mode" validate:"boolean"`
}

type MemorySettings struct {
	Min string
	Max string
}
