package cmd

import (
	"GatherGithubInfo/config"
	"GatherGithubInfo/database"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the database",
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化数据库 root:123456@tcp(127.0.0.1:3306)/github?charset=utf8mb4&parseTime=True&loc=Local
		ip, _ := cmd.Flags().GetString("ip")
		port, _ := cmd.Flags().GetString("port")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		dbname, _ := cmd.Flags().GetString("dbname")
		err := database.InitDB(ip, port, username, password, dbname)
		if err != nil {
			panic("Failed to initialize the database: " + err.Error())
		} else {
			logrus.Info("Database initialized successfully.")
		}

	},
}

func init() {
	initCmd.PersistentFlags().StringVarP(&config.DatabaseIP, "ip", "I", "127.0.0.1", "The Database IP")
	initCmd.PersistentFlags().StringVarP(&config.DatabasePort, "port", "P", "3306", "The Database Port")
	initCmd.PersistentFlags().StringVarP(&config.DatabaseUser, "username", "u", "root", "The Database Username")
	initCmd.PersistentFlags().StringVarP(&config.DatabasePassword, "password", "p", "123456", "The Database Password")
	initCmd.PersistentFlags().StringVarP(&config.DatabaseName, "dbname", "d", "github", "The Database Name")

}
