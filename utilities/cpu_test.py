import multiprocessing
import os


def cpu_intensive_task():
    # An infinite loop that performs some CPU-intensive calculations
    x = 0
    while True:
        x = x ** 2 + x ** 0.5 - x * 0.3  # Random complex calculation to keep the CPU busy


if __name__ == "__main__":
    # Get the number of available CPU cores
    cpu_count = os.cpu_count()

    # Create a process for each CPU core
    processes = []
    for _ in range(cpu_count):
        p = multiprocessing.Process(target=cpu_intensive_task)
        processes.append(p)
        p.start()

    # Wait for all processes to complete (they never will, so you need to manually stop the program)
    for p in processes:
        p.join()
