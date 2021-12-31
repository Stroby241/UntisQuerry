package untis

import (
	"fmt"
	"github.com/Stroby241/UntisAPI"
	"github.com/Stroby241/UntisQuerry/event"
	"github.com/Stroby241/UntisQuerry/state"
	"time"
)

var user *UntisAPI.User

func Init() {
	event.On(event.EventLogin, func(data interface{}) {
		strings := data.([4]string)
		success := login(strings[0], strings[1], strings[2], strings[3])
		event.Go(event.EventLoginResult, success)
	})

	event.On(event.EventLogout, func(data interface{}) {
		logout()
	})

	event.On(event.EventAddTeacher, func(data interface{}) {
		strings := data.([2]string)
		success := addTeacher(strings[0], strings[1])
		if success {
			event.Go(event.EventUpdateQuerryPanel, nil)
			event.Go(event.EventSetPage, state.PageQuerry)
		}
	})

	event.On(event.EventLoadTimeTable, func(data interface{}) {
		times := data.([2]time.Time)
		loadTimetable(times[0], times[1])
	})

	event.On(event.EventQuerryTaecher, func(data interface{}) {
		queryTeacher(data.(*state.Teacher))
	})
}

func login(username string, password string, school string, server string) bool {
	user = UntisAPI.NewUser(username, password, school, server)
	err := user.Login()
	if err != nil {
		user = nil
		return false
	}

	return initCalls()
}

var rooms map[int]UntisAPI.Room
var classes map[int]UntisAPI.Class

func initCalls() bool {
	if user == nil {
		return false
	}
	var err error

	rooms, err = user.GetRooms()
	if err != nil {
		return false
	}

	classes, err = user.GetClasses()
	if err != nil {
		return false
	}

	return true
}

func logout() {
	if user == nil {
		return
	}
	user.Logout()
	user = nil

	event.Go(event.EventSetPage, state.PageStart)
}

var timetable []map[int]UntisAPI.Period
var startDate int
var endDate int

func loadTimetable(startTime time.Time, endTime time.Time) bool {
	newStartDate := UntisAPI.ToUntisDate(startTime)
	newEndDate := UntisAPI.ToUntisDate(endTime)

	if startDate == newStartDate && endDate == newEndDate {
		return true
	}

	startDate = newStartDate
	endDate = newEndDate

	timetable = []map[int]UntisAPI.Period{}
	counter := 0

	event.Go(event.EventStartLoading, "Timetable")
	for _, room := range rooms {
		fmt.Printf("Loading timetable of room: %d of %d. \r", counter, len(rooms))
		event.Go(event.EventUpdateLoading, float64(counter)/float64(len(rooms))*100.0)

		periods, err := user.GetTimeTable(room.Id, 4, startDate, endDate)
		if err != nil {
			return false
		}

		if periods != nil {
			timetable = append(timetable, periods)
		}

		counter++
	}
	event.Go(event.EventSetPage, state.PageQuerry)

	return true
}
