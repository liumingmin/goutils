![gopher](gopher.png "goutils")


[![Go](https://github.com/liumingmin/goutils/actions/workflows/go.yml/badge.svg)](https://github.com/liumingmin/goutils/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/liumingmin/goutils)](https://goreportcard.com/report/github.com/liumingmin/goutils)
[![codecov](https://codecov.io/gh/liumingmin/goutils/graph/badge.svg?token=BQRDOY3CDX)](https://codecov.io/gh/liumingmin/goutils)
![GitHub last commit](https://img.shields.io/github/last-commit/liumingmin/goutils)
![GitHub Tag](https://img.shields.io/github/v/tag/liumingmin/goutils)
![GitHub License](https://img.shields.io/github/license/liumingmin/goutils)

**Read this in other languages: [English](README.md), [中文](README_zh.md).**

# Introduction

This repository aims to provide Golang developers with a series of practical and easy-to-use tools to help them improve development efficiency and work efficiency. These tools cover various fields, including algorithm libraries, container libraries, storage tools, file processing, Http and Websocket networks, NoSql database access, etc.

# Features
- Solution positioning: Based on solving tedious tasks outside the microservice framework, it forms a differentiated complement with the development framework.
- Easy to integrate: low connection, can be integrated with various projects, such as go-zero.
- Function focus: The tool library theme implements tools for various common functions. In order to avoid reinventing the wheel without developing your own SQL ORM, make more use of the powerful ORM solutions of the open source community. Only by standing on the shoulders of giants can you go further.

# Moudles

- [Algorithm Module](algorithm/README.md)
    - Double circle buffer
    - Crc16
    - Descart combination
    - Xor reader and writer
- [Cache Module](cache/README.md)
    - Generics-based function caching
- [YAML Configuration Module](conf/README.md)
- [Container Module](container/README.md)
    - Bitmap
    - Buffer pool
    - Consistent hashing
    - Generics-based sync pool
    - Generics-based common lock wrapper
    - Memory db struct
    - Generics-based Queue
    - Red Black Tree
- [Database Module](db/README.md)
    - Elasticsearch
    - Kafka
    - Mongo
    - Redis
- [Logging Library](log/README.md)
    - Zap wrapper
- [Network Library](net/README.md)
    - Support Http1.x and 2.0 HttpClient
    - Support Http1.x and 2.0 HttpServer    
    - Ip utils
    - Binary Net Packet Protocol
    - Ssh proxy client
- [General Utility Library](utils/README.md)
    - CircuitBreaker
    - Checksum utils
    - Type convert utils
    - Csv and MDB DataTable reader and writer
    - Distributed lock
    - Use gotest generate markdown document utils
    - Finite state machine
    - Http utils
    - Email utils
    - Safe goroutine
    - Snowflake id generater
    - Support timeout synchronous multi-call 
    - Window dll invoke
    - UTF-8 encoding convert
    - File utils
    - Math utils
    - Reflect utils
    - String parser
    - String utils
    - Struct utils
    - Struct tags utils
- [WebSocket Client and Server Library](ws/README.md)
    - Go Websocket client and server(100,000 concurrent 2.3G memory usage)
    - Cpp Websocket client     
    - Ts Websocket client 
    - Js Websocket client

# Roadmap

## Objective:

* Build a high-quality, high-performance Golang project
* Enhance project visibility and influence

## Phase Breakdown:

### Phase 1: Foundation Establishment

* **Goal:** Ensure code quality and establish a solid foundation
* **Key Metrics:**
    - [x] All code is tested with `go test` and `codecover` 
    - [x] GitHub Star exceeds 200
* **Tasks:**
    - [x] Write unit tests and integration tests to ensure code functionality
    - [x] Utilize the `codecover` tool to measure code coverage and continuously improve it
    - [x] Compose clear project documentation and README files

### Phase 2: Performance Optimization and Expansion

* **Goal:** Elevate project performance and scalability
* **Key Metrics:**
    * [x] Major code undergoes performance testing and meets performance expectations
    * [x] Code coverage for all code surpasses 80%
    * GitHub Star exceeds 500
* **Tasks:**
    * [x] Identify performance bottlenecks and implement targeted optimizations
    * [x] Employ load testing tools to evaluate project performance and make continuous improvements
    * [x] Refine code structure and design to enhance code maintainability and extensibility
    * [x] Release new versions and document version change logs

### Phase 3: Community Operation and Promotion

* **Goal:** Grow the project community and expand project influence
* **Key Metrics:**
    * GitHub Star exceeds 1000
    * Establish an independent official website
    * Actively participate in open-source community events
* **Tasks:**
    * Set up community communication platforms, such as forums or Discord
    * Organize online and offline technical exchange events to share project experiences
    * Write blog posts, technical tutorials, etc., to disseminate project knowledge
    * Actively participate in relevant open-source conferences and events to promote the project