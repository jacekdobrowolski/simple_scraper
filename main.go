package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

type wordCount struct {
	word  string
	count int
}

type ScrapeResult struct {
	url string
	wq  []wordCount
}

type WordFreqCache struct {
	urlMap map[string][]wordCount
	mux    sync.RWMutex
}

func (c *WordFreqCache) Set(url string, wordFreq []wordCount) {
	defer c.mux.Unlock()
	c.mux.Lock()
	c.urlMap[url] = wordFreq
}

func (c *WordFreqCache) Get(url string) ([]wordCount, bool) {
	defer c.mux.RUnlock()
	c.mux.RLock()
	wordFreq, ok := c.urlMap[url]
	return wordFreq, ok
}

func NewWordFreqCache() *WordFreqCache {
	return &WordFreqCache{urlMap: make(map[string][]wordCount)}
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
	httpClient := &http.Client{Timeout: 10 * time.Second}
	scrape(urls, httpClient)
}

func scrape(urls []string, httpClient *http.Client) {
	swg := NewSemaphoredWaitGroup(2)
	results := make(chan ScrapeResult, len(urls))
	collectingResultsFinished := make(chan struct{})
	cache := NewWordFreqCache()
	go collectResults(results, cache, collectingResultsFinished)

	for _, url := range urls {
		swg.Add(1)
		go scrapeUrl(url, httpClient, cache, swg, results)
	}
	swg.Wait()
	close(results)
	<-collectingResultsFinished
}

func collectResults(results chan ScrapeResult, cache *WordFreqCache, finished chan struct{}) {
	for {
		result, more := <-results
		cache.Set(result.url, result.wq)
		if more {
			fmt.Println("done", result.url, result.wq[len(result.wq)-3:len(result.wq)])
		} else {
			finished <- struct{}{}
		}
	}
}

func scrapeUrl(url string, httpClient *http.Client, cache *WordFreqCache, swg *SemaphoredWaitGroup, results chan ScrapeResult) {
	defer swg.Done()
	cachedResult, ok := cache.Get(url)
	if ok {
		fmt.Println("from cache", url)
		results <- ScrapeResult{url, cachedResult}
		return
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("requesting", url)
	res, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	result := getWordFreq(res)
	results <- ScrapeResult{url, result}
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
				field = strings.Trim(field, ".,„“”'\";:()[]{}")
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
