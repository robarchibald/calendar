package calendar

import (
	"testing"
	"time"
)

func TestGetDailyOccurrencesWeekdays(t *testing.T) {
	expected := []time.Time{time.Date(2016, 4, 1, 12, 30, 0, 0, time.UTC),
		time.Date(2016, 4, 4, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 5, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 6, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 7, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 8, 12, 30, 0, 0, time.UTC),
		time.Date(2016, 4, 11, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 12, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 13, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 14, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 15, 12, 30, 0, 0, time.UTC),
		time.Date(2016, 4, 18, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 19, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 20, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 21, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 22, 12, 30, 0, 0, time.UTC),
		time.Date(2016, 4, 25, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 26, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 27, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 28, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 29, 12, 30, 0, 0, time.UTC)}
	actual := getDailyOccurrences(time.Date(2010, 1, 1, 12, 30, 0, 0, time.UTC), 1, true, nil, time.Date(2016, 4, 1, 0, 0, 0, 0, time.UTC), time.Date(2016, 5, 1, 0, 0, 0, 0, time.UTC))
	compareTimes(t, expected, actual)
}

func TestGetDailyOccurrencesAllDays(t *testing.T) {
	expected := []time.Time{time.Date(2016, 4, 2, 12, 30, 0, 0, time.UTC),
		time.Date(2016, 4, 5, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 8, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 11, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 14, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 17, 12, 30, 0, 0, time.UTC),
		time.Date(2016, 4, 20, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 23, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 26, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 29, 12, 30, 0, 0, time.UTC)}
	actual := getDailyOccurrences(time.Date(2010, 1, 1, 12, 30, 0, 0, time.UTC), 3, false, nil, time.Date(2016, 4, 1, 0, 0, 0, 0, time.UTC), time.Date(2016, 5, 1, 0, 0, 0, 0, time.UTC))
	compareTimes(t, expected, actual)
}

func TestGetDailyStartTime(t *testing.T) {
	recurrenceStartDate := time.Date(2010, 1, 1, 12, 30, 0, 0, time.UTC)
	timePeriodStart := time.Date(2016, 4, 1, 0, 0, 0, 0, time.UTC)

	// expected dates gathered from Outlook and double-checked with calculator
	actual := getDailyStartTime(recurrenceStartDate, 2, timePeriodStart)
	if actual != time.Date(2016, 4, 1, 12, 30, 0, 0, time.UTC) {
		t.Error("expected to line up with timePeriodStart day", actual)
	}

	actual = getDailyStartTime(recurrenceStartDate, 3, timePeriodStart)
	if actual != time.Date(2016, 4, 2, 12, 30, 0, 0, time.UTC) {
		t.Error("expected to be on 4/2", actual)
	}

	actual = getDailyStartTime(recurrenceStartDate, 5, timePeriodStart)
	if actual != time.Date(2016, 4, 4, 12, 30, 0, 0, time.UTC) {
		t.Error("expected to be on 4/4", actual)
	}
}

func TestGetWeekdayStartTime(t *testing.T) {
	recurrenceStartDate := time.Date(2010, 1, 1, 12, 30, 0, 0, time.UTC)
	timePeriodStart := time.Date(2016, 4, 1, 0, 0, 0, 0, time.UTC)
	recurEvery := 3

	actual := getWeekdayStartTime(recurrenceStartDate, recurEvery, timePeriodStart)
	if actual != time.Date(2016, 4, 5, 12, 30, 0, 0, time.UTC) {
		t.Error("expected correct start date:", actual)
	}
}

func TestGetWeekdays(t *testing.T) {
	// 1972 is 281 weeks + 5 days.  1/1/2010 is a Friday, so that's 281*5+3=
	if actual := getWeekdays(1972, time.Date(2010, 1, 1, 12, 30, 0, 0, time.UTC)); actual != 1408 {
		t.Error("expected 1408 weekdays", actual)
	}

	// 1972 is 281 weeks + 5 days.  1/4/2010 is a Monday, so that's 281*5+4=
	if actual := getWeekdays(1972, time.Date(2010, 1, 4, 12, 30, 0, 0, time.UTC)); actual != 1409 {
		t.Error("expected 1409 weekdays", actual)
	}

	// 1972 is 281 weeks + 5 days.  1/5/2010 is a Tuesday, so that's 281*5+3=
	if actual := getWeekdays(1972, time.Date(2010, 1, 5, 12, 30, 0, 0, time.UTC)); actual != 1408 {
		t.Error("expected 1408 weekdays", actual)
	}
}

