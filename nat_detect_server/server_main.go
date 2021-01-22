package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"natdetect/server"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{

		&cli.IntFlag{
			Name:  "port1,p1",
			Usage: "*Require, port 1",
		},
		&cli.IntFlag{
			Name:  "port2,p2",
			Usage: "*Require, port 2",
		},
		&cli.IntFlag{
			Name:  "port3,p3",
			Usage: "*Require, port 3",
		},
		&cli.IntFlag{
			Name:  "port4,p4",
			Usage: "*Require, port 4",
		},
		&cli.StringFlag{
			Name:  "ip1,1",
			Usage: "*Require, ip address 1",
		},
		&cli.StringFlag{
			Name:  "ip2,2",
			Usage: "*Require, ip address 2",
		},
		&cli.StringFlag{
			Name:  "ip3,3",
			Usage: "ip address 3",
		},
		&cli.StringFlag{
			Name:  "ip4,4",
			Usage: "ip address 4",
		},
		&cli.StringFlag{
			Name:  "ip5,5",
			Usage: "ip address 5",
		},
		&cli.StringFlag{
			Name:  "ip6,6",
			Usage: "ip address 6",
		},
	}

	app.Action = func(c *cli.Context) error {
		ip1 := c.String("ip1")
		if ip1 == "" {
			return errors.New("Error: ip1 is null ")
		}
		ip2 := c.String("ip2")
		if ip2 == "" {
			return errors.New("Error: ip2 is null ")
		}
		ip3 := c.String("ip3")
		ip4 := c.String("ip4")
		ip5 := c.String("ip5")
		ip6 := c.String("ip6")

		port1 := c.Int("port1")
		if port1 == 0 {
			return errors.New("Error: port1 is null ")
		}
		port2 := c.Int("port2")
		if port2 == 0 {
			return errors.New("Error: port2 is null ")
		}
		port3 := c.Int("port3")
		//if port3 == 0 {
		//	return errors.New("Error: port3 is null ")
		//}
		port4 := c.Int("port4")
		//if port4 == 0 {
		//	return errors.New("Error: port4 is null ")
		//}
		fmt.Println(ip1, ip2, ip3, ip4, ip5, ip6, port1, port2, port3, port4)
		return server.ServerRun(ip1, ip2, ip3, ip4, ip5, ip6, port1, port2, port3, port4)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

