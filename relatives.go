package relatives

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type RelativeRegex struct {
	re    *regexp.Regexp
	parse func([]string, time.Duration) (time.Duration, error)
}

var (
	oneYear   = time.Hour * 24 * 365
	oneMonth  = time.Hour * 24 * 30
	oneWeek   = time.Hour * 24 * 7
	oneDay    = time.Hour * 24
	oneHour   = time.Hour
	oneMinute = time.Minute
	oneSecond = time.Second
)

const (
	yearRe   = `[yY](?:ears?)?`
	monthRe  = `months?|M(?:onths?)?`
	weekRe   = `[Ww](?:eeks?)?`
	dayRe    = `[Dd](?:ays?)?`
	hourRe   = `[Hh](?:ours?)?`
	minuteRe = `m(?:inutes?)?|Minutes?`
	secondRe = `[Ss](?:econds?)?`
)

var ReRe = []RelativeRegex{
	{
		re: regexp.MustCompile(`\b(?P<quantity>\d+)\s*(?P<unit>` +
			yearRe + `|` +
			monthRe + `|` +
			weekRe + `|` +
			dayRe + `|` +
			hourRe + `|` +
			minuteRe + `|` +
			secondRe + `)\b`),
		parse: parseQuantityAndUnit,
	},

	{
		re: regexp.MustCompile(`\bago\b`),
		parse: func(_ []string, val time.Duration) (time.Duration, error) {
			if val == 0 {
				return 0, fmt.Errorf("invalid input: ago must be preceded by a quantity")
			}

			return -val, nil
		},
	},

	{
		re:    regexp.MustCompile(`\byesterday\b`),
		parse: func([]string, time.Duration) (time.Duration, error) { return -oneDay, nil },
	},
}

func parseQuantityAndUnit(parts []string, val time.Duration) (time.Duration, error) {
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid number of parts: expected at least 2 but got %d", len(parts))
	}

	quantity, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid quantity %s: %w", parts[0], err)
	}

	unit, err := parseUnit(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid unit %s: %w", parts[1], err)
	}

	return val + time.Duration(quantity)*unit, nil
}

func Parse(in string) (time.Duration, error) {
	fmt.Printf("in: %s\n", in)
	found := false
	val := time.Duration(0)
	var err error

	if len(in) == 0 {
		return 0, nil
	}

	for _, r := range ReRe {
		matches := r.re.FindAllStringSubmatch(in, -1)
		if matches == nil {
			continue
		}

		for _, match := range matches {
			val, err = r.parse(match[1:], val)
			if err != nil {
				return 0, fmt.Errorf("error parsing match: %w", err)
			}

			found = true
		}
	}

	if !found {
		return 0, fmt.Errorf("invalid input: %s", in)
	}

	return val, nil
}

func parseUnit(u string) (time.Duration, error) {
	switch string(u) {
	case "y", "Y", "year", "Year", "years", "Years":
		return oneYear, nil
	case "M", "month", "Month", "months", "Months":
		return oneMonth, nil
	case "w", "week", "Week", "weeks", "Weeks":
		return oneWeek, nil
	case "d", "day", "Day", "days", "Days":
		return oneDay, nil
	case "h", "hour", "Hour", "hours", "Hours":
		return oneHour, nil
	case "m", "minute", "Minute", "minutes", "Minutes":
		return oneMinute, nil
	case "s", "second", "Second", "seconds", "Seconds":
		return oneSecond, nil
	}

	return 0, fmt.Errorf("invalid unit: %s", u)
}

func Format(time.Time) string {
	return ""
}
