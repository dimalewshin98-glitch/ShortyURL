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

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBRepository struct {
	dbDsn        string
	dbConnection *sql.DB
	urls         map[string]string
}

func NewDBRepository(dbDsn string) (*DBRepository, error) {
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
	err = dbRepository.Ping()
	if err != nil {
		return nil, err
	}
	tables, err := dbRepository.GetTables()
	if err != nil {
		return nil, err
	}
	if !slices.Contains(tables, "urls") {
		err = dbRepository.CreateTables()
		if err != nil {
			return nil, err
		}
	}
	return dbRepository, nil
}

func (r *DBRepository) CreateTables() error {
	sqlBytes, err := os.ReadFile("../../migrations/000001_create_urls_table.up.sql")
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = r.dbConnection.ExecContext(ctx, string(sqlBytes))
	if err != nil {
		return err
	}
	return nil
}

func (r *DBRepository) GetTables() ([]string, error) {
	sqlReq := "SELECT table_name FROM information_schema.tables WHERE table_schema='public';"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	rows, err := r.dbConnection.QueryContext(ctx, string(sqlReq))
	if err != nil {
		return nil, err
	}
	var tableName string
	var tables []string
	for rows.Next() {
		rows.Scan(&tableName)
		tables = append(tables, tableName)
	}
	return tables, nil
}

func (r *DBRepository) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := r.dbConnection.PingContext(ctx)
	return err
}

func (r *DBRepository) Store(urlID string, URL string) (string, error) {
	r.urls[urlID] = URL
	return urlID, nil
}

func (r *DBRepository) Get(urlID string) (string, error) {
	URL := r.urls[urlID]
	return URL, nil
}
