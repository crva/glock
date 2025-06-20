package main

import (
	"flag"
	"fmt"

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
	flag.Parse()

	if *version {
		printVersion()
		return
	}

	var durations []float64
	var successCount int

	for i := 0; i < *nbRequestToPerform; i++ {
		duration, err := requester.SendHttpRequest(*url)
		if err == nil {
			durations = append(durations, duration.Seconds())
			successCount++
		}
	}

	stats.PrintHttpReport(durations, successCount)
}
