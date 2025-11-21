package domain

// Пользователь (User) — участник команды с уникальным идентификатором, именем и флагом активности isActive.
type User struct {
	ID       int64  `json:"user_id"`
	Name     string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}
