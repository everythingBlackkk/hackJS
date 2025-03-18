# HackJS

## Introduction

HackJS is a comprehensive web reconnaissance tool designed for security researchers and bug hunters. It efficiently crawls websites to extract JavaScript files, discover subdomains, and search for specific keywords within JavaScript code that might indicate vulnerabilities or sensitive information.

## Features

- **JavaScript File Discovery**: Automatically detects and extracts JavaScript files from target websites
- **Subdomain Enumeration**: Identifies subdomains related to the target domain
- **Keyword Searching**: Searches for specific keywords or patterns within JavaScript files
- **Configurable Crawl Depth**: Set how deep the crawler should traverse the website
- **Multi-threading Support**: Parallel processing for faster reconnaissance
- **Detailed Output**: Comprehensive reporting of findings with context for keyword matches
- **File Export**: Save results to a file for further analysis

## Installation

### Prerequisites

- Python 3.6+
- pip (Python package installer)

### Setup

1. Clone the repository:
```bash
git clone https://github.com/everythingBlackkk/hackJS.git
cd hackJS
```

2. Install required dependencies:
```bash
pip install -r requirements.txt
```

Alternatively, you can install dependencies manually:
```bash
pip install requests beautifulsoup4 colorama
```

## Usage

### Basic Usage

```bash
python3 hackjs.py -u example.com
```

### Command Line Arguments

| Option | Description |
|--------|-------------|
| `-u, --url` | Target URL to scan (required) |
| `-w, --wordlist` | Path to wordlist file for keyword search |
| `-t, --threads` | Number of threads (default: 5) |
| `-d, --depth` | Crawling depth (default: 2) |
| `-o, --output` | Output file to save results |

### Examples

Scan with default settings:
```bash
python3 hackjs.py -u example.com
```

Scan with custom depth and threads:
```bash
python3 hackjs.py -u example.com -d 3 -t 10
```

Scan with keyword search:
```bash
python3 hackjs.py -u example.com -w keywords.txt
```

Save results to a file:
```bash
python3 hackjs.py -u example.com -o results.txt
```

### Creating a Wordlist

Create a text file with keywords or patterns to search for in JavaScript files, one per line:

```
api_key
password
token
secret
credentials
```

## Output

WebRecon provides comprehensive output displaying:
- Discovered JavaScript files
- Identified subdomains
- Keyword matches with surrounding context
- Summary statistics

When saving to a file using the `-o` option, results are formatted in clear sections for easy analysis.


## Use Cases

- Bug bounty hunting
- Penetration testing
- Security assessments
- Asset discovery
- Sensitive information detection


## Author

**Yassin Mohamed** - [@everythingBlackkk](https://github.com/everythingBlackkk)

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Disclaimer

WebRecon is intended for use by security professionals with proper authorization to test target systems. Unauthorized scanning of websites may violate laws and regulations. Always obtain proper permission before scanning websites you don't own.
