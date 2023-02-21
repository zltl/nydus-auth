package db

import (
	"github.com/sirupsen/logrus"
	"github.com/zltl/nydus-auth/internal/pkg/m"
)

// query m.AuthClient by client_id
func (s *State) GetAuthClient(clientId string) (*m.AuthClient, error) {
	var client m.AuthClient
	err := s.QueryRow("SELECT name, secret from auth.client WHERE id=$1",
		clientId).
		Scan(&client.Name, &client.Secret)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	client.Id = clientId

	// get redirect uri list
	rows, err := s.Query("SELECT redirect_uri FROM auth.client_redirect_uri WHERE client_id=$1",
		clientId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var uri string
		err = rows.Scan(&uri)
		if err != nil {
			return nil, err
		}
		client.RedirectUri = append(client.RedirectUri, uri)
	}

	// get scope list
	rows, err = s.Query("SELECT scope FROM auth.client_scope WHERE client_id=$1",
		clientId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var scope string
		err = rows.Scan(&scope)
		if err != nil {
			return nil, err
		}
		client.Scopes = append(client.Scopes, scope)
	}

	return &client, err
}

func (s *State) StoreCode(code *m.AuthCode) error {
	_, err := s.Exec(`INSERT INTO auth.code (client_id, user_id, scope, expire_at, code) 
	values ($1, $2, $3, $4, $5)`,
		code.ClientID,
		code.UserID,
		code.Scope,
		code.ExpireAt,
		code.Code)
	return err
}

func (s *State) QueryCode(clientid, codestr string) (code *m.AuthCode, err error) {
	code = &m.AuthCode{}
	err = s.QueryRow(`SELECT client_id, user_id, scope, expire_at, code FROM auth.code 
	WHERE code=$1 and client_id=$2 limit 1`,
		codestr, clientid).
		Scan(&code.ClientID, &code.UserID, &code.Scope, &code.ExpireAt, &code.Code)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return code, nil
}

func (s *State) SaveRefreshToken(ref *m.RefreshToken) error {
	_, err := s.Exec(`INSERT INTO auth.refresh_token 
	(client_id, user_id, scope, expire_at, refresh_token) 
	values ($1, $2, $3, $4, $5)`,
		ref.ClientID,
		ref.UserID,
		ref.Scope,
		ref.ExpireAt,
		ref.RefreshToken)
	return err
}

func (s *State) QueryRefreshToken(clientId, refreshToken string) (ref []m.RefreshToken, err error) {
	rows, err := s.Query(`SELECT client_id, user_id, scope, expire_at, refresh_token 
	FROM auth.refresh_token 
	WHERE refresh_token=$1 and client_id=$2`,
		refreshToken, clientId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var r m.RefreshToken
		err = rows.Scan(&r.ClientID, &r.UserID, &r.Scope, &r.ExpireAt, &r.RefreshToken)
		if err != nil {
			return nil, err
		}
		ref = append(ref, r)
	}
	return ref, nil
}
