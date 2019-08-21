package cmd

import (
	"context"
	"fmt"
	"github.com/s4kibs4mi/movie-pie/api"
	"github.com/s4kibs4mi/movie-pie/config"
	"github.com/s4kibs4mi/movie-pie/log"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve starts http and gRPC server",
	Run:   serve,
}

func serve(cmd *cobra.Command, args []string) {
	addr := fmt.Sprintf("%s:%d", config.App().Base, config.App().Port)

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	hServer := http.Server{
		Addr:    addr,
		Handler: api.Router(),
	}

	go func() {
		log.Log().Infoln("Http server has been started on", addr)
		if err := hServer.ListenAndServe(); err != nil {
			log.Log().Errorln("Failed to start http server on :", err)
			os.Exit(-1)
		}
	}()

	<-stop

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	hServer.Shutdown(ctx)

	log.Log().Infoln("Http server has been shutdown gracefully")
}
