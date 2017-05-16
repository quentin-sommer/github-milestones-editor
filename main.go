package main

import (
	"flag"
	"time"
)

func main() {
	title := flag.String("title", "", "milestone title")
	description := flag.String("desc", "", "milestone desription")
	date := flag.String("date", "", "milestone due date")
	rm := flag.Bool("remove", false, "remove milestone")

	flag.Parse()

	if *rm {
		if *title == "" {
			println("title is required when removing")
			return
		}
		RemoveMilestone(*title)
		return
	}
	if *title == "" || *description == "" || *date == "" {
		println("title, desc and due are required")
		return
	}

	t, err := time.Parse("2006-01-02", *date)
	if err != nil {
		println("Error : date should be formated like this : yyyy-mm-dd")
		return
	}
	CreateMilestone(*title, *description, t)
}
