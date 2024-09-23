# Raft Distributed Consensus

## Overview

This project implements the Raft consensus algorithm, which is used for managing a replicated log across a distributed system. Raft is designed to be understandable and is often used in systems requiring fault tolerance and consistency.


## Architecture

The architecture consists of the following components:

- **Node**: Represents a single instance in the distributed system.
- **Raft Protocol**: Implements the core functionalities of leader election and log management.
- **Client Interface**: Allows clients to interact with the Raft cluster.

## Getting Started

### Prerequisites

- Go (version 1.16 or higher)
- Redis (for storing logs and states)
- Any necessary dependencies listed in `go.mod`

### Installation

1. Clone the repository:
   git clone https://github.com/tnlmao/sarva.git
   cd your-repo-name
2. Install Dependencies
    go mod tidy
3. Start Application
    go run main.go