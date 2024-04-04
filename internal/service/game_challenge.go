package service

import (
	"github.com/elabosak233/cloudsdale/internal/model"
	"github.com/elabosak233/cloudsdale/internal/model/request"
	"github.com/elabosak233/cloudsdale/internal/model/response"
	"github.com/elabosak233/cloudsdale/internal/repository"
	"github.com/elabosak233/cloudsdale/internal/utils/calculate"
	"github.com/mitchellh/mapstructure"
)

type IGameChallengeService interface {
	Find(req request.GameChallengeFindRequest) (challenges []response.GameChallengeResponse, err error)
	Create(req request.GameChallengeCreateRequest) (err error)
	Update(req request.GameChallengeUpdateRequest) (err error)
	Delete(req request.GameChallengeDeleteRequest) (err error)
}

type GameChallengeService struct {
	gameRepository          repository.IGameRepository
	gameChallengeRepository repository.IGameChallengeRepository
	noticeRepository        repository.INoticeRepository
}

func NewGameChallengeService(appRepository *repository.Repository) IGameChallengeService {
	return &GameChallengeService{
		gameRepository:          appRepository.GameRepository,
		gameChallengeRepository: appRepository.GameChallengeRepository,
		noticeRepository:        appRepository.NoticeRepository,
	}
}

func (g *GameChallengeService) Find(req request.GameChallengeFindRequest) (challenges []response.GameChallengeResponse, err error) {
	games, _, _ := g.gameRepository.Find(request.GameFindRequest{
		ID: req.GameID,
	})
	game := games[0]
	gameChallenges, err := g.gameChallengeRepository.Find(req)
	for _, gameChallenge := range gameChallenges {
		var challenge response.GameChallengeResponse
		_ = mapstructure.Decode(gameChallenge, &challenge)
		pts := calculate.GameChallengePts(
			gameChallenge.MaxPts,
			gameChallenge.MinPts,
			gameChallenge.Challenge.Difficulty,
			len(challenge.Submissions),
			game.FirstBloodRewardRatio,
			game.SecondBloodRewardRatio,
			game.ThirdBloodRewardRatio,
		)
		challenge.Pts = pts
		for index, submission := range challenge.Submissions {
			submission.Pts = calculate.GameChallengePts(
				gameChallenge.MaxPts,
				gameChallenge.MinPts,
				gameChallenge.Challenge.Difficulty,
				int(submission.Rank-1),
				game.FirstBloodRewardRatio,
				game.SecondBloodRewardRatio,
				game.ThirdBloodRewardRatio,
			)
			if req.TeamID != 0 && submission.TeamID != nil && *(submission.TeamID) == req.TeamID {
				sub := submission
				challenge.Solved = sub
				break
			}
			challenge.Submissions[index] = submission
		}
		if req.SubmissionQty > 0 {
			challenge.Submissions = challenge.Submissions[:min(req.SubmissionQty, len(challenge.Submissions))]
		}
		challenges = append(challenges, challenge)
	}
	return challenges, err
}

func (g *GameChallengeService) Create(req request.GameChallengeCreateRequest) (err error) {
	var gameChallenge model.GameChallenge
	err = mapstructure.Decode(req, &gameChallenge)
	err = g.gameChallengeRepository.Create(gameChallenge)
	return err
}

func (g *GameChallengeService) Update(req request.GameChallengeUpdateRequest) (err error) {
	var gameChallenge model.GameChallenge
	err = mapstructure.Decode(req, &gameChallenge)
	err = g.gameChallengeRepository.Update(gameChallenge)
	if gameChallenge.IsEnabled != nil && *(gameChallenge.IsEnabled) {
		_, err = g.noticeRepository.Create(model.Notice{
			Type:        "new_challenge",
			ChallengeID: &gameChallenge.ChallengeID,
			GameID:      &gameChallenge.GameID,
		})
	}
	return err
}

func (g *GameChallengeService) Delete(req request.GameChallengeDeleteRequest) (err error) {
	var gameChallenge model.GameChallenge
	err = mapstructure.Decode(req, &gameChallenge)
	err = g.gameChallengeRepository.Delete(gameChallenge)
	return err
}
