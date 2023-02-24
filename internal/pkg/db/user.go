package db

import (
	"database/sql"

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

func (s *State) GetUserInfo(id id.ID) (u m.UserInfo, err error) {
	var email sql.NullString
	err = s.QueryRow("SELECT name, email FROM auth.user WHERE id=$1",
		id).
		Scan(&u.Name, &email)
	if err != nil {
		logrus.Error(err)
	}
	u.Email = email.String
	u.Sub = id.String()
	return u, err
}
