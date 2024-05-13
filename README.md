# tradebook
This repo implements a trade engine that accepts orders via the REST protocol.
<br>
An order request consists of these information (buy or sell, quantity, market price or limit price). 


## quickstart

0. pre-requirement
```
$ create a key on gcp cloud kms for this app and update SYSTEM_KEY_ID feild in .env
```

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
* add user api group

## todo
* Update to use websocket instead of REST
* Logging library should consider whether this is performance critical to change log format
