package service

import (
	"context"
	"errors"
	"log/slog"
	"math/rand/v2"
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

func (s *PullReqService) CreateTeam(name string, members []*domain.User) (*domain.Team, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	team := &domain.Team{
		Name: name,
	}

	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		if err := s.teamsRepo.Insert(ctx, team); err != nil {
			return err
		}

		for _, member := range members {
			member.TeamName = team.Name

			err := s.usersRepo.CreateUser(ctx, member)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return team, nil
}

func (s *PullReqService) GetTeam(name string) (*domain.Team, []string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	team, err := s.teamsRepo.Get(ctx, name)
	if err != nil {
		return nil, nil, err
	}

	members, err := s.usersRepo.GetAllForTeam(ctx, team.Name)
	if err != nil {
		return nil, nil, err
	}

	return team, members, nil
}

func (s *PullReqService) MergePullReq(id string) (*domain.PR, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	pr := &domain.PR{
		ID:     id,
		Status: "MERGED",
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	reviewers, err := s.assigneReviewers(ctx, authorID)
	if err != nil {
		return nil, err
	}
	pr := &domain.PR{
		ID:                id,
		Name:              name,
		AuthorID:          authorID,
		Status:            "OPEN",
		AssignedReviewers: reviewers,
	}

	err = s.pullReqsRepo.Insert(ctx, pr)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

// i need to query all active users in the same team as author and pick 2
// if there are less than 2 active users in the team, assigne what i have (even 0)
func (s *PullReqService) assigneReviewers(ctx context.Context, authorId string) ([]string, error) {
	user, err := s.usersRepo.Get(ctx, authorId)
	if err != nil {
		return nil, err
	}

	members, err := s.usersRepo.GetAllForTeam(ctx, user.TeamName)
	if err != nil {
		return nil, err
	}

	if len(members) < 2 {
		return members, nil
	}

	rand.Shuffle(len(members), func(i, j int) {
		members[i], members[j] = members[j], members[i]
	})

	return members[:2], nil
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

	possibleReviewers, err := s.assigneReviewers(ctx, pr.AuthorID)
	if err != nil {
		return nil, err
	}

	if possibleReviewers == nil {
		return nil, ErrNoCandidate
	}

	rand.Shuffle(len(possibleReviewers), func(i, j int) {
		possibleReviewers[i], possibleReviewers[j] = possibleReviewers[j], possibleReviewers[i]
	})
	pr.AssignedReviewers[i] = possibleReviewers[0]

	err = s.pullReqsRepo.UpdateReviewers(ctx, pr)
	if err != nil {
		return nil, err
	}

	return pr, nil
}
