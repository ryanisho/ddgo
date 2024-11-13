# DGOS - Distributed System Monitoring Platform

DGOS is a distributed system monitoring platform built in Go (Golang) for efficient monitoring of various system metrics across multiple nodes. It provides a central server for monitoring and control, agents deployed on local nodes, and a web-based frontend for real-time monitoring. The project is useful for monitoring system health and performance in a distributed environment.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Installation](#installation)
- [Usage](#usage)
- [Authors](#authors)

## Overview

DGOS uses a central server to collect data from multiple agents running on different machines in a network. Agents send system metrics (such as CPU usage, memory usage, etc.) to the server, which can be visualized using a web-based frontend.

## Features

- **Centralized Monitoring**: One central server collects data from multiple agents running on distributed nodes.
- **Real-Time Metrics**: Live data visualization of system metrics such as CPU and memory usage.
- **User-Friendly Frontend**: A web-based UI built with Node.js and React to view the system status of all agents.
- **Simple Setup**: Lightweight agents in Go, deployable on any machine with minimal configuration.

## Architecture

The DGOS architecture consists of three primary components:

1. **Central Server**: Manages communication with agents, aggregates data, and serves data to the frontend.
2. **Agent**: Runs on each monitored machine, collects local metrics, and sends them to the central server.
3. **Frontend**: Web-based dashboard that displays real-time metrics from the central server, built with React.

## Installation

### Prerequisites

- **Go**: Ensure that Go is installed. You can download it [here](https://golang.org/dl/).
- **Node.js**: Required to run the frontend. Download from [here](https://nodejs.org/).

### Steps

1. **Clone the repository**:

   ```bash
   git clone https://github.com/yourusername/dgos.git
   cd dgos
   ```

2. **Install frontend dependencies**:
   ```bash
   cd client/ddgo-fe
   npm install
   ```

### Usage

1. **Start the central server**:

   ```bash
   go run cmd/server/main.go -port 8080
   ```

2. **Run the agent**:

   ```bash
   go run cmd/agent/agent.go -server http://<your-ip-here>:8080
   ```

3. **Launch the frontend**
   ```bash
   cd client/dggo-fe
   npm run start
   ```

**To launch DGOS over a network**:

- Launch the central server on one machine (Step 1).
- Run the frontend on the same device as the central server.
- Run agents on remote nodes (local devices) by providing the IP address of the central server.

### Authors

- Ryan Ho
