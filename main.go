package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/brutella/dnssd"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"
)

var instanceFlag = flag.String("name", "Awesome Service", "Service name")
var domainFlag = flag.String("domain", "local", "domain")
var portFlag = flag.Int("port", 12345, "Port")
var ipFlag = flag.String("ip", "", "IP address")
var hostnameFlag = flag.String("hostname", "hap-service", "hostname of service")
var macFlag = flag.String("mac", "", "mac of service")

func main() {
	flag.Parse()
	if len(*instanceFlag) == 0  || len(*domainFlag) == 0  || len(*ipFlag) == 0 || len(*macFlag) == 0 {
		flag.Usage()
		return
	}
	service := "_hap._tcp"

	instance := fmt.Sprintf("%s.%s.%s.", strings.Trim(*instanceFlag, "."), strings.Trim(service, "."), strings.Trim(*domainFlag, "."))

	log.Printf("Registering Service %s port %d\n", instance, *portFlag)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if resp, err := dnssd.NewResponder(); err != nil {
		log.Fatal(err)
	} else {
		cfg := dnssd.Config{
			Name:   *instanceFlag,
			Type:   service,
			Domain: *domainFlag,
			Port:   *portFlag,
			IPs:    []net.IP{net.ParseIP(*ipFlag)},
			Host: 	*hostnameFlag,
			Text: map[string]string{
				"md": *instanceFlag,
				"pv": "1.0",
				"id": *macFlag,
				"c#": "2" ,
				"s#": "1",
				"ff": "0",
				"ci": "2",
				"sf": "1",
			},
		}
		log.Printf("Registering service to ip %s with mac %s", *ipFlag, *macFlag)
		srv, err := dnssd.NewService(cfg)
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			stop := make(chan os.Signal, 1)
			signal.Notify(stop, os.Interrupt)

			select {
			case <-stop:
				cancel()
			}
		}()

		go func() {
			time.Sleep(1 * time.Second)
			handle, err := resp.Add(srv)
			if err != nil {
				log.Fatal(err)
			} else {
				log.Printf("Got a reply for service %s: Name now registered and active\n", handle.Service().ServiceInstanceName())
			}
		}()
		err = resp.Respond(ctx)

		if err != nil {
			log.Fatal(err)
		}
		log.Println("Stopping")
	}
}