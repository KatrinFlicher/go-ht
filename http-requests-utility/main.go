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

type InfoRequests struct {
	Mutex        sync.Mutex
	CountNoResp  int
	RespTime     []time.Duration
	TotalRuntime time.Duration
	URI          string
}

func makeRequest(client *http.Client, info *InfoRequests, wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now()
	resp, err := client.Get(info.URI)
	duration := time.Since(start)
	info.Mutex.Lock()
	if err != nil {
		info.CountNoResp++
	} else {
		defer resp.Body.Close()
	}
	info.RespTime = append(info.RespTime, duration)
	info.TotalRuntime = info.TotalRuntime + duration
	info.Mutex.Unlock()
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
		Timeout: args.TimeOut,
	}
	for _, value := range args.Addresses {
		var wg sync.WaitGroup
		wg.Add(args.NumRequest)
		info := InfoRequests{
			Mutex: sync.Mutex{},
			URI:   value,
		}
		for i := 0; i < args.NumRequest; i++ {
			go makeRequest(client, &info, &wg)
		}
		wg.Wait()
		sort.Slice(info.RespTime, func(i, j int) bool { return info.RespTime[i] < info.RespTime[j] })
		results = append(results, Result{
			Address:      value,
			TotalRuntime: info.TotalRuntime,
			AverageTime:  time.Duration(int(info.TotalRuntime) / len(info.RespTime)),
			MaxRespTime:  info.RespTime[len(info.RespTime)-1],
			MinRespTime:  info.RespTime[0],
			CountNoResp:  info.CountNoResp,
		})
	}
	arguments = args
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
	timeOut := flag.Int("timeout", -1, "Timeout of request")
	flag.Parse()
	if addresses == nil {
		panic("There are not addresses for calling")
	}
	if *numRequest < 0 {
		*numRequest = 20
		fmt.Println("App use default value of number of requests (20)")
	}
	if *timeOut < 0 {
		*timeOut = 200
		fmt.Println("App use default value of timeout (200)")
	}
	return Arguments{
		Addresses:  addresses,
		NumRequest: *numRequest,
		TimeOut:    time.Duration(*timeOut) * time.Millisecond,
	}
}
func main() {
	printResults(utility(parseArgs()))
}
