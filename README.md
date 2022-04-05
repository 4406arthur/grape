# grape
ETH chain scraper


## indexer
used for block discover, will submit relative task into jobqueue

## scraper
consumer for jobqueue, concurrent query ETH Chian by RPC Spec, and batch insert into target RDBMS

## api-service
for client side  fast access chian data


## How to run 

1. run the docker-compose.yaml inside deployments (start local db)
2. setup config file (you could use default value)
3. go run main.go
