# tradebook
This repo implements a trade engine that accepts orders via the protocol (or triggering) you defined.
<br>
An order request at least has this information (buy or sell, quantity, market price or limit price). 


## quickstart
1. spin up local development environment
```
$ make local-dev-up
```

2. start service for local development
```
$ make run
```

3. (Optional) spin down local development environment
```
$ make local-dev-down
```

## Resouce
* [postman file](./tradebook.postman_collection.json) is included for trying the api out

## API DOC
* Local Dev: http://localhost:8000/docs/index.html

## package structure
```
api -> service -> store -> implementations (database, encryption...)
 |        |         |         
 |        V         |
 +----> models <----+
```

## FIXME
* mutex should not only be on codebase level, but also should be on kubernete level to ensure no multitple pods modifying the orders at the same time

## todo
* Update to use websocket instead of REST
* Logging library should consider whether this is performance critical to change log format
