# go-rest-template
This repo provides a simple golang rest template for developing microservices.
<br>
It implements a REST based tradebook service for demo purpose.
<br>
An order in tradebook consists of these information (buy or sell, quantity, market price or limit price). 

## Quickstart

0. pre-requirement
```
$ create a key on gcp cloud kms for this app and update SYSTEM_KEY_ID feild in .env
```

1. spin up local development environment
```
$ make compose-up
```

2. start service for local development
```
$ make run
```

3. (Optional) spin down local development environment
```
$ make compose-down
```

## Commit Requirements

1. Build Document
```
$ make doc
```

2. Build Mocks
```
$ make mocks
```

## API DOC
* Live Doc: http://localhost:8000/docs/index.html
* [postman file](./tradebook.postman_collection.json) is included for trying the api out

## Structure
```
api -> service -> store -> implementations (database, encryption...)
 |        |         |         
 |        V         |
 +----> models <----+
```
