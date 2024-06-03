package db

import (
	"database/sql"
	"fmt"
	"microservices/geo/config"
	"microservices/geo/internal/db/adapter"

	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func NewSqlDB(dbConf config.DB, logger *zap.Logger) (*sqlx.DB, *adapter.SQLAdapter, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password, dbConf.Name)

	var (
		dbRaw *sql.DB
		err   error
	)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	timeoutExceeded := time.After(time.Second + time.Duration(dbConf.Timeout))

	for {
		select {
		case <-timeoutExceeded:
			return nil, nil, fmt.Errorf("db connection failed after %d timeout %s", dbConf.Timeout, err)
		case <-ticker.C:
			dbRaw, err = sql.Open(dbConf.Driver, dsn)
			if err != nil {
				return nil, nil, err
			}
			err = dbRaw.Ping()
			if err == nil {
				db := sqlx.NewDb(dbRaw, dbConf.Driver)
				db.SetMaxOpenConns(50)
				db.SetMaxIdleConns(50)
				sqlAdapter := adapter.NewSQLAdapter(db)
				return db, sqlAdapter, nil
			}
			logger.Error("failed to connect to the database", zap.String("dsn", dsn), zap.Error(err))
		}
	}
}
