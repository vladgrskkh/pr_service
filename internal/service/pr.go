package service

import (
	"log/slog"

	"github.com/vladgrskkh/pr_service/internal/domain"
	"github.com/vladgrskkh/pr_service/internal/repository"
)

type PullReqService struct {
	logger       *slog.Logger // think about why i need this
	pullReqsRepo *repository.PullRequestRepo
	teamsRepo    *repository.TeamRepository
	usersRepo    *repository.UsersRepo
}

func NewPullReqService(
	logger *slog.Logger,
	pullReqsRepo *repository.PullRequestRepo,
	teamsRepo *repository.TeamRepository,
	usersRepo *repository.UsersRepo,
) *PullReqService {
	return &PullReqService{
		logger:       logger,
		pullReqsRepo: pullReqsRepo,
		teamsRepo:    teamsRepo,
		usersRepo:    usersRepo,
	}
}

func (s *PullReqService) SetIsActiveUser(id int64, isActive bool) (*domain.User, error) {
	user, err := s.usersRepo.SetIsActive(id, isActive)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PullReqService) GetReviewByUser(id int64) ([]*domain.PR, error) {
	prs, err := s.pullReqsRepo.GetAllForUser(id)
	if err != nil {
		return nil, err
	}

	return prs, nil
}

func (s *PullReqService) CreateTeam(name string, members []*domain.User) (*domain.Team, error) {
	team := &domain.Team{
		Name:    name,
		Members: members,
	}

	err := s.teamsRepo.Insert(team)
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (s *PullReqService) GetTeam(name string) (*domain.Team, error) {
	team, err := s.teamsRepo.Get(name)
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (s *PullReqService) MergePullReq(id int64) (*domain.PR, error) {
	pr := &domain.PR{
		ID:     id,
		Status: "merged",
	}

	err := s.pullReqsRepo.UpdateStatus(pr)
	if err != nil {
		return nil, err
	}

	return pr, nil
}
