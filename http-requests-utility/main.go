package main

import (
	"flag"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	NegativeValue            = -1
	DefaultValueForRequest   = 50
	DefaultValueForTimeout   = 20
	DefaultValueForAddresses = "https://vk.com/"
)

var initArguments Arguments

type arrayFlags []string

func (f *arrayFlags) String() string {
	return fmt.Sprintf("%d", f)
}

func (f *arrayFlags) Set(value string) error {
	*f = append(*f, value)
	return nil
}

type Arguments struct {
	Addresses  []string
	NumRequest int
	TimeOut    time.Duration
}

type InfoRequests struct {
	Mutex        sync.Mutex
	CountNoResp  int
	RespTime     []time.Duration
	TotalRuntime time.Duration
	URI          string
}

type Result struct {
	TotalRuntime time.Duration
	AverageTime  time.Duration
	MinRespTime  time.Duration
	MaxRespTime  time.Duration
	CountNoResp  int
}

type InfoResult struct {
	Address string
	Result  *Result
}

func utility(args Arguments) (results []InfoResult) {
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
			go makeRequest(*client, &info, &wg)
		}
		wg.Wait()
		results = append(results, InfoResult{
			Address: value,
			Result:  getResult(info, args),
		})
	}
	return
}

func makeRequest(client http.Client, info *InfoRequests, wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now()
	resp, err := client.Get(info.URI)
	duration := time.Since(start)
	info.Mutex.Lock()
	if err != nil {
		info.CountNoResp++
	} else {
		defer resp.Body.Close()
		info.RespTime = append(info.RespTime, duration)
		info.TotalRuntime = info.TotalRuntime + duration
	}
	info.Mutex.Unlock()
}

func getResult(info InfoRequests, arg Arguments) (res *Result) {
	if info.CountNoResp != arg.NumRequest {
		sort.Slice(info.RespTime, func(i, j int) bool { return info.RespTime[i] < info.RespTime[j] })
		res = &Result{
			TotalRuntime: info.TotalRuntime,
			AverageTime:  time.Duration(int(info.TotalRuntime) / len(info.RespTime)),
			MaxRespTime:  info.RespTime[len(info.RespTime)-1],
			MinRespTime:  info.RespTime[0],
			CountNoResp:  info.CountNoResp,
		}
	}
	return
}

func printResults(results []InfoResult) {
	for _, infoResult := range results {
		if infoResult.Result == nil {
			fmt.Println("Work with address ", infoResult.Address, " is not done")
			fmt.Println("All requests were failed")
			return
		}
		fmt.Println("Work with address ", infoResult.Address, " is successfully done")
		fmt.Println("Number of send requests was ", initArguments.NumRequest)
		fmt.Println("Timeout was set in ", initArguments.TimeOut)
		fmt.Println("All requests worked at ", infoResult.Result.TotalRuntime)
		fmt.Println("Average time to request: ", infoResult.Result.AverageTime)
		fmt.Println("Maximum return response time is ", infoResult.Result.MaxRespTime)
		fmt.Println("Minimum return response time is ", infoResult.Result.MinRespTime)
		fmt.Println("The number of answers were not waited: ", infoResult.Result.CountNoResp)
	}
}

func init() {
	//var addresses arrayFlags
	//flag.Var(&addresses, "address", "Addresses for request")
	var addresses []string
	address := flag.String("addresses", "", "Addresses for request")
	numRequest := flag.Int("num", NegativeValue, "Number of requests")
	timeOut := flag.Int("timeout", NegativeValue, "Timeout of request")
	flag.Parse()
	if *address == "" {
		addresses = []string{DefaultValueForAddresses}
		fmt.Println("App use default value of addresses: ", DefaultValueForAddresses)
	} else {
		addresses = strings.Split(*address, ";")
	}
	//if addresses == nil {
	//	addresses = arrayFlags{DefaultValueForAddresses}
	//	fmt.Println("App use default value of addresses: ", DefaultValueForAddresses)
	//}
	if *numRequest < 0 {
		*numRequest = DefaultValueForRequest
		fmt.Println("App use default value of number of requests: ", DefaultValueForRequest)
	}
	if *timeOut < 0 {
		*timeOut = DefaultValueForTimeout
		fmt.Println("App use default value of timeout: ", DefaultValueForTimeout, "s")
	}
	initArguments = Arguments{
		Addresses:  addresses,
		NumRequest: *numRequest,
		TimeOut:    time.Duration(*timeOut) * time.Second,
	}
}

func main() {
	printResults(utility(initArguments))
}
