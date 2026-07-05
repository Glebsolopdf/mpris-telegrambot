package main

import (
	"fmt"
	"path/filepath"
	"time"

	"mpris-tg-status/internal/config"
	"mpris-tg-status/internal/daemon"
)

func run(args []string) error {
	command := "run"
	if len(args) > 1 {
		command = args[1]
	}
	if command == "help" || command == "-h" || command == "--help" {
		printHelp(filepath.Base(args[0]))
		return nil
	}
	if command == "run" || command == "serve" {
		return runService()
	}

	manager, err := newManager()
	if err != nil {
		return err
	}
	switch command {
	case "up", "start":
		options := commandOptions()
		if err := configureAlias(options.UseMSAlias, manager.Executable); err != nil {
			return err
		}
		delay := askStartupDelay(options.StartupDelay, options.PromptTimeout)
		return printResult(manager.UpWithDelay(delay))
	case "down", "stop":
		options := commandOptions()
		restore := askRestoreProfile(options.PromptTimeout)
		return printResult(manager.DownWithRestore(restore))
	case "restart":
		options := commandOptions()
		if err := configureAlias(options.UseMSAlias, manager.Executable); err != nil {
			return err
		}
		restore := askRestoreProfile(options.PromptTimeout)
		delay := askStartupDelay(options.StartupDelay, options.PromptTimeout)
		return printResult(manager.RestartWithOptions(restore, delay))
	case "delete":
		if !askDeleteData(10 * time.Second) {
			fmt.Println("delete cancelled")
			return nil
		}
		return printResult(manager.DeleteData())
	case "status":
		fmt.Println(manager.Status())
		return nil
	case "log", "logs":
		fmt.Println(manager.Paths.Log)
		return nil
	case "help", "-h", "--help":
		printHelp(filepath.Base(args[0]))
		return nil
	default:
		return usage(filepath.Base(args[0]))
	}
}

type commandConfig struct {
	StartupDelay  time.Duration
	PromptTimeout time.Duration
	UseMSAlias    bool
}

func commandOptions() commandConfig {
	cfg, err := config.Load()
	if err != nil {
		return commandConfig{PromptTimeout: 10 * time.Second}
	}
	return commandConfig{
		StartupDelay:  cfg.StartupDelay,
		PromptTimeout: cfg.PromptTimeout,
		UseMSAlias:    cfg.UseMSAlias,
	}
}

func configureAlias(enabled bool, executable string) error {
	message, err := ensureMSAlias(enabled, executable)
	if err != nil {
		return err
	}
	if message != "" {
		fmt.Println(message)
	}
	return nil
}

func newManager() (daemon.Manager, error) {
	dir, err := config.ExecutableDir()
	if err != nil {
		return daemon.Manager{}, err
	}
	executable := currentExecutable()
	return daemon.Manager{Executable: executable, Paths: daemon.NewPaths(dir)}, nil
}

func printResult(message string, err error) error {
	if err != nil {
		return err
	}
	fmt.Println(message)
	return nil
}

func usage(name string) error {
	return fmt.Errorf("unknown command, run %s help", name)
}
