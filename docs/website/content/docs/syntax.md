---
title: "Sysl Syntax"
description: "Understanding the Sysl language."
date: 2018-02-28T10:11:18+11:00
weight: 40
draft: false
bref: ""
toc: true
---

The Sysl syntax is indentation based.
The most basic Sysl file looks like this:

```
HelloWorld:
	...
```
and can be parsed with

	sysl textpb hello.sysl --out hello.textpb

Sysl Syntax Guide
	* Applications are defined on the top level, followed by `:`
	* All endpoints are indented. Use a `tab` or `spaces` to indent.
	* `<:` is used to define the arguments to `Login` endpoint.
	* `!type` is used to define a new data type `LoginData`.
	* Again, `...` is used to show we don't have enough details yet about each endpoint.
	* Attributes

Patterns
    A pattern is `~` followed by a word that means something to you. E.g. `[~tag]`.

Key-Value pair

    As the name suggests, you can associate some data with your application or an endpoint.
  ```
  Application [version="1.1"]:
  ```

Data Types
-------------
Sysl supports following data types out of the box.
  * int, int64, int32
  * float, decimal
  * string
  * bool
  * datetime, date
  * any
  * xml

Transforms Syntax
-------------

* `!view` 
* 	`View definition
* `let
	* Function Declaration

### Go Function Calls

Go function calls can also be used in transforms

* Contains
* Count
* Fields
* FindAllString
* HasPrefix
* HasSuffix
* Join
* LastIndex
* MatchString
* Replace
* Split
* Title
* ToLower
* ToTitle
* ToUpper
* Trim
* TrimLeft
* TrimPrefix
* TrimRight
* TirmSpace
* TrimSuffix
