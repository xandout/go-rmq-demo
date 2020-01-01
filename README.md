# Distributed Work via RabbitMQ

This project is a quick example of how to use a queue to distribute work amongst many, possibly remote workers in go.

## Prep

> 🚨 Disclaimer 🚨

> Majestic owns all rights to the data used in this example.  They provide this data for free, please be kind to their servers.

```
curl -o million.csv http://downloads.majestic.com/majestic_million.csv
```

## Docker Compose

[docker-compose.yml](docker-compose.yml) 


### Build

```
docker-compose build
```

### Run

Because of the way queues are created in RabbitMQ, you need to start the RMQ service & the consumers first.

```
docker-compose up -d rabbitmq us-east-1 us-east-2 us-west-2
```

This will start the RMQ server and 3 consumers as well as setup the `Exchange` called "domains" and registering a `Queue` for each "region".  

Once these containers are up, you can start the "producer" with 

```
docker-compose up producer
```

This container will parse the [Top 1,000,000 Domains](http://downloads.majestic.com/majestic_million.csv) from the CSV and will place them on the `domains` exchange in RMQ.  Each record is placed on one of the 3 "regional" workers, pseudorandomly.

```
producer_1   | time="2020-01-01T14:39:38Z" level=info msg="Adding osaka.jp to us-east-2"
producer_1   | time="2020-01-01T14:39:38Z" level=info msg="Adding yandex.ua to us-west-2"
producer_1   | time="2020-01-01T14:39:38Z" level=info msg="Adding uni-kassel.de to us-west-2"
producer_1   | time="2020-01-01T14:39:38Z" level=info msg="Adding dovepress.com to us-west-2"
producer_1   | time="2020-01-01T14:39:38Z" level=info msg="Adding uillinois.edu to us-west-2"
producer_1   | time="2020-01-01T14:39:38Z" level=info msg="Adding cbdoildiscount.net to us-west-2"
producer_1   | time="2020-01-01T14:39:38Z" level=info msg="Adding hermanmiller.com to us-east-1"
producer_1   | time="2020-01-01T14:39:38Z" level=info msg="Adding segodnya.ua to us-east-1"
producer_1   | time="2020-01-01T14:39:38Z" level=info msg="Adding freeuk.com to us-east-1"
producer_1   | time="2020-01-01T14:39:38Z" level=info msg="Adding crainsdetroit.com to us-east-2"
```

Messages are delivered to the exchange and placed on the correct queue based on the `Routing Key`.  The workers will process(print to stdout) each record and `ACK` back to RMQ.


### RMQ Admin UI
[http://localhost:15672]()
```
 UN: user
 PW: bitnami
```

## Teardown

```
docker-compose down -v
```

