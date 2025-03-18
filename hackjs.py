#!/usr/bin/env python3
"""
WebRecon - A web reconnaissance tool for security researchers and bug hunters
This tool crawls a website to extract JavaScript files, discover subdomains,
and search for specific keywords in JavaScript files.
"""

import argparse
import re
import sys
import requests
from bs4 import BeautifulSoup
from urllib.parse import urlparse, urljoin
import concurrent.futures
from colorama import Fore, Style, init

init(autoreset=True)

# Global variables We Will Use :) 
visited_urls = set()
js_files = set()
subdomains = set()
found_keywords = {}

def print_banner():
    banner = f"""
    
Hack JavaScript                 |     |
                                \\_V_//
                                \/=|=\/
                                 [=v=]
                               __\___/_____
                              /..[  _____  ]
                             /_  [ [  M /] ]
                            /../.[ [ M /@] ]
                           <-->[_[ [M /@/] ]
                          /../ [.[ [ /@/ ] ]
     _________________]\ /__/  [_[ [/@/ C] ]
    <_________________>>0---]  [=\ \@/ C / /
       ___      ___   ]/000o   /__\ \ C / /
          \    /              /....\ \_/ /
       ....\||/....           [___/=\___/
      .    .  .    .          [...] [...]
     .      ..      .         [___/ \___]
     .    0 .. 0    .         <---> <--->
  /\/\.    .  .    ./\/\      [..]   [..]
 / / / .../|  |\... \ \ \    _[__]   [__]_
/ / /       \/       \ \ \  [____>   <____]

    {Fore.YELLOW}[-u URL] [-w WORDLIST] [-t THREADS] [-d DEPTH]
    {Fore.GREEN} @Github : everthingBlackkk
    {Fore.GREEN} #Dev : Yassin Mohamed
    """
    print(banner)

def parse_args():
    parser = argparse.ArgumentParser(description='Web Reconnaissance Tool')
    parser.add_argument('-u', '--url', required=True, help='Target URL to scan')
    parser.add_argument('-w', '--wordlist', help='Path to wordlist file for keyword search')
    parser.add_argument('-t', '--threads', type=int, default=5, help='Number of threads (default: 5)')
    parser.add_argument('-d', '--depth', type=int, default=2, help='Crawling depth (default: 2)')
    parser.add_argument('-o', '--output', help='Output file to save results')
    return parser.parse_args()

def normalize_url(url):
    if not url.startswith(('http://', 'https://')):
        url = 'https://' + url
    return url

def extract_domain(url):
    parsed_url = urlparse(url)
    return parsed_url.netloc

def is_valid_url(url, base_domain):
    try:
        parsed = urlparse(url)
        return parsed.netloc == base_domain or parsed.netloc.endswith('.' + base_domain)
    except:
        return False

def extract_js_files(url, base_domain):
    global js_files
    try:
        headers = {
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
        }
        response = requests.get(url, headers=headers, timeout=10)
        if response.status_code != 200:
            return
        
        soup = BeautifulSoup(response.text, 'html.parser')
        
        for script in soup.find_all('script'):
            src = script.get('src')
            if src:
                full_url = urljoin(url, src)
                parsed_url = urlparse(full_url)
                
                if parsed_url.path.endswith('.js'):
                    js_files.add(full_url)
                    print(f"{Fore.GREEN}[+] Found JavaScript file: {full_url}")
        
        # Also look for JavaScript files in HTML content
        js_pattern = re.compile(r'(https?://[^\s\'">]+\.js)')
        for match in js_pattern.finditer(response.text):
            js_url = match.group(1)
            js_files.add(js_url)
            print(f"{Fore.GREEN}[+] Found JavaScript file: {js_url}")
        
        extract_subdomains(response.text, base_domain)
        
    except Exception as e:
        print(f"{Fore.RED}[!] Error processing {url}: {str(e)}")

def extract_subdomains(content, base_domain):
    global subdomains
    
    # Match subdomains pattern
    subdomain_pattern = re.compile(r'(?:https?://)?([a-zA-Z0-9][-a-zA-Z0-9]*\.)+' + re.escape(base_domain))
    
    for match in subdomain_pattern.finditer(content):
        subdomain = match.group(0).lower()
        # Clean up the subdomain
        if subdomain.startswith('http://'):
            subdomain = subdomain[7:]
        elif subdomain.startswith('https://'):
            subdomain = subdomain[8:]
        
        if subdomain not in subdomains and subdomain != base_domain:
            subdomains.add(subdomain)
            print(f"{Fore.CYAN}[+] Found subdomain: {subdomain}")

def search_keywords(js_url, keywords):
    """Search for keywords in JavaScript files"""
    global found_keywords
    try:
        headers = {
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
        }
        response = requests.get(js_url, headers=headers, timeout=10)
        if response.status_code != 200:
            return
        
        content = response.text
        
        for keyword in keywords:
            if keyword.lower() in content.lower():
                # Find the context around the keyword
                pattern = re.compile(r'.{0,50}' + re.escape(keyword) + r'.{0,50}', re.IGNORECASE)
                matches = pattern.findall(content)
                
                if js_url not in found_keywords:
                    found_keywords[js_url] = {}
                
                if keyword not in found_keywords[js_url]:
                    found_keywords[js_url][keyword] = []
                
                for match in matches:
                    found_keywords[js_url][keyword].append(match.strip())
                
                print(f"{Fore.YELLOW}[+] Found keyword '{keyword}' in {js_url}")
                
    except Exception as e:
        print(f"{Fore.RED}[!] Error searching keywords in {js_url}: {str(e)}")

