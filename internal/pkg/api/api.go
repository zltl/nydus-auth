package api

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zltl/nydus-auth/pkg/id"
)

type State struct {
	idGenerator   *id.Instance
	csrfJwtSecret []byte
	csrfJwtAlgo   string
	csrfJwtTTL    int
	csrfJwtKeyId  string

	sessionJwtSecret        []byte
	sessionJwtAlgo          string
	sessionJwtTTL           int
	sessionJwtKeyId         string
	sessionExpire           int
	sessionJwtRefreshExpire int

	externUri string
}

func (s *State) Start() {
	s.idGenerator = id.NewInstance(viper.GetInt("server.id"))
	s.csrfJwtAlgo = viper.GetString("server.csrf_jwt_algo")
	s.csrfJwtSecret = []byte(viper.GetString("server.csrf_jwt_secret"))
	s.csrfJwtTTL = viper.GetInt("server.csrf_jwt_ttl")
	s.csrfJwtKeyId = viper.GetString("server.csrf_jwt_kid")
	s.sessionExpire = viper.GetInt("server.session_expire")
	s.sessionJwtAlgo = viper.GetString("server.session_jwt_algo")
	s.sessionJwtSecret = []byte(viper.GetString("server.session_jwt_secret"))
	s.sessionJwtTTL = viper.GetInt("server.session_jwt_ttl")
	s.sessionJwtKeyId = viper.GetString("server.session_jwt_kid")
	s.sessionJwtRefreshExpire = viper.GetInt("server.session_jwt_refresh_expire")

	s.externUri = viper.GetString("server.extern_uri")
	if s.externUri == "" {
		logrus.Panicf("server.extern_uri is empty")
	}
	if s.externUri[len(s.externUri)-1] == '/' {
		s.externUri = s.externUri[:len(s.externUri)-1]
	}

	tmpl, err := loadTmpl(viper.GetString("server.tmpl"))
	if err != nil {
		logrus.Panic(err)
	}

	r := gin.Default()
	r.SetHTMLTemplate(tmpl)

	// serve static files in ./static
	r.Static("/auth/static", viper.GetString("server.static"))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello",
		})
	})

	// signup pages
	r.GET("/auth/signup", s.handleGetAuthSignup)
	r.POST("/auth/signup", s.handlePostAuthSignup)

	// OAuth2 authorize
	r.GET("/auth/authorize", s.handleGetAuthAuthorize)
	r.POST("/auth/authorize", s.handlePostAuthAuthorize)

	// r.POST("/auth/token")

	r.POST("/auth/token", s.handleAuthToken)
	r.GET("/auth/userinfo", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "oauth2/userinfo",
		})
	})
	r.GET("/auth/keys", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "oauth2/keys",
		})
	})

	listenAddr := viper.GetString("server.addr") + ":" + viper.GetString("server.port")
	r.Run(listenAddr)
}

func loadTmpl(dir string) (*template.Template, error) {
	t := template.New("")
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".tmpl" {
			return nil
		}
		_, err = t.ParseFiles(path)
		return err
	})
	return t, err
}

func (s *State) externFullUri(path string) string {
	return s.externUri + path
}
