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
	txMap        map[string]*sql.Tx
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
		txMap:        make(map[string]*sql.Tx)}
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
		sqlReqBytes, err = os.ReadFile("../../migrations/000002_add_unique_index_to_orig_url.up.sql")
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

func (r *DBRepository) Store(ctx context.Context, userID string, urlID string, URL string) (string, error) {
	var ctxUUID string
	var singleReq bool
	var isLastReq bool
	var err error
	var tx *sql.Tx
	if ctx.Value("UUID") == nil {
		ctxUUID = ""
		singleReq = true
		isLastReq = true
	} else {
		ctxUUID = ctx.Value("UUID").(string)
		singleReq = false
		if ctx.Value("isLastReq") == nil {
			isLastReq = false
		} else {
			isLastReq = ctx.Value("isLastReq").(bool)
		}
	}
	if singleReq {
		tx = nil
	} else {
		tx = r.txMap[ctxUUID]
	}
	if tx == nil {
		tx, err = r.dbConnection.BeginTx(ctx, nil)
		if err != nil {
			return urlID, err
		}
	}
	if !singleReq {
		r.txMap[ctxUUID] = tx
	}
	select {
	case <-ctx.Done():
		rbErr := tx.Rollback()
		if !singleReq {
			delete(r.txMap, ctxUUID)
		}
		return urlID, rbErr
	default:
		sqlInsert := "INSERT INTO urls (uuid, short_url, original_url) VALUES ($1, $2, $3) ON CONFLICT (original_url) DO NOTHING RETURNING uuid"
		var UUID string
		isExistsURL := false
		row := tx.QueryRowContext(ctx, sqlInsert, uuid.New().String(), urlID, URL)
		err = row.Scan(&UUID)
		if err != nil {
			if err == sql.ErrNoRows {
				isExistsURL = true
			} else {
				if rbErr := tx.Rollback(); rbErr != nil {
					return urlID, rbErr
				}
				return urlID, err
			}
		}
		if isExistsURL {
			sqlSelect := "SELECT short_url FROM urls WHERE original_url = $1"
			row := r.dbConnection.QueryRowContext(ctx, sqlSelect, URL)
			err = row.Scan(&urlID)
			if err != nil {
				if rbErr := tx.Rollback(); rbErr != nil {
					return urlID, rbErr
				}
				return urlID, err
			}
			err = errors.New("Short URL already exists")
		}
		if isLastReq {
			commitErr := tx.Commit()
			if !singleReq {
				delete(r.txMap, ctxUUID)
			}
			if commitErr == nil {
				return urlID, err
			} else {
				return urlID, commitErr
			}
		}
		return urlID, err
	}
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
