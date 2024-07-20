package main

import (
	"bufio"
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

var sensitiveWords []string

func main() {
	printBanner()

	if len(os.Args) < 3 {
		fmt.Println("Usage: hackJS -u <URL> or hackJS -l <file> [-w <wordlist>]")
		return
	}

	wordlistProvided := false
	if len(os.Args) == 5 && os.Args[3] == "-w" {
		loadWordlist(os.Args[4])
		wordlistProvided = true
	} else {
		if !loadDefaultWordlist() {
			fmt.Println("\033[31mWarning: The file WordList.txt is missing. Please download it from GitHub.\033[0m")
		}
	}

	if os.Args[1] == "-u" {
		processURL(os.Args[2], wordlistProvided)
	} else if os.Args[1] == "-l" {
		processFile(os.Args[2], wordlistProvided)
	} else {
		fmt.Println("Invalid option. Use -u for a URL or -l for a file.")
	}
}

func loadWordlist(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening wordlist file: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sensitiveWords = append(sensitiveWords, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading wordlist file: %v\n", err)
	}
}

func loadDefaultWordlist() bool {
	fileName := "WordList.txt"
	file, err := os.Open(fileName)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sensitiveWords = append(sensitiveWords, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading default wordlist file: %v\n", err)
	}
	return true
}

func processURL(targetUrl string, wordlistProvided bool) {
	fmt.Printf("\nStarting %s...\n", "hackJS")

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
	var sensitiveData []string

	for _, jsFile := range jsFiles {
		jsContent, err := fetchJSContent(jsFile)
		if err != nil {
			fmt.Printf("Error fetching JS file %s: %v\n", jsFile, err)
			continue
		}

		results = append(results, filterLinks(extractLinks(jsContent, targetUrl), targetUrl)...)
		subdomains = append(subdomains, filterSubdomains(extractSubdomains(jsContent, targetUrl), targetUrl)...)
		if wordlistProvided {
			sensitiveData = append(sensitiveData, findSensitiveData(jsContent, jsFile)...)
		}
	}

	results = removeDuplicates(results)
	subdomains = removeDuplicates(subdomains)
	jsFiles = removeDuplicates(jsFiles)
	sensitiveData = removeDuplicates(sensitiveData)

	printResults("Links", results, "\033[32m")
	printResults("Subdomains", subdomains, "\033[36m")
	printResults("JS Files", jsFiles, "\033[33m")
	if len(sensitiveData) > 0 {
		printResults("Sensitive Data", sensitiveData, "\033[31m")
	} else {
		fmt.Println("\n\033[31mNo sensitive data found.\033[0m")
	}

	saveResults(targetUrl, results, subdomains, jsFiles, sensitiveData)
}

func processFile(fileName string, wordlistProvided bool) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		processURL(scanner.Text(), wordlistProvided)
		fmt.Println("_____________________________________________________________________________________________")
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
}

func printBanner() {
	fmt.Println("\033[32m")
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
	fmt.Println("          # hackJS , Coded By Yassin Abd-elrazik")
	fmt.Println("          Made By <3 github : everythingBlackkk")
	fmt.Println("\033[0m")
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
		jsFile = cleanURL(jsFile)
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
	re := regexp.MustCompile(`https?://[^\s"<>()']+`)
	for _, line := range lines {
		lineMatches := re.FindAllString(line, -1)
		for _, match := range lineMatches {
			if strings.Contains(match, baseDomain) && !strings.HasSuffix(match, ".js") {
				matches = append(matches, cleanURL(match))
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

func findSensitiveData(jsContent, jsFile string) []string {
	var matches []string
	for _, word := range sensitiveWords {
		if strings.Contains(jsContent, word) {
			matches = append(matches, fmt.Sprintf("🔹 %s ➔ %s", word, jsFile))
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

func saveResults(targetUrl string, results, subdomains, jsFiles, sensitiveData []string) {
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

	writeSection(file, "===Links===", results)
	writeSection(file, "===Subdomains===", subdomains)
	writeSection(file, "===JS Files===", jsFiles)
	writeSection(file, "===Sensitive Data===", sensitiveData)

	fmt.Printf("\n\033[32mResults saved to %s\n\033[0m", fileName)
}

func writeSection(file *os.File, title string, results []string) {
	file.WriteString(fmt.Sprintf("\n%s:\n", title))
	for _, result := range results {
		file.WriteString(fmt.Sprintf("%s\n", result))
	}
}

func printResults(title string, results []string, color string) {
	fmt.Printf("\n%s%s:\n", color, title)
	for _, result := range results {
		fmt.Printf("%s%s\n", color, result)
	}
	fmt.Printf("\033[0mTotal %s found: %d\n", title, len(results))
}

func extractDomain(targetUrl string) string {
	parsedUrl, err := url.Parse(targetUrl)
	if err != nil {
		return ""
	}
	return parsedUrl.Host
}

func cleanURL(rawURL string) string {
	trimChars := []string{`"`, `'`, `>`, `<`, `)`}
	for _, char := range trimChars {
		if idx := strings.Index(rawURL, char); idx != -1 {
			rawURL = rawURL[:idx]
		}
	}
	return rawURL
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
