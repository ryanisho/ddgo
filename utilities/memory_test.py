import time


def consume_memory():
    memory_hog = []
    try:
        while True:
            # Append a large string to the list to consume memory
            memory_hog.append(' ' * 10**9)  # 1MB per iteration
            # time.sleep(0.1)  # Slight delay to slow down the consumption rate
    except MemoryError:
        print("Memory is full!")


if __name__ == "__main__":
    consume_memory()
