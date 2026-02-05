package main

import (
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



                                                                                                    
                                                  @@@@@@@@@@@@@@@@@@@@@                             
                                      @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@                   
                              @@@@@@@@@@@@@@@@@@@@@      @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@             
                        @@@@@@@@@@@@@         @@@@@       @@   @@@@@@@@@@   @@@@@@@@@@@@@@          
                    @@@@@@@@@               @@@@@@             @@@@@@@@@@@        @@@@@@@@@@@       
               @@@@@@@@                    @@@@@@@           @@@@@@@@@@@@@@           @@@@@@@@@     
            @@@@@@                        @@@@@@@@          @@@@@@@@@@@@@@@@             @@@@@@@@   
         @@@@@                            @@@@@@@@@@    @@@ @@@@@@@@@@@@@@@@                @@@@@@  
      @@@@                                @@@@@@@@@@@  @@@@@@@@@@@@@@@@@@@@@@                @@@@@@ 
    @@@                                  @@@@@@@@@@@@@    @@@@@@@@@@@@@@@@@@@                 @@@@@ 
  @@                                     @@@@@@@@@@@@@@@@@ @@@@@@@@@@@@@@@@@@                  @@@@ 
 @                                        @@@@@@@@@@@@@@@@@@@@       @@@@@@@@                  @@@@ 
                                          @@@@@@@@@@@@@@@@@@@            @@@                   @@@  
                                           @@@@@@@@@@@@@@@@@@           @@@@                  @@@   
                                            @@@@@@@@@@@@@@@@@@         @@@@                   @@    
                                             @@@@@@@@@@@@@@@@@@@      @@@@                  @@      
                                              @@@@@@@@@@@@@@@@@@    @@@@                            
                                                @@@@@@@@@@@@@@@   @@@@                              
                                                   @@@@@@@@@@@@ @@@@                                
                                                        @@@@@@                                      
                                                                                                    

Hack the Gibson. Hack the Planet!
`

func main() {
	fmt.Printf(banner)
	fmt.Println("Welcome to Echelon - SS13 Log Analysis Tool")

	ckey := "beypazarifan"

	results, err := parser.FilterByCkey()
	if err != nil {
		fmt.Println("Error filtering logs:", err)
		return
	}

	fmt.Printf("Found %d log entries for ckey %s\n", len(results), ckey)

	//Write results to output file
	outputFile, err := os.Create("filtered_" + ckey + ".log")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	for _, entry := range results {
		fmt.Fprintf(outputFile, "[%s] %s\n", entry.Timestamp, entry.Message)
	}

	fmt.Printf("Results written to %s\n", "filtered_%s.log\n", ckey)
}
