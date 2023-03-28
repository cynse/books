package CLI

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

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
		port := 8080
		requestURL := fmt.Sprintf("http://%s:%d", hostname, port)

		if resp, err := http.Get(requestURL); err == nil {
			if respBody, err := io.ReadAll(resp.Body); err != nil {
				fmt.Printf("client: could not read response body: %s\n. Please try again", err)
				os.Exit(1)
			} else {
				fmt.Printf(string(respBody))
			}
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
	rootCmd.PersistentFlags().StringP("host", "h", "127.0.0.1", "the hostname or ip address where the server can be found")
}