func TestGetWeeklyOccurrences(t *testing.T) {
	recurrenceStartDate := time.Date(2010, 1, 1, 12, 30, 0, 0, time.UTC)
	timePeriodStart := time.Date(2016, 4, 1, 0, 0, 0, 0, time.UTC)
	timePeriodEnd := time.Date(2016, 5, 1, 0, 0, 0, 0, time.UTC)
	expected := []time.Time{time.Date(2016, 4, 1, 12, 30, 0, 0, time.UTC),
		time.Date(2016, 4, 4, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 6, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 8, 12, 30, 0, 0, time.UTC),
		time.Date(2016, 4, 11, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 13, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 15, 12, 30, 0, 0, time.UTC),
		time.Date(2016, 4, 18, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 20, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 22, 12, 30, 0, 0, time.UTC),
		time.Date(2016, 4, 25, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 27, 12, 30, 0, 0, time.UTC), time.Date(2016, 4, 29, 12, 30, 0, 0, time.UTC)}
	// 42 = MWF weekly meeting
	actual := getWeeklyOccurrences(recurrenceStartDate, 1, getIncludedWeeklyDays(42), nil, timePeriodStart, timePeriodEnd)
	compareTimes(t, expected, actual)
}

func TestGetWeeklyOccurrencesEndsMidWeek(t *testing.T) {
	recurrenceStartDate := time.Date(2010, 1, 1, 12, 30, 0, 0, time.UTC)
	timePeriodStart := time.Date(2016, 5, 1, 0, 0, 0, 0, time.UTC)
	timePeriodEnd := time.Date(2016, 6, 1, 0, 0, 0, 0, time.UTC)
	expected := []time.Time{time.Date(2016, 5, 1, 12, 30, 0, 0, time.UTC),
		time.Date(2016, 5, 8, 12, 30, 0, 0, time.UTC), time.Date(2016, 5, 15, 12, 30, 0, 0, time.UTC),
		time.Date(2016, 5, 22, 12, 30, 0, 0, time.UTC), time.Date(2016, 5, 29, 12, 30, 0, 0, time.UTC)}
	// 64 = SUN weekly meeting
	actual := getWeeklyOccurrences(recurrenceStartDate, 1, getIncludedWeeklyDays(64), nil, timePeriodStart, timePeriodEnd)
	compareTimes(t, expected, actual)
}

func TestGetIncludedDays(t *testing.T) {
	days := getIncludedDays([]time.Weekday{time.Sunday, time.Tuesday, time.Friday},
		time.Date(2016, 5, 29, 12, 30, 0, 0, time.UTC),
		time.Date(2016, 5, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2016, 5, 31, 23, 59, 59, 999, time.UTC))
	if len(days) != 2 || days[0] != time.Date(2016, 5, 29, 12, 30, 0, 0, time.UTC) || days[1] != time.Date(2016, 5, 31, 12, 30, 0, 0, time.UTC) {
		t.Error("expected Sunday and Tuesday", days)
	}
}

func TestGetWeeklyStartTime(t *testing.T) {
	recurrenceStartDate := time.Date(2010, 1, 1, 12, 30, 0, 0, time.UTC)
	timePeriodStart := time.Date(2016, 4, 1, 0, 0, 0, 0, time.UTC)

	actual := getWeeklyStartTime(recurrenceStartDate, 3, timePeriodStart)
	if actual != time.Date(2016, 4, 3, 12, 30, 0, 0, time.UTC) {
		t.Error("expected correct start date 1:", actual)
	}

	actual = getWeeklyStartTime(recurrenceStartDate, 1, timePeriodStart)
	if actual != time.Date(2016, 3, 27, 12, 30, 0, 0, time.UTC) {
		t.Error("expected correct start date 2:", actual)
	}
}

func TestGetWeeklyStartTimeInSameWeek(t *testing.T) {
	recurrenceStartDate := time.Date(2016, 3, 30, 12, 30, 0, 0, time.UTC)
	timePeriodStart := time.Date(2016, 4, 1, 0, 0, 0, 0, time.UTC)
	recurEvery := 3

	actual := getWeeklyStartTime(recurrenceStartDate, int(recurEvery), timePeriodStart)
	if actual != time.Date(2016, 3, 27, 12, 30, 0, 0, time.UTC) {
		t.Error("expected correct start date:", actual)
	}
}

