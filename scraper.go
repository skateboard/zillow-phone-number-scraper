package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/data-harvesters/goapify"
)

type scraper struct {
	actor  *goapify.Actor
	input  *input
	client tls_client.HttpClient
}

func newScraper(input *input, actor *goapify.Actor) (*scraper, error) {
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Chrome_124),
		tls_client.WithNotFollowRedirects(),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return nil, err
	}

	return &scraper{
		actor:  actor,
		input:  input,
		client: client,
	}, nil
}

func (s *scraper) Run() {
	fmt.Println("beginning scrapping...")

	var wg sync.WaitGroup
	for _, zpid := range s.input.ZPids {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r, err := s.scrapePhoneNumber(zpid)
			if err != nil {
				fmt.Printf("%s: failed to scrape phone number: %v\n", zpid, err)
				return
			}

			err = s.actor.Output(map[string]string{
				"zpid":          zpid,
				"display_name":  r.PropertyInfo.AgentInfo.DisplayName,
				"business_name": r.PropertyInfo.AgentInfo.BusinessName,
				"phone_number":  r.PropertyInfo.AgentInfo.PhoneNumber,
			})
			if err != nil {
				fmt.Printf("%s: failed to send output: %v\n", zpid, err)
				return
			}
		}()
	}
	wg.Wait()
	fmt.Println("succesfully scraped all phone numbers")
}

func (s *scraper) scrapePhoneNumber(zPid string) (*response, error) {
	payload := fmt.Sprintf("{\"zpid\":\"%s\",\"pageType\":\"BDP\",\"isInstantTourEnabled\":false,\"isCachedInstantTourAvailability\":true,\"tourTypes\":[]}", zPid)

	req, err := http.NewRequest("POST", "https://www.zillow.com/rentals/api/rcf/v1/rcf", strings.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Host", "www.zillow.com")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:131.0) Gecko/20100101 Firefox/131.0")
	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("accept-language", "en-US,en;q=0.5")
	req.Header.Add("content-type", "application/json;charset=utf-8")
	req.Header.Add("origin", "https://www.zillow.com")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("te", "trailers")

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to scrape reviews: %d %s", res.StatusCode, string(b))
	}

	var response response
	err = json.Unmarshal(b, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *scraper) getApiKey(roomId string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.airbnb.com/rooms/%s/reviews", roomId), nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Host", "www.airbnb.com")
	req.Header.Add("Sec-Ch-Ua", "\"Chromium\";v=\"127\", \"Not)A;Brand\";v=\"99\"")
	req.Header.Add("Accept-Language", "en-US")
	req.Header.Add("Sec-Ch-Ua-Platform-Version", "\"\"")
	req.Header.Add("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Add("Device-Memory", "8")
	req.Header.Add("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.6533.100 Safari/537.36")
	req.Header.Add("Viewport-Width", "1536")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Ect", "3g")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Priority", "u=1, i")

	res, err := s.client.Do(req)
	if err != nil {
		return "", err

	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get api key: %d", res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	matches := apiKeyRe.FindAllStringSubmatch(string(b), -1)

	if len(matches) == 0 {
		return "", errors.New("failed to find api key")
	}
	apiKey := matches[0][1]

	return apiKey, nil
}
