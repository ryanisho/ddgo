function updateCharts() {
    fetch('/api/metrics')
        .then(response => response.json())
        .then(data => {
            // Process and update charts with the data
            updateCPUChart(data.filter(m => m.Name === 'cpu_usage'));
            // updateMemoryChart(data.filter(m => m.Name === 'memory_usage'));
            // updateDiskChart(data.filter(m => m.Name === 'disk_usage'));
        });
}

// Initialize charts and start updating
document.addEventListener('DOMContentLoaded', () => {
    // Initialize Chart.js charts here
    // Set up periodic updates
    setInterval(updateCharts, 30000);
    updateCharts();
});