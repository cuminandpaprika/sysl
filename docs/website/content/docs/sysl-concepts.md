---
title: "Sysl Concepts"
description: "Understanding the concepts behind Sysl"
date: 2018-02-28T10:11:18+11:00
weight: 40
draft: false
bref: ""
toc: false
---

# Overview

Sysl allows you to specify all the connections and dependencies around your applications.


## Applications
An application is an independent entity that provides services via its various `endpoints`.

Here is how an application is defined in sysl.
```
MobileApp:
  ...
```
`MobileApp` is a user-defined Application that does not have any endpoints yet. We will design this app as we move along.

## Endpoints

Endpoints are the services that an application offers. Let's add endpoints to our `MobileApp`.
```
MobileApp:
  Login: ...
  Search: ...
  Order: ...
```

## Data
You will have various kinds of data passing through your systems. Sysl allows you to express ownership, information classification and other attributes of your data in one place.

Continuing with the previous example, let's define a `Server` that expects `LoginData` for the `Login` Flow.
```
Server:
  Login (request <: Server.LoginData): ...

  !type LoginData:
    username <: string
    password <: string
```

## Statements

Statements describe the behaviour of applications.

```
MobileApp:
  Login:
    Server <- Login

Server:
  Login(data <: LoginData):
    build query
    DB <- Query
    check result
    return Server.LoginResponse

  !type LoginData:
    username <: string
    password <: string

  !type LoginResponse:
    message <: string

DB:
  Query:
    lookup data
    return data
  Save:
    ...
```


## Transforms

A transform is a mapping from one type of data to another. In Sysl, transforms define what each generated file contains and "transforms" sysl into code. Currently, only Go is supported but there are plans to extend the output code to Swift and Kotlin.


