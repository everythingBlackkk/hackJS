import sys
import psutil
import winreg
import subprocess
import logging
import hashlib
import time
from datetime import datetime
from prettytable import PrettyTable
from watchdog.observers import Observer
from watchdog.events import FileSystemEventHandler

# Logging configuration
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(message)s')
logger = logging.getLogger(__name__)

class FileMonitorHandler(FileSystemEventHandler):
    """Monitor file system changes"""
    def on_modified(self, event):
        logger.info(f"File modified: {event.src_path}")

    def on_created(self, event):
        logger.info(f"File created: {event.src_path}")

    def on_deleted(self, event):
        logger.info(f"File deleted: {event.src_path}")

def validate_process_exists(target_process):
    """Check if the process exists"""
    proc = subprocess.Popen(['where', target_process], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    proc.communicate()
    return proc.returncode == 0

def get_process_details(pid):
    """Retrieve process details"""
    process = psutil.Process(pid)
    details = {
        "Process Name": process.name(),
        "Process ID": pid,
        "Executable Path": process.exe(),
        "Username": process.username(),
        "Started At": datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
        "CPU Usage (%)": process.cpu_percent(interval=1.0),
        "Memory Usage (MB)": process.memory_info().rss / (1024 * 1024),
        "Threads": process.num_threads(),
        "Open Files": len(process.open_files()),
        "Status": process.status()
    }
    return details

def calculate_file_hash(filepath):
    """Calculate SHA256 hash of the file"""
    sha256_hash = hashlib.sha256()
    with open(filepath, "rb") as f:
        for byte_block in iter(lambda: f.read(4096), b""):
            sha256_hash.update(byte_block)
    return sha256_hash.hexdigest()

def monitor_network_connections(pid):
    """Monitor network connections for the process"""
    process = psutil.Process(pid)
    connections = process.connections()
    
    if connections:
        print("\n=== Network Connections ===")
        table = PrettyTable()
        table.field_names = ["Local Port", "Remote IP"]
        table.align = "l"
        for conn in connections:
            remote_ip = conn.raddr.ip if conn.raddr else 'Not Connected'
            table.add_row([conn.laddr.port, remote_ip])
        print(table)

def monitor_registry_changes():
    """Monitor registry changes dynamically"""
    try:
        registry_hives = [
            winreg.HKEY_LOCAL_MACHINE,
            winreg.HKEY_CURRENT_USER,
            winreg.HKEY_CLASSES_ROOT,
            winreg.HKEY_USERS,
            winreg.HKEY_CURRENT_CONFIG
        ]

        for hive in registry_hives:
            key = winreg.ConnectRegistry(None, hive)
            reg_key = winreg.OpenKey(key, "", 0, winreg.KEY_NOTIFY)
            winreg.QueryInfoKey(reg_key)
            logger.info(f"Monitoring registry hive: {hive}")

    except Exception as e:
        logger.error(f"Error monitoring registry: {e}")

def monitor_file_system(path="C:\\"):
    """Monitor file system changes"""
    event_handler = FileMonitorHandler()
    observer = Observer()
    observer.schedule(event_handler, path, recursive=True)
    observer.start()
    logger.info(f"Monitoring file system changes in: {path}")
    return observer

def monitor_process(target_process):
    """Monitor the process for suspicious activity"""
    if not validate_process_exists(target_process):
        logger.error(f"Cannot find process: {target_process}")
        sys.exit(1)
    
    process = subprocess.Popen(target_process)
    pid = None
    for attempt in range(20):
        for proc in psutil.process_iter(['name']):
            if proc.info['name'] == target_process:
                pid = proc.pid
                break
        if pid:
            break
    
    if not pid:
        logger.error("Failed to detect process after launching")
        sys.exit(1)

    process_details = get_process_details(pid)
    print("\n=== Process Details ===")
    table = PrettyTable()
    table.field_names = ["Detail", "Value"]
    table.align = "l"
    for key, value in process_details.items():
        table.add_row([key, value])
    print(table)
    
    file_hash = calculate_file_hash(process_details["Executable Path"])
    print(f"\n=== File SHA256 Hash ===\n{file_hash}\n\n")

    file_observer = monitor_file_system()
    
    while True:
        if not psutil.pid_exists(pid):
            print("\n=== Process Terminated ===")
            file_observer.stop()
            file_observer.join()
            break
        monitor_network_connections(pid)
        monitor_registry_changes()
        time.sleep(5)

def main():
    if len(sys.argv) < 2:
        logger.error("Please specify the program name to monitor")
        sys.exit(1)
    
    target_process = sys.argv[1]
    monitor_process(target_process)

if __name__ == "__main__":
    main()
