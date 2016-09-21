package calendar

import (
	"math"
	"time"
)

type Recurrence struct {
	StartDate             time.Time  // Date to start Recurrence. Note that time and time zone information is NOT used in calculations
	RecurrencePatternCode string     // D for daily, W for weekly, M for monthly or Y for yearly
	RecurEvery            int16      // number of days, weeks, months or years between occurrences
	YearlyMonth           *int16     // month of the year to recur (applies only to RecurrencePatternCode: Y)
	MonthlyWeekOfMonth    *int16     // week of the month to recur. used together with MonthlyDayOfWeek (applies only to RecurrencePatternCode: M or Y)
	MonthlyDayOfWeek      *int16     // day of the week to recur. used together with MonthlyWeekOfMonth (applies only to RecurrencePatternCode: M or Y)
	MonthlyDay            *int16     // day of the month to recur (applies only to RecurrencePatternCode: M or Y)
	WeeklyDaysIncluded    *int16     // integer representing binary values AND'd together for 1000000-64 (Sun), 0100000-32 (Mon), 0010000-16 (Tu), 0001000-8 (W), 0000100-4 (Th), 0000010-2 (F), 0000001-1 (Sat). (applies only to RecurrencePatternCode: M or Y)
	DailyIsOnlyWeekday    *bool      // indicator that daily recurrences should only be on weekdays (applies only to RecurrencePatternCode: D)
	EndByDate             *time.Time // date by which all occurrences must end by. Note that time and time zone information is NOT used in calculations
}

func (r *Recurrence) GetOccurrences(timePeriodStart, timePeriodEnd time.Time) []time.Time {
	// Remove all time and time zone information from the recurrence start and end dates
	startDate := time.Date(r.StartDate.Year(), r.StartDate.Month(), r.StartDate.Day(), 0, 0, 0, 0, time.UTC)
	var endDate *time.Time
	if r.EndByDate != nil {
		end := time.Date(r.EndByDate.Year(), r.EndByDate.Month(), r.EndByDate.Day(), 0, 0, 0, 0, time.UTC)
		endDate = &end
	}
	switch {
	case r.RecurrencePatternCode == "D":
		dailyIsOnlyWeekday := false
		if r.DailyIsOnlyWeekday != nil {
			dailyIsOnlyWeekday = *r.DailyIsOnlyWeekday
		}
		return getDailyOccurrences(startDate, int(r.RecurEvery), dailyIsOnlyWeekday, endDate, timePeriodStart, timePeriodEnd)
	case r.RecurrencePatternCode == "W":
		var weeklyDaysIncluded int16 = 127 // all days
		if r.WeeklyDaysIncluded != nil {
			weeklyDaysIncluded = *r.WeeklyDaysIncluded
		}
		return getWeeklyOccurrences(startDate, int(r.RecurEvery), getIncludedWeeklyDays(weeklyDaysIncluded), endDate, timePeriodStart, timePeriodEnd)
	case r.RecurrencePatternCode == "M":
		return getMonthlyOccurrences(startDate, int(r.RecurEvery), r.MonthlyDay, r.MonthlyDayOfWeek, r.MonthlyWeekOfMonth, endDate, timePeriodStart, timePeriodEnd)
	case r.RecurrencePatternCode == "Y":
		return getYearlyOccurrences(startDate, int(r.RecurEvery), r.YearlyMonth, r.MonthlyDay, r.MonthlyDayOfWeek, r.MonthlyWeekOfMonth, endDate, timePeriodStart, timePeriodEnd)
	}
	return []time.Time{}
}

func (r *Recurrence) IsValidOccurrenceDate(occurrenceDate time.Time) bool {
	// Remove all time and time zone information from the occurrenceDate
	date := time.Date(occurrenceDate.Year(), occurrenceDate.Month(), occurrenceDate.Day(), 0, 0, 0, 0, time.UTC)
	occurrences := r.GetOccurrences(date, date)
	return len(occurrences) == 1 && occurrences[0] == date
}

