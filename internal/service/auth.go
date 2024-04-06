package service

import (
	"github.com/elabosak233/cloudsdale/internal/model"
	"github.com/elabosak233/cloudsdale/internal/repository"
)

type IAuthService interface {
	CanModifyUser(user *model.User, targetUserID uint) bool
	CanModifyTeam(user *model.User, targetTeamID uint) bool
}

type AuthService struct {
	userRepository repository.IUserRepository
	teamRepository repository.ITeamRepository
}

func NewAuthService(appRepository *repository.Repository) IAuthService {
	return &AuthService{
		userRepository: appRepository.UserRepository,
		teamRepository: appRepository.TeamRepository,
	}
}

func (a *AuthService) CanModifyUser(user *model.User, targetUserID uint) bool {
	return user.Group.Name == "admin" || user.ID == targetUserID
}

func (a *AuthService) CanModifyTeam(user *model.User, targetTeamID uint) bool {
	isCaptain := func() bool {
		for _, team := range user.Teams {
			if team.ID == targetTeamID && team.CaptainID == user.ID {
				return true
			}
		}
		return false
	}
	return user.Group.Name == "admin" || isCaptain()
}