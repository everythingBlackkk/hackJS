
```
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
                                                                                                                              
```

# hackJS

**hackJS** is a Python-based tool that extracts JavaScript file links from a website, searches for endpoints inside them, and optionally checks for specific keywords.

## Features

- Extracts all JavaScript files from a given URL.
- Searches JS files for hidden endpoints (e.g., API paths).
- Optional keyword matching inside JavaScript files.
- Output results to console or file.

## Usage

```bash
python3 hackjs.py -u <target_url> [-o <output_file>] [-w <keywords_file>]
````

### Parameters

* `-u` / `--url` (required): Target website URL.
* `-o` / `--output` (optional): Save results to a file.
* `-w` / `--wordlist` (optional): File containing keywords to search inside JS files.

### Example

```bash
python3 hackjs.py -u https://example.com -o result.txt -w WordList.txt
```

## Requirements

Install dependencies with:

```bash
pip install -r requirements.txt
```

## Author

Developed by [Yassin Mohamed](https://github.com/everythingBlackkk)

