package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	RA "github.com/MohaAmmar/DISYSAsgn4/proto"
	"google.golang.org/grpc"
)

type peer struct {
	RA.UnimplementedRAServiceServer
	id               int32
	name             string
	clients          map[int32]RA.RAServiceClient
	ctx              context.Context
	lamportTimestamp int32
	requesting       bool // is peer requesting access to the critical section?
	inCS             bool // is peer in the critical section?
	myReq            *RA.Request
}

func main() {
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	my_port := int32(arg1) + 3000

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	names := [3]string{"Alice", "Bob", "Charlie"}

	p := &peer{
		id:               my_port - 3000,
		name:             names[my_port-3000],
		clients:          make(map[int32]RA.RAServiceClient),
		ctx:              ctx,
		lamportTimestamp: 0,
		requesting:       false,
		inCS:             false,
	}

	// Create listener tcp on port my_port
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", my_port))
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}

	// Clear log.txt file when a new server is started
	if err := os.Truncate("log.txt", 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}

	// Connect to log file
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// Create server
	grpcServer := grpc.NewServer()
	RA.RegisterRAServiceServer(grpcServer, p)

	go func() {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to server %v", err)
		}
	}()

	for i := 0; i < 3; i++ {
		port := int32(3000) + int32(i)

		if port == my_port {
			continue
		}

		var conn *grpc.ClientConn
		fmt.Printf("Trying to dial: %v\n", port)
		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		defer conn.Close()
		c := RA.NewRAServiceClient(conn)
		p.clients[port] = c
	}

	for {
		randomInt := rand.Intn(1000)

		if randomInt < 150 {
			// Sending request to other peers
			p.RequestCritical()

			// Enter the critical section
			p.inCS = true
			fmt.Printf("%s (id: %d) enters critical section at lamport time %d\n", p.name, p.id, p.lamportTimestamp)
			log.Printf("%s (id: %d) enters critical section at lamport time %d\n", p.name, p.id, p.lamportTimestamp)

			// Sleep for 3 seconds
			time.Sleep(3 * time.Second)

			// Exit the critical section
			fmt.Printf("%s (id: %d) exiting critical section.\n", p.name, p.id)
			log.Printf("%s (id: %d) exiting critical section.\n", p.name, p.id)
			p.inCS = false
			p.requesting = false
		}
	}

}

func (p *peer) RA(ctx context.Context, r *RA.Request) (*RA.Reply, error) {
	p.lamportTimestamp++
	fmt.Printf("%s (id: %d) received a request from %s (id: %d) at lamport time %d\n", p.name, p.id, r.Name, r.Id, r.LamportTimestamp)
	log.Printf("%s (id: %d) received a request from %s (id: %d) at lamport time %d\n", p.name, p.id, r.Name, r.Id, r.LamportTimestamp)

	// This is a while loop that ensures that the loop continues to run
	// until the boolean returned from ReplyCritical is true
	for !p.ReplyCritical(r) {

	}

	fmt.Printf("%s (id: %d) replies to %s (id: %d)'s request at lamport time %d\n", p.name, p.id, r.Name, r.Id, r.LamportTimestamp)
	log.Printf("%s (id: %d) replies to %s (id: %d)'s request at lamport time %d\n", p.name, p.id, r.Name, r.Id, r.LamportTimestamp)
	reply := &RA.Reply{Msg: "free"}
	return reply, nil
}

func (p *peer) RequestCritical() {
	r := &RA.Request{Id: p.id, Name: p.name, LamportTimestamp: p.lamportTimestamp}
	fmt.Printf("%s (id: %d) is sending a request at lamport time: %d\n", p.name, p.id, p.lamportTimestamp)
	log.Printf("%s (id: %d) is sending a request at lamport time: %d\n", p.name, p.id, p.lamportTimestamp)
	p.requesting = true
	p.myReq = r

	for id, client := range p.clients {
		reply, err := client.RA(p.ctx, r)
		if err != nil {
			fmt.Printf("ERROR")
		}
		fmt.Printf("Reply from {Name Missing} (id: %d): Critical section is %s\n", id-3000, reply.Msg)
		log.Printf("Reply from {Name Missing} (id: %d): Critical section is %s\n", id-3000, reply.Msg)
	}
	p.lamportTimestamp++

}

func (p *peer) ReplyCritical(r *RA.Request) bool {
	if p.inCS {
		return false
	}

	if p.requesting {
		if p.myReq.LamportTimestamp < r.LamportTimestamp {
			return false
		}
	}

	return true

}
