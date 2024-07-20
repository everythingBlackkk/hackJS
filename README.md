# hackJS
hackJS is a tool designed to analyze JavaScript files from a specified URL or a list of URLs to extract useful information such as links, subdomains, and sensitive data. It can also work with wordlists to search for specific keywords within JavaScript files. This tool helps in discovering hidden resources and potential security issues on websites.

____
## Features

- Extract Links: Finds and filters all links in JavaScript files.
- Extract Subdomains: Identifies subdomains mentioned in JavaScript files.
- Find Sensitive Data: Searches for sensitive words in JavaScript files using a provided or default wordlist.
- Output Results: Saves results to a file and displays them on the console.
- Handle Multiple URLs: Can process a single URL or multiple URLs from a file.

## Installation

1. Clone the repository to your local machine:

    ```bash
    git clone https://github.com/everythingBlackkk/hackJS.git
    ```

2. Navigate to the hackJS directory:

    ```bash
    cd hackJS
    ```
3. run tool:
   ```bash
    go run main.go
   ```
   
## Usage


   ```Basic Usage
    go run main.go -u <URL>
   ```
   ```Scan Multiple URLs from a File:
    go run main.go -l <file>
   ```
   ```bash
    go run main.go -l url.txt -w sensitive-world-list.txt
   ```

## Options
- -u <URL>: Specifies the URL to scan.
- -l <file>: Specifies a file containing a list of URLs to scan.
- -w <wordlist>: Specifies a custom wordlist file to search for sensitive data.

## Output
The results are categorized and saved into a result directory. Each category includes:

#Sample Output

```golang
===Links===
https://example.com/page1
https://example.com/page2
Total Links found: 2

===Subdomains===
sub.example.com
api.example.com
Total Subdomains found: 2

===JS Files===
https://example.com/script1.js
https://example.com/script2.js
Total JS Files found: 2

===Sensitive Data===
🔹 api_key ➔ https://example.com/script1.js
🔹 token ➔ https://example.com/script2.js
Total Sensitive Data found: 2

```

##Contact
For any questions or feedback, please contact:

Name: Yassin Abd-elrazik
GitHub: everythingBlackkk

## Contributing
Contributions are welcome! Fork the repository, make your changes, and submit a pull request.
