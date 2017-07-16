package main

import (
	"fmt"
	"os"
	"net/http"
	"sync"
	
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "xCurl"
	app.Usage = "Curl made async"
	app.Action = func(c *cli.Context) error {
		nargs := c.NArg()
		switch {
		case nargs == 0:
			break
		case nargs == 1:
			url := c.Args().Get(0)
			if res, err := http.Get(url); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("I got %d from %q\n", res.StatusCode, url)
			}
		case nargs > 1:
			urls := c.Args()
			var wg sync.WaitGroup
			for _, url := range urls {
				wg.Add(1)
				go func(url string) {
					res, err := http.Get(url)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Printf("I got %d from %q\n", res.StatusCode, url)
					wg.Done()
					
				}(url)
			}
			wg.Wait()
		}
		return nil
	}
	app.Run(os.Args)
}
