package db

import "github.com/zltl/nydus-auth/pkg/id"

func (s *State) AddUser(i id.ID, username, bcryptPassword string) error {
	_, err := s.Exec("INSERT INTO auth.user (id, name, password) VALUES ($1, $2, $3)",
		i, username, bcryptPassword)
	return err
}
