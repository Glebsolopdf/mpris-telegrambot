package mpris

import "strings"

func matchesPreferred(name string, preferred []string) bool {
	for _, player := range preferred {
		if strings.Contains(strings.ToLower(name), strings.ToLower(player)) {
			return true
		}
	}
	return false
}

func normalizedPreferred(preferred string) []string {
	parts := strings.Split(preferred, ",")
	players := make([]string, 0, len(parts))
	for _, part := range parts {
		player := strings.TrimSpace(part)
		if player == "" {
			continue
		}
		if strings.EqualFold(player, "all") {
			return nil
		}
		players = append(players, player)
	}
	return players
}
