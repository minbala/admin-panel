package commands

import (
	"admin-panel/docs"
	userModule "admin-panel/internal/user/port/public/http"
	commonImport "admin-panel/pkg/common"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type ServeCommand struct {
	container *commonImport.Container
}

func (s *ServeCommand) RunE(cmd *cobra.Command, args []string) error {
	f, err := os.OpenFile(s.container.Config.App.LogFilePath+"/"+s.container.Config.App.LogFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	gin.DefaultWriter = io.MultiWriter(f)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	corsPolicy := cors.DefaultConfig()
	corsPolicy.AllowHeaders = append(corsPolicy.AllowHeaders, "Authorization")
	corsPolicy.AllowAllOrigins = true
	r.Use(cors.New(corsPolicy))
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/ping", func(c *gin.Context) {

		c.String(200, "pong")
	})

	//r.Use(middlewares.MaintenanceMode(s.container.Config))
	//r.Use(sessions.Sessions("mysession", s.container.Session.Store))
	////r.Use(csrf.Middleware(csrf.Options{
	////	Secret: "SHsHZ28711587148418",
	////	ErrorFunc: func(c *gin.Context) {
	////		c.String(400, "CSRF token mismatch")
	////		c.Abort()
	////	},
	////}))
	//

	//prometheus.MustRegister()
	//r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	////p := ginprometheus.NewPrometheus("gin")
	////p.Use(r)
	r.Use(s.container.Operation.LimitBodySize(), s.container.Operation.AddDataToLogger())
	userModule.SetupAPI(r, s.container)
	//
	//if s.container.Config.App.Maintenance == false {
	//	_, errScheduler := s.container.Scheduler.Run()
	//	if errScheduler != nil {
	//		s.container.Log.Error(errScheduler.Error())
	//	}
	//}
	//
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%v", s.container.Config.App.APIPORT),
		Handler:           r,
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		MaxHeaderBytes:    8 * 1024, // 8KiB

	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *ServeCommand) NewCommand(container *commonImport.Container) *cobra.Command {
	s.container = container
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the server",
		Long:  "This command starts the server.",
		RunE:  s.RunE,
	}
}
