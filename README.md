# Web-Digger

## Getting started.

To make it easy for you to get to know with Web-digger, here's a some suggestions.

Web-Digger follows Hexagonal pattern (Port and Adapter) to be easily maintainable and scalable.
For testing phase, it has been considered not to use code quality tools like Linter, Conform or maybe Go-lines. The main reason of it is that reviewer can obtain a better vision about style and the hand write of developer.   


The main functionality of this application is placed in [service Layer](./internal/core/services) Where the Analyzer service receive data from HTTP adapter and forward it to maybe lower layers of source code like [infrastructure](./internal/core/infrastructure) layer. But if you want to get to know better with this source code, it is suggested to flow the cycle from [first place (main.go)](./main.go). This tour can help you to understand what happens deep into the architecture. 

## installation.
install Web-digger locally by running the following command:

```shell
# Please make .env ready.
cp .env.example .env

docker-compose -f docker-compose-local.yml up --build -d
# OR if your machine has Make installed:
make up
# OR you can run it on your own machine:
go run main.go
```

## Improvements that can be pointed at.

It would be better to cache responses according to URL addresses. It is apparent that processing and health checking modern huge web pages require noticeable amount of resources. for This case, personally, I think providing a driven adapter for using cache databases can be helpful.
Therefore, writing test cases in each layer as it is possible is always welcomed. For example, writing test cases for pkg layers can assure higher level of code coverage or other places with same situation too.