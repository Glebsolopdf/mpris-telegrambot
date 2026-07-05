# MPRIS Telegram Business Status

Small Linux service that reads your active MPRIS media player and updates a Telegram Business account name, bio, and optional profile photo.

It is designed for Spotify, but it can work with any MPRIS player. If `PREFERRED_PLAYER` is set, every other player is ignored, so browser videos do not accidentally become your Telegram bio.

## Features

- Updates the Telegram Business bio with the currently playing track.
- Adds configurable music emoji placeholders to the profile first name and bio.
- Optionally generates an album-cover profile photo.
- Restores your default name, bio, and avatar when playback stops or the service shuts down.
- Waits for `business_connection_id` instead of exiting when the bot has not received it yet.
- Keeps service logs in RAM with a 1 MiB cap.
- Keeps avatar cooldown state only in RAM.
- Provides a small `ms` command alias for background service control.

## Requirements

- Linux desktop session with D-Bus and MPRIS support.
- Go 1.23 or newer.
- Telegram Business account.
- Telegram bot connected to that Business account.
- Bot permissions to edit the Business account name, bio, and profile photo.

## Installation

Fedora:

```sh
sudo dnf install golang git dbus-x11 file
```

Debian or Ubuntu:

```sh
sudo apt update
sudo apt install golang-go git dbus-user-session file
```

Arch Linux:

```sh
sudo pacman -S go git dbus file
```

Void Linux:

```sh
sudo xbps-install -S go git dbus file
```

openSUSE:

```sh
sudo zypper install go git dbus-1 file
```

`playerctl` is optional, but useful for checking what your desktop exposes through MPRIS. Install it with your package manager if you want this diagnostic command:

```sh
playerctl metadata
```

After installing dependencies, clone the repository:

```sh
git clone https://github.com/Glebsolopdf/mpris-telegrambot
cd mpris-telegrambot
```

## Telegram Setup

1. Create a bot with `@BotFather` and copy the bot token.
2. Open Telegram Business settings.
3. Connect the bot to your Business account.
4. Allow the bot to edit your profile name, bio, and profile photo.
5. Get your numeric Telegram user ID from a trusted ID bot or Telegram API tool.

The app discovers `business_connection_id` from Telegram `business_connection` updates and saves it in `business_connection.json`. If the update is missing or the Business connection is disabled, the service keeps running and retries on the configured interval.

## Build

Build the binary into this project `bin` directory:

```sh
chmod +x ./build.sh
./build.sh
```

The output binary is:

```sh
./bin/mpris-tg-status
```

Manual build command:

```sh
mkdir -p ./bin
go build -buildvcs=false -o ./bin/mpris-tg-status ./cmd/mpris-tg-status
```

## Configuration

On first start, the app creates `.env` next to the binary and exits. Fill in at least:

```sh
TELEGRAM_BOT_TOKEN=123456:token
TELEGRAM_TARGET_USER_ID=123456789
```

Full example:

```sh
TELEGRAM_BOT_TOKEN=123456:token
TELEGRAM_TARGET_USER_ID=123456789
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
```

## Important Options

`PREFERRED_PLAYER` matches the MPRIS bus name. Use `spotify` for Spotify. When it is set to a player name, Firefox, Chrome, YouTube, and other players are ignored. Set `PREFERRED_PLAYER=all` or leave it empty to accept any MPRIS player.

`DEFAULT_BIO`, `DEFAULT_FIRST_NAME`, `DEFAULT_LAST_NAME`, and `DEFAULT_AVATAR_PATH` are restored when playback stops or when you choose to restore the profile on shutdown.

`ACTIVE_BIO_TEMPLATE` controls the bio while music is playing. `ACTIVE_FIRST_NAME_TEMPLATE` controls the active profile first name. Available placeholders are `{emoji_name}`, `{emoji_bio}`, `{emoji}`, `{artist}`, `{title}`, and `{default_first_name}`. `{emoji}` is kept for older configs and resolves to the active template's emoji.

`ACTIVE_NAME_EMOJIS` and `ACTIVE_BIO_EMOJIS` customize the emoji pool for each template. Leave them empty to use the same random music emoji for name and bio. Use comma-separated values such as `🎧,🎵,🔥`. Set either value to `false` to make that emoji variable empty. Example: `ACTIVE_BIO_TEMPLATE={emoji_bio} Listening: {artist} — {title}`.

`GENERATED_AVATAR_ENABLED=false` disables generated album-cover avatars while keeping name and bio updates active.

`NO_PLAYER_POLL_INTERVAL` is used after no matching MPRIS player is found. The normal `POLL_INTERVAL` resumes when the player appears again.

`AVATAR_MIN_UPDATE_INTERVAL` throttles profile photo changes in memory. It is not written to disk.

`MEMORY_LIMIT_MB` and `GC_PERCENT` tune Go runtime memory use for background mode.

`STARTUP_DELAY` is offered by `up` and `restart`. If you do not answer within `PROMPT_TIMEOUT`, startup continues without waiting.

`BUSINESS_CONNECTION_WAIT_INTERVAL` controls how often the service retries Telegram Business connection discovery.

`USE_MS_ALIAS=true` makes the app configure `~/.local/bin/ms` automatically when the alias is missing. Existing unrelated `ms` commands are left unchanged.

## Commands

Run in the foreground:

```sh
./mpris-tg-status run
```

Manage the background service:

```sh
./mpris-tg-status up
./mpris-tg-status down
./mpris-tg-status restart
./mpris-tg-status status
./mpris-tg-status logs
./mpris-tg-status delete
./mpris-tg-status help
```

`down` asks whether to restore the default profile. If you do not answer within `PROMPT_TIMEOUT`, it restores the profile.

`up` asks whether to wait for `STARTUP_DELAY` before the first requests. If you do not answer within `PROMPT_TIMEOUT`, it starts immediately.

## Alias

Create a short `ms` alias by symlinking the binary into a directory from your `PATH`:

```sh
mkdir -p "$HOME/.local/bin"
ln -sf "$PWD/bin/mpris-tg-status" "$HOME/.local/bin/ms"
```

Then use:

```sh
ms up
ms status
ms restart
ms down
ms logs
ms delete
ms help
```

## Autostart

Create a desktop autostart entry:

```sh
mkdir -p "$HOME/.config/autostart"
cat > "$HOME/.config/autostart/mpris-tg-status.desktop" <<EOF
[Desktop Entry]
Type=Application
Name=MPRIS Telegram Business Status
Exec=$PWD/bin/mpris-tg-status up
Path=$PWD
Terminal=false
X-GNOME-Autostart-enabled=true
EOF
```

Because autostart has no interactive terminal, startup prompts use their default answers.

## Logs

Logs are kept in RAM and capped at 1 MiB. Print the current log path with:

```sh
ms logs
```

On most systems the path is under `XDG_RUNTIME_DIR`, for example `/run/user/1000/mpris-tg-status.log`.

## Runtime Files

- `.env` stores local configuration.
- `business_connection.json` stores the Telegram Business connection ID.
- `mpris-tg-status.pid` stores the background service PID.
- `shutdown_restore.txt` is a temporary control file used during shutdown.

## Delete Local Data

Remove local data created by the app without deleting source files or binaries:

```sh
ms delete
```

The command restores the default profile, stops the background service, removes `.env`, `business_connection.json`, runtime control files, RAM logs, known `ms` symlinks, and the desktop autostart entry.
