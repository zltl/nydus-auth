package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type State struct {
	db *sql.DB
}

var Ctx = State{}

func (s *State) Open(url string) error {
	if url == "" {
		url = fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
			viper.GetString("db.username"),
			viper.Get("db.password"),
			viper.GetString("db.host"),
			viper.GetInt("db.port"),
			viper.GetString("db.dbname"))
	}
	log.Tracef("openning db: %s", url)

	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Error(err)
		return err
	}
	s.db = db
	return nil
}

func (s *State) Query(query string, args ...any) (*sql.Rows, error) {
	log.Debugf("query: %s, args: %v", query, args)
	r, err := s.db.Query(query, args...)
	if err != nil {
		log.Error(err)
	}
	return r, err
}

func (s *State) QueryRow(query string, args ...any) *sql.Row {
	log.Debugf("query row: %s, args: %v", query, args)
	r := s.db.QueryRow(query, args...)
	return r
}

func (s *State) Exec(query string, args ...any) (sql.Result, error) {
	log.Debugf("exec: %s, args: %v", query, args)
	r, err := s.db.Exec(query, args...)
	if err != nil {
		log.Error(err)
	}
	return r, err
}
