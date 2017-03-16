package main

import (
	"context"
	"github.com/google/go-github/github"
)

import (
	"golang.org/x/oauth2"
	"regexp"
	"time"
)

const accessToken string = "5c76b22dd318e9973a5f7b742bb4b7be365ccabd"
const matcher = "opentrends"

var client *github.Client = nil
var ctx context.Context

func InitClient() {
	if client != nil {
		return
	}
	ctx = context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)

	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)
}
func getOwnedRepos() []*github.Repository {
	var ret []*github.Repository
	opt := &github.RepositoryListOptions{
		Type: "owner",
	}
	repos, _, err := client.Repositories.List(ctx, "", opt)
	if err != nil {
		println(err)
		return ret
	}
	for _, r := range repos {
		if matched, _ := regexp.Match("opentrends.*", []byte(*r.Name)); matched {
			ret = append(ret, r)
		}
	}
	return ret
}

func getMemberRepos() []*github.Repository {
	var ret []*github.Repository
	opt := &github.RepositoryListOptions{
		Type: "member",
	}
	repos, _, err := client.Repositories.List(ctx, "", opt)
	if err != nil {
		println(err)
		return ret
	}
	for _, r := range repos {
		if matched, _ := regexp.Match("opentrends.*", []byte(*r.Name)); matched {
			ret = append(ret, r)
		}
	}

	return ret
}
func timePtr(t time.Time) *time.Time { return &t }

func CreateMilestone(title string, desc string) {
	InitClient()
	var allRepos []*github.Repository

	allRepos = append(getOwnedRepos(), getMemberRepos()...)

	m := &github.Milestone{
		Title:       github.String(title),
		Description: github.String(desc),
		DueOn:       timePtr(time.Date(2017, 03, 25, 0, 0, 0, 0, time.UTC)),
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
	InitClient()
	var allRepos []*github.Repository

	allRepos = append(getOwnedRepos(), getMemberRepos()...)
	for _, r := range allRepos {
		milestones, _, err := client.Issues.ListMilestones(ctx, *r.Owner.Login, *r.Name, nil)
		if err != nil {
			println(err.Error())
			return
		}
		for _, milestone := range milestones {
			if *milestone.Title == title {
				println("Removed milestone", *milestone.Title, "from repository", *r.Name)
				_, err := client.Issues.DeleteMilestone(ctx, *r.Owner.Login, *r.Name, *milestone.Number)

				if err != nil {
					println(err.Error())
				}
			}
		}
	}
}
