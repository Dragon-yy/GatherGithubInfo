package crawler

import (
	"fmt"
	"testing"
)

func TestCrawlGithubUsers(t *testing.T) {
	usersList := CrawlGithubUsers(100000, 20, 100040)
	if len(usersList) == 0 {
		t.Error("CrawlGithubUsers Error: ", usersList)
	}
	for _, user := range usersList {
		t.Log(user)
	}
}

func TestCrawlGithubRepoEmail(t *testing.T) {
	email := CrawlGithubRepoEmail("mojombo", 1)
	fmt.Println(email)
}
