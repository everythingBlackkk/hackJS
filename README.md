# hackJS
To use hackJS, specify a target URL that contains JavaScript files you want to analyze. The tool will fetch these JavaScript files, extract embedded links and subdomains.
____
## Features

- JavaScript Extraction: Fetches JavaScript files linked in HTML pages.
- Subdomain Extraction: Identifies subdomains referenced in JavaScript files.
- Results Organization: Outputs and saves extracted data into categorized result files.

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

To use hackJS, specify a target URL that contains JavaScript files you want to analyze. The tool will fetch these JavaScript files, extract embedded links and subdomains
    ```
    go run main.go -u https://example.com
    ```
    
## Contributing
Contributions are welcome! Fork the repository, make your changes, and submit a pull request.


