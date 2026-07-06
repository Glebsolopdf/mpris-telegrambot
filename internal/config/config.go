package config

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"mpris-tg-status/internal/status"
)

const (
	defaultPollInterval = 20 * time.Second
	defaultNoPlayerPoll = 3 * time.Minute
	defaultTelegramMin  = 90 * time.Second
	defaultAvatarMin    = 24 * time.Hour
	defaultHTTPTimeout  = 15 * time.Second
	defaultMemoryLimit  = 64
	defaultGCPercent    = 50
	defaultPromptWait   = 10 * time.Second
	defaultBusinessWait = 30 * time.Second
)

type Config struct {
	TelegramBotToken     string
	TargetUserID         int64
	PreferredPlayer      string
	HTTPTimeout          time.Duration
	LogLevel             string
	MemoryLimitMB        int
	GCPercent            int
	StartupDelay         time.Duration
	PromptTimeout        time.Duration
	BusinessWaitInterval time.Duration
	UseMSAlias           bool
	EnvPath              string
	BusinessStatePath    string
	ShutdownRestorePath  string
	StatusSettings       status.Settings
}

var ErrEnvCreated = errors.New("env file created")

func Load() (Config, error) {
	dir, err := ExecutableDir()
	if err != nil {
		return Config{}, err
	}
	return LoadFromDir(dir)
}

func LoadFromDir(dir string) (Config, error) {
	envPath := filepath.Join(dir, ".env")
	values, err := readEnvFile(envPath)
	if errors.Is(err, os.ErrNotExist) {
		return Config{}, createDefaultEnv(envPath)
	}
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		TelegramBotToken:     env(values, "TELEGRAM_BOT_TOKEN"),
		TargetUserID:         int64Env(values, "TELEGRAM_TARGET_USER_ID"),
		PreferredPlayer:      env(values, "PREFERRED_PLAYER"),
		HTTPTimeout:          durationEnv(values, "HTTP_TIMEOUT", defaultHTTPTimeout),
		LogLevel:             logLevel(values),
		MemoryLimitMB:        intEnv(values, "MEMORY_LIMIT_MB", defaultMemoryLimit),
		GCPercent:            intEnv(values, "GC_PERCENT", defaultGCPercent),
		StartupDelay:         durationEnv(values, "STARTUP_DELAY", 0),
		PromptTimeout:        durationEnv(values, "PROMPT_TIMEOUT", defaultPromptWait),
		BusinessWaitInterval: durationEnv(values, "BUSINESS_CONNECTION_WAIT_INTERVAL", defaultBusinessWait),
		UseMSAlias:           boolEnv(values, "USE_MS_ALIAS", false),
		EnvPath:              envPath,
		BusinessStatePath:    filepath.Join(dir, "business_connection.json"),
		ShutdownRestorePath:  filepath.Join(dir, "shutdown_restore.txt"),
		StatusSettings: status.Settings{
			PollInterval:            durationEnv(values, "POLL_INTERVAL", defaultPollInterval),
			NoPlayerPollInterval:    durationEnv(values, "NO_PLAYER_POLL_INTERVAL", defaultNoPlayerPoll),
			TelegramMinInterval:     durationEnv(values, "TELEGRAM_MIN_UPDATE_INTERVAL", defaultTelegramMin),
			AvatarMinInterval:       durationEnv(values, "AVATAR_MIN_UPDATE_INTERVAL", defaultAvatarMin),
			DisableGeneratedAvatar:  !boolEnv(values, "GENERATED_AVATAR_ENABLED", true),
			IdleBio:                 defaultBio(values),
			ActiveBioTemplate:       env(values, "ACTIVE_BIO_TEMPLATE"),
			ActiveFirstNameTemplate: env(values, "ACTIVE_FIRST_NAME_TEMPLATE"),
			ActiveNameEmojis:        env(values, "ACTIVE_NAME_EMOJIS"),
			ActiveBioEmojis:         env(values, "ACTIVE_BIO_EMOJIS"),
			DefaultAvatarPath:       env(values, "DEFAULT_AVATAR_PATH"),
			DefaultFirstName:        env(values, "DEFAULT_FIRST_NAME"),
			DefaultLastName:         env(values, "DEFAULT_LAST_NAME"),
		},
	}
	cfg.StatusSettings.DefaultAvatarPath = resolvePath(dir, cfg.StatusSettings.DefaultAvatarPath)

	if err := validate(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func resolvePath(dir string, path string) string {
	if path == "" || filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(dir, path)
}
