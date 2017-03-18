package main

import (
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
	"time"
)

var client *github.Client = nil
var ctx context.Context

func InitClient() {
	if client != nil {
		return
	}
	token, err := ioutil.ReadFile("accessToken.txt")
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	ctx = context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: string(token)},
	)

	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)
}

func getOwnedRepos() []*github.Repository {
	InitClient()
	var ret []*github.Repository
	opt := &github.RepositoryListOptions{
		Type: "owner",
	}
	repos, _, err := client.Repositories.List(ctx, "", opt)
	if err != nil {
		println(err.Error())
		return ret
	}
	for _, r := range repos {
		ret = append(ret, r)
	}
	return ret
}

func timePtr(t time.Time) *time.Time { return &t }

func CreateMilestone(title string, desc string, date time.Time) {
	var allRepos []*github.Repository

	allRepos = getOwnedRepos()
	m := &github.Milestone{
		Title:       github.String(title),
		Description: github.String(desc),
		DueOn:       timePtr(date),
	}
	for _, r := range allRepos {
		m, _, err := client.Issues.CreateMilestone(ctx, *r.Owner.Login, *r.Name, m)
		if err != nil {
			println(err.Error())
		} else {
			println("Created milestone", *m.Title, "at", *m.HTMLURL)
		}
	}
}

func RemoveMilestone(title string) {
	for _, repo := range getOwnedRepos() {
		milestones, _, err := client.Issues.ListMilestones(ctx, *repo.Owner.Login, *repo.Name, nil)
		if err != nil {
			println(err.Error())
			return
		}
		for _, milestone := range milestones {
			if *milestone.Title == title {
				println("Removed milestone", *milestone.Title, "n°", milestone.GetNumber(), "from repository", *repo.Name)
				_, err := client.Issues.DeleteMilestone(ctx, *repo.Owner.Login, *repo.Name, *milestone.Number)
				if err != nil {
					println(err.Error())
				}
			}
		}
	}
}
