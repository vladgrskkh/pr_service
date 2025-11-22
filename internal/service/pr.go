package service

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"time"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/vladgrskkh/pr_service/internal/domain"
	"github.com/vladgrskkh/pr_service/internal/repository"
)

var (
	ErrUserNotAssigned = errors.New("reviwer is not assigned to this PR")
	ErrNoCandidate     = errors.New("no active replacement candidate in team")
	ErrMergedPRChange  = errors.New("cannot reassign on merged PR")
)

type PullReqService struct {
	logger       *slog.Logger // think about why i need this
	pullReqsRepo *repository.PullRequestRepo
	teamsRepo    *repository.TeamRepository
	usersRepo    *repository.UsersRepo
	trManager    *manager.Manager
}

func NewPullReqService(
	logger *slog.Logger,
	pullReqsRepo *repository.PullRequestRepo,
	teamsRepo *repository.TeamRepository,
	usersRepo *repository.UsersRepo,
	trManager *manager.Manager,
) *PullReqService {
	return &PullReqService{
		logger:       logger,
		pullReqsRepo: pullReqsRepo,
		teamsRepo:    teamsRepo,
		usersRepo:    usersRepo,
		trManager:    trManager,
	}
}

func (s *PullReqService) SetIsActiveUser(id string, isActive bool) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user, err := s.usersRepo.SetIsActive(ctx, id, isActive)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PullReqService) GetReviewByUser(id string) ([]*domain.PR, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	prs, err := s.pullReqsRepo.GetAllForUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return prs, nil
}

func (s *PullReqService) CreateTeam(name string, members []string) (*domain.Team, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	team := &domain.Team{
		Name: name,
	}

	err := s.teamsRepo.Insert(ctx, team)
	if err != nil {
		return nil, err
	}

	// add logic for adding users to team

	return team, nil
}

func (s *PullReqService) GetTeam(name string) (*domain.Team, []string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	team, err := s.teamsRepo.Get(ctx, name)
	if err != nil {
		return nil, nil, err
	}

	// implement logic here
	members := make([]string, 0)

	return team, members, nil
}

func (s *PullReqService) MergePullReq(id string) (*domain.PR, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	pr := &domain.PR{
		ID:     id,
		Status: "merged",
	}

	err := s.pullReqsRepo.UpdateStatus(ctx, pr)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *PullReqService) CreatePullReq(id, name, authorID string) (*domain.PR, error) {
	// problem here is when i query users and than someone change status(false)
	// i guess its fine cause they were active the moment i query them
	// it no matter if afterwards there no active
	//
	// another thing to consider is do i need to use transaction
	// i dont change any state in users so i gueess i dont
	pr := &domain.PR{
		ID:                id,
		Name:              name,
		AuthorID:          authorID,
		Status:            "open",
		AssignedReviewers: s.assigneReviewers(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.pullReqsRepo.Insert(ctx, pr)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

// i need to query all active users in the same team as author and pick 2
// if there are less than 2 active users in the team, assigne what i have (even 0)
func (s *PullReqService) assigneReviewers() []string {
	return make([]string, 2)
}

func (s *PullReqService) ReassignReviewer(prID, userID string) (*domain.PR, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	pr, err := s.pullReqsRepo.GetByID(ctx, prID)
	if err != nil {
		return nil, err
	}

	if pr.Status != "open" {
		return nil, ErrMergedPRChange
	}

	if pr.AuthorID == userID {
		return nil, ErrUserNotAssigned
	}

	usersID := append([]string{}, pr.AssignedReviewers...)

	i := slices.Index(usersID, userID)
	if i == -1 {
		return nil, ErrUserNotAssigned
	}

	possibleReviewers := s.assigneReviewers()
	if possibleReviewers == nil {
		return nil, ErrNoCandidate
	}

	// change that to random
	pr.AssignedReviewers[i] = possibleReviewers[0]

	err = s.pullReqsRepo.UpdateReviewers(ctx, pr)
	if err != nil {
		return nil, err
	}

	return pr, nil
}
