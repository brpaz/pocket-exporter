// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

// TODO configure viper to store consumer key and access token
// ADD COLOR output to logs.
// Handle authentication errors on Retrieve method.
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	pocket "github.com/brpaz/pocket-exporter/pocket"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Pocket access token
var consumerKey string

// The output file.
var outputFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pocket-exporter",
	Short: "Pocket data exporter",
	Long:  "Command line tool that allow you to export your Pocket articles as json format.",
	Run: func(cmd *cobra.Command, args []string) {

		httpClient := &http.Client{}
		pocketClient := pocket.NewClient(consumerKey, httpClient)
		_, err := pocketClient.Authenticate()

		// TODO save accessToken for later use

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		log.Println("Retrieving articles")
		articles, err := pocketClient.Retrieve()

		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		articlesJSON, _ := json.Marshal(articles)

		ioutil.WriteFile(outputFile, articlesJSON, 0644)

		log.Println("Export with success")
	},
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
	currentDir, _ := os.Getwd()
	fileName := fmt.Sprintf("pocket_export_%s.json", time.Now().Format("2006_01_02"))

	rootCmd.Flags().StringVarP(&consumerKey, "consumerKey", "k", "some-key", "The consumer key obtained from Pocket website")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", path.Join(currentDir, fileName), "The output")
	rootCmd.MarkFlagRequired("consumerKey")

	viper.BindPFlag("consumerKey", rootCmd.PersistentFlags().Lookup("consumer_key"))
}
