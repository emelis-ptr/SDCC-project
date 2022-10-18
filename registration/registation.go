package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"main/util"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type Registry struct {
	Peer util.Registration
}

var members = 0
var registration util.Registration

// RegisterMember adds a new member and returns in res the memberhip when all members are registered
func (registry *Registry) RegisterMember(arg util.Peer, res *util.Registration) error {

	registration.Index = members
	members++ // increment the number of registered members

	//*res = waitForAll(registry, arg) // put the member in wait for others and give him the membership
	registration.Peer = append(registration.Peer, arg)
	*res = registration
	registry.Peer = *res
	return nil
}

func (registry *Registry) RetrieveMember(bool bool, res *util.Registration) error {
	registration = registry.Peer

	*res = registration
	return nil
}

func main() {
	var err error

	fmt.Println("Registration service is up")

	util.OpenEnv()

	registry := new(Registry)

	//read configuration
	var conf util.Conf
	conf.ReadConf()

	//server := rpc.NewServer()
	err = rpc.Register(registry)
	if err != nil {
		log.Fatalf("Error in register server name", err)
	}

	//init variables and signal handler for shutdown
	sigs := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//cancel context on shutdown
	go func() {
		<-sigs
		cancel()
	}()

	fmt.Println("Registration is starting...")

	//set up listening for incoming connection
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(conf.RegPort))
	if err != nil {
		log.Fatalln("Listening failed with error: ", err)
	}

	log.Printf("Serving rpc sulla porta %s", strconv.Itoa(conf.RegPort))
	log.Printf("Address reg %s", conf.RegIP)

	//expose api on open port
	err = rpc.Register(new(Registry))
	if err != nil {
		log.Fatal("Error registry", err)
	}
	rpc.HandleHTTP()
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return http.Serve(lis, nil)
	})

	//close listener on shutdown
	g.Go(func() error {
		<-gCtx.Done()
		return lis.Close()
	})
	if err := g.Wait(); err != nil {
		fmt.Println("\nRegistration service shutdown")
	}
}
