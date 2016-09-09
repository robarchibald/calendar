# calendar
A simple Go library to calculate recurring appointment dates within a time period. Supports all types of recurrences in Outlook and Google calendar

## Getting Started
    go get https://github.com/robarchibald/calendar

 1. Create a Recurrence Struct (normally pulled from database)
 2. call GetOccurrences(startTime, endTime)


```
// Every 4th weekday
startTime := time.Date(2016, 1, 1, 12, 30, 0, 0, time.UTC)
dailyIsOnlyWeekday := true
r := Recurrence{
	StartDateTime:         startTime,
	RecurrencePatternCode: "D",
	RecurEvery:            4,
	DailyIsOnlyWeekday:    &dailyIsOnlyWeekday}
occurrences := r.GetOccurences(startTime, startTime.AddDate(0, 1, 0))
```
For additional examples, see recurrence_test.go

## Notes about the Recurrence Struct
The Recurrence struct is modeled after the recurring schedule data model used by both Microsoft Outlook and Google Calendar for recurring appointments. Just like Outlook, you can pick from Daily ("D"), Weekly ("W"), Monthly ("M") and Yearly ("Y") recurrence pattern codes. Each of those recurrence patterns then require the corresponding information to be filled in.

**All recurrences:**

 - StartDateTime - start time of the appointment. Should be set to the first desired occurence of the recurring appointment
 - RecurrencePatternCode - D: daily, W: weekly, M: monthly or Y: yearly
 - RecurEvery - number defining how many days, weeks, months or years to wait between recurrences
 - EndByDate (optional) - date by which recurrences must be done by 
 - NumberOfOccurrences (optional) - data for UI which can be used to store the number of recurrences. Has no effect in calculations though. EndByDate must be calculated based on NumberOfOccurrences

**Recurrence Pattern Code D (daily)**

 - DailyIsOnlyWeekday (optional) - ensure that daily occurrences only fall on weekdays (M, T, W, Th, F)

**Recurrence Pattern Code W (weekly)**

 - WeeklyDaysIncluded - binary value (converted to int16) to indicate days included (e.g. 0101010 or decimal 42 would be MWF). Each of the individual days are bitwise AND'd together to get the value.
	 - Sunday - 64 (1000000)
	 - Monday - 32 (0100000)
	 - Tuesday - 16 (0010000)
	 - Wednesday - 8 (0001000)
	 - Thursday - 4 (0000100)
	 - Friday - 2 (0000010)
	 - Saturday - 1 (0000001)

**Recurrence Pattern Code M (monthly)**

 - MonthlyWeekOfMonth - week of the month to recur on. e.g. Thanksgiving is always on the 4th week of the month. Must be used together with MonthlyDayOfWeek
 - MonthlyDayOfWeek - day of the week to recur on (0=Sunday, 1=Monday, 2=Tuesday, 3=Wednesday, 4=Thursday, 5=Friday, 6=Saturday). Must be used together with MonthlyWeekOfMonth
 **OR**
 - MonthlyDay - day of the month to recur on. e.g. 5 would recur on the 5th of every month

**Recurrence Pattern Code Y (yearly)**

 - YearlyMonth - month of the year to recur on (1=January, 2=February, 3=March, 4=April, 5=May, 6=June, 7=July)
 - MonthlyWeekOfMonth - week of the month to recur on. e.g. Thanksgiving is always on the 4th week of the month. Must be used together with MonthlyDayOfWeek
 - MonthlyDayOfWeek - day of the week to recur on (0=Sunday, 1=Monday, 2=Tuesday, 3=Wednesday, 4=Thursday, 5=Friday, 6=Saturday). Must be used together with MonthlyWeekOfMonth
 **OR**
 - MonthlyDay - day of the month to recur on. e.g. 5 would recur on the 5th of every month


![Outlook Recurrence Setup](https://robarchibald.github.io/calendar/images/outlookrecurrence.jpg)