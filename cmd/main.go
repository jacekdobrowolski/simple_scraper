package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/jacekdobrowolski/simple_scrapper/pkg/cache"
	"golang.org/x/net/html"
	"golang.org/x/sync/errgroup"
)

func main() {
	urls := []string{
		"https://pl.wikipedia.org/wiki/Polska",
		"https://pl.wikipedia.org/wiki/Niemcy",
		"https://pl.wikipedia.org/wiki/Francja",
		"https://pl.wikipedia.org/wiki/Polska",
		"https://pl.wikipedia.org/wiki/Polska",
		"https://pl.wikipedia.org/wiki/Polska",
		"https://pl.wikipedia.org/wiki/Polska",
	}

	s, ctx := NewScraper(context.Background(), &http.Client{Timeout: 10 * time.Second}, 2)
	s.scrape(urls)
	<-ctx.Done()
}

type wordCount struct {
	word  string
	count int
}

type ScrapeResult struct {
	url string
	wq  []wordCount
}

type Scraper struct {
	httpClient *http.Client
	cache      *cache.Cache[[]wordCount]
	group      *errgroup.Group
	results    chan ScrapeResult
}

func NewScraper(ctx context.Context, httpClient *http.Client, maxConnections int) (*Scraper, context.Context) {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(maxConnections)
	return &Scraper{
		httpClient: httpClient,
		group:      g,
		results:    make(chan ScrapeResult),
		cache:      cache.New[[]wordCount](),
	}, ctx
}

func (s *Scraper) scrape(urls []string) {
	collectingResultsFinished := make(chan struct{})
	go s.collectResults(collectingResultsFinished)

	for _, url := range urls {
		url := url
		s.group.Go(func() error { return s.scrapeUrl(url) })
	}
	err := s.group.Wait()
	if err != nil {
		log.Fatal(err)
	}
	close(s.results)
	<-collectingResultsFinished
}

func (s *Scraper) collectResults(finished chan struct{}) {
	for {
		result, more := <-s.results
		if more {
			fmt.Println("done", result.url, result.wq[len(result.wq)-5:len(result.wq)])
		} else {
			finished <- struct{}{}
		}
	}
}

func (s *Scraper) scrapeUrl(url string) error {
	cachedResult, ok := s.cache.Get(url)
	if ok {
		fmt.Println("from cache", url)
		s.results <- ScrapeResult{url, cachedResult}
		return nil
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	fmt.Println("requesting", url)
	res, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	result := getWordFreq(res)
	s.cache.Set(url, result)
	s.results <- ScrapeResult{url, result}
	return nil
}

func getWordFreq(res *http.Response) []wordCount {
	tokenizer := html.NewTokenizer(res.Body)
	wordFreq := map[string]int{}
	for {
		tokenType := tokenizer.Next()
		tagName, _ := tokenizer.TagName()
		tagNameString := string(tagName)
		if tagNameString == "script" || tagNameString == "noscript" {
			for {
				tokenType := tokenizer.Next()
				tagName, _ := tokenizer.TagName()
				endTagNameString := string(tagName)
				if tokenType == html.EndTagToken && endTagNameString == tagNameString {
					break
				}
			}
		}
		if tokenType == html.TextToken {
			text := strings.TrimSpace(tokenizer.Token().String())
			text = strings.ToLower(text)
			fields := strings.Fields(text)
			for _, field := range fields {
				field = strings.Trim(field, ".,„“”'\";:()[]{}–↑")
				if field != "" {
					wordFreq[field]++
				}
			}
		} else if tokenType == html.ErrorToken {
			break
		}
	}
	result := make([]wordCount, 0, len(wordFreq))
	for word, freq := range wordFreq {
		result = append(result, wordCount{word, freq})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].count < result[j].count
	})
	return result
}