func getDailyOccurrences(recurrenceStartDate time.Time, recurEvery int, dailyIsOnlyWeekday bool, recurrenceEndByDate *time.Time, timePeriodStart, timePeriodEnd time.Time) []time.Time {
	recurrences := []time.Time{}
	currentDate := recurrenceStartDate
	if currentDate.Before(timePeriodStart) {
		if dailyIsOnlyWeekday {
			currentDate = getWeekdayStartTime(recurrenceStartDate, recurEvery, timePeriodStart)
		} else {
			currentDate = getDailyStartTime(recurrenceStartDate, recurEvery, timePeriodStart)
		}
	}
	for currentDate.Before(timePeriodEnd) || currentDate.Equal(timePeriodEnd) {
		recurrences = append(recurrences, currentDate)
		if dailyIsOnlyWeekday {
			currentDate = addWeekdays(int(recurEvery), currentDate)
		} else {
			currentDate = currentDate.AddDate(0, 0, int(recurEvery))
		}
	}
	return recurrences
}

func getDailyStartTime(recurrenceStartDate time.Time, recurEvery int, timePeriodStart time.Time) time.Time {
	days := getDays(recurrenceStartDate, timePeriodStart)
	return recurrenceStartDate.AddDate(0, 0, getStartAdder(days, recurEvery)+days)
}

/**********************************************************************************************************
Daily recurring meetings only on weekdays is supported by Outlook UI, but not Google calendar although you
   can get the same fuctionality in Google with a weekly meeting on M,T,W,Th,F

Daily recurring meetings only on weekdays that recur every N number of days is not creatable by either
   Outlook or Google calendar UI's, but can be viewed by both since they both support the ICalendar spec.
   The below rule yields a meeting every other weekday

FREQ=DAILY;BYDAY=MO,TU,WE,TH,FR;INTERVAL=2

From https://www.ietf.org/rfc/rfc2445.txt
   The BYDAY rule part specifies a COMMA character (US-ASCII decimal 44) separated list of days of the week;
   MO indicates Monday; TU indicates Tuesday; WE indicates Wednesday; TH indicates Thursday; FR indicates
   Friday; SA indicates Saturday; SU indicates Sunday.

   BYxxx rule parts modify the recurrence in some manner. BYxxx rule parts for a period of time which is the
   same or greater than the frequency generally reduce or limit the number of occurrences of the recurrence
   generated. For example, "FREQ=DAILY;BYMONTH=1" reduces the number of recurrence instances from all days
   (if BYMONTH tag is not present) to all days in January
************************************************************************************************************/
func getWeekdayStartTime(recurrenceStartDate time.Time, recurEvery int, timePeriodStart time.Time) time.Time {
	days := getDays(recurrenceStartDate, timePeriodStart)
	weekdays := getWeekdays(days, recurrenceStartDate)
	startDateTime := recurrenceStartDate.AddDate(0, 0, getStartAdder(weekdays, recurEvery)+days)
	if startDateTime.Weekday() == time.Sunday || startDateTime.Weekday() == time.Saturday {
		startDateTime = startDateTime.AddDate(0, 0, 2) // add 2 days either way.  if Saturday, we need two to get to Monday, if Sunday then 2 to make up for Saturday which is also not a weekday
	}
	return startDateTime
}

func getDays(recurrenceStartDate, timePeriodStart time.Time) int {
	return int(math.Ceil(timePeriodStart.Sub(recurrenceStartDate).Hours() / 24)) // include timePeriodStart even though it is midnight
}

func getWeekdays(days int, firstOccurrence time.Time) int {
	weeks := days / 7
	weekdays := weeks * 5
	extradays := days - weeks*7 // number of days past the full weeks (add extra weekdays below)
	for i := 1; i <= extradays; i++ {
		tmpDate := firstOccurrence.AddDate(0, 0, i)
		if tmpDate.Weekday() != time.Sunday && tmpDate.Weekday() != time.Saturday {
			weekdays++ // date is a weekday so add it to the total weekdays
		}
	}
	return weekdays
}

func addWeekdays(weekdays int, startDate time.Time) time.Time {
	endTime := startDate
	for i := 1; weekdays > 0; i++ {
		endTime = startDate.AddDate(0, 0, i)
		if endTime.Weekday() != time.Sunday && endTime.Weekday() != time.Saturday {
			weekdays-- // date is a weekday so add it to the total weekdays
		}
	}
	return endTime
}

