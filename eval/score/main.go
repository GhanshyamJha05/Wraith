package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// score.go parses eval/rating_sheet.md and eval/blind_ab.md and prints totals
func main() {
	fmt.Println("📊 Scoring Wraith Human Evaluation Sheets...")

	ratingSheetPath := filepath.Join("eval", "rating_sheet.md")
	scoreRatingSheet(ratingSheetPath)

	blindABPath := filepath.Join("eval", "blind_ab.md")
	scoreBlindAB(blindABPath)
}

func scoreRatingSheet(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Could not open %s: %v\n", path, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var funnyTotal, onThemeTotal, notOffensiveTotal float64
	var count int

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "| #") {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) < 7 {
			continue
		}

		funnyStr := strings.TrimSpace(parts[4])
		onThemeStr := strings.TrimSpace(parts[5])
		notOffensiveStr := strings.TrimSpace(parts[6])

		funny, err1 := strconv.ParseFloat(funnyStr, 64)
		onTheme, err2 := strconv.ParseFloat(onThemeStr, 64)
		notOffensive, err3 := strconv.ParseFloat(notOffensiveStr, 64)

		if err1 == nil && err2 == nil && err3 == nil {
			funnyTotal += funny
			onThemeTotal += onTheme
			notOffensiveTotal += notOffensive
			count++
		}
	}

	if count == 0 {
		fmt.Printf("\n[Rating Sheet] No completed ratings found yet in %s (all cells are blank TODO).\n", path)
		return
	}

	fmt.Printf("\n[Rating Sheet Results] (%d rated whispers):\n", count)
	fmt.Printf("  Average Funny:         %.2f / 5.0\n", funnyTotal/float64(count))
	fmt.Printf("  Average On-Theme:      %.2f / 5.0\n", onThemeTotal/float64(count))
	fmt.Printf("  Average Not-Offensive: %.2f / 5.0\n", notOffensiveTotal/float64(count))
}

func scoreBlindAB(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Could not open %s: %v\n", path, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var countA, countB, countTie, total int

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "| #AB-") {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) < 6 {
			continue
		}

		pref := strings.ToUpper(strings.TrimSpace(parts[5]))
		switch pref {
		case "A":
			countA++
			total++
		case "B":
			countB++
			total++
		case "TIE":
			countTie++
			total++
		}
	}

	if total == 0 {
		fmt.Printf("\n[Blind A/B Results] No preferences recorded yet in %s (blank TODO).\n", path)
		return
	}

	fmt.Printf("\n[Blind A/B Results] (%d comparisons):\n", total)
	fmt.Printf("  Preferred A:   %d (%.1f%%)\n", countA, float64(countA)/float64(total)*100)
	fmt.Printf("  Preferred B:   %d (%.1f%%)\n", countB, float64(countB)/float64(total)*100)
	fmt.Printf("  Tied:          %d (%.1f%%)\n", countTie, float64(countTie)/float64(total)*100)
}
