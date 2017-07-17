package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"net/http"
	"errors"
	
	"github.com/urfave/cli"
)

func main() {	
	app := cli.NewApp()
	app.Name = "xCurl"
	app.Usage = "Curl made async"
	app.Action = func(c *cli.Context) error {
		emitters := [](func() interface{}){}
		
		nargs := c.NArg()
		if nargs == 0 {
			return errors.New("Argument not provided")
		}
		
		urls := c.Args()
		for _, url := range urls {
			em := func() interface{} {
				res, err := http.Get(url)
				if err != nil {
					return err
				}

				defer res.Body.Close()
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					return err
				}
				return string(body)
			}
			
			emitters = append(emitters, em)
		}
		rxq := NewQueue(emitters...)
		err, done := rxq.Subscribe(func(i interface{}) {
			fmt.Printf("%q\n\n", i)
		})

		select {
		case e := <-err:
			fmt.Println(e)
		case d := <-done:
			fmt.Println(d)
		}
		
		return nil
	}
	app.Run(os.Args)
}
