package postgresAdapter

import (
	"admin-panel/internal/user/domain"
	"admin-panel/pkg/postgres"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"sync"
)

type UserLogRepository struct {
	client *postgres.Queries
	db     *pgxpool.Pool
}

var userLogRepository domain.UserLogRepositoryInterface
var setUpOnceUserLogRepository sync.Once

func ProvideUserLogRepositoryInterface(client *postgres.Queries, db *pgxpool.Pool) domain.UserLogRepositoryInterface {
	setUpOnceUserLogRepository.Do(func() {
		userLogRepository = &UserLogRepository{
			client: client,
			db:     db,
		}
	})
	return userLogRepository
}

func (b UserLogRepository) Create(ctx context.Context, userLog domain.UserLog) error {
	return b.client.CreateUserLog(ctx, postgres.CreateUserLogParams{
		UserID:       postgres.Int64ToPgxInt(userLog.UserID),
		Event:        userLog.Event,
		RequestUrl:   userLog.RequestUrl,
		Data:         userLog.Data,
		Status:       userLog.Status,
		ErrorMessage: postgres.StringToNullString(userLog.ErrorMessage),
	})
}

func (b UserLogRepository) Get(ctx context.Context, query domain.UserLogQuery) ([]domain.UserLog, error) {
	logs, err := b.client.GetUserLogs(ctx, postgres.GetUserLogsParams{
		UserID: postgres.Int64ToPgxInt(query.UserID),
		Limit:  query.Limit,
		Offset: query.Offset,
	})
	return unmarshalUserLogFromDB(logs), err
}

func (b UserLogRepository) Count(ctx context.Context, query domain.UserLogQuery) (int64, error) {
	return b.client.GetUserLogsCount(ctx, postgres.Int64ToPgxInt(query.UserID))
}

func unmarshalUserLogFromDB(input []postgres.UserLogs) []domain.UserLog {
	response := make([]domain.UserLog, len(input))
	for i := range input {
		response[i] = domain.UserLog{
			ID:           input[i].ID,
			CreatedAt:    input[i].CreatedAt.Time,
			UpdatedAt:    input[i].UpdatedAt.Time,
			UserID:       input[i].UserID.Int64,
			Event:        input[i].Event,
			RequestUrl:   input[i].RequestUrl,
			Data:         input[i].Data,
			Status:       input[i].Status,
			ErrorMessage: input[i].ErrorMessage.String,
		}

	}
	return response
}
