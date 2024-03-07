package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/jacekdobrowolski/simple_scrapper/pkg/cache"
	"golang.org/x/net/html"
)

type Scraper struct {
	httpClient *http.Client
	cache      *cache.Cache[[]wordCount]
	swg        *SemaphoredWaitGroup
	results    chan ScrapeResult
}

type wordCount struct {
	word  string
	count int
}

type ScrapeResult struct {
	url string
	wq  []wordCount
}

type SemaphoredWaitGroup struct {
	wg  sync.WaitGroup
	sem chan struct{}
}

func (s *SemaphoredWaitGroup) Add(delta int) {
	s.wg.Add(delta)
	s.sem <- struct{}{}
}
func (s *SemaphoredWaitGroup) Done() {
	<-s.sem
	s.wg.Done()
}

func (s *SemaphoredWaitGroup) Wait() {
	s.wg.Wait()
}

func NewSemaphoredWaitGroup(maxConcurrency int) *SemaphoredWaitGroup {
	return &SemaphoredWaitGroup{
		sem: make(chan struct{}, maxConcurrency),
	}
}

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

	s := Scraper{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		swg:        NewSemaphoredWaitGroup(2),
		results:    make(chan ScrapeResult, len(urls)),
		cache:      cache.New[[]wordCount](),
	}
	s.scrape(urls)
}

func (s *Scraper) scrape(urls []string) {
	collectingResultsFinished := make(chan struct{})
	go s.collectResults(collectingResultsFinished)

	for _, url := range urls {
		s.swg.Add(1)
		go s.scrapeUrl(url)
	}
	s.swg.Wait()
	close(s.results)
	<-collectingResultsFinished
}

func (s *Scraper) collectResults(finished chan struct{}) {
	for {
		result, more := <-s.results
		s.cache.Set(result.url, result.wq)
		if more {
			fmt.Println("done", result.url, result.wq[len(result.wq)-5:len(result.wq)])
		} else {
			finished <- struct{}{}
		}
	}
}

func (s *Scraper) scrapeUrl(url string) {
	defer s.swg.Done()
	cachedResult, ok := s.cache.Get(url)
	if ok {
		fmt.Println("from cache", url)
		s.results <- ScrapeResult{url, cachedResult}
		return
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("requesting", url)
	res, err := s.httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	result := getWordFreq(res)
	s.results <- ScrapeResult{url, result}
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
