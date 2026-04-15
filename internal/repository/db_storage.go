package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBRepository struct {
	dbDsn        string
	dbConnection *sql.DB
	urls         map[string]string
}

func NewDBRepository(dbDsn string) (*DBRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	hostPort := strings.Split(dbDsn, ":")
	var ps string
	if len(hostPort) == 1 {
		ps = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
			hostPort[0], `postgres`, `admin`, `postgres`)
	} else if len(hostPort) == 2 {
		ps = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			hostPort[0], hostPort[1], `postgres`, `admin`, `postgres`)
	} else {
		return nil, errors.New("Database host/port wrong format")
	}
	db, err := sql.Open("pgx", ps)
	if err != nil {
		return nil, err
	}
	dbRepository := &DBRepository{
		dbDsn:        dbDsn,
		dbConnection: db,
		urls:         make(map[string]string)}
	err = dbRepository.Ping(ctx)
	if err != nil {
		return nil, err
	}
	err = dbRepository.CreateTables(ctx)
	if err != nil {
		return nil, err
	}
	return dbRepository, nil
}

func (r *DBRepository) CreateTables(ctx context.Context) error {
	tx, err := r.dbConnection.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	sqlReq := "SELECT table_name FROM information_schema.tables WHERE table_schema='public';"
	rows, err := tx.QueryContext(ctx, string(sqlReq))
	if err != nil {
		return nil
	}
	var tableName string
	var tables []string
	for rows.Next() {
		rows.Scan(&tableName)
		tables = append(tables, tableName)
	}
	if !slices.Contains(tables, "urls") {
		sqlReqBytes, err := os.ReadFile("../../migrations/000001_create_urls_table.up.sql")
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, string(sqlReqBytes))
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return rbErr
			}
			return err
		}
	}
	return tx.Commit()
}

func (r *DBRepository) Ping(ctx context.Context) error {
	err := r.dbConnection.PingContext(ctx)
	return err
}

func (r *DBRepository) Store(ctx context.Context, urlID string, URL string) (string, error) {
	tx, err := r.dbConnection.BeginTx(ctx, nil)
	if err != nil {
		return urlID, err
	}
	UUID := uuid.New().String()
	sqlInsert := "INSERT INTO urls (uuid, short_url, original_url) VALUES ($1, $2, $3)"
	_, err = tx.ExecContext(ctx, sqlInsert, UUID, urlID, URL)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return urlID, rbErr
		}
		return urlID, err
	}
	return urlID, tx.Commit()
}

func (r *DBRepository) Get(ctx context.Context, urlID string) (string, error) {
	sqlSelect := "SELECT original_url FROM urls WHERE short_url = $1"
	var originalURL string
	row := r.dbConnection.QueryRowContext(ctx, sqlSelect, urlID)
	err := row.Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return originalURL, nil
}
