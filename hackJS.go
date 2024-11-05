package main

import (
    "bufio"
    "crypto/tls"
    "flag"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "regexp"
    "sort"
    "strings"
    "time"
)

var (
    urlsFile      string
    wordlistFile  string
    timeout       int
    outputDir     string
    saveResults   bool
    sensitiveWords []string
)

func main() {
    parseCommandLineArgs()
    printBanner()
    loadWordlist()
    processInputURLs()
}

func parseCommandLineArgs() {
    flag.StringVar(&urlsFile, "i", "", "File containing a list of URLs to analyze")
    flag.StringVar(&wordlistFile, "w", "", "File containing a list of sensitive words")
    flag.IntVar(&timeout, "t", 30, "Timeout for HTTP requests (in seconds)")
    flag.StringVar(&outputDir, "o", "", "Output directory for results (default is $HOME/hackJS_results)")
    flag.BoolVar(&saveResults, "s", true, "Save results to files (default is true)")
    flag.Parse()
}

func loadWordlist() {
    if wordlistFile != "" {
        file, err := os.Open(wordlistFile)
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
    } else {
        loadDefaultWordlist()
    }
}

func loadDefaultWordlist() {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        fmt.Printf("Error getting home directory: %v\n", err)
        return
    }
    fileName := filepath.Join(homeDir, "bin", "WordList.txt")
    file, err := os.Open(fileName)
    if err != nil {
        fmt.Println("\033[31mWarning: The file WordList.txt is missing. Please download it from GitHub.\033[0m")
        return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        sensitiveWords = append(sensitiveWords, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        fmt.Printf("Error reading default wordlist file: %v\n", err)
    }
}

func processInputURLs() {
    if urlsFile == "" {
        fmt.Println("Please provide a file containing the URLs to analyze.")
        return
    }

    file, err := os.Open(urlsFile)
    if err != nil {
        fmt.Printf("Error opening URLs file: %v\n", err)
        return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        targetURL := scanner.Text()
        fmt.Printf("\nProcessing URL: %s\n", targetURL)
        processURL(targetURL)
        fmt.Println("_____________________________________________________________________________________________")
    }

    if err := scanner.Err(); err != nil {
        fmt.Printf("Error reading URLs file: %v\n", err)
    }
}

func processURL(targetURL string) {
    resp, err := httpGet(targetURL, timeout)
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

    jsFiles := extractJSFiles(string(body), targetURL)
    if len(jsFiles) == 0 {
        fmt.Println("No JavaScript files found.")
        return
    }

    var results []string
    var subdomains []string
    var sensitiveData []string

    for _, jsFile := range jsFiles {
        jsContent, err := fetchJSContent(jsFile, timeout)
        if err != nil {
            fmt.Printf("Error fetching JS file %s: %v\n", jsFile, err)
            continue
        }

        results = append(results, filterLinks(extractLinks(jsContent, targetURL), targetURL)...)
        subdomains = append(subdomains, filterSubdomains(extractSubdomains(jsContent, targetURL), targetURL)...)
        sensitiveData = append(sensitiveData, findSensitiveData(jsContent, jsFile)...)
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

    if saveResults {
        saveResultsToFiles(targetURL, results, subdomains, jsFiles, sensitiveData)
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

func httpGet(targetURL string, timeout int) (*http.Response, error) {
    customTransport := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{
        Transport: customTransport,
        Timeout:   time.Duration(timeout) * time.Second,
    }
    return client.Get(targetURL)
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

func fetchJSContent(jsFile string, timeout int) (string, error) {
    resp, err := httpGet(jsFile, timeout)
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
    encountered := make(map[string]bool)
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
    encountered := make(map[string]bool)
    for _, subdomain := range subdomains {
        if !encountered[subdomain] && strings.HasSuffix(subdomain, baseDomain) {
            encountered[subdomain] = true
            filteredSubdomains = append(filteredSubdomains, subdomain)
        }
    }
    return filteredSubdomains
}

func removeDuplicates(elements []string) []string {
    encountered := make(map[string]bool)
    var result []string

    for _, v := range elements {
        if !encountered[v] {
            encountered[v] = true
            result = append(result, v)
        }
    }

    sort.Strings(result)
    return result
}

func cleanURL(dirtyURL string) string {
    cleanURL, err := url.Parse(dirtyURL)
    if err != nil {
        return dirtyURL
    }
    cleanURL.Fragment = ""
    return cleanURL.String()
}

func extractDomain(rawURL string) string {
    parsedURL, err := url.Parse(rawURL)
    if err != nil {
        return ""
    }

    host := parsedURL.Hostname()
    parts := strings.Split(host, ".")
    if len(parts) >= 2 {
        return parts[len(parts)-2] + "." + parts[len(parts)-1]
    }

    return host
}

func printResults(label string, results []string, colorCode string) {
    if len(results) > 0 {
        fmt.Printf("\n%s%s:\033[0m\n", colorCode, label)
        for _, result := range results {
            fmt.Println(result)
        }
    }
}

func saveResultsToFiles(targetURL string, links, subdomains, jsFiles, sensitiveData []string) {
    domain := extractDomain(targetURL)
    if domain == "" {
        fmt.Println("Invalid URL provided.")
        return
    }

    if outputDir == "" {
        homeDir, err := os.UserHomeDir()
        if err != nil {
            fmt.Printf("Error getting user home directory: %v\n", err)
            return
        }
        outputDir = filepath.Join(homeDir, "hackJS_results")
    }

    resultsDir := filepath.Join(outputDir, domain)
    if err := os.MkdirAll(resultsDir, 0755); err != nil {
        fmt.Printf("Error creating results directory: %v\n", err)
        return
    }

    saveToFile(filepath.Join(resultsDir, "links.txt"), links)
    saveToFile(filepath.Join(resultsDir, "subdomains.txt"), subdomains)
    saveToFile(filepath.Join(resultsDir, "jsfiles.txt"), jsFiles)
    if len(sensitiveData) > 0 {
        saveToFile(filepath.Join(resultsDir, "sensitive.txt"), sensitiveData)
    }

    fmt.Printf("Results saved to: %s\n", resultsDir)
}

func saveToFile(fileName string, data []string) {
    file, err := os.Create(fileName)
    if err != nil {
        fmt.Printf("Error creating file %s: %v\n", fileName, err)
        return
    }
    defer file.Close()

    for _, line := range data {
        _, err := file.WriteString(line + "\n")
        if err != nil {
            fmt.Printf("Error writing to file %s: %v\n", fileName, err)
            return
        }
    }
}