// Recurrence makes it so that we skip days in the calendar and may not start
// at the beginning of the time period we're looking at, so calculate how
// many we need to add to get to our first recurrence after the start
func getStartAdder(days, recurEvery int) int {
	startAdder := int(math.Mod(float64(days), float64(recurEvery)))
	if startAdder != 0 {
		startAdder = recurEvery - startAdder
	}
	return startAdder
}

func getIncludedWeeklyDays(weeklyDaysIncluded int16) []time.Weekday {
	var days []time.Weekday
	if weeklyDaysIncluded&64 != 0 {
		days = append(days, time.Sunday)
	}
	if weeklyDaysIncluded&32 != 0 {
		days = append(days, time.Monday)
	}
	if weeklyDaysIncluded&16 != 0 {
		days = append(days, time.Tuesday)
	}
	if weeklyDaysIncluded&8 != 0 {
		days = append(days, time.Wednesday)
	}
	if weeklyDaysIncluded&4 != 0 {
		days = append(days, time.Thursday)
	}
	if weeklyDaysIncluded&2 != 0 {
		days = append(days, time.Friday)
	}
	if weeklyDaysIncluded&1 != 0 {
		days = append(days, time.Saturday)
	}
	return days
}

func getWeeklyOccurrences(recurrenceStartDate time.Time, recurEvery int, daysIncluded []time.Weekday, recurrenceEndByDate *time.Time, timePeriodStart, timePeriodEnd time.Time) []time.Time {
	recurrences := []time.Time{}
	currentDate := recurrenceStartDate
	if currentDate.Before(timePeriodStart) {
		currentDate = getWeeklyStartTime(recurrenceStartDate, recurEvery, timePeriodStart)
	} else {
		currentDate = currentDate.AddDate(0, 0, -1*int(currentDate.Weekday())) // turn into beginning of week
	}
	for (currentDate.Before(timePeriodEnd) || currentDate.Equal(timePeriodEnd)) && (recurrenceEndByDate == nil || currentDate.Before(*recurrenceEndByDate)) {
		recurrences = append(recurrences, getIncludedDays(daysIncluded, currentDate, timePeriodStart, timePeriodEnd)...)
		currentDate = currentDate.AddDate(0, 0, 7*(recurEvery))
	}
	return recurrences
}

func getIncludedDays(daysIncluded []time.Weekday, startDate, timePeriodStart, timePeriodEnd time.Time) []time.Time {
	days := []time.Time{}
	for _, day := range daysIncluded {
		date := startDate.AddDate(0, 0, int(day))
		if date.After(timePeriodEnd) {
			break
		}
		if date.After(timePeriodStart) || date.Equal(timePeriodStart) {
			days = append(days, date)
		}
	}
	return days
}

func getWeeklyStartTime(recurrenceStartDate time.Time, recurEvery int, timePeriodStart time.Time) time.Time {
	weekStartDate := recurrenceStartDate.AddDate(0, 0, -1*int(recurrenceStartDate.Weekday())) // turn into beginning of week
	weeks := getWeeks(weekStartDate, timePeriodStart)
	adder := getStartAdder(weeks, recurEvery)
	return weekStartDate.AddDate(0, 0, 7*(adder+weeks))
}

func getWeeks(fromDate, toDate time.Time) int {
	return int(math.Floor(toDate.Sub(fromDate).Hours() / 24 / 7)) // include toDate even though it is midnight
}

func getMonthlyOccurrences(recurrenceStartDate time.Time, recurEvery int, monthlyDay, monthlyDayOfWeek, monthlyWeekOfMonth *int16, recurrenceEndByDate *time.Time, timePeriodStart, timePeriodEnd time.Time) []time.Time {
	recurrences := []time.Time{}
	currentDate := recurrenceStartDate
	if currentDate.Before(timePeriodStart) {
		currentDate = getMonthlyStartTime(recurrenceStartDate, recurEvery, timePeriodStart)
	}
	for (currentDate.Before(timePeriodEnd) || currentDate.Equal(timePeriodEnd)) && (recurrenceEndByDate == nil || currentDate.Before(*recurrenceEndByDate)) {
		recurrences = append(recurrences, getMonthOccurrence(currentDate, timePeriodStart, timePeriodEnd, monthlyDay, monthlyDayOfWeek, monthlyWeekOfMonth)...)
		currentDate = currentDate.AddDate(0, recurEvery, 0)
	}
	return recurrences
}

