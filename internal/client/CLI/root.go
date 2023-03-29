package CLI

import (
	"books/internal/client"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "books",
	Short: "Retrieves a list of books",
	Long: `Books is an application that will reach to a server, 
           and displays a list of books, sorted depending on the user's preference'`,

	Run: func(cmd *cobra.Command, args []string) {
		hostname, _ := cmd.PersistentFlags().GetString("host")
		searchTerm, _ := cmd.PersistentFlags().GetString("search")
		sortBy, _ := cmd.PersistentFlags().GetString("sort")
		page, _ := cmd.PersistentFlags().GetInt("page")

		if sortBy != "title" && sortBy != "author" {
			fmt.Print("sortBy must be by either title or author!")
			return
		}
		if page <= 0 {
			fmt.Print("Page number must be 1 or greater!")
		}
		client.Execute(hostname, searchTerm, sortBy, page)
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
	rootCmd.PersistentFlags().Int("page", 1, "Choose which page of the result to display")
	rootCmd.PersistentFlags().String("sort", "title", "Sorts the results by the specified field 'author' or 'title'")
	rootCmd.PersistentFlags().StringP("search", "s", "", "Search the Goodreads' API and display the results on screen")
	rootCmd.PersistentFlags().StringP("host", "h", "127.0.0.1", "The hostname or ip address where the server can be found")
}
