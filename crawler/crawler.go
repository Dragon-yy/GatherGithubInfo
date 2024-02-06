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
)

// SetHeaders sets common headers for simulating a browser request
func SetHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	if config.GithubToken != "" {
		logrus.Info("Github Token Set: " + config.GithubToken)
		req.Header.Add("Authorization", "Bearer "+config.GithubToken)
	}
}

func MakeRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	SetHeaders(req)
	return http.DefaultClient.Do(req)
}

func CrawlGithubUsers(startNum, perPage, endNum int) []models.User {
	// 爬取 Github 用户信息
	// 1. 获取 Github 用户列表
	var usersList = []models.User{}
	for startNum < endNum {
		targetUrl := fmt.Sprintf(config.GithubUsersApi, startNum, perPage)
		fmt.Println(targetUrl)
		resp, err := MakeRequest("GET", targetUrl, nil)

		if err != nil {
			logrus.Error("CrawlGithubUsers Error: ", err)
		}
		if resp.StatusCode != http.StatusOK {
			logrus.Error("CrawlGithubUsers StatusCode Error: ", resp.StatusCode)
			return nil
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logrus.Error("CrawlGithubUsers ReadAll Error: ", err)
		}
		//fmt.Println(string(body))
		var users []models.User
		//decoder := json.NewDecoder(resp.Body)
		//err = decoder.Decode(&users)
		err = json.Unmarshal(body, &users)
		if err != nil {
			fmt.Println("CrawlGithubUsers decoding JSON Error: ", err)
			return nil
		}

		// Print the decoded User data
		//for _, user := range users {
		//	fmt.Println(user)
		//}

		for _, user := range users {
			// 2. 获取用户详细信息
			user = CrawlGithubUserDetail(user.Login)
			user.EmailRepo = CrawlGithubRepoEmail(user.Login, 1)
			usersList = append(usersList, user)
		}
		startNum += perPage
		//return usersList
	}
	return usersList
}

func CrawlGithubUserDetail(user string) models.User {
	// 2. 获取用户详细信息
	targetUrl := fmt.Sprintf(config.GithubUserDetailApi, user)
	resp, err := MakeRequest("GET", targetUrl, nil)
	if err != nil {
		logrus.Error("CrawlGithubUserDetail Error: ", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("CrawlGithubUserDetail ReadAll Error: ", err)
	}
	var userDetail models.User
	//fmt.Println(string(body))
	err = json.Unmarshal(body, &userDetail)
	if err != nil {
		logrus.Error("CrawlGithubUserDetail decoding JSON Error: ", err)
		return models.User{}
	}
	return userDetail
}

func CrawlGithubRepos(user string, perPage int) string {
	// 爬取 Github 仓库信息
	// 1. 获取 Github 仓库列表
	targetUrl := fmt.Sprintf(config.GithubUserReposApi, user)
	logrus.Println(targetUrl)
	resp, err := MakeRequest("GET", targetUrl, nil)
	if err != nil {
		logrus.Error("CrawlGithubRepos Error: ", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("CrawlGithubRepos ReadAll Error: ", err)
	}
	// 提取json中的commit author email
	// 解析json数据
	var data []models.UserRepo
	err = json.Unmarshal(body, &data)
	if err != nil {
		logrus.Error("Error parsing JSON:", err)
		return ""
	}
	//fmt.Println(data)
	if len(data) == 0 {
		logrus.Error("CrawlGithubRepos No Repo Data Error: ", data)
		return ""
	}
	var commitsUrl = ""
	var flag = false
	for _, repo := range data {
		//fmt.Println(repo)
		//fmt.Println(repo.Fork, repo.CommitsURL)
		if !repo.Fork {
			commitsUrl = repo.CommitsURL
			break
		}
	}
	commitsUrl, flag = strings.CutSuffix(commitsUrl, "{/sha}")
	if !flag {
		logrus.Error("CrawlGithubRepos Unable To Find Commits_URL Error: ", flag)
	}
	//fmt.Println(commitsUrl)
	return commitsUrl
}

func CrawlGithubRepoEmail(user string, perPage int) string {
	// 爬取 Github 仓库提交者的邮箱
	commitsUrl := CrawlGithubRepos(user, perPage)
	if commitsUrl == "" {
		return ""
	}
	// 2. 获取仓库信息，主要为了获取提交者的邮箱，注意要获取fork为false的仓库
	resp, err := MakeRequest("GET", commitsUrl, nil)
	if err != nil {
		logrus.Error("CrawlGithubRepos Error: ", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("CrawlGithubRepos ReadAll Error: ", err)
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
