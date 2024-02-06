package cmd

import (
	"GatherGithubInfo/config"
	"GatherGithubInfo/crawler"
	"GatherGithubInfo/database"
	"github.com/spf13/cobra"
)

var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Start crawling github users",
	Run: func(cmd *cobra.Command, args []string) {
		// 调用爬虫函数
		startNum, _ := cmd.Flags().GetInt("startnum")
		perPage, _ := cmd.Flags().GetInt("perpage")
		endNum, _ := cmd.Flags().GetInt("endnum")
		user := cmd.Flag("user").Value.String()
		if user != "" {
			user := crawler.CrawlGithubUserDetail(user)
			user.EmailRepo = crawler.CrawlGithubRepoEmail(user.Login, 1)
			database.SaveUser(&user)
			return
		}

		if startNum < 0 || perPage < 0 || endNum < 0 {
			panic("Invalid number")
		}
		userList := crawler.CrawlGithubUsers(startNum, perPage, endNum)
		for _, user := range userList {
			database.SaveUser(&user)
		}

	},
}

func init() {
	crawlCmd.Flags().IntP("startnum", "s", 0, "The number of users you want to start crawling")
	crawlCmd.Flags().IntP("perpage", "p", 0, "The number of users you want to crawl perpage")
	crawlCmd.Flags().IntP("endnum", "e", 0, "The number of users you want to end crawling")

	crawlCmd.Flags().StringP("user", "u", "", "The single github user you want to crawl")
	crawlCmd.PersistentFlags().StringVarP(&config.GithubToken, "token", "t", "", "The github token")
}
