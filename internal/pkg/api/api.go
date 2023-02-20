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

	sessionJwtSecret []byte
	sessionJwtAlgo   string
	sessionJwtTTL    int
	sessionJwtKeyId  string
	sessionExpire    int
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

	tmpl, err := loadTmpl(viper.GetString("server.tmpl"))
	if err != nil {
		logrus.Panic(err)
	}

	r := gin.Default()
	r.SetHTMLTemplate(tmpl)

	// serve static files in ./static
	r.Static("/auth/static", viper.GetString("server.static"))

	// signup pages
	r.GET("/auth/signup", s.handleGetAuthSignup)
	r.POST("/auth/signup", s.handlePostAuthSignup)

	// Oauth2
	r.GET("/auth/authorize", func(c *gin.Context) {
		// tmpl.Ctx.ExecuteTemplate(c.Writer, "authorize", nil)

		c.JSON(http.StatusOK, gin.H{
			"message": "auth/authorize",
		})
	})
	r.POST("/auth/token", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "auth/token",
		})
	})
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

	r.Run(":8080")
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
