package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func askRestoreProfile(timeout time.Duration) bool {
	return askYesNo(
		"Restore the default profile?",
		true,
		timeout,
	)
}

func askStartupDelay(delay time.Duration, timeout time.Duration) time.Duration {
	if delay <= 0 {
		return 0
	}
	wait := askYesNo(
		fmt.Sprintf("Wait %s before the first requests?", delay),
		false,
		timeout,
	)
	if !wait {
		return 0
	}
	return delay
}

func askDeleteData(timeout time.Duration) bool {
	return askYesNo(
		"Delete all local app data, config, aliases, and autostart entry?",
		false,
		timeout,
	)
}

func askYesNo(question string, fallback bool, timeout time.Duration) bool {
	if !stdinIsTerminal() || timeout <= 0 {
		return fallback
	}

	fmt.Printf("%s Auto-select in %s: %s\n> ", question, timeout, yesNoText(fallback))
	answers := make(chan string, 1)
	go func() {
		line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		answers <- strings.TrimSpace(line)
	}()

	select {
	case answer := <-answers:
		return parseYesNo(answer, fallback)
	case <-time.After(timeout):
		fmt.Println(yesNoText(fallback))
		return fallback
	}
}

func stdinIsTerminal() bool {
	info, err := os.Stdin.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

func parseYesNo(answer string, fallback bool) bool {
	switch strings.ToLower(answer) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		return fallback
	}
}

func yesNoText(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}
