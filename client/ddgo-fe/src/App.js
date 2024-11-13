import React, { useState, useEffect } from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Cell } from 'recharts';
import { AreaChart, Area } from 'recharts';
import { formatBytes } from './utils/formatters';

const MetricCard = ({ title, children, className = '' }) => (
    <div className={`bg-white rounded-lg shadow-sm p-4 ${className}`}>
        <h3 className="text-sm font-medium text-gray-700 mb-3">{title}</h3>
        {children}
    </div>
);

const DetailRow = ({ label, value, unit = '' }) => (
    <div className="flex justify-between items-center py-1 text-sm">
        <span className="text-gray-600">{label}</span>
        <span className="font-medium">{value}{unit}</span>
    </div>
);

function CPUCoresChart({ cores }) {
    return (
        <div className="h-[200px]">
            <ResponsiveContainer width="100%" height="100%">
                <BarChart data={cores} margin={{ top: 5, right: 5, left: 5, bottom: 5 }}>
                    <CartesianGrid strokeDasharray="3 3" opacity={0.1} />
                    <XAxis dataKey="core" />
                    <YAxis domain={[0, 100]} />
                    <Tooltip
                        formatter={(value) => value.toFixed(2) + '%'}
                        labelFormatter={(label) => `Core ${label}`}
                    />
                    <Bar dataKey="usage" fill="#3B82F6">
                        {cores.map((entry, index) => (
                            <Cell
                                key={`cell-${index}`}
                                fill={entry.usage > 80 ? '#EF4444' : entry.usage > 60 ? '#F59E0B' : '#3B82F6'}
                            />
                        ))}
                    </Bar>
                </BarChart>
            </ResponsiveContainer>
        </div>
    );
}

