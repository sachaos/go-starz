/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/jedib0t/go-pretty/table"
	"github.com/sachaos/go-starz/lib"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var version string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "go-starz [username]",
	Short:   "GitHub stars",
	Args:    cobra.ExactArgs(1),
	Version: version,
	Run:     run,
}

func run(cmd *cobra.Command, args []string) {
	token := os.Getenv("GITHUB_TOKEN")
	username := args[0]
	client := lib.NewClient(token)

	limit, err := cmd.Flags().GetInt("limit")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	thresh, err := cmd.Flags().GetInt("thresh")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	list, err := client.GetStarzList(context.Background(), username)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	total := 0
	filtered := make([]*lib.Starz, 0, len(list))
	for _, starz := range list {
		total += starz.StargazersCount

		if starz.StargazersCount >= thresh {
			filtered = append(filtered, starz)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].StargazersCount > filtered[j].StargazersCount
	})

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)

	t.AppendHeader(table.Row{"Name", "Star Count"})

	var cnt = 0
	for _, starz := range filtered {
		if limit != -1 && cnt >= limit {
			continue
		}

		t.AppendRow(table.Row{starz.Name, starz.StargazersCount})
		cnt++
	}

	t.AppendFooter(table.Row{"Total", total})

	t.Render()
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-starz.yaml)")

	rootCmd.Flags().IntP("limit", "l", -1, "Only show this many repos. -1 means infinite")
	rootCmd.Flags().IntP("thresh", "t", 1, "Only show repos above this threshold.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".go-starz" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".go-starz")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
