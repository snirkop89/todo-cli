package todo

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
	"time"
	"unicode"
)

func (tm *TodoManager) parse(task string) Todo {
	var s scanner.Scanner
	r := strings.NewReader(task)
	s.Init(r)

	s.IsIdentRune = func(ch rune, i int) bool {
		return (ch == ':' && i > 0) ||
			unicode.IsLetter(ch) || unicode.IsDigit(ch) && i > 0
	}

	t := Todo{
		Task:        task,
		Completed:   false,
		CreatedAt:   time.Now().UTC(),
		CompletedAt: time.Time{},
	}

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		// If it's a special symbol with meaning for todo,
		// read the next ident and parse into the struct
		switch tok {
		case '@':
			_ = s.Scan()
			// todo s.TokenText()
		case '+':
			_ = s.Scan()
			t.Project = s.TokenText()
		case '#':
			_ = s.Scan()
			fmt.Printf("Tags: %s\n", s.TokenText())
		}

		switch s.TokenText() {
		case "due:":
			_ = s.Scan()
			due := tm.parseDueDate(s.TokenText())
			t.DueDate = due
		case "pri:":
			_ = s.Scan()
			t.SetPriority(s.TokenText())
		}
	}

	return t
}

func (tm *TodoManager) parseDueDate(d string) time.Time {
	switch d {
	case "tom", "tomorrow":
		return time.Now().AddDate(0, 0, 1).UTC()
	case "sun", "sunday":
		return nextAfter(time.Sunday)
	case "mon", "monday":
		return nextAfter(time.Monday)
	case "thu", "thursday":
		return nextAfter(time.Thursday)
	case "wed", "wednesday":
		return nextAfter(time.Wednesday)
	case "tue", "tuesday":
		return nextAfter(time.Tuesday)
	case "fri", "friday":
		return nextAfter(time.Friday)
	case "sat", "saturday":
		return nextAfter(time.Saturday)
	case "eow":
		// End of week
		return endOfWeek(tm.eow)
	case "eom":
		// TODO End of month
	}

	// If user specific a number, add this to the current date
	toAdd, err := strconv.Atoi(d)
	if err != nil {
		return time.Time{}
	}
	return time.Now().AddDate(0, 0, toAdd)
}

func nextAfter(weekday time.Weekday) time.Time {
	diff := int(weekday) - int(time.Now().Weekday())
	if diff <= 0 {
		diff += 7
	}
	return time.Now().AddDate(0, 0, diff)
}

func endOfWeek(endDay time.Weekday) time.Time {
	// Get next configured end of the week
	return nextAfter(endDay)
}
