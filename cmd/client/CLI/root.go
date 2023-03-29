package CLI

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

const port = 8080
const searchQueryKey = "q"
const sortQueryKey = "s"
const pageQueryKey = "p"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "books",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {
		hostname, _ := cmd.PersistentFlags().GetString("host")
		requestURL := fmt.Sprintf("http://%s:%d", hostname, port)
		u, err := url.Parse(requestURL)
		if err != nil {
			log.Printf("client: could not parse URL: %s\n. Please try again", err)
			return
		}

		searchTerm, _ := cmd.PersistentFlags().GetString("search")
		sortBy, _ := cmd.PersistentFlags().GetString("sort")

		queryParams := u.Query()
		queryParams.Set(searchQueryKey, searchTerm)
		queryParams.Set(sortQueryKey, sortBy)
		u.RawQuery = queryParams.Encode()

		if resp, err := http.Get(u.String()); err == nil {
			if respBody, err := io.ReadAll(resp.Body); err != nil {
				log.Printf("client: could not read response body: %s\n. Please try again", err)
				return
			} else {
				fmt.Printf(string(respBody))
			}
		} else {
			if resp.StatusCode >= 500 {
				log.Printf("Internal server error: Please try again")
			}
			return
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Remove the -h help shorthand, since we want to use -h for host
	rootCmd.PersistentFlags().BoolP("help", "", false, "help for this command")

	rootCmd.PersistentFlags().String("sort", "title", "Sorts the results by the specified field 'author' or 'title'")
	rootCmd.PersistentFlags().StringP("search", "s", "", "Search the Goodreads' API and display the results on screen")
	rootCmd.PersistentFlags().StringP("host", "h", "127.0.0.1", "The hostname or ip address where the server can be found")
}
