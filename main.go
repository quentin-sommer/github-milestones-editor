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
	mask := flag.String("mask", "", "regex mask to match repos")

	flag.Parse()

	if *rm {
		if *title == "" {
			println("title is required when removing")
			return
		}
		RemoveMilestone(*title, *mask)
		return
	}
	if *title == "" || *description == "" {
		println("title and desc")
		return
	}

	if *date != "" {
		t, err := time.Parse("2006-01-02", *date)
		if err != nil {
			println("Error : date should be formated like this : yyyy-mm-dd")
			return
		}
		CreateMilestone(*title, *description, &t, *mask)
	} else {
		CreateMilestone(*title, *description, nil, *mask)
	}

}
