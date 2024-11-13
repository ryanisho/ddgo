import React, { useEffect, useState } from 'react';
import { Bar } from 'react-chartjs-2';
import axios from 'axios';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';

ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend
);

const ChartOne = () => {
  const [cpuData, setCpuData] = useState<number[]>([]);
  const [labels, setLabels] = useState<string[]>([]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get('http://localhost:8080/api/metrics');
        const responseData = response.data;
        const key = Object.keys(responseData)[0];
        const data = responseData[key]?.metrics?.cpu?.cores;

        if (!data) {
          throw new Error('Invalid data structure');
        }

        const usageData = data.map((core: any) => core.usage);
        // Adjust coreLabels to start from 1 instead of 0
        const coreLabels = data.map((core: any) => `Core ${core.core + 1}`);
        setCpuData(usageData);
        setLabels(coreLabels);
      } catch (error) {
        console.error('Error fetching CPU data:', error);
      }
    };

    // Initial fetch
    fetchData();

    // Set interval for real-time updates
    const intervalId = setInterval(fetchData, 5000); // Fetch data every 5 seconds

    // Clear interval on component unmount
    return () => clearInterval(intervalId);
  }, []);

  const chartData = {
    labels: labels,
    datasets: [
      {
        label: 'CPU Usage (%)',
        data: cpuData,
        backgroundColor: 'rgba(75,192,192,0.6)',
        borderColor: 'rgba(75,192,192,1)',
        borderWidth: 1,
      },
    ],
  };

  const options = {
    responsive: true,
    maintainAspectRatio: false,
    scales: {
      y: {
        beginAtZero: true,
        max: 100,
        title: {
          display: true,
          text: 'Usage (%)',
        },
      },
      x: {
        title: {
          display: true,
          text: 'CPU Cores',
        },
      },
    },
  };

  return (
    <div className="col-span-12 bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md border border-gray-200 dark:border-gray-700">
      <h2 className="text-2xl font-semibold text-gray-800 dark:text-gray-100 mb-4">CPU Usage</h2>
      <div className="w-full h-80">
        <Bar data={chartData} options={options} />
      </div>
    </div>
  );
};

export default ChartOne;
