package api

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/zltl/nydus-auth/internal/pkg/db"
	"github.com/zltl/nydus-auth/internal/pkg/m"
	"golang.org/x/crypto/bcrypt"
)

// oauth2 authorize
func (s *State) handleGetAuthAuthorize(c *gin.Context) {
	responseType := c.Query("response_type")
	clientId := c.Query("client_id")
	redirectUri := c.Query("redirect_uri")
	state := c.Query("state")
	scope := c.Query("scope")
	if responseType != "code" {
		c.HTML(http.StatusBadRequest, "error", gin.H{
			"error": "response_type must be 'code'",
		})
		return
	}

	client, err := db.Ctx.GetAuthClient(clientId)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error", gin.H{
			"error": "client not found",
		})
		return
	}

	// check redirect_uri
	ok := checkRedirectUri(client, redirectUri)
	if !ok {
		c.HTML(http.StatusBadRequest, "error", gin.H{
			"error": "redirect_uri not allow",
		})
		return
	}

	// check if element of scope is in client.Scopes
	if !checkScopeList(client.Scopes, scope) {
		c.HTML(http.StatusBadRequest, "error", gin.H{
			"error": "scope not allow",
		})
		logrus.Errorf("scope not allow: got %+v, expected %+v", scope, client.Scopes)
		return
	}

	thisUri := s.externFullUri("/auth/authorize")
	csrfToken, err := s.genCsrfJwtToken([]string{thisUri}, []string{c.ClientIP()})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error", gin.H{
			"error": "could not generate csrf token:" + err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "/auth/authorize", gin.H{
		"client_id":    client.Id,
		"redirect_uri": redirectUri,
		"state":        state,
		"scope":        scope,
		"csrf_token":   csrfToken,
	})
}

func (s *State) handlePostAuthAuthorize(c *gin.Context) {
	var body m.AuthorizeLoginReq
	err := c.Bind(&body)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error", gin.H{
			"error": "invalid request body",
		})
		return
	}

	// check csrf token
	_, err = s.validateCsrfJwtToken(body.CSRFToken,
		s.externFullUri("/auth/authorize"),
		c.ClientIP())
	if err != nil {
		c.HTML(http.StatusBadRequest, "error", gin.H{
			"error": "invalid csrf token",
		})
		return
	}

	client, err := db.Ctx.GetAuthClient(body.ClientId)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error", gin.H{
			"error": "client not found",
		})
		return
	}

	// verify redirect_uri
	ok := checkRedirectUri(client, body.RedirectURI)
	if !ok {
		c.HTML(http.StatusBadRequest, "error", gin.H{
			"error": "redirect_uri not allow",
		})
		return
	}

	// verify scopes
	if !checkScopeList(client.Scopes, body.Scope) {
		c.HTML(http.StatusBadRequest, "error", gin.H{
			"error": "scope not allow",
		})
		logrus.Errorf("scope not allow: got %+v, expected %+v", body.Scope, client.Scopes)
		return
	}

	// check username and password
	u, err := db.Ctx.GetUserPassword(body.Username)
	if err != nil {
		c.HTML(http.StatusBadRequest, "/auth/authorize", gin.H{
			"error":        "invalid username or password",
			"client_id":    body.ClientId,
			"redirect_uri": body.RedirectURI,
			"state":        body.State,
			"scope":        body.Scope,
			"csrf_token":   body.CSRFToken,
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(body.Password))
	if err != nil {
		logrus.Error(err)
		c.HTML(http.StatusBadRequest, "/auth/authorize", gin.H{
			"error":        "invalid username or password",
			"client_id":    body.ClientId,
			"redirect_uri": body.RedirectURI,
			"state":        body.State,
			"scope":        body.Scope,
			"csrf_token":   body.CSRFToken,
		})
		return
	}

	// generate string length 64
	code := make([]byte, 32)
	_, err = rand.Read(code)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error", gin.H{
			"error": "could not generate code",
		})
		return
	}
	codeStr := hex.EncodeToString(code)
	var co m.AuthCode
	co.ClientID = body.ClientId
	co.Code = codeStr
	co.Scope = body.Scope
	co.UserID = u.ID
	co.ExpireAt = time.Now().Add(time.Hour)
	err = db.Ctx.StoreCode(&co)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error", gin.H{
			"error": "could not store code",
		})
		return
	}

	// redirect to redirect_uri
	redirectUri, err := url.Parse(body.RedirectURI)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error", gin.H{
			"error": "could not parse redirect_uri",
		})
		return
	}
	q := redirectUri.Query()
	q.Set("code", codeStr)
	q.Set("state", body.State)
	redirectUri.RawQuery = q.Encode()
	c.Redirect(http.StatusFound, redirectUri.String())
}

func checkRedirectUri(client *m.AuthClient, redirectUriReq string) bool {
	// check redirect_uri
	ok := false
	for _, uri := range client.RedirectUri {
		if strings.HasPrefix(redirectUriReq, uri) {
			ok = true
			break
		}
	}
	return ok
}

func checkScopeList(clientScopes []string, scopes string) (ok bool) {
	scopeList := strings.Split(scopes, ",")
	ok = true
	// if every element in scopes are in client.Scopes, ok = true
	for _, scopeReq := range scopeList {
		oneOk := false
		for _, s := range clientScopes {
			if strings.HasPrefix(scopeReq, s) {
				oneOk = true
				break
			}
		}
		if !oneOk {
			ok = false
			break
		}
	}
	return ok
}
