import React, { useState, useEffect } from 'react';
import { Box, Container, Grid, Heading, Select, Text } from '@chakra-ui/react';
import { BarChart, Bar, AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';

const formatBytes = (bytes) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / k ** i).toFixed(2)) + ' ' + sizes[i];
};

const MetricCard = ({ title, children }) => (
    <Box
        bg="white"
        p={4}
        borderRadius="lg"
        boxShadow="sm"
        border="1px"
        borderColor="gray.200"
    >
        <Text fontSize="sm" fontWeight="medium" mb={2} color="gray.500">
            {title}
        </Text>
        {children}
    </Box>
);

const CPUChart = ({ data }) => (
    <ResponsiveContainer width="100%" height={200}>
        <BarChart data={data.cores}>
            <CartesianGrid strokeDasharray="3 3" opacity={0.1} />
            <XAxis dataKey="core" label="Core" />
            <YAxis domain={[0, 100]} label="Usage %" />
            <Tooltip />
            <Bar
                dataKey="usage"
                fill="#3182ce"
                radius={[4, 4, 0, 0]}
            />
        </BarChart>
    </ResponsiveContainer>
);

const MemoryChart = ({ data }) => {
    const memoryData = {
        used: data.virtual.used,
        free: data.virtual.free,
        total: data.virtual.total,
        usage: data.virtual.usage
    };

    return (
        <Box>
            <Text fontSize="2xl" fontWeight="bold" mb={2}>
                {memoryData.usage.toFixed(1)}%
            </Text>
            <Box h="4" bg="gray.100" borderRadius="full" overflow="hidden">
                <Box
                    h="100%"
                    bg="blue.500"
                    borderRadius="full"
                    transition="width 0.3s ease"
                    width={`${memoryData.usage}%`}
                />
            </Box>
            <Grid templateColumns="repeat(2, 1fr)" gap={4} mt={4}>
                <Box>
                    <Text fontSize="sm" color="gray.500">Used</Text>
                    <Text fontWeight="medium">{formatBytes(memoryData.used)}</Text>
                </Box>
                <Box>
                    <Text fontSize="sm" color="gray.500">Total</Text>
                    <Text fontWeight="medium">{formatBytes(memoryData.total)}</Text>
                </Box>
            </Grid>
        </Box>
    );
};

const DiskIOChart = ({ data }) => {
    const [historicalData, setHistoricalData] = useState([]);

    useEffect(() => {
        setHistoricalData(prev => [...prev, {
            time: new Date().toLocaleTimeString(),
            read: data.io.read_bytes,
            write: data.io.write_bytes
        }].slice(-20));
    }, [data]);

    return (
        <ResponsiveContainer width="100%" height={200}>
            <AreaChart data={historicalData}>
                <CartesianGrid strokeDasharray="3 3" opacity={0.1} />
                <XAxis dataKey="time" />
                <YAxis />
                <Tooltip formatter={(value) => formatBytes(value)} />
                <Area
                    type="monotone"
                    dataKey="read"
                    stroke="#3182ce"
                    fill="#bee3f8"
                    stackId="1"
                    name="Read"
                />
                <Area
                    type="monotone"
                    dataKey="write"
                    stroke="#38a169"
                    fill="#c6f6d5"
                    stackId="1"
                    name="Write"
                />
            </AreaChart>
        </ResponsiveContainer>
    );
};

function App() {
    const [metrics, setMetrics] = useState(null);
    const [selectedAgent, setSelectedAgent] = useState('all');
    const [agents, setAgents] = useState([]);

    useEffect(() => {
        const fetchMetrics = async () => {
            try {
                const response = await fetch('http://localhost:8080/api/metrics');
                const data = await response.json();
                
                setAgents(Object.keys(data));
                if (selectedAgent === 'all') {
                    setMetrics(data);
                } else {
                    setMetrics({ [selectedAgent]: data[selectedAgent] });
                }
            } catch (error) {
                console.error('Error fetching metrics:', error);
            }
        };

        fetchMetrics();
        const interval = setInterval(fetchMetrics, 2000);
        return () => clearInterval(interval);
    }, [selectedAgent]);

    if (!metrics) {
        return <Text p={8}>Loading metrics...</Text>;
    }

    return (
        <Container maxW="container.xl" py={8}>
            <Box mb={8}>
                <Heading size="lg" mb={4}>System Metrics Dashboard</Heading>
                <Select
                    value={selectedAgent}
                    onChange={(e) => setSelectedAgent(e.target.value)}
                    maxW="300px"
                >
                    <option value="all">All Agents</option>
                    {agents.map(agentId => (
                        <option key={agentId} value={agentId}>
                            {metrics[agentId].hostname}
                        </option>
                    ))}
                </Select>
            </Box>

            {Object.entries(metrics).map(([agentId, agentMetrics]) => (
                <Box key={agentId} mb={8}>
                    <Heading size="md" mb={4}>{agentMetrics.hostname}</Heading>
                    <Grid templateColumns="repeat(12, 1fr)" gap={4}>
                        <Box gridColumn="span 8">
                            <MetricCard title="CPU Usage">
                                <CPUChart data={agentMetrics.metrics.cpu} />
                            </MetricCard>
                        </Box>

                        <Box gridColumn="span 4">
                            <MetricCard title="Memory Usage">
                                <MemoryChart data={agentMetrics.metrics.memory} />
                            </MetricCard>
                        </Box>

                        <Box gridColumn="span 6">
                            <MetricCard title="Load Average">
                                <Grid templateColumns="repeat(3, 1fr)" gap={4}>
                                    <Box>
                                        <Text fontSize="sm" color="gray.500">1m</Text>
                                        <Text fontSize="xl" fontWeight="bold">
                                            {agentMetrics.metrics.cpu.load["1m"].toFixed(2)}
                                        </Text>
                                    </Box>
                                    <Box>
                                        <Text fontSize="sm" color="gray.500">5m</Text>
                                        <Text fontSize="xl" fontWeight="bold">
                                            {agentMetrics.metrics.cpu.load["5m"].toFixed(2)}
                                        </Text>
                                    </Box>
                                    <Box>
                                        <Text fontSize="sm" color="gray.500">15m</Text>
                                        <Text fontSize="xl" fontWeight="bold">
                                            {agentMetrics.metrics.cpu.load["15m"].toFixed(2)}
                                        </Text>
                                    </Box>
                                </Grid>
                            </MetricCard>
                        </Box>

                        <Box gridColumn="span 6">
                            <MetricCard title="System Info">
                                <Grid templateColumns="repeat(2, 1fr)" gap={4}>
                                    <Box>
                                        <Text fontSize="sm" color="gray.500">Processes</Text>
                                        <Text fontSize="xl" fontWeight="bold">
                                            {agentMetrics.metrics.cpu.info.process_count}
                                        </Text>
                                    </Box>
                                    <Box>
                                        <Text fontSize="sm" color="gray.500">Threads</Text>
                                        <Text fontSize="xl" fontWeight="bold">
                                            {agentMetrics.metrics.cpu.info.thread_count}
                                        </Text>
                                    </Box>
                                </Grid>
                            </MetricCard>
                        </Box>

                        <Box gridColumn="span 12">
                            <MetricCard title="Disk I/O">
                                <DiskIOChart data={agentMetrics.metrics.disk} />
                            </MetricCard>
                        </Box>
                    </Grid>
                </Box>
            ))}
        </Container>
    );
}

export default App;