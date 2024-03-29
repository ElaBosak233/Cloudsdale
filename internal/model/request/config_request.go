package request

type Site struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type Container struct {
	ParallelLimit int64 `json:"parallel_limit,omitempty"`
	RequestLimit  int64 `json:"request_limit,omitempty"`
}

type User struct {
	AllowRegistration bool `json:"allow_registration,omitempty"`
}

type ConfigUpdateRequest struct {
	Site      Site      `json:"site,omitempty"`
	Container Container `json:"container,omitempty"`
	User      User      `json:"user,omitempty"`
}
