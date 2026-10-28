package theme

// ASCII and styled ceremonies for Wraith CLI

// GhostBanner returns the multi-line ASCII ghost banner
func GhostBanner(plain bool) string {
	banner := `
          .---.
         /     \
        | () () |     W R A I T H
         \  ^  /      The Haunted Terminal
          |||||       (runs 100% offline via Ollama)
          |||||
`
	if plain {
		return banner
	}
	// Lavender truecolor or bright magenta ANSI escape
	return "\033[38;2;179;157;219m" + banner + "\033[0m"
}

// FarewellCeremony returns quiet goodbye text on exorcism
func FarewellCeremony(plain bool, path string) string {
	if plain {
		return "[wraith] The spirit untethers from " + path + ". Silence returns to your terminal.\n"
	}
	return "\033[90m🪦 [wraith] The spirit untethers from " + path + "...\033[0m\n\033[38;2;179;157;219;3m   A gentle breeze extinguishes the candle. Silence returns to your terminal.\033[0m\n"
}

// StatusHeader returns a themed status banner
func StatusHeader(plain bool, persona string, hushed bool) string {
	statusTag := "AWAKE"
	if hushed {
		statusTag = "HUSHED"
	}

	if plain {
		return "=== WRAITH SPECTRAL STATUS [" + statusTag + "] ==="
	}

	colorCode := "\033[38;2;179;157;219;1m" // Lavender bold
	if hushed {
		colorCode = "\033[90;1m" // Dim gray
	}
	return colorCode + "👻 === WRAITH SPECTRAL STATUS [" + statusTag + "] ===" + "\033[0m"
}

