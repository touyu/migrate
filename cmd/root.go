package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate for mysql with environment flag",
	Long: "",
	Run: func(cmd *cobra.Command, args []string) {
		dbConfig := viper.GetStringMapString(Env)
		username := dbConfig["username"]
		password := dbConfig["password"]
		host := dbConfig["host"]
		port := dbConfig["port"]
		database := dbConfig["database"]

		if DryRun {
			args := strings.Fields(fmt.Sprintf("-u %s -p %s -h %s -P %s %s --dry-run --file schema.sql", username, password, host, port, database))
			out, _ := exec.Command("mysqldef", args...).CombinedOutput()
			fmt.Print(string(out))
		} else {
			args := strings.Fields(fmt.Sprintf("-u %s -p %s -h %s -P %s %s --file schema.sql", username, password, host, port, database))
			out, _ := exec.Command("mysqldef", args...).CombinedOutput()
			fmt.Print(string(out))
		}
	},
}

var (
	Env string
	DryRun bool
)

func init() {
	viper.SetConfigName("dbconfig")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

	rootCmd.PersistentFlags().StringVarP(&Env, "env", "e", "development", "Environment")
	rootCmd.PersistentFlags().BoolVarP(&DryRun, "dry-run", "", false, "Dry run")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
