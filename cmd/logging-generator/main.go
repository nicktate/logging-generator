package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nicktate/logging-generator/pkg/generator"
)

const (
	leaderElectionResource = "floating-ip-pool-controller-leader-election"
)

var ctx context.Context
var log logrus.FieldLogger

var killPeriod time.Duration

var messageSize int64
var messageTotal int64
var messagePeriod time.Duration

var rootCmd = &cobra.Command{
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if isatty.IsTerminal(os.Stdout.Fd()) {
			logrus.SetFormatter(&logrus.TextFormatter{})
		}
	},
	RunE: runMain,
}

func init() {
	viper.SetEnvPrefix("logging-generator")

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	log = logrus.New()
	ctx = signalContext(context.Background(), log)
}

func main() {
	rootCmd.Flags().DurationVar(&killPeriod, "kill-period", generator.DefaultKillPeriod, "duration until the logging generator exits")
	rootCmd.Flags().Int64Var(&messageSize, "message-size", generator.DefaultMessageSize, "size of each individual message")
	rootCmd.Flags().Int64Var(&messageTotal, "message-total", generator.DefaultMessageTotal, "total number of messages to send")
	rootCmd.Flags().DurationVar(&messagePeriod, "message-period", generator.DefaultMessagePeriod, "time between log output")
	rootCmd.Execute()
}

func signalContext(ctx context.Context, log logrus.FieldLogger) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Info("got interrupt signal; shutting down")
		cancel()
		<-c
		log.Info("got second interrupt signal; unclean shutdown")
		os.Exit(1) // exit hard for the impatient
	}()

	return ctx
}

func runMain(cmd *cobra.Command, args []string) error {
	if time.Duration(messageTotal)*messagePeriod > killPeriod {
		return errors.New("message-total * message-period must be less than kill-period")
	}

	s := generator.NewGenerator(
		generator.WithLog(log),
		generator.WithKillPeriod(killPeriod),
		generator.WithMessageSize(messageSize),
		generator.WithMessageTotal(messageTotal),
		generator.WithMessagePeriod(messagePeriod),
	)

	return s.Run(ctx)
}
