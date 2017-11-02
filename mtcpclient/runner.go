package mtcpclient

import (
	"github.com/dachad/check-max-tcp-connections/tcpclient"
	"sync"
	"fmt"
	"time"
	"io"
	"strconv"
)

func MultiTCPConnect(numberConnections int, delay int, host string, port int,
		connStatusCh chan tcpclient.Connection, debugOut io.Writer) {
	var wg sync.WaitGroup
	wg.Add(numberConnections)
	for runner := 0; runner < numberConnections; runner++ {
		fmt.Fprintln(debugOut, "Initiating runner # " + strconv.Itoa(runner))
		// TODO: closure channel should be created by the closure thread and passed into here
		go tcpclient.TCPConnect(runner, host, port, &wg, debugOut, connStatusCh,  make(chan bool))
		time.Sleep(time.Duration(delay) * time.Millisecond)
		fmt.Fprintln(debugOut, "Runner " + strconv.Itoa(runner) +
			" initated. Remaining: " + strconv.Itoa(numberConnections - runner))
	}
	fmt.Fprintln(debugOut, "Waiting runners to finish")
	wg.Wait()
}