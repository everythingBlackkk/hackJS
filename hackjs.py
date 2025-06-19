import argparse
import requests
import re
from bs4 import BeautifulSoup
from urllib.parse import urljoin
import sys

def print_banner():
    banner = r"""  
Hack JavaScript                

         _       _    _                    _             _                  _        _        
        / /\    / /\ / /\                /\ \           /\_\               /\ \     / /\      
       / / /   / / // /  \              /  \ \         / / /  _            \ \ \   / /  \     
      / /_/   / / // / /\ \            / /\ \ \       / / /  /\_\          /\ \_\ / / /\ \__  
     / /\ \__/ / // / /\ \ \          / / /\ \ \     / / /__/ / /         / /\/_// / /\ \___\ 
    / /\ \___\/ // / /  \ \ \        / / /  \ \_\   / /\_____/ / _       / / /   \ \ \ \/___/ 
   / / /\/___/ // / /___/ /\ \      / / /    \/_/  / /\_______/ /\ \    / / /     \ \ \       
  / / /   / / // / /_____/ /\ \    / / /          / / /\ \ \    \ \_\  / / /  _    \ \ \      
 / / /   / / // /_________/\ \ \  / / /________  / / /  \ \ \   / / /_/ / /  /_/\__/ / /      
/ / /   / / // / /_       __\ \_\/ / /_________\/ / /    \ \ \ / / /__\/ /   \ \/___/ /       
\/_/    \/_/ \_\___\     /____/_/\/____________/\/_/      \_\_\\/_______/     \_____\/        
                                                                                              
                                    @Github : everthingBlackkk
                                      #Dev : Yassin Mohamed
    """
    print(banner)

HEADERS = {
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 "
                  "(KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
    "Accept": "*/*",
    "Accept-Language": "en-US,en;q=0.9",
    "Connection": "keep-alive",
}

ENDPOINT_REGEX = re.compile(
    r"""(?:"|')                                  
    (
        (?:
            (?:\/[a-zA-Z0-9_\-\/\.\?\=\&\%\#]+)  
            |
            (?:https?:\/\/[^\s"'<>]+)            
        )
    )
    (?:"|')                                      
    """, re.VERBOSE
)

def is_valid_endpoint(ep):
    ep = ep.strip()
    if not ep or ep in ('/', '//'):
        return False
    if len(ep) < 3:
        return False
    if all(c in '/:.' for c in ep):
        return False
    return True

def fetch_html(url):
    try:
        response = requests.get(url, headers=HEADERS, timeout=10)
        response.raise_for_status()
        return response.text
    except Exception as e:
        print(f"[ERROR] Failed to fetch {url}: {e}")
        return None

def extract_js_files(base_url, html):
    soup = BeautifulSoup(html, "html.parser")
    js_files = set()
    for script in soup.find_all("script"):
        src = script.get("src")
        if src:
            full_url = urljoin(base_url, src)
            js_files.add(full_url)
    return sorted(js_files)

def extract_endpoints(js_url):
    endpoints = set()
    content = fetch_html(js_url)
    if content:
        matches = ENDPOINT_REGEX.findall(content)
        for match in matches:
            cleaned = match.strip()
            if is_valid_endpoint(cleaned):
                endpoints.add(cleaned)
    return sorted(endpoints)

def search_keywords_in_js(js_url, keywords):
    found = []
    content = fetch_html(js_url)
    if content:
        for word in keywords:
            if word.lower() in content.lower():
                found.append(word)
    return found

def load_keywords(filepath):
    try:
        with open(filepath, 'r') as f:
            return [line.strip() for line in f if line.strip()]
    except Exception as e:
        print(f"[ERROR] Failed to read keywords file: {e}")
        return []

def output_line(text, file=None):
    print(text)
    if file:
        file.write(text + '\n')

def main():
    print_banner()
    parser = argparse.ArgumentParser(description="Extract JS files and endpoints from a website")
    parser.add_argument("-u", "--url", required=True, help="Target website URL")
    parser.add_argument("-o", "--output", help="Optional output file to save results")
    parser.add_argument("-w", "--wordlist", help="Keyword list file to search inside JS files")
    args = parser.parse_args()

    base_url = args.url.strip()
    output_path = args.output
    wordlist_path = args.wordlist
    out_file = open(output_path, "w") if output_path else None

    keywords = []
    if wordlist_path:
        output_line(f"[*] Loading keywords from: {wordlist_path}", out_file)
        keywords = load_keywords(wordlist_path)
        if not keywords:
            output_line("[!] No valid keywords loaded.", out_file)

    output_line("[*] Fetching website HTML...", out_file)
    html = fetch_html(base_url)
    if not html:
        output_line("[!] Failed to fetch base URL.", out_file)
        if out_file: out_file.close()
        return

    output_line("[*] Extracting JavaScript files...", out_file)
    js_files = extract_js_files(base_url, html)
    if not js_files:
        output_line("[!] No JavaScript files found.", out_file)
    else:
        output_line(f"[+] Found {len(js_files)} JavaScript file(s):", out_file)
        for i, js in enumerate(js_files, 1):
            output_line(f"    [{i}] {js}", out_file)

    all_endpoints = set()
    output_line("\n[*] Searching for endpoints and keywords in JS files...\n", out_file)

    for js_index, js in enumerate(js_files, 1):
        output_line(f"[{js_index}] Checking: {js}", out_file)
        endpoints = extract_endpoints(js)
        if endpoints:
            output_line(f"    Found {len(endpoints)} endpoint(s):", out_file)
            for idx, ep in enumerate(endpoints, 1):
                output_line(f"        [{idx}] {ep}", out_file)
                all_endpoints.add(ep)
        else:
            output_line("    No endpoints found.", out_file)

        if keywords:
            found_words = search_keywords_in_js(js, keywords)
            if found_words:
                output_line(f"    Found keyword(s): {', '.join(found_words)}", out_file)

        output_line("-" * 50, out_file)

    if not all_endpoints:
        output_line("\n[!] No valid endpoints found in any JavaScript files.", out_file)
    else:
        output_line(f"\n[✓] Total unique valid endpoints found: {len(all_endpoints)}", out_file)

    if out_file:
        out_file.close()
        output_line(f"\n[✓] Results saved to: {output_path}")

if __name__ == "__main__":
    main()
