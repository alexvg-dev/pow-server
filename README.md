# PoW-server

Requirements:
````
Design and implement “Word of Wisdom” tcp server.
* TCP server should be protected from DDOS attacks with the Proof of Work (https://en.wikipedia.org/wiki/Proof_of_work),
the challenge-response protocol should be used.
* The choice of the POW algorithm should be explained.
* After Proof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other
collection of the quotes.
* Docker file should be provided both for the server and for the client that solves the POW challenge
````

Server: 
1. Incoming requests limited by MAX_CONNECTIONS env
2. Connection time limited by SESSION_TTL_SEC env
3. Metrics server available

Client:
1. Client pkg extracted
2. Reply policy available if server temporary unavailable (backoff algo)

# PoW

- Challenge-response protocol implemented based on SCRYPT hashing algorithm.
- SCRYPT, compared to SHA-256 and Hashcash, is more ASIC-resistant due to memory-hard structure
  - you can increase memory consumption by increasing Scrypt params


## Commands:

### Dev
```
make lint 
make test
make local_run_server
make local_run_client

```

### 1. Start server
````
make build_server
make start_server
````

### 2. Start client
````
make build_client
make start_client
````
