package crawler

import (
	"GatherGithubInfo/config"
	"GatherGithubInfo/models"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
	"sync"
)

// CrawlGithubUsers starts multiple goroutines to crawl GitHub users concurrently.
func CrawlGithubUsersMulti(startNum, perPage, numWorkers int) []models.User {
	targetUrl := fmt.Sprintf(config.GithubUsersApi, startNum, perPage)
	fmt.Println(targetUrl)
	resp, err := http.Get(targetUrl)
	if err != nil {
		logrus.Error("CrawlGithubUsers Error: ", err)
		return nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("CrawlGithubUsers ReadAll Error: ", err)
		return nil
	}

	var users []models.User
	err = json.Unmarshal(body, &users)
	if err != nil {
		logrus.Error("CrawlGithubUsers decoding JSON Error: ", err)
		return nil
	}

	usersList := make([]models.User, len(users))
	userChan := make(chan models.User)
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for user := range userChan {
				userDetails := CrawlGithubUserDetailMulti(user.Login)
				userDetails.EmailRepo = CrawlGithubRepoEmailMulti(userDetails.Login, 1)
				usersList = append(usersList, userDetails)
			}
		}()
	}

	// Send users to the workers
	for _, user := range users {
		userChan <- user
	}
	close(userChan)

	wg.Wait() // Wait for all workers to finish
	return usersList
}

func CrawlGithubUserDetailMulti(user string) models.User {
	// 2. 获取用户详细信息
	targetUrl := fmt.Sprintf(config.GithubUserDetailApi, user)
	resp, err := http.Get(targetUrl)
	if err != nil {
		logrus.Error("CrawlGithubUserDetailMulti Error: ", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("CrawlGithubUserDetailMulti ReadAll Error: ", err)
	}
	var userDetail models.User
	//fmt.Println(string(body))
	err = json.Unmarshal(body, &userDetail)
	if err != nil {
		fmt.Println("CrawlGithubUserDetailMulti decoding JSON Error: ", err)
		return models.User{}
	}
	return userDetail
}

func CrawlGithubReposMulti(user string, perPage int) string {
	// 爬取 Github 仓库信息
	// 1. 获取 Github 仓库列表
	targetUrl := fmt.Sprintf(config.GithubUserReposApi, user, perPage)
	logrus.Println(targetUrl)
	resp, err := http.Get(targetUrl)
	if err != nil {
		logrus.Error("CrawlGithubReposMulti Error: ", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("CrawlGithubReposMulti ReadAll Error: ", err)
	}
	// 提取json中的commit author email
	// 解析json数据
	var data []map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return ""
	}
	//fmt.Println(data)
	commitsUrl, flag := strings.CutSuffix(data[0]["commits_url"].(string), "{/sha}")
	if !flag {
		logrus.Error("CrawlGithubReposMulti Unable To Find Commits_URL Error: ", flag)
	}
	//fmt.Println(commitsUrl)
	return commitsUrl
}

func CrawlGithubRepoEmailMulti(user string, perPage int) string {
	// 爬取 Github 仓库提交者的邮箱
	commitsUrl := CrawlGithubReposMulti(user, perPage)
	// 2. 获取仓库信息，主要为了获取提交者的邮箱
	resp, err := http.Get(commitsUrl + "?per_page=1")
	if err != nil {
		logrus.Error("CrawlGithubReposMulti Error: ", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("CrawlGithubReposMulti ReadAll Error: ", err)
		return ""
	}
	var data []map[string]interface{}

	// 提取json中的commit author email
	// 解析json数据
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return ""
	}
	emailRepo := data[0]["commit"].(map[string]interface{})["author"].(map[string]interface{})["email"]
	return emailRepo.(string)

}
