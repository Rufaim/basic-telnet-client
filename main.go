package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	timeout := flag.String("timeout", "10s", "timeout for a connection")
	flag.Parse()
	if len(flag.Args()) == 0 {
		usage()
		return
	}
	hostPort := flag.Args()[0]

	timeoutDuration, err := time.ParseDuration(*timeout)
	if err != nil {
		usage()
		panic(err.Error())
	}

	telnetClient := NewTelnetClient(hostPort, timeoutDuration, os.Stdin, os.Stdout)
	err = telnetClientConnectAndAttachStandartStreams(telnetClient)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	<-telnetClient.Done()
	telnetClient.Close()
}

func usage() {
	fmt.Printf("Usage: %s [--timeout 10s] hostname:port\n", os.Args[0])
}

func telnetClientConnectAndAttachStandartStreams(tc TelnetClient) error {
	if err := tc.Connect(); err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-tc.Done():
				return
			default:
			}
			tc.Receive()
		}
	}()

	go func() {
		for {
			select {
			case <-tc.Done():
				return
			default:
			}
			tc.Send()
		}
	}()

	return nil
}
