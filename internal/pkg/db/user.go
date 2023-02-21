package db

import (
	"github.com/sirupsen/logrus"
	"github.com/zltl/nydus-auth/internal/pkg/m"
	"github.com/zltl/nydus-auth/pkg/id"
)

func (s *State) AddUser(i id.ID, username, bcryptPassword string) error {
	_, err := s.Exec("INSERT INTO auth.user (id, name, password) VALUES ($1, $2, $3)",
		i, username, bcryptPassword)
	return err
}

func (s *State) GetUserPassword(username string) (u m.UserPassword, err error) {
	err = s.QueryRow("SELECT id, name, password FROM auth.user WHERE name=$1",
		username).
		Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		logrus.Error(err)
	}
	return u, err
}
