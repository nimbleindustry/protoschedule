package protoschedule

// Package protoschedule implements functionality to define a weekly prototypical schedule in JSON format. You can then
// query the schedule to determine if a given time is within that schedule. Additionally, you can retreive all
// intervals within the schedule that a supplied time matches.

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

// The 'time' standard library has this wacky construct in which date parsing relies on specific
// numeral constants in order to indicate appropriate parsing. Yep, 1504 indicates 'military time' encoding.
const ANSIC_FORMAT = "1504"

// Struct used in the decoding of the JSON-encoded schedule
type IntervalJSONEncoding struct {
	Start    string `json:"start"`
	Duration string `json:"duration"`
	Label    string `json:"label"`
}

// Identifies unique intervals in the ScheduleDefinition
type IntervalValue struct {
	start time.Time
	end   time.Time
	label string
}

// Declarations and functions required for IntervalValue sorting
type ByStart []IntervalValue

func (a ByStart) Len() int           { return len(a) }
func (a ByStart) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByStart) Less(i, j int) bool { return a[i].start.Before(a[j].start) }

func (iv IntervalValue) String() string {
	return fmt.Sprintf("%s: %s ---> %s", iv.label, iv.start, iv.end)
}

// A ScheduleDefinition represents intervals which comprise a weekly schedule, calculated from a prototypical definition.
// Multiple entries covering the same intervals can be declared, and intervals can overflow into other intervals.
type ScheduleDefinition struct {
	Description string `json:"description"`
	Schedule    struct {
		Mon []IntervalJSONEncoding `json:"mon"`
		Tue []IntervalJSONEncoding `json:"tue"`
		Wed []IntervalJSONEncoding `json:"wed"`
		Thu []IntervalJSONEncoding `json:"thu"`
		Fri []IntervalJSONEncoding `json:"fri"`
		Sat []IntervalJSONEncoding `json:"sat"`
		Sun []IntervalJSONEncoding `json:"sun"`
	}
	currentJSON string
	first       time.Time
	last        time.Time
	intervals   []IntervalValue
}

// Construct a new ScheduleDefinition from the supplied json string
func New(jsonString string) (sd *ScheduleDefinition, err error) {
	return NewFromTime(jsonString, time.Now())
}

// Construct a new ScheduleDefinition, using the supplied json string and a specific time
func NewFromTime(jsonString string, t time.Time) (sd *ScheduleDefinition, err error) {
	sd = new(ScheduleDefinition)
	if sd.currentJSON != jsonString {
		sd.intervals = nil
		sd.currentJSON = jsonString
		err = json.Unmarshal([]byte(sd.currentJSON), sd)
		if err != nil {
			return
		}
	}
	err = sd.loadNew(t)
	return
}

// Loads new intervals from the defined prototype schedule, based on the time supplied
func (sd *ScheduleDefinition) loadNew(t time.Time) (err error) {
	epoch := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	for epoch.Weekday() != time.Monday {
		epoch = epoch.AddDate(0, 0, -1)
	}
	for _, value := range [...][]IntervalJSONEncoding{sd.Schedule.Mon, sd.Schedule.Tue, sd.Schedule.Wed, sd.Schedule.Thu, sd.Schedule.Fri, sd.Schedule.Sat, sd.Schedule.Sun} {
		for _, interval := range value {
			var start, base time.Time
			base, _ = time.Parse(ANSIC_FORMAT, "0000")
			start, err = time.Parse(ANSIC_FORMAT, interval.Start)
			if err != nil {
				sd = nil
				return
			}
			seconds := start.Sub(base)
			shiftStart := epoch.Add(seconds)
			var duration time.Duration
			duration, err = time.ParseDuration(interval.Duration)
			if err != nil {
				sd = nil
				return
			}
			duration -= (1 * time.Second)
			shiftEnd := shiftStart.Add(duration)
			sd.intervals = append(sd.intervals, IntervalValue{shiftStart, shiftEnd, interval.Label})
		}
		epoch = epoch.Add(time.Duration(24 * time.Hour))
	}
	sort.Sort(ByStart(sd.intervals))
	sd.first = sd.intervals[0].start
	sd.last = sd.intervals[len(sd.intervals)-1].end
	return
}

// Stringer implemention for printing
func (sd *ScheduleDefinition) String() string {
	var s string
	s = fmt.Sprintf("Schedule Defintion: %s, %d intervals, \nspan %s -> %s\n", sd.Description, len(sd.intervals), sd.first, sd.last)
	for _, v := range sd.intervals {
		s += fmt.Sprintf("%s\n", v)
	}
	return s
}

// Determines if the supplied time is within the defined prototypical schedule, using Second precision. This function
// allocates new intervals if the supplied time is outside the intervals currently calculated by previous lookups.
// TODO: Implement a binary search to speed this up
func (sd *ScheduleDefinition) Within(t time.Time) bool {
	if t.After(sd.last) || t.Before(sd.first) {
		sd.loadNew(t)
	}
	for _, v := range sd.intervals {
		if (t.After(v.start) || t.Unix() == v.start.Unix()) && (t.Before(v.end) || t.Unix() == v.end.Unix()) {
			return true
		}
	}
	return false
}

// Returns all matching intervals in the weekly schedule based on the supplied time. Protoschedule allows overlapping
// intervals (shifts)–you can differentiate intervals using the "Label" attribute.
func (sd *ScheduleDefinition) MatchingIntervals(t time.Time) []IntervalValue {
	var a []IntervalValue
	for _, v := range sd.intervals {
		if (t.After(v.start) || t.Unix() == v.start.Unix()) && (t.Before(v.end) || t.Unix() == v.end.Unix()) {
			a = append(a, v)
		}
	}
	return a
}
