package domain

// Пользователь (User) — участник команды с уникальным идентификатором, именем и флагом активности isActive.
type User struct {
	ID       string `json:"user_id"`
	Name     string `json:"username"`
	TeamName string `json:"team_name,omitempty"`
	IsActive bool   `json:"is_active"`
}
