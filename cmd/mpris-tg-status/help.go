package main

import "fmt"

func printHelp(name string) {
	fmt.Printf(`Usage: %[1]s <command>

Commands:
  up, start       Start background service
  down, stop      Stop service and ask whether to restore default profile
  restart         Stop, then start service again
  status          Show service state and pid
  logs, log       Print RAM log file path
  run, serve      Run service in the foreground
  delete          Delete local config, runtime data, aliases, and autostart
  help            Show this help

Examples:
  %[1]s up
  %[1]s status
  %[1]s down
  %[1]s delete
`, name)
}
