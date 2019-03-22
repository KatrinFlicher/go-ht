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
	Addresses  []string
	NumRequest int
	TimeOut    time.Duration
}

type Result struct {
	Address      string
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
	value = "https://" + value
	*f = append(*f, value)
	return nil
}

func utility(args Arguments) (results []Result, arguments Arguments) {
	client := &http.Client{
		Timeout: args.TimeOut * time.Second,
	}
	for _, value := range args.Addresses {
		var wg sync.WaitGroup
		wg.Add(args.NumRequest)
		var mutex = &sync.Mutex{}
		var countNoResp int
		var respTime []time.Duration
		var totalRuntime time.Duration
		for i := 0; i < args.NumRequest; i++ {
			go func(group *sync.WaitGroup) {
				defer wg.Done()
				start := time.Now()
				resp, err := client.Get(value)
				defer resp.Body.Close()
				duration := time.Since(start)
				mutex.Lock()
				if err != nil {
					countNoResp++
				}
				respTime = append(respTime, duration)
				totalRuntime = totalRuntime + duration
				mutex.Unlock()
			}(&wg)
		}
		wg.Wait()
		sort.Slice(respTime, func(i, j int) bool { return respTime[i] < respTime[j] })
		results = append(results, Result{
			Address:      value,
			TotalRuntime: totalRuntime,
			AverageTime:  time.Duration(int(totalRuntime) / len(respTime)),
			MaxRespTime:  respTime[len(respTime)-1],
			MinRespTime:  respTime[0],
			CountNoResp:  countNoResp,
		})
		fmt.Println(results)
	}
	return
}

func printResults(results []Result, arguments Arguments) {
	for _, result := range results {
		fmt.Println("Work with address ", result.Address, " is successfully done")
		fmt.Println("Number of send requests was ", arguments.NumRequest)
		fmt.Println("Timeout was set in ", arguments.TimeOut)
		fmt.Println("All requests worked at ", result.TotalRuntime)
		fmt.Println("Average time to request: ", result.AverageTime)
		fmt.Println("Maximum return response time is ", result.MaxRespTime)
		fmt.Println("Minimum return response time is ", result.MinRespTime)
		fmt.Println("The number of answers were not waited: ", result.CountNoResp)
	}
}

func parseArgs() Arguments {
	var addresses arrayFlags
	flag.Var(&addresses, "address", "Addresses for request")
	numRequest := flag.Int("num", -1, "Number of requests")
	timeOut := time.Duration(*flag.Int("timeOut", -1, "Timeout of request"))
	flag.Parse()
	if addresses == nil {
		panic("There are not addresses for calling")
	}
	if *numRequest < 0 {
		*numRequest = 20
		fmt.Println("App use default value of number of requests (20)")
	}
	if timeOut < 0 {
		timeOut = time.Duration(200)
		fmt.Println("App use default value of timeout (200)")
	}
	return Arguments{
		Addresses:  addresses,
		NumRequest: *numRequest,
		TimeOut:    timeOut,
	}
}
func main() {
	printResults(utility(parseArgs()))
}
