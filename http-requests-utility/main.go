package main

import (
	"flag"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

type Arguments struct {
	Adr        string
	Addresses  []string
	NumRequest int
	TimeOut    time.Duration
}

type Result struct {
	TotalRuntime time.Duration
	AverageTime  time.Duration
	MinRespTime  time.Duration
	MaxRespTime  time.Duration
	CountNoResp  int
}

type arrayFlags []string

func (f *arrayFlags) String() string {
	return fmt.Sprintf("%d", f)
}

func (f *arrayFlags) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func utility(args Arguments) (results []Result) {
	client := &http.Client{
		Timeout: args.TimeOut * time.Second,
	}
	//for _, value := range args.Addresses {
	var wg sync.WaitGroup
	wg.Add(args.NumRequest)
	var countNoResp int
	respTime := make([]time.Duration, args.NumRequest)
	var totalRuntime time.Duration
	fmt.Println("Start", args)
	for i := 0; i <= args.NumRequest; i++ {
		go func(group *sync.WaitGroup) {
			defer wg.Done()
			start := time.Now()
			resp, err := client.Get(args.Adr)
			duration := time.Since(start)
			if err != nil {
				countNoResp++
			}
			defer resp.Body.Close()
			fmt.Println(duration)
			respTime = append(respTime, duration)
			totalRuntime = totalRuntime + duration
			fmt.Println(totalRuntime)
		}(&wg)
	}
	sort.Slice(respTime, func(i, j int) bool { return respTime[i] < respTime[j] })
	fmt.Println("Finish")
	results = append(results, Result{
		TotalRuntime: totalRuntime,
		AverageTime:  time.Duration(int(totalRuntime) / len(respTime)),
		MaxRespTime:  respTime[len(respTime)-1],
		MinRespTime:  respTime[0],
		CountNoResp:  countNoResp,
	})
	fmt.Println(results)
	//}
	return
}

func printResults(results []Result) {
	for _, result := range results {
		fmt.Println("All requests worked at ", result.TotalRuntime)
		fmt.Println("Average time to request: ", result.AverageTime)
		fmt.Println("Maximum return response time is ", result.MaxRespTime,
			"Minimum return response time is ", result.MinRespTime)
		fmt.Println("The number of answers were not waited: ", result.CountNoResp)
	}
}

func parseArgs() Arguments {
	//var addresses arrayFlags
	//flag.Var(&addresses, "address", "Addresses for request")
	adr := *flag.String("address", "https://www.google.com/", "Addresses for request")
	numRequest := *flag.Int("numRequest", 20, "Number of requests")
	//timeOut := time.Duration(*flag.Int("timeOut", 200, "Timeout of request"))
	timeOut := time.Duration(20)
	flag.Parse()
	if adr == "" {
		panic("There are not addresses for calling")
	}

	if numRequest == 5 {
		fmt.Println("App use default value of number of requests")
	}
	if timeOut == 0 {
		fmt.Println("App use default value of timeout")
	}
	return Arguments{
		Adr:        adr,
		NumRequest: numRequest,
		TimeOut:    timeOut,
	}
}
func main() {
	printResults(utility(parseArgs()))
}
