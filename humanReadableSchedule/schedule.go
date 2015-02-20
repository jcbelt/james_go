package humanReadableSchedule

import (
    "time"
    "sort"
    "github.com/nowaitapp/guestApi"
)

type dayNum int

// Funky way of doing enums in go
const (
    S dayNum = iota
    M
    T
    W
    Th
    F
    Sa
)

var dayStrings = [...]string {
     "S",
     "M",
     "T",
     "W",
     "Th",
     "F",
     "S",
}



// period of time when restaurant is open
// we allow endtimes past midnight
type OpenHours struct {
    StartTime time.Time
    EndTime time.Time
} 

// array of open hour periods for a single day
type DailyOpenHours []OpenHours


// array of open hours indexed by day of week
type HumanReadableSchedule [7]DailyOpenHours


// Sorting interface implementation for DailyOpenHours
func (s DailyOpenHours) Len() int {
    return len(s)
}
func (s DailyOpenHours) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
func (s DailyOpenHours) Less(i, j int) bool {
    return s[i].StartTime.Before(s[j].StartTime)
}



// Get HumanReadableSchedule for a guestApi.Schedule
func New(apiSchedule guestApi.Schedule) (schedule HumanReadableSchedule){

    for _, period := range apiSchedule.Periods {

        // specify format for parsing time strings
        form := "15:04:00"

        startDay := period.Open.Day
        startTime, _ := time.Parse(form, period.Open.Time)

        endDay := period.Close.Day
        endTime, _ := time.Parse(form, period.Close.Time)

        // Determine if the period spans multiple days
        // Handle wrap arround from sat to sunday
        var daySpan int
        if (endDay >= startDay){
            daySpan = (endDay - startDay)
        } else {
            daySpan = (7 - startDay + endDay)
        }

        // Determine if we need to split this period into multiple daily hours
        // We allow the end times to go past midnight for a given day (i.e M:10pm - 1am)
        // but if we go past 4am we will consider that the start of a new day
        startOfNewDay, _ := time.Parse(form, "04:00:00")

        // split periods if necesary until we are no longer spaning multiple days
        var openHours OpenHours
        for daySpan >=0 {
            if ( daySpan < 1 || (daySpan == 1 && endTime.Before(startOfNewDay))){
                // this is the normal case where the start and end times are on the same day
                // or the end time is "late night" hours on the next day
                openHours = OpenHours{startTime, endTime}
                daySpan = 0
            } else {
                // since the period spans past the start of a new day we will split it up
                // and carry the remander of the period over to a new day
                openHours = OpenHours{startTime, startOfNewDay}
            }

            schedule[startDay] = append(schedule[startDay], openHours)

            startTime = startOfNewDay
            startDay = (startDay + 1) % 7 // wrap sat to sun
            daySpan -= 1
        }
    }

    // sort the open hours for each day by start time
    for _, dailyOpenHours := range schedule {
        sort.Sort(dailyOpenHours)
    }

    return
}

// Convert the schedule to a string
func (schedule *HumanReadableSchedule) String() string{
    format := "3:04 PM"
    var scheduleString string
    for dayNum, dailyOpenHours := range schedule {
        scheduleString += dayStrings[dayNum] + ": "
        for _,openHours := range dailyOpenHours {
            scheduleString += openHours.StartTime.Format(format) + " - " + openHours.EndTime.Format(format)
        }
        scheduleString += "\n"

    }
    
    return scheduleString
}

