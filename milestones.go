package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
	"regexp"
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

func repoMatchesMask(repoName, mask string) bool {
	matched, err := regexp.MatchString(mask, repoName)
	if err != nil {
		fmt.Println("Error in Regex parsing", err)
		return false
	}
	return matched
}

func CreateMilestone(title string, desc string, date time.Time, mask string) {
	var repositories []*github.Repository

	repositories = getOwnedRepos()
	m := &github.Milestone{
		Title:       github.String(title),
		Description: github.String(desc),
		DueOn:       timePtr(date),
	}
	for _, r := range repositories {
		if repoMatchesMask(r.GetName(), mask) {
			m, _, err := client.Issues.CreateMilestone(ctx, *r.Owner.Login, *r.Name, m)
			if err != nil {
				println(err.Error())
			} else {
				println("Created milestone", *m.Title, "n°", m.GetNumber(), "at ", *m.HTMLURL)
			}
		}
	}
}

func RemoveMilestone(title, mask string) {
	for _, r := range getOwnedRepos() {
		if repoMatchesMask(r.GetName(), mask) {
			milestones, _, err := client.Issues.ListMilestones(ctx, *r.Owner.Login, *r.Name, nil)
			if err != nil {
				println(err.Error())
				return
			}
			for _, milestone := range milestones {
				if *milestone.Title == title {
					println("Removed milestone", *milestone.Title, "n°", milestone.GetNumber(), "from repository", *r.Name)
					_, err := client.Issues.DeleteMilestone(ctx, *r.Owner.Login, *r.Name, *milestone.Number)
					if err != nil {
						println(err.Error())
					}
				}
			}
		}
	}
}
