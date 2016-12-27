package githubstat

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
)

const (
	DimensionOverall = "Overall"
	DimensionWeek    = "Week"
)

type Load interface {
}

var Config Configuration

type User struct {
	Name                  string
	StackalyticsDeviation int
}
type Configuration struct {
	StatBeginTime    time.Time
	AccessToken      string
	Users            []User
	Repos            []string
	Metrics          string
	Dimension        string
	WeekFirstDay     time.Weekday
	ThisWeekFirstDay time.Time
}

func getWeekFirstDay(t time.Time) time.Time {
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	for {
		weekday := t.Weekday()
		if weekday != Config.WeekFirstDay {
			t = t.Add(-24 * time.Hour)
		} else {
			break
		}
	}
	return t

}

// read config file
var _ = func() int {
	if _, err := toml.DecodeFile("./config.toml", &Config); err != nil {
		panic(err)
	}
	Config.ThisWeekFirstDay = getWeekFirstDay(time.Now())
	fmt.Printf("this week first day is : %v\n", Config.ThisWeekFirstDay)
	fmt.Printf("stat begin time is : %v\n", Config.StatBeginTime)
	return 0
}()
