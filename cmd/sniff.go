/*
Copyright © 2022 ndzn me@x3.lol
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

// TODO:

// sniffCmd represents the sniff command
var sniffCmd = &cobra.Command{
	Use:   "sniff",
	Short: "Sniffs out information about a website.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:
Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// add flag to more aff domains from crt.sh associated with the certificates
		// filter out cloudflare cert domains if it becomes an issue (probably not)
		if len(args) == 0 {
			fmt.Println("Please provide a domain to sniff.")
			return
		}
		domain := args[0]
		fmt.Println("Sniffing out information about", domain)

		resp, err := http.Get("http://crt.sh/?q=" + args[0])
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		re := regexp.MustCompile(`(?m)<TD>([a-zA-Z0-9\.\-]+)\.` + args[0] + `</TD>`)
		matches := re.FindAllStringSubmatch(string(body), -1)

		seen := make(map[string]bool)
		for _, match := range matches {
			if !seen[match[1]] {
				seen[match[1]] = true
				fmt.Println(match[1] + "." + args[0])
			}
		}
		if len(matches) == 0 {
			fmt.Println("No subdomains found.")
		}

		outputFiles, _ := cmd.Flags().GetString("output")
		if outputFiles != "" {
			outputFile(seen, outputFiles)
		}

	},
}

func init() {
	rootCmd.AddCommand(sniffCmd)

	sniffCmd.Flags().StringP("output", "o", "", "Name of the output file")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sniffCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sniffCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


// fix this
func outputFile(i map[string]bool, name string) {
	file, fileErr := os.Create(name + ".txt")
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}
	fmt.Fprintf(file, "%v\n", i)
}
