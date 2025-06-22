package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/crva/glock/requester"
	"github.com/crva/glock/stats"
)

func printAscii() {
	glock := `                                                                                                    
              +xxxXXxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxXXXxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx               
              Xx++xxxxxxXxxxx+x+++++++++++++++++++++++++++++++++++++++++++++x+++xxx+xxxxxxxxxxxxxx              
              XxxxXxxxxXXxxxxxxxxxxxXXXXXxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx              
              XxxxXxxxxxXxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxXxXxXxXxXxXxX              
              XxxXXXxxxXXXXXXxXxxXXXXxxxxxxxxXxxxxxXXXXxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxXxXxXxXxXxXxX              
             +xxxxxxxxxxxx+++x+++x$XxXx++++++x+++++xxxxx++;;;+;+xXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX              
            ;+xxX$Xx+xxxxxxxxx+xxxXXXxx+xxxxxx++xxXxxxxx++xxx+++XXxxxxxxxxxxxxxxxxxxxXxXxXxXxXxXxX              
                :;xXx;xXXxxxxxxxxxxxxxxxxxxxxXxxxxX++xXXXXXXXXXXXXXxXxXxxXXxxxXXXXXXxxxxxxxxxxxxxx              
                  ;xX+++xxXXXXXXXxxxxxxXXXXX$$$Xxx+;;:;xXXXXx                                       
                   +XXx+xx++++++++++xxXX$$$$$XXx        ;xx                                         
                   +XXx+xx+;;;;+++++xXX$X+xxxXx+        :xx                                         
                  ;xXXxxxxxXXXXXXXXxxX$$x   xXXxx       :xx                                         
                  xXXxxXXXXxXXxxXXxX$$$Xxx    +xxX     :xxx                                         
                 +xXXxxxxxXxxXxxxxXX$XXx+x++++++++++++xxxxxx                                        
                +xXXxxxxxxxxxxxxXxXXXx     ;xx++xxx+++                                              
               +xXXxxxxxxxxxxxxxxXXXx                                                               
              +xXXXxxxxxxxxxxxxxXXXx                                                                
             ;xXXXxxxxxxxxxxxxxxXXXX;                                                               
            +xXXXxxxxxxxxxxxxxxXXX+                                                                 
           +XXXxxxxxxxxxxxxxxXXXXx                                                                  
          +xXXxxxxxxXXxxxxXXxXXXX+                                                                  
         +xXXxxxXXXXXXxxxxXXxXXXx+                                                                  
        ;+XXxxxxxXXXXXxxxxxxXXX+                                                                    
        +XXxxxxxxxxxxxxXxxxXXXx                                                                     
       +XXx++xxxxxxxxxxxxxxxXX;                                                                     
       +XXx+xxxxxxxxxxxxxxxxXX+                                                                     
         +++xxxxxxxxxxxxxxxXXx`
	fmt.Println(glock)
}

func printVersion() {
	printAscii()
	fmt.Println("Glock - HTTP Request Performance Tester")
	fmt.Println("Version: 0.0.1")
}

func main() {
	version := flag.Bool("v", false, "Print version information")
	url := flag.String("u", "", "URL to send the HTTP request to")
	nbRequestToPerform := flag.Int("n", 1, "Number of requests to send")
	nbGoroutines := flag.Int("c", 1, "Number of goroutines to use for sending requests")
	flag.Parse()

	if *version {
		printVersion()
		return
	}

	if *nbGoroutines < 1 || *nbGoroutines > *nbRequestToPerform {
		fmt.Println("Error: Number of goroutines must be greater than 0 and less than the number of requests to perform.")
		return
	}

	var durations []float64
	var successCount int
	var requestWg sync.WaitGroup
	nbRequestToPerformPerGoroutine := *nbRequestToPerform / *nbGoroutines
	nbRemainingRequests := *nbRequestToPerform % *nbGoroutines

	// Handle all the evenly distributed requests
	for i := 0; i < *nbGoroutines; i++ {
		requestWg.Add(1)

		go func() {
			defer requestWg.Done()

			for j := 0; j < nbRequestToPerformPerGoroutine; j++ {
				duration, err := requester.SendHttpRequest(*url)
				if err == nil {
					successCount++
				}
				durations = append(durations, duration.Seconds())
			}
		}()
	}

	// Handle the remaining requests in a separate goroutine
	if nbRemainingRequests > 0 {
		requestWg.Add(1)
		go func() {
			defer requestWg.Done()

			for j := 0; j < nbRemainingRequests; j++ {
				duration, err := requester.SendHttpRequest(*url)
				if err == nil {
					successCount++
				}
				durations = append(durations, duration.Seconds())
			}
		}()
	}

	start := time.Now()
	requestWg.Wait()
	totalDuration := time.Since(start).Seconds()

	stats.PrintHttpReport(durations, successCount)
	stats.PrintTotalDuration(totalDuration)
}
