package types

type CreateServerData struct {
	ContainerID  string
	ServerID     string
	ServerConfig *CreateServerConfig
	ImageName    string
}
type CreateServerConfig struct {
	Name          string `json:"name" validate:"required,min=3,max=64"`
	ServerVersion string `json:"server_version" validate:"required"`
	ServerType    string `json:"server_type" validate:"required"`
	Region        string `json:"region" validate:"required"`
	PlayerCount   int    `json:"player_count" validate:"required,min=1,max=100"`
	RamPlan       string `json:"ram_plan" validate:"required"`
	Difficulty    string `json:"difficulty" validate:"required,oneof=peaceful easy normal hard hardcore"`
	OnlineMode    bool   `json:"online_mode"`
}

type RecommendationServerParams struct {
	PlayerCount int    `query:"player_count" validate:"required,min=1,max=100"`
	ServerType  string `query:"server_type" validate:"required"`
	Region      string `query:"region" validate:"required"`
}

type MemorySettings struct {
	Min string
	Max string
}

type RecommendationResponse struct {
	Recommendation string `json:"recommendation"`
}
