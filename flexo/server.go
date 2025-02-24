/*
Copyright © 2021 Bren 'fraq' Briggs (code@fraq.io)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package flexo

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/SECCDC/flexo/util"
)

type Config struct {
	DBUser string
	DBPass string
	DBAddr string
	DBName string
	DBssl  string
	Secret string
}

type Server struct {
	Router *gin.Engine
	DB     *gorm.DB
}

func Migrate(c Config) error {
	fmt.Println("Running migrations...")
	err := util.DBinit(c.DBUser, c.DBPass, c.DBAddr, c.DBName, c.DBssl)
	return err
}

func Run(c Config) {
	fmt.Println("Starting Flexo...")
	fmt.Println("Hello from version v0.1.19.4")
	s := Server{
		Router: gin.New(),
		DB:     util.DBconnect(c.DBUser, c.DBPass, c.DBAddr, c.DBName, c.DBssl),
	}

	s.Router.Use(
		gin.LoggerWithWriter(gin.DefaultWriter, "/healthz"),
		gin.Recovery(),
		cors.New(cors.Config{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{"authorization", "origin", "content-type", "accept"},
		}),
	)

	migrateErr := Migrate(c)
	if migrateErr != nil {
		fmt.Println(migrateErr)
	} else {
		fmt.Println("Migrations completed successfully!")
	}

	authorized := s.Router.Group("/")
	authorized.Use(util.SecretProvided(c.Secret))
	{
		s.Router.GET("/targets", s.getTargets)
		s.Router.POST("/target", s.postTarget)

		s.Router.GET("/teams", s.getTeams)
		s.Router.POST("/team", s.postTeam)

		s.Router.GET("/categories", s.getCategories)
		s.Router.POST("/category", s.postCategory)

		s.Router.GET("/events", s.getEvents)
		s.Router.POST("/event", s.postEvent)

		s.Router.GET("/ecom", s.getEcomEvents)
		s.Router.POST("/ecom", s.postEcomEvent)

		s.Router.GET("/report/team/:ID", s.teamReport)
	}

	s.Router.GET("/healthz", s.healthCheck)
	s.Router.Run()
	fmt.Println("Hello from new version! (possible v0.1.19.3)")

	defer fmt.Println("Goodbye!")
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
