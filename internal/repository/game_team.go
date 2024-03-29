package repository

import (
	"github.com/elabosak233/cloudsdale/internal/model"
	"gorm.io/gorm"
)

type IGameTeamRepository interface {
	Insert(gameTeam model.GameTeam) (err error)
	Update(gameTeam model.GameTeam) (err error)
	Delete(gameTeam model.GameTeam) (err error)
	Find(gameTeam model.GameTeam) (gameTeams []model.GameTeam, err error)
	FindByGameID(gameID uint) (gameTeams []model.GameTeam, err error)
	FindByTeamID(teamID uint) (gameTeams []model.GameTeam, err error)
}

type GameTeamRepository struct {
	db *gorm.DB
}

func NewGameTeamRepository(db *gorm.DB) IGameTeamRepository {
	return &GameTeamRepository{db: db}
}

func (g *GameTeamRepository) Insert(gameTeam model.GameTeam) (err error) {
	result := g.db.Table("game_teams").Create(&gameTeam)
	return result.Error
}

func (g *GameTeamRepository) Delete(gameTeam model.GameTeam) (err error) {
	result := g.db.Table("game_teams").Delete(&gameTeam)
	return result.Error
}

func (g *GameTeamRepository) Update(gameTeam model.GameTeam) (err error) {
	result := g.db.Table("game_teams").Where(
		"game_id = ? AND team_id = ?", gameTeam.GameID, gameTeam.TeamID,
	).Updates(&gameTeam)
	return result.Error
}

func (g *GameTeamRepository) Find(gameTeam model.GameTeam) (gameTeams []model.GameTeam, err error) {
	result := g.db.Table("game_teams").
		Where(&gameTeam).
		Preload("Team", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Captain", func(db *gorm.DB) *gorm.DB {
				return db.Select([]string{"id", "nickname", "username", "email"})
			}).Preload("Users", func(db *gorm.DB) *gorm.DB {
				return db.Select([]string{"id", "nickname", "username", "email"})
			})
		}).
		Preload("Game").
		Find(&gameTeams)
	return gameTeams, result.Error
}

func (g *GameTeamRepository) FindByGameID(gameID uint) (gameTeams []model.GameTeam, err error) {
	result := g.db.Table("game_teams").Where("game_id = ?", gameID).Find(&gameTeams)
	return gameTeams, result.Error
}

func (g *GameTeamRepository) FindByTeamID(teamID uint) (gameTeams []model.GameTeam, err error) {
	result := g.db.Table("game_teams").Where("team_id = ?", teamID).Find(&gameTeams)
	return gameTeams, result.Error
}
