package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const defaultEnv = `TELEGRAM_BOT_TOKEN=
TELEGRAM_TARGET_USER_ID=
PREFERRED_PLAYER=spotify
DEFAULT_BIO=
DEFAULT_AVATAR_PATH=default_avatar.png
DEFAULT_FIRST_NAME=
DEFAULT_LAST_NAME=
ACTIVE_BIO_TEMPLATE={emoji_bio} Listening now: {artist} -- {title}
ACTIVE_FIRST_NAME_TEMPLATE={emoji_name} {default_first_name}
ACTIVE_NAME_EMOJIS=
ACTIVE_BIO_EMOJIS=
GENERATED_AVATAR_ENABLED=true
POLL_INTERVAL=20s
NO_PLAYER_POLL_INTERVAL=3m
TELEGRAM_MIN_UPDATE_INTERVAL=90s
AVATAR_MIN_UPDATE_INTERVAL=24h
HTTP_TIMEOUT=15s
LOG_LEVEL=info
MEMORY_LIMIT_MB=64
GC_PERCENT=50
STARTUP_DELAY=5m
PROMPT_TIMEOUT=10s
BUSINESS_CONNECTION_WAIT_INTERVAL=30s
USE_MS_ALIAS=false
`

func createDefaultEnv(path string) error {
	if err := os.WriteFile(path, []byte(defaultEnv), 0o600); err != nil {
		return err
	}
	return fmt.Errorf("%w: fill %s and restart", ErrEnvCreated, path)
}

func readEnvFile(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	values := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parseEnvLine(scanner.Text(), values)
	}
	return values, scanner.Err()
}

func parseEnvLine(line string, values map[string]string) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		return
	}

	key, value, ok := strings.Cut(trimmed, "=")
	if !ok {
		return
	}
	values[strings.TrimSpace(key)] = cleanEnvValue(value)
}

func cleanEnvValue(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, `"`)
	return strings.Trim(value, `'`)
}
