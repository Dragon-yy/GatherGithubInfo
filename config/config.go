package config

// 获取用户列表
var GithubUsersApi = "https://api.github.com/users?since=%d&per_page=%d"

// 获取用户详细信息
var GithubUserDetailApi = "https://api.github.com/users/%s"

// 获取用户Repo列表（commits_url字段)
// var GithubUserReposApi = "https://api.github.com/users/%s/repos?per_page=%d"
var GithubUserReposApi = "https://api.github.com/users/%s/repos"

// 新用户就算是使用token速率也被限制在60
var GithubToken = ""

var DatabaseIP = ""
var DatabasePort = ""
var DatabaseUser = ""
var DatabasePassword = ""
var DatabaseName = ""