func TestGetIncludedWeekdays(t *testing.T) {
	expected := []time.Weekday{time.Monday, time.Wednesday, time.Friday}
	compareWeekdays(t, expected, getIncludedWeeklyDays(32+8+2))

	expected = []time.Weekday{time.Sunday, time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday}
	compareWeekdays(t, expected, getIncludedWeeklyDays(64+32+16+8+4+2+1))

	expected = []time.Weekday{time.Sunday, time.Tuesday, time.Thursday, time.Saturday}
	compareWeekdays(t, expected, getIncludedWeeklyDays(64+16+4+1))
}

func TestGetMonthlyOccurrences(t *testing.T) {
	recurrenceStartDate := time.Date(2010, 1, 1, 12, 30, 0, 0, time.UTC)
	timePeriodStart := time.Date(2016, 4, 1, 0, 0, 0, 0, time.UTC)
	timePeriodEnd := time.Date(2016, 6, 1, 0, 0, 0, 0, time.UTC)
	var monthlyDay, monthlyDayOfWeek, monthlyWeekOfMonth int16
	expected := []time.Time{time.Date(2016, 4, 15, 12, 30, 0, 0, time.UTC), time.Date(2016, 5, 15, 12, 30, 0, 0, time.UTC)}
	monthlyDay = 15 // 15th of every month
	actual := getMonthlyOccurrences(recurrenceStartDate, 1, &monthlyDay, nil, nil, nil, timePeriodStart, timePeriodEnd)
	compareTimes(t, expected, actual)

	monthlyDayOfWeek = 4   // Thursday
	monthlyWeekOfMonth = 3 // 3rd week
	expected = []time.Time{time.Date(2016, 4, 21, 12, 30, 0, 0, time.UTC), time.Date(2016, 5, 19, 12, 30, 0, 0, time.UTC)}
	actual = getMonthlyOccurrences(recurrenceStartDate, 1, nil, &monthlyDayOfWeek, &monthlyWeekOfMonth, nil, timePeriodStart, timePeriodEnd)
	compareTimes(t, expected, actual)
}

func TestGetMonthOccurrence(t *testing.T) {
	startDate := time.Date(2016, 5, 1, 12, 30, 0, 0, time.UTC)
	timePeriodStart := time.Date(2016, 4, 1, 0, 0, 0, 0, time.UTC)
	timePeriodEnd := time.Date(2016, 6, 1, 0, 0, 0, 0, time.UTC)
	var monthlyDay, monthlyDayOfWeek, monthlyWeekOfMonth int16
	monthlyDay = 15
	date := getMonthOccurrence(startDate, timePeriodStart, timePeriodEnd, &monthlyDay, nil, nil)
	if len(date) != 1 || date[0] != time.Date(2016, 5, 15, 12, 30, 0, 0, time.UTC) {
		t.Error("expected 5/15/2016", date)
	}

	monthlyDayOfWeek = 4   // Thursday
	monthlyWeekOfMonth = 3 // 3rd week
	date = getMonthOccurrence(startDate, timePeriodStart, timePeriodEnd, nil, &monthlyDayOfWeek, &monthlyWeekOfMonth)
	if len(date) != 1 || date[0] != time.Date(2016, 5, 19, 12, 30, 0, 0, time.UTC) {
		t.Error("expected 5/19/2016", date)
	}

	monthlyDayOfWeek = 4   // Thursday
	monthlyWeekOfMonth = 5 // 5th week
	date = getMonthOccurrence(startDate, timePeriodStart, timePeriodEnd, nil, &monthlyDayOfWeek, &monthlyWeekOfMonth)
	if len(date) != 0 {
		t.Error("expected no 5th Thursday", date)
	}

	monthlyDayOfWeek = 2   // Thursday
	monthlyWeekOfMonth = 5 // 5th week
	date = getMonthOccurrence(startDate, timePeriodStart, timePeriodEnd, nil, &monthlyDayOfWeek, &monthlyWeekOfMonth)
	if len(date) != 1 || date[0] != time.Date(2016, 5, 31, 12, 30, 0, 0, time.UTC) {
		t.Error("expected 5/31", date)
	}

	// no valid date.  Starts & ends on same day
	date = getMonthOccurrence(startDate, timePeriodStart, timePeriodStart, nil, &monthlyDayOfWeek, &monthlyWeekOfMonth)
	if len(date) != 0 {
		t.Error("expected empty list", date)
	}
}

