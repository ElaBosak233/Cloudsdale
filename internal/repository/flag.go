package repository

import (
	"github.com/elabosak233/cloudsdale/internal/model"
	"gorm.io/gorm"
)

type IFlagRepository interface {
	Insert(flag model.Flag) (f model.Flag, err error)
	Update(flag model.Flag) (f model.Flag, err error)
	Delete(flag model.Flag) (err error)
	FindByChallengeID(challengeIDs []uint) (flags []model.Flag, err error)
}

type FlagRepository struct {
	db *gorm.DB
}

func NewFlagRepository(db *gorm.DB) IFlagRepository {
	return &FlagRepository{db: db}
}

func (t *FlagRepository) Insert(flag model.Flag) (f model.Flag, err error) {
	result := t.db.Table("flags").Create(&flag)
	return flag, result.Error
}

func (t *FlagRepository) Update(flag model.Flag) (f model.Flag, err error) {
	result := t.db.Table("flags").Model(&flag).Updates(&flag)
	return flag, result.Error
}

func (t *FlagRepository) Delete(flag model.Flag) (err error) {
	result := t.db.Table("flags").Delete(&flag)
	return result.Error
}

func (t *FlagRepository) FindByChallengeID(challengeIDs []uint) (flags []model.Flag, err error) {
	result := t.db.Table("flags").
		Where("challenge_id IN ?", challengeIDs).
		Find(&flags)
	return flags, result.Error
}
