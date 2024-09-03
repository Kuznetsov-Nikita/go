//go:build !solution

package hotelbusiness

import "sort"

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

type Event struct {
	Date        int
	GuestChange int
}

func ComputeLoad(guests []Guest) []Load {
	if len(guests) == 0 {
		return []Load{}
	}

	var events []Event

	for _, guest := range guests {
		events = append(events, Event{Date: guest.CheckInDate, GuestChange: 1})
		events = append(events, Event{Date: guest.CheckOutDate, GuestChange: -1})
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Date < events[j].Date
	})

	var result []Load
	var date, load = events[0].Date, 0

	for _, event := range events {
		if event.Date != date {
			if len(result) == 0 || result[len(result)-1].GuestCount != load {
				result = append(result, Load{StartDate: date, GuestCount: load})
			}
			date = event.Date
		}
		load += event.GuestChange
	}
	if len(result) == 0 || result[len(result)-1].GuestCount != load {
		result = append(result, Load{StartDate: date, GuestCount: load})
	}

	return result
}
