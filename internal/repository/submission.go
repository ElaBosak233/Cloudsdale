package repository

import (
	"github.com/elabosak233/pgshub/internal/model"
	"github.com/elabosak233/pgshub/internal/model/dto/request"
	"github.com/elabosak233/pgshub/internal/model/dto/response"
	"gorm.io/gorm"
)

type ISubmissionRepository interface {
	Insert(submission model.Submission) (err error)
	Delete(id uint) (err error)
	Find(req request.SubmissionFindRequest) (submissions []response.SubmissionResponse, count int64, err error)
	BatchFind(req request.SubmissionBatchFindRequest) (submissions []response.SubmissionResponse, err error)
}

type SubmissionRepository struct {
	Db *gorm.DB
}

func NewSubmissionRepository(Db *gorm.DB) ISubmissionRepository {
	return &SubmissionRepository{Db: Db}
}

func (t *SubmissionRepository) Insert(submission model.Submission) (err error) {
	result := t.Db.Table("submissions").Create(&submission)
	return result.Error
}

func (t *SubmissionRepository) Delete(id uint) (err error) {
	result := t.Db.Table("submissions").Delete(&model.Submission{ID: id})
	return result.Error
}

func (t *SubmissionRepository) Find(req request.SubmissionFindRequest) (submissions []response.SubmissionResponse, count int64, err error) {
	applyFilters := func(q *gorm.DB) *gorm.DB {
		if req.UserID != 0 && req.TeamID == nil && req.GameID == nil {
			q = q.Where("user_id = ?", req.UserID)
		}
		if req.ChallengeID != 0 {
			q = q.Where("challenge_id = ?", req.ChallengeID)
		}
		if req.TeamID != nil {
			q = q.Where("team_id = ?", *(req.TeamID))
		}
		if req.GameID != nil {
			q = q.Where("game_id = ?", *(req.GameID))
		}
		if req.Status != 0 {
			q = q.Where("status = ?", req.Status)
		}
		if req.IsDetailed == 0 {
			q = q.Omit("flag")
		}
		return q
	}
	db := applyFilters(t.Db.Table("submissions"))
	ct := applyFilters(t.Db.Table("submissions"))
	result := ct.Model(&model.Submission{}).Count(&count)
	if len(req.SortBy) > 0 {
		sortKey := req.SortBy[0]
		sortOrder := req.SortBy[1]
		if sortOrder == "asc" {
			db = db.Order("submissions." + sortKey + " ASC")
		} else if sortOrder == "desc" {
			db = db.Order("submissions." + sortKey + " DESC")
		}
	} else {
		db = db.Order("submissions.id DESC") // 默认采用 IDs 降序排列
	}
	if req.Page != 0 && req.Size > 0 {
		offset := (req.Page - 1) * req.Size
		db = db.Offset(offset).Limit(req.Size)
	}
	db = db.Joins("INNER JOIN users ON submissions.user_id = users.id").
		Joins("INNER JOIN challenges ON submissions.challenge_id = challenges.id").
		Joins("LEFT JOIN teams ON submissions.team_id = teams.id").
		Joins("LEFT JOIN games ON submissions.game_id = games.id")
	result = db.Find(&submissions)
	return submissions, count, result.Error
}

func (t *SubmissionRepository) BatchFind(req request.SubmissionBatchFindRequest) (submissions []response.SubmissionResponse, err error) {
	applyFilters := func(q *gorm.DB) *gorm.DB {
		if req.UserID != 0 {
			q = q.Where("submissions.user_id = ?", req.UserID)
		}
		if req.TeamID != nil {
			q = q.Where("submissions.team_id = ?", *(req.TeamID))
		}
		if req.GameID != nil {
			q = q.Where("submissions.game_id = ?", *(req.GameID))
		}
		if req.Status != 0 {
			q = q.Where("submissions.status = ?", req.Status)
		}
		return q
	}
	db := applyFilters(t.Db.Table("submissions"))
	if len(req.SortBy) > 0 {
		sortKey := req.SortBy[0]
		sortOrder := req.SortBy[1]
		if sortOrder == "asc" {
			db = db.Order("submissions." + sortKey + " ASC")
		} else if sortOrder == "desc" {
			db = db.Order("submissions." + sortKey + " DESC")
		}
	}
	db = db.Joins("INNER JOIN users ON submissions.user_id = users.id").
		Joins("LEFT JOIN teams ON submissions.team_id = teams.id").
		Joins("LEFT JOIN challenges ON submissions.challenge_id = challenges.id").
		Where("submissions.challenge_id IN ?", req.ChallengeID)
	_ = db.Find(&submissions)
	return submissions, err
}