package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	fmt.Println("banan")

	res, err := http.Get("https://www.ipdb.org/machine.cgi?id=20")
	if err != nil {
		fmt.Printf("%s", errors.New(err.Error()))
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Printf("Response status code was %d\n", res.StatusCode)
	}

	ct := res.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "text/html") {
		fmt.Printf("Response content type was %s, not text/html\n", ct)
	}

	tokenizer := html.NewTokenizer(res.Body)
	
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				break
			}
			fmt.Println(tokenizer.Err())
		}

		if tokenType == html.StartTagToken {
			token := tokenizer.Token()
			if token.Data == "title" {
				tokenType = tokenizer.Next()
				if tokenType == html.TextToken {
					fmt.Println(tokenizer.Token().Data)
					break
				}
			}
		}
	}
}
