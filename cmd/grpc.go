/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/donech/pearl/internal"
	"github.com/donech/pearl/internal/common"

	"github.com/donech/pearl/pkg/log"

	"github.com/spf13/cobra"
)

// grpcCmd represents the grpc command
var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "erp grpc api service",
	Long:  `erp grpc api service`,
	Run: func(cmd *cobra.Command, args []string) {
		common.InitGlobal()
		entry := internal.NewGrpcEntry(common.C().Grpc)
		if err := entry.Run(); err != nil {
			log.S(context.Background()).Fatalf("entry run failed, %#v", err)
		}

		signalCh := make(chan os.Signal, 4)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	WAIT:
		s := <-signalCh
		if s == syscall.SIGHUP {
			//todo implement
			log.S(context.Background()).Info("reload config")
			goto WAIT
		}
		log.S(context.Background()).Info("Shutdown Server ...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := entry.Stop(ctx); err != nil {
			log.S(context.Background()).Fatal("Entry stop err :", err)
		}
		log.S(context.Background()).Info("Entry exiting")
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(grpcCmd)
}
