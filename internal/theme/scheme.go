package theme

// WindowsTerminalSchemeJSON provides the exact JSON definition matching Wraith's spectral aesthetic
const WindowsTerminalSchemeJSON = `{
  "name": "Wraith Spectral",
  "background": "#0F0E17",
  "foreground": "#E0DEF4",
  "cursorColor": "#EBBCBA",
  "selectionBackground": "#2A283E",
  "black": "#191724",
  "red": "#EB6F92",
  "green": "#00FF66",
  "yellow": "#FFB000",
  "blue": "#9CCFD8",
  "purple": "#B39DDB",
  "cyan": "#31748F",
  "white": "#E0DEF4",
  "brightBlack": "#6E6A86",
  "brightRed": "#EB6F92",
  "brightGreen": "#56F996",
  "brightYellow": "#FFA726",
  "brightBlue": "#9CCFD8",
  "brightPurple": "#C4A7E7",
  "brightCyan": "#EBBCBA",
  "brightWhite": "#FFFFFF"
}`

// ThemeInstructions returns step-by-step guidance for Windows Terminal
func ThemeInstructions() string {
	return `To install the "Wraith Spectral" color scheme in Windows Terminal:

1. Open Windows Terminal settings (press Ctrl + ,).
2. Click "Open JSON file" at the bottom left (or press Ctrl + Shift + ,).
3. Under the "schemes" array in settings.json, paste the JSON snippet below.
4. Save settings.json and set "colorScheme": "Wraith Spectral" in your profile default.
`
}

