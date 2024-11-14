import React, { useEffect, useState } from "react";
import axios from "axios";

interface CoreTime {
  core: string;
  user: number;
  system: number;
  idle: number;
  iowait: number;
  irq: number;
}

const TableOne = () => {
  const [cpuTimes, setCpuTimes] = useState<CoreTime[]>([]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get("http://localhost:8080/api/metrics"); // Adjust the API endpoint as needed
        const responseData = response.data;
        const key = Object.keys(responseData)[0];
        const data = responseData[key]?.metrics?.cpu?.times;

        if (!data) {
          throw new Error("Invalid data structure");
        }

        // Format and sort the data by core number
        const formattedData: CoreTime[] = Object.keys(data)
          .map((coreKey) => {
            const coreNumber = parseInt(coreKey.replace("cpu", ""), 10) + 1; // Convert "cpu0" to "Core 1", etc.
            return {
              core: `Core ${coreNumber}`,
              user: data[coreKey].user,
              system: data[coreKey].system,
              idle: data[coreKey].idle,
              iowait: data[coreKey].iowait,
              irq: data[coreKey].irq,
            };
          })
          .sort((a, b) => parseInt(a.core.replace("Core ", "")) - parseInt(b.core.replace("Core ", ""))); // Sort by core number

        setCpuTimes(formattedData);
      } catch (error) {
        console.error("Error fetching CPU data:", error);
      }
    };

    fetchData();
  }, []);

  return (
    <div className="rounded-sm border border-stroke bg-white px-4 pb-2 pt-4 shadow-default dark:border-strokedark dark:bg-boxdark sm:px-6 xl:pb-1">
      <h4 className="mb-4 text-lg font-semibold text-black dark:text-white">
        CPU Core Times
      </h4>

      <div className="flex flex-col">
        <div className="grid grid-cols-5 bg-gray-100 dark:bg-meta-4 text-xs font-medium text-gray-700 dark:text-gray-300">
          <div className="p-2 text-center">Core</div>
          <div className="p-2 text-center">User (s)</div>
          <div className="p-2 text-center">System (s)</div>
          <div className="p-2 text-center">Idle (s)</div>
          <div className="p-2 text-center">IRQ + IOWait</div>
        </div>

        {cpuTimes.map((core, index) => (
          <div
            className={`grid grid-cols-5 items-center text-xs ${
              index === cpuTimes.length - 1
                ? ""
                : "border-b border-gray-200 dark:border-strokedark"
            }`}
            key={index}
          >
            <div className="p-2 text-center text-gray-800 dark:text-gray-100">{core.core}</div>
            <div className="p-2 text-center text-gray-800 dark:text-gray-100">{core.user.toFixed(2)}</div>
            <div className="p-2 text-center text-gray-800 dark:text-gray-100">{core.system.toFixed(2)}</div>
            <div className="p-2 text-center text-gray-800 dark:text-gray-100">{core.idle.toFixed(2)}</div>
            <div className="p-2 text-center text-gray-800 dark:text-gray-100">{(core.irq + core.iowait).toFixed(2)}</div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default TableOne;
