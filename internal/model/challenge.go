package model

import (
	"gorm.io/gorm"
)

// Challenge is the challenge for Jeopardy-style CTF game.
type Challenge struct {
	ID            uint          `json:"id"`                                                     // The challenge's id. As primary key.
	Title         string        `gorm:"type:varchar(32);not null;" json:"title"`                // The challenge's title.
	Description   string        `gorm:"type:text;not null;" json:"description"`                 // The challenge's description.
	CategoryID    uint          `gorm:"not null;" json:"category_id"`                           // The challenge's category.
	Category      *Category     `json:"category,omitempty"`                                     // The challenge's category.
	HasAttachment *bool         `gorm:"not null;default:false" json:"has_attachment"`           // Whether the challenge has attachment.
	AttachmentURL string        `gorm:"type:varchar(255);" json:"attachment_url,omitempty"`     // The challenge's attachment URL.
	IsPracticable *bool         `gorm:"not null;default:false" json:"is_practicable,omitempty"` // Whether the challenge is practicable. (Is the practice field visible.)
	IsDynamic     *bool         `gorm:"default:false" json:"is_dynamic"`                        // Whether the challenge is based on dynamic container.
	Difficulty    int64         `gorm:"default:1" json:"difficulty"`                            // The degree of difficulty. (From 1 to 5)
	PracticePts   int64         `gorm:"default:200" json:"practice_pts,omitempty"`              // The points will be given when the challenge is solved in practice field.
	Duration      int64         `gorm:"default:1800" json:"duration,omitempty"`                 // The duration of container maintenance in the initial state. (Seconds)
	CreatedAt     int64         `gorm:"autoUpdateTime:milli" json:"created_at,omitempty"`       // The challenge's creation time.
	UpdatedAt     int64         `gorm:"autoUpdateTime:milli" json:"updated_at,omitempty"`       // The challenge's last update time.
	Flags         []*Flag       `json:"flags,omitempty"`
	Hints         []*Hint       `json:"hints,omitempty"`
	Images        []*Image      `json:"images,omitempty"`
	Solved        *Submission   `json:"solved,omitempty"`
	Submissions   []*Submission `json:"submissions,omitempty"`
}

func (c *Challenge) BeforeDelete(db *gorm.DB) (err error) {
	db.Table("flags").Where("challenge_id = ?", c.ID).Delete(&Flag{})
	db.Table("hints").Where("challenge_id = ?", c.ID).Delete(&Hint{})

	var images []*Image
	db.Table("images").Where("challenge_id = ?", c.ID).Find(&images)
	for _, image := range images {
		db.Table("images").Delete(image)
	}

	db.Table("submissions").Where("challenge_id = ?", c.ID).Delete(&Submission{})
	db.Table("game_challenges").Where("challenge_id = ?", c.ID).Delete(&GameChallenge{})
	return nil
}
