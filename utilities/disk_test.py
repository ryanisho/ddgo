import os


def fill_disk(path):
    try:
        with open(path, 'wb') as f:
            while True:
                f.write(b'0' * 1024 * 1024)  # Write 1MB of data
    except OSError as e:
        print(f"Disk is full: {e}")


if __name__ == "__main__":
    fill_disk("./loc/disk.txt")
