# hackJS
To use hackJS, specify a target URL that contains JavaScript files you want to analyze. The tool will fetch these JavaScript files, extract embedded links and subdomains and sensitive data in js file " -w ".
____
## Features

- JavaScript Extraction: Fetches JavaScript files linked in HTML pages.
- Subdomain Extraction: Identifies subdomains referenced in JavaScript files.
- Results Organization: Outputs and saves extracted data into categorized result files.
- Get sensitive data in javascript files like Token,Api,..... " But You Must put your word list '-w' "

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

To use hackJS, specify a target URL that contains JavaScript files you want to analyze. The tool will fetch these JavaScript files, extract embedded links and subdomains and 
sensitive data in all js file using " -w " flag.

   ```bash
    go run main.go -u https://everythingBlackkk.com
   ```
   ```bash
    go run main.go -l urls.txt
   ```
   ```bash
    go run main.go -l url.txt -w sensitive-world-list.txt
   ```
    
    
    
## Contributing
Contributions are welcome! Fork the repository, make your changes, and submit a pull request.


