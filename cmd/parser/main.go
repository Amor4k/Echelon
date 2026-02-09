package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Amor4k/Echelon/internal/parser"
)

const banner = `
                  ███████╗ ██████╗██╗  ██╗███████╗██╗      ██████╗ ███╗   ██╗
                  ██╔════╝██╔════╝██║  ██║██╔════╝██║     ██╔═══██╗████╗  ██║
                  █████╗  ██║     ███████║█████╗  ██║     ██║   ██║██╔██╗ ██║
                  ██╔══╝  ██║     ██╔══██║██╔══╝  ██║     ██║   ██║██║╚██╗██║
                  ███████╗╚██████╗██║  ██║███████╗███████╗╚██████╔╝██║ ╚████║
                  ╚══════╝ ╚═════╝╚═╝  ╚═╝╚══════╝╚══════╝ ╚═════╝ ╚═╝  ╚═══╝
`

func main() {
	//CLI Input flag handler
	//========================
	ckey := flag.String("ckey", "", "Player ckey to filter")
	cleanMobIds := flag.Bool("clean-mob-ids", true, "Remove mob IDs from output")
	noBanner := flag.Bool("no-banner", false, "Disable Echelon banner display")
	flag.Parse()
	logFiles := flag.Args()

	if *ckey == "" {
		fmt.Println("Error: Please provide a ckey using the --ckey flag.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if len(logFiles) == 0 {
		fmt.Println("Error: Please provide at least one log file.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	//If multiple logfiles, check if they are of the same round
	if len(logFiles) > 1 {
		if err := parser.ValidateRoundIDs(logFiles); err != nil {
			fmt.Println("Warning: ", err)
			fmt.Println("Continue anyway? (Y/n):")
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Aborted")
				os.Exit(1)
			}
		}
	}
	//========================

	if !*noBanner {
		fmt.Print(banner)
		fmt.Println("Welcome to Echelon - SS13 Log Analysis Tool")
	}

	opts := parser.FilterOptions{
		CleanMobIds: *cleanMobIds,
	}

	results, err := parser.FilterByCkey(logFiles, *ckey, opts)
	if err != nil {
		fmt.Println("Error filtering logs:", err)
		return
	}

	fmt.Printf("Found %d log entries for ckey %s\n", len(results), *ckey)

	//Write results to output file
	outputFile, err := os.Create("filtered_" + *ckey + ".log")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	for _, entry := range results {
		fmt.Fprintf(outputFile, "[%s] %s\n", entry.Timestamp, entry.Message)
	}

	fmt.Printf("Results written to filtered_%s.log\n", *ckey)
}
