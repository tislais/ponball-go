package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type machine struct {
	title string
	ipdbId int
	manufacturer string
}

func main() {
	
	fetchedHTML, err := fetchHTML("https://www.ipdb.org/machine.cgi?id=20")
	if err != nil {
		fmt.Println("ouchies")
	}
	defer fetchedHTML.Close()
	
	
	// m, err := processTokens(fetchedHTML)
	// if err != nil {
	// 	fmt.Printf("%v\n", err)
	// }

	m2, err := parseHTML(fetchedHTML)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Title: %v\n", m2.title)
}

func fetchHTML(URL string) (io.ReadCloser, error) {
	res, err := http.Get(URL)
	if err != nil {
		fmt.Printf("hi %s", errors.New(err.Error()))
	}

	if res.StatusCode != http.StatusOK {
		fmt.Printf("Response status code was %d\n", res.StatusCode)
	}

	ct := res.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "text/html") {
		err = errors.New("big bad: response content was not text/html, instead got " + ct)
	}
	return res.Body, err
}

func parseHTML(body io.ReadCloser) (machine, error) {
	s, err := ioutil.ReadAll(body)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%s", s)

	text := string(s)

	dom, err := html.Parse(strings.NewReader(text))
	fmt.Println("dom ", dom.LastChild)

	return machine{}, err
}

func processTokens(body io.ReadCloser) (machine, error) {
	// fmt.Println(dom, err)
	// if err != nil {
	// 	fmt.Printf("Error from html.Parse:%v", err)
	// }
	// var f func(*html.Node)
	// f = func(n *html.Node) {
	// 	if n.Type == html.ElementNode && n.Data == "a" {
	// 		for _, a := range n.Attr {
	// 			if a.Key == "href" {
	// 				fmt.Println(a.Val)
	// 				break
	// 			}
	// 		}
	// 		for c := n.FirstChild; c != nil; c = c.NextSibling {
	// 			f(c)
	// 		}
	// 	}
	// 	f(dom)
	// }
	// z = tokenizer
	var err error
	z := html.NewTokenizer(body)
	m := machine{}
	
	for {
		// tt = tokenType
		tt := z.Next()
		if tt == html.ErrorToken {
			err := z.Err()
			if err == io.EOF {
				err = nil
				break
			}
			fmt.Println(z.Err())
		}

		if tt == html.StartTagToken {
			token := z.Token()
			if token.Data == "title" {
				tt = z.Next()
				if tt == html.TextToken {
					m.title = z.Token().Data
					fmt.Println(z.Token().Data)
				}
			} 
		}
	}
	return m, err
}
