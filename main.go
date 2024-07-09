package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
)

const toolName = "hackJS"

func main() {
	printBanner()

	if len(os.Args) < 3 {
		fmt.Println("Usage: hackJS -u <URL>")
		return
	}

	targetUrl := os.Args[2]
	fmt.Printf("Starting %s...\n", toolName)

	resp, err := httpGet(targetUrl)
	if err != nil {
		fmt.Printf("Error fetching the URL: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading the response body: %v\n", err)
		return
	}

	jsFiles := extractJSFiles(string(body), targetUrl)
	if len(jsFiles) == 0 {
		fmt.Println("No JavaScript files found.")
		return
	}

	var results []string
	var subdomains []string

	for _, jsFile := range jsFiles {
		jsContent, err := fetchJSContent(jsFile)
		if err != nil {
			fmt.Printf("Error fetching JS file %s: %v\n", jsFile, err)
			continue
		}

		results = append(results, filterLinks(extractLinks(jsContent, targetUrl), targetUrl)...)
		subdomains = append(subdomains, filterSubdomains(extractSubdomains(jsContent, targetUrl), targetUrl)...)
	}

	results = removeDuplicates(results)
	subdomains = removeDuplicates(subdomains)

	printResults("Links", results)
	printResults("Subdomains", subdomains)

	saveResults(targetUrl, results, subdomains)
}

func printBanner() {
	fmt.Println(`
 __                            __           _____   ______  
/  |                          /  |         /     | /      \ 
$$ |____    ______    _______ $$ |   __    $$$$$ |/$$$$$$  |
$$      \  /      \  /       |$$ |  /  |      $$ |$$ \__$$/ 
$$$$$$$  | $$$$$$  |/$$$$$$$/ $$ |_/$$/  __   $$ |$$      \ 
$$ |  $$ | /    $$ |$$ |      $$   $$<  /  |  $$ | $$$$$$  |
$$ |  $$ |/$$$$$$$ |$$ \_____ $$$$$$  \ $$ \__$$ |/  \__$$ |
$$ |  $$ |$$    $$ |$$       |$$ | $$  |$$    $$/ $$    $$/ 
$$/   $$/  $$$$$$$/  $$$$$$$/ $$/   $$/  $$$$$$/   $$$$$$/  
                                                            
                                                            
                                                            
`)
	fmt.Println("Created by: everythingBlackkk\n")
}

func httpGet(targetUrl string) (*http.Response, error) {
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: customTransport}
	return client.Get(targetUrl)
}

func extractJSFiles(html, baseURL string) []string {
	re := regexp.MustCompile(`src="([^"]+\.js)"`)
	matches := re.FindAllStringSubmatch(html, -1)

	var jsFiles []string
	for _, match := range matches {
		jsFile := match[1]
		if !strings.HasPrefix(jsFile, "http") {
			jsFile = baseURL + "/" + jsFile
		}
		// Remove any trailing HTML tags or attributes
		jsFile = strings.Split(jsFile, `"`)[0]
		jsFiles = append(jsFiles, jsFile)
	}
	return jsFiles
}

func fetchJSContent(jsFile string) (string, error) {
	resp, err := httpGet(jsFile)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func extractLinks(jsContent string, baseURL string) []string {
	lines := strings.Split(jsContent, "\n")
	baseDomain := extractDomain(baseURL)
	var matches []string
	re := regexp.MustCompile(`https?://[^\s"<]+`)
	for _, line := range lines {
		lineMatches := re.FindAllString(line, -1)
		for _, match := range lineMatches {
			if strings.Contains(match, baseDomain) {
				matches = append(matches, match)
			}
		}
	}
	return matches
}

func extractSubdomains(jsContent string, baseURL string) []string {
	lines := strings.Split(jsContent, "\n")
	baseDomain := extractDomain(baseURL)
	var matches []string
	re := regexp.MustCompile(`\b(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z]{2,6}\b`)
	for _, line := range lines {
		lineMatches := re.FindAllString(line, -1)
		for _, match := range lineMatches {
			if strings.Contains(match, baseDomain) {
				matches = append(matches, match)
			}
		}
	}
	return matches
}

func filterLinks(links []string, baseURL string) []string {
	baseDomain := extractDomain(baseURL)
	var filteredLinks []string
	encountered := map[string]bool{}
	for _, link := range links {
		if !encountered[link] && strings.Contains(link, baseDomain) {
			encountered[link] = true
			filteredLinks = append(filteredLinks, link)
		}
	}
	return filteredLinks
}

func filterSubdomains(subdomains []string, baseURL string) []string {
	baseDomain := extractDomain(baseURL)
	var filteredSubdomains []string
	encountered := map[string]bool{}
	for _, subdomain := range subdomains {
		if !encountered[subdomain] && strings.Contains(subdomain, baseDomain) {
			encountered[subdomain] = true
			filteredSubdomains = append(filteredSubdomains, subdomain)
		}
	}
	return filteredSubdomains
}

func saveResults(targetUrl string, results, subdomains []string) {
	domain := extractDomain(targetUrl)
	if domain == "" {
		fmt.Println("Invalid URL provided.")
		return
	}

	dir := "result"
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creating result directory: %v\n", err)
		return
	}

	fileName := fmt.Sprintf("%s/%s.txt", dir, domain)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error creating result file: %v\n", err)
		return
	}
	defer file.Close()

	writeSection(file, "Links", results)
	writeSection(file, "Subdomains", subdomains)

	fmt.Printf("Results saved to %s\n", fileName)
}

func writeSection(file *os.File, title string, results []string) {
	file.WriteString(fmt.Sprintf("\n%s:\n", title))
	for _, result := range results {
		file.WriteString(result + "\n")
	}
}

func printResults(title string, results []string) {
	fmt.Printf("\n%s:\n", title)
	for _, result := range results {
		fmt.Println(result)
	}
	fmt.Printf("Total %s found: %d\n", title, len(results))
}

func extractDomain(targetUrl string) string {
	parsedUrl, err := url.Parse(targetUrl)
	if err != nil {
		return ""
	}
	return parsedUrl.Host
}

func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for _, v := range elements {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}

	sort.Strings(result)
	return result
}