def crawl_website(url, depth, base_domain, max_depth):
    """Crawl website to find links and JavaScript files"""
    global visited_urls
    
    if url in visited_urls or depth > max_depth:
        return
    
    visited_urls.add(url)
    print(f"{Fore.BLUE}[*] Crawling: {url} (Depth: {depth}/{max_depth})")
    
    try:
        headers = {
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
        }
        response = requests.get(url, headers=headers, timeout=10)
        if response.status_code != 200:
            return
        
        extract_js_files(url, base_domain)
        
        soup = BeautifulSoup(response.text, 'html.parser')
        links = soup.find_all('a', href=True)
        
        for link in links:
            href = link['href']
            full_url = urljoin(url, href)
            
            if is_valid_url(full_url, base_domain) and full_url not in visited_urls:
                crawl_website(full_url, depth + 1, base_domain, max_depth)
                
    except Exception as e:
        print(f"{Fore.RED}[!] Error crawling {url}: {str(e)}")

def load_wordlist(wordlist_path):
    try:
        with open(wordlist_path, 'r') as f:
            return [line.strip() for line in f if line.strip()]
    except Exception as e:
        print(f"{Fore.RED}[!] Error loading wordlist: {str(e)}")
        sys.exit(1)

def save_results(output_file):
    try:
        with open(output_file, 'w') as f:
            f.write("=== WEB RECONNAISSANCE RESULTS ===\n\n")
            
            f.write("=== DISCOVERED JAVASCRIPT FILES ===\n")
            for js_file in js_files:
                f.write(f"{js_file}\n")
            f.write("\n")
            
            f.write("=== DISCOVERED SUBDOMAINS ===\n")
            for subdomain in subdomains:
                f.write(f"{subdomain}\n")
            f.write("\n")
            
            f.write("=== KEYWORD MATCHES ===\n")
            for js_url, keywords in found_keywords.items():
                f.write(f"\nFile: {js_url}\n")
                for keyword, contexts in keywords.items():
                    f.write(f"  Keyword: {keyword}\n")
                    for context in contexts:
                        f.write(f"    Context: ...{context}...\n")
                f.write("\n")
            
        print(f"{Fore.GREEN}[+] Results saved to {output_file}")
    except Exception as e:
        print(f"{Fore.RED}[!] Error saving results: {str(e)}")

def main():
    print_banner()
    args = parse_args()
    
    # Normalize URL
    target_url = normalize_url(args.url)
    base_domain = extract_domain(target_url)
    
    print(f"{Fore.BLUE}[*] Target URL: {target_url}")
    print(f"{Fore.BLUE}[*] Base Domain: {base_domain}")
    
    keywords = []
    if args.wordlist:
        print(f"{Fore.BLUE}[*] Loading wordlist: {args.wordlist}")
        keywords = load_wordlist(args.wordlist)
        print(f"{Fore.GREEN}[+] Loaded {len(keywords)} keywords")
    
    # Crawl the website
    print(f"{Fore.BLUE}[*] Starting crawl with depth {args.depth}")
    crawl_website(target_url, 1, base_domain, args.depth)
    
    print(f"\n{Fore.GREEN}[+] Crawling complete")
    print(f"{Fore.GREEN}[+] Found {len(js_files)} JavaScript files")
    print(f"{Fore.GREEN}[+] Found {len(subdomains)} subdomains")
    
    if keywords:
        print(f"\n{Fore.BLUE}[*] Searching for keywords in JavaScript files")
        with concurrent.futures.ThreadPoolExecutor(max_workers=args.threads) as executor:
            futures = [executor.submit(search_keywords, js_url, keywords) for js_url in js_files]
            concurrent.futures.wait(futures)
    
    print(f"\n{Fore.GREEN}=== SUMMARY ===")
    print(f"{Fore.GREEN}[+] Total JavaScript files: {len(js_files)}")
    print(f"{Fore.GREEN}[+] Total subdomains: {len(subdomains)}")
    print(f"{Fore.GREEN}[+] Files with keyword matches: {len(found_keywords)}")
    
    if args.output:
        save_results(args.output)
    
    if subdomains:
        print(f"\n{Fore.CYAN}=== SUBDOMAINS ===")
        for subdomain in sorted(subdomains):
            print(f"{Fore.CYAN}[+] {subdomain}")
    
    if js_files:
        print(f"\n{Fore.GREEN}=== JAVASCRIPT FILES ===")
        for js_file in sorted(js_files):
            print(f"{Fore.GREEN}[+] {js_file}")
    
    if found_keywords:
        print(f"\n{Fore.YELLOW}=== KEYWORD MATCHES ===")
        for js_url, keywords in found_keywords.items():
            print(f"\n{Fore.YELLOW}File: {js_url}")
            for keyword, contexts in keywords.items():
                print(f"{Fore.YELLOW}  Keyword: {keyword}")
                for context in contexts[:3]: 
                    print(f"{Fore.YELLOW}    Context: ...{context}...")
                if len(contexts) > 3:
                    print(f"{Fore.YELLOW}    ... and {len(contexts) - 3} more matches")

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print(f"\n{Fore.RED}[!] Execution interrupted by user")
        sys.exit(0)
