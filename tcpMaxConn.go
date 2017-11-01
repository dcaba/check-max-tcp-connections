package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/dachad/check-max-tcp-connections/mtcpclient"
	"github.com/spf13/pflag"
)

func main() {
	hostPtr := pflag.StringP("host", "h", "", "Host you want to open tcp connections against")
	// according to https://en.wikipedia.org/wiki/List_of_TCP_and_UDP_port_numbers, you are probably not using this
	portPtr := pflag.StringP("port", "p", "", "Port you want to open tcp connections against")
	numberConnectionsPtr := pflag.IntP("connections", "c", 100, "Number of connections you want to open")
	delayPtr := pflag.IntP("sleep", "s", 10, "Time you want to sleep between connections, in ms")
	debugPtr := pflag.BoolP("debug", "d", false, "Print debugging information to the standard error")
	reportingIntervalPtr := pflag.IntP("interval", "i", 1, "Interval, in seconds, between stats updates")
	assumeyesPtr := pflag.BoolP("assume-yes", "y", false, "Force execution without asking for confirmation")
	pflag.Parse()

	// Target host and port are mandatory to run the TCP check
	if *hostPtr == "" || *portPtr == "" {
		pflag.PrintDefaults()
		os.Exit(1)
	}
	port, _ := strconv.Atoi(*portPtr)

	var debugOut io.Writer = ioutil.Discard
	if *debugPtr {
		debugOut = os.Stderr
	}

	if *assumeyesPtr == false {
		*assumeyesPtr = askForUserConfirmation(*hostPtr, *portPtr, *numberConnectionsPtr)
	}

	if *assumeyesPtr {
		connStatusCh := mtcpclient.StartReportingLogic(*numberConnectionsPtr, *reportingIntervalPtr)
		mtcpclient.MultiTCPConnect(*numberConnectionsPtr, *delayPtr, *hostPtr, port, connStatusCh, debugOut)
		fmt.Println("\n*** Execution Terminated ***")
	} else {
		fmt.Println("\n*** Execution aborted as prompted by the user ***")
	}

}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func askForUserConfirmation(host string, port string, connections int) bool {

	fmt.Println("****************************** WARNING ******************************")
	fmt.Println("* You are going to  run a TCP stress check with these arguments:")
	fmt.Println("*	- Target: " + host)
	fmt.Println("*	- TCP Port: " + port)
	fmt.Println("*	- # of concurrent connections: " + strconv.Itoa(connections))
	fmt.Println("*********************************************************************")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Do you want to continue? (y/N)")
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Response not processed")
			os.Exit(1)
		}

		response = strings.TrimSuffix(response, "\n")
		response = strings.ToLower(response)
		switch {
		case stringInSlice(response, []string{"yes", "y"}):
			return true
		case stringInSlice(response, []string{"no", "n", ""}):
			return false
		default:
			fmt.Println("\nSorry, response not recongized. Try again, please")
		}
	}
}
