package crawler

import "testing"

func TestCrawlGithubUsersMulti(t *testing.T) {
	usersList := CrawlGithubUsersMulti(0, 10, 3)
	if len(usersList) == 0 {
		t.Error("CrawlGithubUsersMulti Error: ", usersList)
	}
	for _, user := range usersList {
		t.Log(user)
	}
}
