package domain

type Team struct {
	ID      int64   `json:"-"`
	Name    string  `json:"team_name"`
	Members []*User `json:"members"`
}
