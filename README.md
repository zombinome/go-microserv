# go-microserv
Small microservice library I'm writing to learn go. Inspired by .net core approach.

## What is covered (at least for some extent):
* JSON & reflection
* http package

## TODO:
* Tests
* More examples & scenarios covered
* Cookies & sessions
* Authentication and authorization groundwork

Base idea is:
You have local host which can contain one or more services as hosted services. Additionally you can register one or more remote services with endpoints. If services share name - it is considered that you have same services deployed to different hosts (sharding).
Each service contains several endpoints.
Each endpoint is basically a function to be called.
You can register a set of middlewares to preform additional filtering/parsing/request modification before or after service endpoint is called.
You can call another service endpoint from current endpoint handler (if several instances of the same service is known, you can use custom service selector to select which service should actually be called).
