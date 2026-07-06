package text

func fitBioTemplate(template string, vars templateVars) string {
	rendered := renderTemplate(template, vars)
	if runeLen(rendered) <= telegramBioLimit {
		return rendered
	}

	limited := shrinkTrackText(template, vars)
	if runeLen(limited) <= telegramBioLimit {
		return limited
	}
	return truncateRunes(limited, telegramBioLimit)
}

func shrinkTrackText(template string, vars templateVars) string {
	artistLimit := runeLen(vars.Artist)
	titleLimit := runeLen(vars.Title)

	for artistLimit > 0 || titleLimit > 0 {
		if artistLimit >= titleLimit && artistLimit > 0 {
			artistLimit--
		} else if titleLimit > 0 {
			titleLimit--
		}

		candidate := vars
		candidate.Artist = truncateRunes(vars.Artist, artistLimit)
		candidate.Title = truncateRunes(vars.Title, titleLimit)
		rendered := renderTemplate(template, candidate)
		if runeLen(rendered) <= telegramBioLimit {
			return rendered
		}
	}

	vars.Artist = ""
	vars.Title = ""
	return renderTemplate(template, vars)
}

func runeLen(value string) int {
	return len([]rune(value))
}