function App() {
    const [metrics, setMetrics] = useState(null);
    const [historicalIO, setHistoricalIO] = useState([]);
    const [selectedCoreDetails, setSelectedCoreDetails] = useState(null);

    useEffect(() => {
        const fetchMetrics = async () => {
            try {
                const response = await fetch('http://localhost:8080/api/metrics');
                const data = await response.json();
                setMetrics(data);
                
                setHistoricalIO(prev => [...prev, {
                    time: new Date().toLocaleTimeString(),
                    read: data.disk.io.read_bytes,
                    write: data.disk.io.write_bytes
                }].slice(-20));
            } catch (error) {
                console.error('Error fetching metrics:', error);
            }
        };

        fetchMetrics();
        const interval = setInterval(fetchMetrics, 2000);
        return () => clearInterval(interval);
    }, []);

    if (!metrics) return <div className="p-4">Loading metrics...</div>;

    return (
        <div className="min-h-screen bg-gray-50 p-6">
            <div className="max-w-[1600px] mx-auto">
                {/* Header */}
                <div className="flex justify-between items-center mb-6">
                    <h1 className="text-2xl font-bold text-gray-900">DGoS</h1>
                    <div className="text-sm text-gray-500">
                        Last updated: {new Date(metrics.time).toLocaleTimeString()}
                    </div>
                </div>

                {/* Main Grid */}
                <div className="grid grid-cols-12 gap-4">
                    {/* CPU Overview */}
                    <div className="col-span-12 lg:col-span-8">
                        <MetricCard title="CPU Usage by Core">
                            <CPUCoresChart cores={metrics.cpu.cores} />
                        </MetricCard>
                    </div>

                    {/* CPU Info */}
                    <div className="col-span-12 lg:col-span-4">
                        <MetricCard title="CPU Information">
                            <DetailRow label="Physical Cores" value={metrics.cpu.info.physical_cores} />
                            <DetailRow label="Logical Cores" value={metrics.cpu.info.logical_cores} />
                            <DetailRow label="Total Processes" value={metrics.cpu.info.process_count} />
                            <DetailRow label="Total Threads" value={metrics.cpu.info.thread_count} />
                            <div className="mt-2 border-t pt-2">
                                <DetailRow label="Load Average (1m)" value={metrics.cpu.load["1m"].toFixed(2)} />
                                <DetailRow label="Load Average (5m)" value={metrics.cpu.load["5m"].toFixed(2)} />
                                <DetailRow label="Load Average (15m)" value={metrics.cpu.load["15m"].toFixed(2)} />
                            </div>
                        </MetricCard>
                    </div>

                    {/* CPU Times */}
                    <div className="col-span-12">
                        <MetricCard title="CPU Times">
                            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
                                {Object.entries(metrics.cpu.times)
                                    // Sort CPUs by numeric index extracted from the "cpu" key
                                    .sort(([a], [b]) => {
                                        const numA = parseInt(a.replace('cpu', ''));
                                        const numB = parseInt(b.replace('cpu', ''));
                                        return numA - numB;
                                    })
                                    .map(([cpu, times]) => (
                                        <div key={cpu} className="bg-gray-50 p-3 rounded-lg">
                                            <h4 className="font-medium mb-2">CPU {cpu.replace('cpu', '')}</h4>
                                            <DetailRow label="User" value={times.user.toFixed(2)} unit="s" />
                                            <DetailRow label="System" value={times.system.toFixed(2)} unit="s" />
                                            <DetailRow label="Idle" value={times.idle.toFixed(2)} unit="s" />
                                            <DetailRow label="I/O Wait" value={times.iowait.toFixed(2)} unit="s" />
                                            <DetailRow label="IRQ" value={times.irq.toFixed(2)} unit="s" />
                                        </div>
                                    ))}
                            </div>
                        </MetricCard>
                    </div>

                    {/* Memory Stats */}
                    <div className="col-span-12 lg:col-span-6">
                        <MetricCard title="Memory Usage">
                            <div className="mb-4">
                                <div className="h-2 bg-gray-200 rounded-full">
                                    <div 
                                        className="h-2 bg-blue-500 rounded-full transition-all duration-300"
                                        style={{ width: `${metrics.memory.virtual.usage}%` }}
                                    />
                                </div>
                                <div className="mt-1 text-sm text-gray-600 text-right">
                                    {metrics.memory.virtual.usage.toFixed(2)}%
                                </div>
                            </div>
                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <DetailRow 
                                        label="Total Memory" 
                                        value={formatBytes(metrics.memory.virtual.total)} 
                                    />
                                    <DetailRow 
                                        label="Used Memory" 
                                        value={formatBytes(metrics.memory.virtual.used)} 
                                    />
                                    <DetailRow 
                                        label="Free Memory" 
                                        value={formatBytes(metrics.memory.virtual.free)} 
                                    />
                                </div>
                                <div>
                                    <DetailRow 
                                        label="Cached" 
                                        value={formatBytes(metrics.memory.virtual.cached)} 
                                    />
                                    <DetailRow 
                                        label="Available" 
                                        value={formatBytes(metrics.memory.virtual.available)} 
                                    />
                                    <DetailRow 
                                        label="Usage" 
                                        value={`${metrics.memory.virtual.usage.toFixed(2)}%`} 
                                    />
                                </div>
                            </div>
                        </MetricCard>
                    </div>

                    {/* Swap Memory */}
                    <div className="col-span-12 lg:col-span-6">
                        <MetricCard title="Swap Memory">
                            <div className="grid grid-cols-2 gap-4">
                                <DetailRow 
                                    label="Total Swap" 
                                    value={formatBytes(metrics.memory.swap.total)} 
                                />
                                <DetailRow 
                                    label="Used Swap" 
                                    value={formatBytes(metrics.memory.swap.used)} 
                                />
                                <DetailRow 
                                    label="Free Swap" 
                                    value={formatBytes(metrics.memory.swap.free)} 
                                />
                                <DetailRow 
                                    label="Swap Usage" 
                                    value={`${metrics.memory.swap.usage.toFixed(2)}%`} 
                                />
                            </div>
                        </MetricCard>
                    </div>

                    {/* Disk Usage */}
                    <div className="col-span-12 lg:col-span-6">
                        <MetricCard title="Disk Usage">
                            <div className="mb-4">
                                <div className="h-2 bg-gray-200 rounded-full">
                                    <div 
                                        className="h-2 bg-green-500 rounded-full transition-all duration-300"
                                        style={{ width: `${metrics.disk.usage}%` }}
                                    />
                                </div>
                                <div className="mt-1 text-sm text-gray-600 text-right">
                                    {metrics.disk.usage.toFixed(2)}%
                                </div>
                            </div>
                            <div className="grid grid-cols-2 gap-4">
                                <DetailRow 
                                    label="Total Space" 
                                    value={formatBytes(metrics.disk.total)} 
                                />
                                <DetailRow 
                                    label="Free Space" 
                                    value={formatBytes(metrics.disk.free)} 
                                />
                            </div>
                        </MetricCard>
                    </div>

                    {/* Disk I/O */}
                    <div className="col-span-12 lg:col-span-6">
                        <MetricCard title="Disk I/O Statistics">
                            <div className="grid grid-cols-2 gap-4">
                                <DetailRow 
                                    label="Read Operations" 
                                    value={metrics.disk.io.read_count.toLocaleString()} 
                                />
                                <DetailRow 
                                    label="Write Operations" 
                                    value={metrics.disk.io.write_count.toLocaleString()} 
                                />
                                <DetailRow 
                                    label="Bytes Read" 
                                    value={formatBytes(metrics.disk.io.read_bytes)} 
                                />
                                <DetailRow 
                                    label="Bytes Written" 
                                    value={formatBytes(metrics.disk.io.write_bytes)} 
                                />
                            </div>
                            <div className="mt-4 h-[150px]">
                                <ResponsiveContainer width="100%" height="100%">
                                    <AreaChart data={historicalIO}>
                                        <CartesianGrid strokeDasharray="3 3" opacity={0.1} />
                                        <XAxis dataKey="time" />
                                        <YAxis />
                                        <Tooltip formatter={(value) => formatBytes(value)} />
                                        <Area 
                                            type="monotone" 
                                            dataKey="read" 
                                            stroke="#3B82F6" 
                                            fill="#93C5FD" 
                                            name="Read"
                                        />
                                        <Area 
                                            type="monotone" 
                                            dataKey="write" 
                                            stroke="#10B981" 
                                            fill="#6EE7B7" 
                                            name="Write"
                                        />
                                    </AreaChart>
                                </ResponsiveContainer>
                            </div>
                        </MetricCard>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default App;