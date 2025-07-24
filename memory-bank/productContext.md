# Product Context

Describe the product.

## Overview

Provide a high-level overview of the project.

## Core Features

- Feature 1
- Feature 2

## Technical Stack

- Tech 1
- Tech 2

## Project Description

Ctrl-Alt-Play Agent is a lightweight game server management agent that communicates with the Ctrl-Alt-Play panel to manage Docker-based game servers on remote Linux systems. It functions similarly to the "Wings" system in Pelican Panel/Pterodactyl, providing seamless integration with the Panel system for efficient server management and real-time communication. Currently requires protocol updates to align with Panel's new command format.



Ctrl-Alt-Play Agent is a lightweight game server management agent that communicates with the Ctrl-Alt-Play panel to manage Docker-based game servers on remote Linux systems. It functions similarly to the "Wings" system in Pelican Panel/Pterodactyl, providing seamless integration with the Panel system for efficient server management and real-time communication.



## Architecture

Panel+Agent distributed architecture where the Agent runs on remote nodes to manage Docker containers via Docker API, communicating with the central Panel via WebSocket for real-time bidirectional messaging. Uses bearer token authentication and supports automatic reconnection with heartbeat mechanism. Recently updated with Issue #27 breaking changes requiring new command protocol format.



Panel+Agent distributed architecture where the Agent runs on remote nodes to manage Docker containers via Docker API, communicating with the central Panel via WebSocket for real-time bidirectional messaging. Uses bearer token authentication and supports automatic reconnection with heartbeat mechanism.



## Technologies

- Go 1.21+
- Docker Engine
- WebSocket
- Linux
- JSON
- Bearer Token Authentication
- Docker API
- Health Check HTTP Server



- Go
- Docker
- WebSocket
- Linux
- JSON
- Bearer Token Authentication



## Libraries and Dependencies

- github.com/gorilla/websocket
- github.com/docker/docker
- context
- net/http
- encoding/json
- time
- log
- os
- syscall
- runtime
- strconv
- sync
- net/url



- github.com/gorilla/websocket
- github.com/docker/docker
- context
- net/http
- encoding/json
- time
- log
- os
- syscall

