package relatives

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	tests := []struct {
		in          string
		want        time.Duration
		assertError func(*testing.T, error)
	}{
		{},
		{in: "1y", want: oneYear},
		{in: "1year", want: oneYear},
		{in: "1years", want: oneYear},
		{in: "1Year", want: oneYear},
		{in: "1Years", want: oneYear},

		{in: "1M", want: oneMonth},
		{in: "1month", want: oneMonth},
		{in: "1months", want: oneMonth},
		{in: "1Month", want: oneMonth},
		{in: "1Months", want: oneMonth},

		{in: "1m", want: oneMinute},

		{in: "3 Hours", want: 3 * oneHour},
		{in: "1 day 2 hours", want: oneDay + 2*oneHour},
		{in: "1 year 30s", want: oneYear + 30*oneSecond},
		{in: "10 y 1   M 1m", want: 10*oneYear + oneMonth + oneMinute},

		{in: "1 day ago", want: -1 * oneDay},
		{in: "2d ago", want: -2 * oneDay},
		{in: "2 days ago", want: -2 * oneDay},
		{in: "yesterday", want: -1 * oneDay},
	}

	for _, test := range tests {
		t.Run(string(test.in), func(t *testing.T) {
			got, err := Parse(test.in)

			if test.assertError != nil {
				test.assertError(t, err)
			} else if err != nil {
				t.Errorf("expected no error but got %#v", err)
			}

			if got != test.want {
				t.Errorf("want %v but got %v", test.want, got)
			}
		})
	}

}
