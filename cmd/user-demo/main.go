package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	Run()
}

var (
	host       = "127.0.0.1"
	apiPort    = "8080"
	streamPort = "9111"
)

func Run() {
	// init users
	var (
		wg        = sync.WaitGroup{}
		userOne   = newUser()
		userTwo   = newUser()
		userThree = newUser()
		userfour  = newUser()
		userfive  = newUser()
	)

	println("going to interact with stream service...")

	wg.Add(3)

	go userOne.StartWatching(userOne.Register(), "zato01", &wg)      // case when stream of channel 'zato01'
	go userTwo.StartWatching(userTwo.Register(), "zato02", &wg)      // case when stream of channel 'zato02'
	go userThree.StartWatching(userThree.Register(), "unknown", &wg) // case when no stream should be received
	go userfour.StartWatching(userfour.Register(), "zato01", &wg)    // case when stream of channel 'zato02'
	go userfive.StartWatching(userfive.Register(), "zato02", &wg)    // case when stream of channel 'zato02'

	wg.Wait()
}

var lastHTTPUserID = 1

func newUser() *user {
	uid := "user-" + strconv.Itoa(lastHTTPUserID)
	lastHTTPUserID++

	f, err := os.Create(fmt.Sprintf("/home/%s/test/"+uid+".ts", os.Getenv("USER")))
	if err != nil {
		panic(err)
	}

	return &user{
		id:     uid,
		output: f,
	}
}

type user struct {
	id     string
	output io.Writer
}

func (u *user) Register() <-chan struct{} {
	done := make(chan struct{})

	go func() {
		endpoint := "http://" + host + ":" + apiPort + "/register"
		client := http.Client{
			Timeout: 2 * time.Second,
		}

		for {
			resp, err := client.Post(endpoint, "text/plain", strings.NewReader(u.id))
			if err == nil && resp.StatusCode == http.StatusOK {
				break
			}

			println(u.id + ": registration failed. Will retry soon...")
			time.Sleep(500 * time.Millisecond)
		}

		println(u.id + ": registration succeed. Will stream soon...")
		done <- struct{}{}
	}()

	return done
}

func (u *user) StartWatching(start <-chan struct{}, channelID string, done *sync.WaitGroup) {
	defer done.Done()
	<-start

	conn, err := net.Dial("tcp", host+":"+streamPort)
	if err != nil {
		println("failed to establish tcp connection: " + err.Error())
		return
	}

	_, err = conn.Write([]byte(u.id + " " + channelID + "\n"))
	if err != nil {
		println("failed to write data: " + err.Error())
		return
	}

	reader := bufio.NewReader(conn)
	for {
		_, err := reader.WriteTo(u.output)
		if err != nil {
			println("error during stream: " + err.Error())
			break
		}
	}
}