func TestGetMonthlyStartTime(t *testing.T) {
	recurrenceStartDate := time.Date(2010, 1, 1, 12, 30, 0, 0, time.UTC)
	timePeriodStart := time.Date(2016, 4, 1, 0, 0, 0, 0, time.UTC)

	actual := getMonthlyStartTime(recurrenceStartDate, 7, timePeriodStart)
	if actual != time.Date(2016, 6, 1, 12, 30, 0, 0, time.UTC) {
		t.Error("expected correct start date 1:", actual)
	}
}

func TestGetMonths(t *testing.T) {
	fromDate := time.Date(2010, 1, 1, 12, 30, 0, 0, time.UTC)
	toDate := time.Date(2016, 4, 1, 0, 0, 0, 0, time.UTC)
	actual := getMonths(fromDate, toDate)
	if actual != 75 {
		t.Error("expected to return 75 months")
	}

	actual = getMonths(time.Date(2010, 5, 1, 12, 30, 0, 0, time.UTC), time.Date(2011, 1, 1, 12, 30, 0, 0, time.UTC))
	if actual != 8 {
		t.Error("expected to return 8 months", actual)
	}

	actual = getMonths(fromDate, fromDate)
	if actual != 0 {
		t.Error("expected to return 0 months", actual)
	}

	actual = getMonths(toDate, fromDate)
	if actual != -75 {
		t.Error("expected to return -75 months", actual)
	}

	// months are different, but only actually 2 days apart, so should return 0
	actual = getMonths(time.Date(2010, 1, 30, 0, 0, 0, 0, time.UTC), time.Date(2010, 2, 1, 0, 0, 0, 0, time.UTC))
	if actual != 0 {
		t.Error("expected to return 0 months", actual)
	}
}

func TestGetYearlyOccurrences(t *testing.T) {
	recurrenceStartDate := time.Date(2010, 1, 1, 12, 30, 0, 0, time.UTC)
	timePeriodStart := time.Date(2016, 4, 1, 0, 0, 0, 0, time.UTC)
	timePeriodEnd := time.Date(2018, 4, 1, 0, 0, 0, 0, time.UTC)
	var yearlyMonth, monthlyDay, monthlyDayOfWeek, monthlyWeekOfMonth int16
	expected := []time.Time{time.Date(2017, 2, 14, 12, 30, 0, 0, time.UTC), time.Date(2018, 2, 14, 12, 30, 0, 0, time.UTC)}
	yearlyMonth = 2
	monthlyDay = 14 // 14th of every month
	actual := getYearlyOccurrences(recurrenceStartDate, 1, &yearlyMonth, &monthlyDay, nil, nil, nil, timePeriodStart, timePeriodEnd)
	compareTimes(t, expected, actual)

	monthlyDayOfWeek = 4   // Thursday
	monthlyWeekOfMonth = 3 // 3rd week
	expected = []time.Time{time.Date(2017, 2, 16, 12, 30, 0, 0, time.UTC), time.Date(2018, 2, 15, 12, 30, 0, 0, time.UTC)}
	actual = getYearlyOccurrences(recurrenceStartDate, 1, &yearlyMonth, nil, &monthlyDayOfWeek, &monthlyWeekOfMonth, nil, timePeriodStart, timePeriodEnd)
	compareTimes(t, expected, actual)
}

/*********************************************************************************************/

func compareTimes(t *testing.T, expected []time.Time, actual []time.Time) {
	if len(expected) != len(actual) {
		t.Log("expected:", expected)
		t.Log("actual:", actual)
		t.Fatalf("expected matching lengths.  Expected:%d, Actual:%d", len(expected), len(actual))
	}
	for i, item := range actual {
		if item != expected[i] {
			t.Errorf("expected[%d] %v vs actual[%d] %v", i, expected[i], i, actual[i])
		}
	}
}

func compareWeekdays(t *testing.T, expected []time.Weekday, actual []time.Weekday) {
	if len(expected) != len(actual) {
		t.Log("expected:", expected)
		t.Log("actual:", actual)
		t.Fatalf("expected matching lengths.  Expected:%d, Actual:%d", len(expected), len(actual))
	}
	for i, item := range actual {
		if item != expected[i] {
			t.Errorf("expected[%v] %v vs actual[%v] %v", i, expected[i], i, actual[i])
		}
	}
}
