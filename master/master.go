package main

import (
	"fmt"
	"log"
	kmeans2 "main/kmeans"
	"main/util"
	"net/rpc"
	"os"
	"strconv"
	"time"
)

func main() {
	fmt.Println("Master is up")
	var err error

	var conf util.Conf
	conf.ReadConf()

	util.OpenEnv()
	numWorker, _ := strconv.Atoi(os.Getenv(util.NumWorker))
	numPoint, _ := strconv.Atoi(os.Getenv(util.NumPoint))
	numCentroid, _ := strconv.Atoi(os.Getenv(util.NumCluster))

	var registrations util.Registration

	var bools bool
	client, err := rpc.DialHTTP("tcp", conf.RegIP+":"+strconv.Itoa(conf.RegPort))

	time.Sleep(time.Second)
	err = client.Call("Registry.RetrieveMember", bools, &registrations)
	for err != nil {
		log.Fatal("Error call registry", err.Error())
	}

	for i := range registrations.Peer {
		log.Printf("Port %s", strconv.Itoa(registrations.Peer[i].Port))
		log.Printf("Address %s ", registrations.Peer[i].Address)
	}

	clients := make([]*rpc.Client, len(registrations.Peer))
	calls := make([]*rpc.Call, len(registrations.Peer))

	for i := 0; i < len(registrations.Peer); i++ {
		clients[i], err = rpc.DialHTTP("tcp", registrations.Peer[i].Address+":"+strconv.Itoa(registrations.Peer[i].Port))
		if err != nil {
			fmt.Println("Errore di connessione , retring ... ")
		}
	}

	numClient := len(clients)
	if numClient < numWorker {
		workerCrash := numWorker - numClient
		log.Fatalf("Errore! %d worker hanno subito un crash", workerCrash)
	}
	data := GeneratePoint(numPoint)
	points := CreateClusteredPoint(data)

	fmt.Println(" ** Llyod **")
	kmeans2.Llyod(numClient, numCentroid, points, util.Llyod, clients, calls)

	fmt.Println(" ** Standard KMeans **")
	kmeans2.StandardKMeans(numClient, numCentroid, points, util.Standard, clients, calls)

	fmt.Println(" ** KMeans++ **")
	kmeans2.KMeansPlusPlus(numClient, numCentroid, points, util.KmeansPlusPlus, clients, calls)

}
