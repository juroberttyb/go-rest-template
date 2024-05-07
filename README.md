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

## package structure
```
api -> service -> store -> implementations (database, encryption...)
 |        |         |         
 |        V         |
 +----> models <----+
```

## todo
* Update to use websocket isntead of REST
