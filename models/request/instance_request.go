package request

type InstanceCreateRequest struct {
	ChallengeId int64 `binding:"required" json:"challenge_id"`
	TeamId      int64 `json:"team_id"`
	GameId      int64 `json:"game_id"`
	UserId      int64 `json:"-"`
}

type InstanceFindRequest struct {
	ChallengeId int64 `json:"challenge_id"`
	UserId      int64 `json:"-"`
	TeamId      int64 `json:"team_id"`
	GameId      int64 `json:"game_id"`
	IsAvailable int   `json:"is_available"`
	Page        int   `json:"page"`
	Size        int   `json:"size"`
}

type InstanceRemoveRequest struct {
	InstanceId int64 `binding:"required" json:"id"`
	TeamId     int64 `json:"team_id"`
	GameId     int64 `json:"game_id"`
	UserId     int64 `json:"-"`
}

type InstanceRenewRequest struct {
	InstanceId int64 `binding:"required" json:"id"`
	TeamId     int64 `json:"team_id"`
	GameId     int64 `json:"game_id"`
	UserId     int64 `json:"-"`
}