func getMonthOccurrence(startDate, timePeriodStart, timePeriodEnd time.Time, monthlyDay, monthlyDayOfWeek, monthlyWeekOfMonth *int16) []time.Time {
	var occurrence time.Time
	if monthlyDay != nil {
		occurrence = time.Date(startDate.Year(), startDate.Month(), int(*monthlyDay), startDate.Hour(), startDate.Minute(), startDate.Second(), startDate.Nanosecond(), startDate.Location())
	} else if monthlyDayOfWeek != nil && monthlyWeekOfMonth != nil {
		weekAdder := *monthlyWeekOfMonth
		if *monthlyDayOfWeek >= int16(startDate.Weekday()) { // first of my desired day of week occurs in first week
			weekAdder--
		}
		occurrence = startDate.AddDate(0, 0, int(7*weekAdder+*monthlyDayOfWeek)-int(startDate.Weekday()))
	}
	if (occurrence.Before(timePeriodEnd) || occurrence.Equal(timePeriodEnd)) && (occurrence.After(timePeriodStart) || occurrence.Equal(timePeriodStart)) {
		return []time.Time{occurrence}
	}
	return []time.Time{}
}

func getMonthlyStartTime(recurrenceStartDate time.Time, recurEvery int, timePeriodStart time.Time) time.Time {
	monthStartDate := recurrenceStartDate.AddDate(0, 0, -1*int(recurrenceStartDate.Day()-1)) // turn into beginning of month
	months := getMonths(monthStartDate, timePeriodStart)
	adder := getStartAdder(months, recurEvery)
	return monthStartDate.AddDate(0, adder+months, 0)
}

func getMonths(fromDate, toDate time.Time) int {
	years := toDate.Year() - fromDate.Year()
	months := int(toDate.Month() - fromDate.Month())
	days := int(toDate.Day() - fromDate.Day())
	if days < 0 {
		months--
	}

	return years*12 + months
}

func getYearlyOccurrences(recurrenceStartDate time.Time, recurEvery int, yearlyMonth, monthlyDay, monthlyDayOfWeek, monthlyWeekOfMonth *int16, recurrenceEndByDate *time.Time, timePeriodStart, timePeriodEnd time.Time) []time.Time {
	recurrences := []time.Time{}
	currentDate := recurrenceStartDate
	if currentDate.Before(timePeriodStart) {
		currentDate = getYearlyStartTime(recurrenceStartDate, yearlyMonth, recurEvery, timePeriodStart)
	}
	for (currentDate.Before(timePeriodEnd) || currentDate.Equal(timePeriodEnd)) && (recurrenceEndByDate == nil || currentDate.Before(*recurrenceEndByDate)) {
		recurrences = append(recurrences, getMonthOccurrence(currentDate, timePeriodStart, timePeriodEnd, monthlyDay, monthlyDayOfWeek, monthlyWeekOfMonth)...)
		currentDate = time.Date(currentDate.Year()+recurEvery, time.Month(*yearlyMonth), 1, currentDate.Hour(), currentDate.Minute(), currentDate.Second(), currentDate.Nanosecond(), currentDate.Location())
	}
	return recurrences
}

func getYearlyStartTime(recurrenceStartDate time.Time, yearlyMonth *int16, recurEvery int, timePeriodStart time.Time) time.Time {
	yearStartDate := time.Date(recurrenceStartDate.Year(), time.Month(*yearlyMonth), 1, recurrenceStartDate.Hour(), recurrenceStartDate.Minute(), recurrenceStartDate.Second(), recurrenceStartDate.Nanosecond(), recurrenceStartDate.Location())
	years := getYears(yearStartDate, timePeriodStart)
	adder := getStartAdder(years, recurEvery)
	return yearStartDate.AddDate(adder+years, 0, 0)
}

func getYears(fromDate, toDate time.Time) int {
	years := toDate.Year() - fromDate.Year()
	months := int(toDate.Month() - fromDate.Month())
	if months < 0 {
		years--
	}

	return years
}
