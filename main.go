package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
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

// handleJsonFlag processes the JSON string provided via the -j flag.
// It unmarshals the JSON string into a map and then marshals it back to bytes.
func handleJsonFlag(jsonString string) ([]byte, error) {
	if jsonString == "" {
		return nil, nil
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonString), &jsonData); err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return nil, fmt.Errorf("error marshalling JSON: %w", err)
	}

	return jsonBytes, nil
}

func main() {
	version := flag.Bool("v", false, "Print version information")
	url := flag.String("u", "", "URL to send the HTTP request to")
	nbRequestToPerform := flag.Int("n", 1, "Number of requests to send")
	nbGoroutines := flag.Int("c", 1, "Number of goroutines to use for sending requests")
	method := flag.String("m", "GET", "HTTP method to use (GET, POST, PUT, DELETE, etc.)")
	jsonString := flag.String("j", "", "JSON data to send with the request")
	flag.Parse()

	if *version {
		printVersion()
		return
	}

	if *nbGoroutines < 1 || *nbGoroutines > *nbRequestToPerform {
		fmt.Println("Error: Number of goroutines must be greater than 0 and less than the number of requests to perform.")
		return
	}

	jsonBytes, err := handleJsonFlag(*jsonString)
	if err != nil {
		fmt.Printf("Error parsing JSON data: %v\n", err)
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
				// We must create a new bytes.Reader for each request because io.Reader is stateful:
				// After the first read, the reader's position is at the end, so reusing it would result in empty bodies for subsequent requests.
				// By creating a new bytes.Reader from the same jsonBytes for each request, we ensure the full JSON body is sent every time.
				var body io.Reader
				if len(jsonBytes) > 0 {
					body = bytes.NewReader(jsonBytes)
				}
				duration, err := requester.SendHttpRequest(*url, *method, body)
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
				// Same reason as above: create a new bytes.Reader for each request to avoid empty bodies.
				var body io.Reader
				if len(jsonBytes) > 0 {
					body = bytes.NewReader(jsonBytes)
				}
				duration, err := requester.SendHttpRequest(*url, *method, body)
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
