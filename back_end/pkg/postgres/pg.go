package postgres

import (
	"admin-panel/pkg/config"
	"context"
	_ "embed"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log"
	"sync"
)

//
//go:embed migration/schema.sql

var ddl string

var client *Queries
var db *pgxpool.Pool
var setDbOnce sync.Once

func ProvideClient(config *config.Config) *Queries {
	setDbOnce.Do(func() {
		connectPostgres(config)
	})
	return client
}

func connectPostgres(config *config.Config) {
	// urlExample := "postgres_adapter://username:password@localhost:5432/database_name"
	conn, err := pgxpool.New(context.Background(), config.Database.DBURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	client = New(conn)
	db = conn
	log.Println("connected to postgres")
	log.Println(db.Ping(context.Background()))
}

func ProvideSqlx(config *config.Config) *pgxpool.Pool {
	setDbOnce.Do(func() {
		connectPostgres(config)
	})
	return db
}

func ProvideDBWithInitialData(dbClient *Queries, db *pgxpool.Pool, hash func(password string) (string, error)) {
	ctx := context.Background()
	_, err := db.Exec(ctx, ddl)
	if err != nil {
		log.Fatal(err)
	}
	fakePassword, err := hash("minbala33")
	if err != nil {
		log.Fatal(err)
	}
	_, err = dbClient.CreateUser(ctx, CreateUserParams{
		Email:    "minbala33@gmail.com",
		Name:     "minbala",
		Password: fakePassword,
		UserRole: "admin",
	})
	if err != nil {
		log.Fatal(err)
	}

}

func WithTx(ctx context.Context, db *pgxpool.Pool, client *Queries, txFn func(*Queries) error) error {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	q := client.WithTx(tx)
	err = txFn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			err = fmt.Errorf("tx failed: %v, unable to rollback: %v", err, rbErr)
		}
	} else {
		err = tx.Commit(ctx)
	}
	return err
}
