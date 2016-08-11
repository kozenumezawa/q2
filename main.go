package main

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/html"
)

type Page struct {
	Title       string
	Description string
}

func isDescription(attrs []html.Attribute) bool {
	for _, attr := range attrs {
		if attr.Key == "name" && attr.Val == "description" {
			return true
		}
	}
	return false
}

func Get(url string) (*Page, error) {
	gotPage := &Page

	resp, err := http.Get(url)
	if err != nil {
		// 実際にはちゃんとエラー処理をしましょう
		panic(err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		// 実際にはちゃんとエラー処理しましょう
		panic(err)
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			gotPage.Title = n.FirstChild.Data
		}
		if n.Type == html.ElementNode && n.Data == "meta" {
			if isDescription(n.Attr) {
				for _, attr := range n.Attr {
					fmt.Println(attr)
					gotPage.Description = attr.Val
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return gotPage, nil
}

func main() {
	// Getを利用しているとします
	p, err := Get("http://voyagegroup.com")
	if err != nil && err != io.EOF {
		panic(err)
	}
	fmt.Println("%#v", p)
}
