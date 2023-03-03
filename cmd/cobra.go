package cmd

import (
	"context"
	"fmt"
	gocron "github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tapd-notice/common/dingding"
	"tapd-notice/common/logger"
	"tapd-notice/common/orm"
	"tapd-notice/common/tapd"
	"tapd-notice/config"
	"tapd-notice/internal/cron"
	"tapd-notice/model/migrate"
	"tapd-notice/router"
	"time"
)

var (
	hostFlag  string
	portFlag  string
	helpFlag  bool
	configYml string
	goCron    *gocron.Cron
	command   = &cobra.Command{
		Use:     "tapd-notice",
		Short:   "Start TAPD-Notice Server",
		Example: "tapd-notice -h 127.0.0.1 -p 8080",
		PreRun: func(cmd *cobra.Command, args []string) {
			config.Init(hostFlag, portFlag, configYml)
			orm.Setup(migrate.WithMigrate)
			logger.Setup()
		},
		Run: func(cmd *cobra.Command, args []string) {
			runCron()
			runAPI()
		},
	}
)

func init() {
	command.PersistentFlags().StringVarP(&hostFlag, "host", "h", "", "Start server with specific host")
	command.PersistentFlags().StringVarP(&portFlag, "port", "p", "", "Start server with specific port")
	command.PersistentFlags().StringVarP(&configYml, "config", "c", "config/settings.yml", "Start server with specific config file")
	command.PersistentFlags().BoolVarP(&helpFlag, "help", "", false, "Help default flag")

	goCron = gocron.New()
}

func Execute() {
	if err := command.Execute(); err != nil {
		os.Exit(-1)
	}
}

func runCron() {
	// dingdingClient
	tapdAgentConfig := config.GetDingdingConfig().TapdAgent
	dingdingClient := dingding.NewClient(tapdAgentConfig.AppKey, tapdAgentConfig.AppSecret, "")
	// tapdClient
	tapdConfig := config.GetTAPDConfig()
	tapdClient := tapd.NewClient(tapdConfig.CompanyId, tapdConfig.ApiUser, tapdConfig.ApiPassword)

	deptCronJob := cron.NewDingdingDeptCronJob(orm.DB, dingdingClient)
	userCronJob := cron.NewDingdingUserCronJob(orm.DB, dingdingClient)

	projectCronJob := cron.NewTAPDProjectCronJob(orm.DB, tapdClient)

	goCron.AddJob("15 0 * * *", deptCronJob)
	goCron.AddJob("30 0 * * *", userCronJob)
	goCron.AddJob("0 0 * * *", projectCronJob)
	goCron.Start()
}

func runAPI() {
	applicationConfig := config.GetApplicationConfig()
	ginRouter := router.NewRouter()
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", applicationConfig.Host, applicationConfig.Port),
		Handler:      ginRouter,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("api server启动失败, err: %s\n", err)
			os.Exit(-1)
		}
	}()
	fmt.Printf("api server listen on %s\n", server.Addr)

	// 优雅关闭
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	fmt.Println("closing http server gracefully ...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalln("closing http server gracefully failed: ", err)
	}
	goCron.Stop()
}
