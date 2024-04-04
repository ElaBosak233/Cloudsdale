package service

import (
	"errors"
	"github.com/elabosak233/cloudsdale/internal/model"
	"github.com/elabosak233/cloudsdale/internal/model/request"
	"github.com/elabosak233/cloudsdale/internal/repository"
)

type IUserTeamService interface {
	Join(req request.TeamUserJoinRequest) (err error)
	Create(req request.TeamUserCreateRequest) (err error)
	Delete(req request.TeamUserDeleteRequest) (err error)
}

type UserTeamService struct {
	userTeamRepository repository.IUserTeamRepository
	teamRepository     repository.ITeamRepository
	userRepository     repository.IUserRepository
}

func NewUserTeamService(appRepository *repository.Repository) IUserTeamService {
	return &UserTeamService{
		userTeamRepository: appRepository.UserTeamRepository,
		teamRepository:     appRepository.TeamRepository,
		userRepository:     appRepository.UserRepository,
	}
}

func (t *UserTeamService) Join(req request.TeamUserJoinRequest) (err error) {
	user, err := t.userRepository.FindById(req.UserID)
	team, err := t.teamRepository.FindById(req.TeamID)
	if err != nil || user.ID == 0 || team.ID == 0 {
		return errors.New("用户或团队不存在")
	}
	if team.InviteToken != req.InviteToken {
		return errors.New("邀请码错误")
	}
	err = t.userTeamRepository.Create(model.UserTeam{
		TeamID: team.ID,
		UserID: req.UserID,
	})
	return err
}

func (t *UserTeamService) Create(req request.TeamUserCreateRequest) (err error) {
	user, err := t.userRepository.FindById(req.UserID)
	team, err := t.teamRepository.FindById(req.TeamID)
	if err != nil || user.ID == 0 || team.ID == 0 {
		return errors.New("用户或团队不存在")
	}
	err = t.userTeamRepository.Create(model.UserTeam{
		TeamID: team.ID,
		UserID: req.UserID,
	})
	return err
}

func (t *UserTeamService) Delete(req request.TeamUserDeleteRequest) (err error) {
	user, err := t.userRepository.FindById(req.UserID)
	team, err := t.teamRepository.FindById(req.TeamID)
	if err != nil || user.ID == 0 || team.ID == 0 {
		return errors.New("用户或团队不存在")
	}
	err = t.userTeamRepository.Delete(model.UserTeam{
		TeamID: team.ID,
		UserID: req.UserID,
	})
	return err
}
