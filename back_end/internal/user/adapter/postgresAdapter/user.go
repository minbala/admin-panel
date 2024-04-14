package postgresAdapter

import (
	"admin-panel/internal/user/domain"
	"admin-panel/pkg/postgres"
	"admin-panel/pkg/resources"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"sync"
)

type UserRepository struct {
	client *postgres.Queries
	db     *pgxpool.Pool
}

var userRepository domain.UserRepositoryInterface
var setUpOnceUserRepository sync.Once

func ProvideUserRepositoryInterface(client *postgres.Queries, db *pgxpool.Pool) domain.UserRepositoryInterface {
	setUpOnceUserRepository.Do(func() {
		userRepository = &UserRepository{
			client: client,
			db:     db,
		}
	})
	return userRepository
}

func (b UserRepository) Create(ctx context.Context, user domain.User) (int64, error) {
	return b.client.CreateUser(ctx, postgres.CreateUserParams{
		Email:    user.Email,
		Name:     user.Name,
		Password: user.Password,
		UserRole: postgres.UserRoleType(user.UserRole),
	})
}

func (b UserRepository) Delete(ctx context.Context, userID int64) error {
	return b.client.DeleteUser(ctx, userID)
}

func (b UserRepository) Update(ctx context.Context, user domain.User) error {
	return b.client.UpdateUser(ctx, postgres.UpdateUserParams{
		Email:    user.Email,
		Name:     user.Name,
		UserRole: postgres.UserRoleType(user.UserRole),
		Password: user.Password,
		ID:       user.ID,
	})
}

func (b UserRepository) Get(ctx context.Context, query domain.UserQuery) ([]domain.User, error) {
	users, err := b.client.GetUsers(ctx, postgres.GetUsersParams{
		ID:       query.ID,
		Email:    query.Email,
		UserRole: postgres.UserRoleType(postgres.ContainQuery(string(query.UserRole))),
		Name:     postgres.ContainQuery(query.Name),
		Limit:    query.Limit,
		Offset:   query.Offset,
	})
	return unmarshalUsersFromDB(users), err
}

func (b UserRepository) Count(ctx context.Context, query domain.UserQuery) (int64, error) {
	return b.client.GetUsersCount(ctx, postgres.GetUsersCountParams{
		ID:       query.ID,
		Email:    query.Email,
		UserRole: postgres.UserRoleType(postgres.ContainQuery(string(query.UserRole))),
		Name:     postgres.ContainQuery(query.Name),
	})

}

func unmarshalUsersFromDB(input []postgres.Users) []domain.User {
	response := make([]domain.User, len(input))
	for i := range input {
		response[i] = domain.User{
			ID:        input[i].ID,
			CreatedAt: input[i].CreatedAt.Time,
			UpdatedAt: input[i].UpdatedAt.Time,
			Email:     input[i].Email,
			Name:      input[i].Name,
			Password:  input[i].Password,
			UserRole:  resources.UserRoleType(input[i].UserRole),
		}

	}
	return response
}
