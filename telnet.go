package main

import (
	"bufio"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//TelnetClient provides main telnet operations
type TelnetClient interface {
	Connect() error
	Send() error
	Receive() error
	Close() error
	Done() <-chan struct{}
}

type basicTelnetClient struct {
	address    string
	connection net.Conn
	timeout    time.Duration
	in         io.ReadCloser
	out        io.Writer
	done       chan struct{}
}

//NewTelnetClient returns a telnet client instance
func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan struct{})
	go func() {
		<-sigs
		close(done)
	}()
	return &basicTelnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
		done:    done,
	}
}

func (btc *basicTelnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", btc.address, btc.timeout)
	if err != nil {
		return err
	}
	btc.connection = conn
	return nil
}

func (btc *basicTelnetClient) Send() error {
	buffer, err := bufio.NewReader(btc.in).ReadBytes(byte('\n'))
	switch {
	case endOftransmissionCheck(err):
		tryToCloseChannel(btc.done)
		return nil
	case err != nil:
		return err
	default:
	}
	_, err = btc.connection.Write(buffer)
	return err
}

func (btc *basicTelnetClient) Receive() error {
	buffer, err := bufio.NewReader(btc.connection).ReadBytes(byte('\n'))
	switch {
	case endOftransmissionCheck(err):
		tryToCloseChannel(btc.done)
		return nil
	case err != nil:
		return err
	default:
	}
	_, err = btc.out.Write(buffer)
	return err
}

func (btc *basicTelnetClient) Close() error {
	return btc.connection.Close()
}

func (btc *basicTelnetClient) Done() <-chan struct{} {
	return btc.done
}

func endOftransmissionCheck(err error) bool {
	if err == io.EOF {
		return true
	}
	return false
}

func tryToCloseChannel(c chan struct{}) {
	select {
	case <-c:
	default:
		close(c)
	}
}
