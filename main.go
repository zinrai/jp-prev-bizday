package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

// HolidayResponse represents the response from the Japanese holiday API
type HolidayResponse struct {
	Holiday bool   `json:"holiday"`
	Name    string `json:"name"`
}

// weekdayJP maps weekdays to Japanese characters
var weekdayJP = map[time.Weekday]string{
	time.Sunday:    "日",
	time.Monday:    "月",
	time.Tuesday:   "火",
	time.Wednesday: "水",
	time.Thursday:  "木",
	time.Friday:    "金",
	time.Saturday:  "土",
}

func main() {
	var (
		dateStr string
		verbose bool
		help    bool
	)

	flag.StringVar(&dateStr, "date", "", "基準日 (YYYY-MM-DD形式、デフォルト: 今日)")
	flag.BoolVar(&verbose, "verbose", false, "詳細表示モード")
	flag.BoolVar(&help, "help", false, "ヘルプを表示")
	flag.Parse()

	if help {
		printHelp()
		os.Exit(0)
	}

	// Set base date
	var baseDate time.Time
	var err error
	if dateStr == "" {
		// Get today's date in JST
		jst, _ := time.LoadLocation("Asia/Tokyo")
		baseDate = time.Now().In(jst)
	} else {
		baseDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: 無効な日付形式です: %s\n", dateStr)
			os.Exit(1)
		}
	}

	// Find the previous business day
	prevBizDay, err := findPreviousBusinessDay(baseDate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// Output results
	if verbose {
		outputVerbose(baseDate, prevBizDay)
	} else {
		outputSimple(prevBizDay)
	}
}

// printHelp displays the help message
func printHelp() {
	fmt.Println("jp-prev-bizday - 日本の直前の営業日を取得するツール")
	fmt.Println()
	fmt.Println("使用方法:")
	fmt.Println("  jp-prev-bizday [オプション]")
	fmt.Println()
	fmt.Println("説明:")
	fmt.Println("  指定された日付（デフォルトは今日）から遡って、")
	fmt.Println("  最初の営業日（土日祝日を除く平日）を返します。")
	fmt.Println("  日本の祝日に対応しています。")
	fmt.Println()
	fmt.Println("オプション:")
	fmt.Println("  -date string    基準日 (YYYY-MM-DD形式、デフォルト: 今日)")
	fmt.Println("  -verbose        詳細表示モード")
	fmt.Println("  -help           このヘルプを表示")
	fmt.Println()
	fmt.Println("例:")
	fmt.Println("  jp-prev-bizday")
	fmt.Println("  jp-prev-bizday -date 2025-07-24")
	fmt.Println("  jp-prev-bizday -verbose")
}

// findPreviousBusinessDay returns the first business day before the specified date
func findPreviousBusinessDay(from time.Time) (time.Time, error) {
	// Search up to 30 days back (sufficient for practical use)
	maxDays := 30
	current := from.AddDate(0, 0, -1) // Start from one day before

	for i := 0; i < maxDays; i++ {
		isBizDay, err := isBusinessDay(current)
		if err != nil {
			return time.Time{}, fmt.Errorf("営業日判定エラー: %w", err)
		}

		if isBizDay {
			return current, nil
		}

		current = current.AddDate(0, 0, -1)
	}

	return time.Time{}, fmt.Errorf("営業日が見つかりませんでした")
}

// isBusinessDay checks if the specified date is a business day
func isBusinessDay(date time.Time) (bool, error) {
	// Check if it's a weekend
	if isWeekend(date) {
		return false, nil
	}

	// Check if it's a holiday
	isHoliday, _, err := checkHoliday(date)
	if err != nil {
		return false, err
	}

	return !isHoliday, nil
}

// isWeekend checks if the date is Saturday or Sunday
func isWeekend(date time.Time) bool {
	weekday := date.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// checkHoliday calls the Japanese holiday API to check if the date is a holiday
func checkHoliday(date time.Time) (bool, string, error) {
	// Build API endpoint URL
	url := fmt.Sprintf("https://jp-holiday.net/api/v1/holiday/%d/%02d/%02d",
		date.Year(), date.Month(), date.Day())

	// Configure HTTP client with timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Make API request
	resp, err := client.Get(url)
	if err != nil {
		return false, "", fmt.Errorf("API呼び出しエラー: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return false, "", fmt.Errorf("APIエラー: ステータスコード %d", resp.StatusCode)
	}

	// Parse JSON response
	var holiday HolidayResponse
	if err := json.NewDecoder(resp.Body).Decode(&holiday); err != nil {
		return false, "", fmt.Errorf("JSONパースエラー: %w", err)
	}

	return holiday.Holiday, holiday.Name, nil
}

// outputSimple outputs only the date
func outputSimple(date time.Time) {
	fmt.Printf("%s\n", date.Format("2006-01-02"))
}

// outputVerbose outputs detailed information
func outputVerbose(baseDate, businessDay time.Time) {
	fmt.Printf("基準日: %s (%s)\n",
		baseDate.Format("2006-01-02"),
		weekdayJP[baseDate.Weekday()])
	fmt.Printf("直前の営業日: %s (%s)\n",
		businessDay.Format("2006-01-02"),
		weekdayJP[businessDay.Weekday()])
}
