package main

import (
	"github.com/gocraft/web"
	"fmt"
	"net/http"
	"strings"
	"log"
	"time"
	"os"
	"strconv"
	pb "github.com/garrywright/tibial/messaging"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Context struct {
	HelloCount int
	logger     *log.Logger
	conn *grpc.ClientConn
	sender pb.SenderServiceClient
}

func  initOnce() *Context {
	c := new(Context)
	c.HelloCount = 30
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	s := pb.NewSenderServiceClient(conn)
	c.conn = conn
	c.sender = s
	c.logger = log.New(os.Stdout, "[tibial] ", 0)
	t := time.Now()
	c.logger.Printf("Created Logger at 3000 %d-%02d-%02d %02d:%02d:%02d\n", t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	return c
}

func (c *Context) init(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	c.HelloCount = 30
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	s := pb.NewSenderServiceClient(conn)
	c.sender = s
	c.logger = log.New(os.Stdout, "[tibial] ", 0)
	t := time.Now()
	c.logger.Printf("Created Logger at 3000 %d-%02d-%02d %02d:%02d:%02d\n", t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	next(rw, req)
}

func (c *Context) SayHello(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprint(rw, strings.Repeat("Hello ", c.HelloCount), "World!")
}

func (c *Context) SayHelloNum(rw web.ResponseWriter, req *web.Request) {
	num := req.PathParams["num"]
	count, err := strconv.Atoi(num)
	c.HelloCount = count
	t := time.Now()
	c.logger.Printf("Counter set to %d  %d-%02d-%02d %02d:%02d:%02d\n", c.HelloCount, t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	if err != nil {
		count = 1
	}
	fmt.Fprint(rw, strings.Repeat("Hello ", count), "World!")
}
func (c *Context) SayHelloRemote(rw web.ResponseWriter, req *web.Request) {

	name := "Garry"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	r, err := c.sender.SendMessage(context.Background(), &pb.Message{Body: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	//reply := r.Reply
	c.logger.Printf("Greeting: %s", r.Reply)
	fmt.Fprint(rw, "Greeting: ", r.Reply, "\n")
}
func main() {
        c := initOnce()
	router := web.New(*c).// Create your router
	Middleware(web.LoggerMiddleware).// Use some included middleware
	Middleware(web.ShowErrorsMiddleware).// ...
	//Middleware((*Context).init).// Your own middleware!
	Get("/", (*c).SayHello).// Add a router
	Get("/:num", (*c).SayHelloNum).               // Add a router
	Get("/remote", (*c).SayHelloRemote)               // Add a router
	t := time.Now()
	c.logger.Printf("Listening on 3000 %d-%02d-%02d %02d:%02d:%02d\n", t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	defer c.conn.Close()
	http.ListenAndServe("localhost:3000", router)   // Start the server!
}

