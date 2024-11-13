# dgos - distributed system monitoring platform made in Golang

# instructions 
- run central server `go run cmd/server/main.go -port 8080`
- run agent (local) `go run cmd/agent/agent.go -server http://<your-ip-here>:8080`
-   IP here must be the IP of central server machine
- run frontend (same device as central server)(`cd/client/ddgo-fe && npm run start`)

# to use over a network:
- launch central server on one machine (step #1)
- run frontend on same device as central server
- run agents on nodes (local devices) (with IP of central server)
  
# authors
- Ryan Ho
- Junkai Zheng
