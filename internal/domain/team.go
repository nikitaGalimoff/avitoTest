package domain

// TeamMember представляет участника команды
type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

// Team представляет команду
type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

