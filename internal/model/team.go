package model

import (
	"gorm.io/gorm"
)

type Team struct {
	ID          uint    `json:"id"`                                                // The team's id. As primary key.
	Name        string  `gorm:"type:varchar(36);not null" json:"name"`             // The team's name.
	Description string  `gorm:"type:text" json:"description"`                      // The team's description.
	Email       string  `gorm:"type:varchar(64);" json:"email,omitempty"`          // The team's email.
	CaptainID   uint    `gorm:"not null" json:"captain_id,omitempty"`              // The captain's id.
	Captain     *User   `json:"captain,omitempty"`                                 // The captain's user.
	IsLocked    *bool   `gorm:"not null;default:false" json:"is_locked,omitempty"` // Whether the team is locked. (true/false)
	CreatedAt   int64   `gorm:"autoUpdateTime:milli" json:"created_at,omitempty"`  // The team's creation time.
	UpdatedAt   int64   `gorm:"autoUpdateTime:milli" json:"updated_at,omitempty"`  // The team's last update time.
	Users       []*User `gorm:"many2many:user_teams;" json:"users,omitempty"`      // The team's users.
}

func (t *Team) BeforeDelete(db *gorm.DB) (err error) {
	db.Table("user_teams").Where("team_id = ?", t.ID).Delete(&UserTeam{})
	db.Table("game_teams").Where("team_id = ?", t.ID).Delete(&GameTeam{})
	return nil
}
