package postgresAdapter

import (
	"admin-panel/internal/user/domain"
	"admin-panel/pkg/postgres"
	"admin-panel/pkg/resources"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"sync"
)

type SessionRepository struct {
	client *postgres.Queries
	db     *pgxpool.Pool
}

var sessionRepository domain.SessionRepositoryInterface
var setUpOnceSessionRepository sync.Once

func (u SessionRepository) Create(ctx context.Context, session domain.Session) error {
	return u.client.CreateSession(ctx, postgres.CreateSessionParams{
		AccessToken: session.AccessToken,
		UserID:      session.UserID,
	})
}

func (u SessionRepository) Get(ctx context.Context, userID int64, accessToken string) (domain.Session, error) {
	session, err := u.client.GetSession(ctx,
		postgres.GetSessionParams{
			UserID:      userID,
			AccessToken: accessToken,
		})
	if err != nil {
		return domain.Session{}, err
	}
	if session.ID == 0 {
		return domain.Session{}, resources.ErrNoRecordFound
	}
	return convertEmployeeSessionToDomainModel(session), nil
}

func convertEmployeeSessionToDomainModel(session postgres.Sessions) domain.Session {
	return domain.Session{
		ID:          session.ID,
		AccessToken: session.AccessToken,
		UserID:      session.UserID,
		CreatedAt:   session.CreatedAt.Time,
		UpdatedAt:   session.UpdatedAt.Time,
	}
}

func (u SessionRepository) Delete(ctx context.Context, session domain.Session) error {
	return u.client.DeleteSession(ctx, postgres.DeleteSessionParams{
		UserID:      session.UserID,
		AccessToken: session.AccessToken,
	})
}

func ProvideUserSessionRepository(client *postgres.Queries, db *pgxpool.Pool) domain.SessionRepositoryInterface {
	setUpOnceSessionRepository.Do(func() {
		sessionRepository = &SessionRepository{
			client: client,
			db:     db,
		}
	})
	return sessionRepository
}
