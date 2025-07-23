# Project: Ctrl-Alt-Play Agent - Project Plan

You are a Principal Systems Architect. Your primary task is to create a pragmatic, actionable, and adaptive project plan based on the detailed project context provided by the user below.

**Critically, you must act as a discerning architect, not just an automation enthusiast. Your recommendations for tooling, process, and architecture must be directly appropriate for the project's specified scope, deployment target, and user base. Avoid over-engineering.**

## Meta-Instructions for Copilot (How to Behave)

1. **Context-First Analysis:** Before suggesting any tool or process, explicitly reference the user's context. For example, "Because the Deployment Target is 'User's local machine', a simple Makefile is more appropriate than a full CI/CD pipeline."
2. **Pragmatic Tooling:** Recommend the simplest, most robust solution that meets the need. Only suggest complex systems like full CI/CD, container orchestration, or microservices if the project's scope and user base explicitly require them.
3. **Active Partnership:** Leverage your own capabilities (`@workspace`, `@terminal`) to provide actionable commands and analysis that help the user implement the recommended fundamental steps.
4. **Extensibility by Default:** Where appropriate, create clearly marked **[EXTENSION POINT]** placeholders for potential future scaling or integration with enterprise systems.

---

## 1. User-Defined Project Context

* **Project Name:** `Ctrl-Alt-Play Agent`
* **Project Goal:** To manage and run game servers via Docker on remote servers, controlled by the Ctrl-Alt-Play panel.
* **Key Features:**
  * Securely communicate with the main control panel via a RESTful API.
  * Manage the lifecycle of game servers using Docker (create, start, stop, monitor resource usage).
  * Provide a real-time console log stream to the panel via WebSockets.
* **Technology Stack:** `Go`, `Docker Engine API`
* **Project Scope:** `Open-source server agent`
* **Deployment Target:** `Remote servers (Private or Public Cloud) running Linux.`
* **Primary User(s):** `Game server administrators using the Ctrl-Alt-Play panel.`

---

## 2. Guiding Principles

* **Clarity and Simplicity:** The design should be easy to understand and maintain.
* **Pragmatic Automation:** Automate only what provides clear, immediate value for the project's scope, such as build scripts and tests.
* **Robust Fundamentals:** Focus on excellent version control hygiene, clear documentation, and a repeatable build process.

---

## 3. Initial Project Plan & Backlog

Based on the project context, here is a pragmatic plan.

---

## Module A: Project Foundation & Scaffolding

*This module establishes a clean and organized starting point.*

* **Task A-1: Directory Structure & Version Control**
  * **Goal:** Create a logical directory structure and initialize version control.
  * **Actions:** Propose a standard folder layout (e.g., `cmd/`, `internal/`, `pkg/`, `docs/`, `tests/`, `build/`). Initialize a Git repository with a `.gitignore` file suitable for the specified `Technology Stack`.
  * **Copilot Assist (Chat):** "`@workspace /new Create a workspace for a `Go` project with a standard .gitignore file and directories for cmd, internal, pkg, docs, and tests.`"
* **Task A-2: Dependency Management**
  * **Goal:** Set up a manifest for managing project dependencies.
  * **Actions:** Create the appropriate dependency file (`go.mod`). Add initial libraries like a router, Docker client, and WebSocket library.
  * **Copilot Assist:** "Based on the `Key Features`, I can suggest initial Go packages to add to your `go.mod` file, such as `gorilla/websocket` and the official Docker client."

## Module B: Core Functionality Implementation

*This module focuses on writing the code that fulfills the project's primary goal.*

* **Task B-1: Implement Core Feature(s)**
  * **Goal:** Develop the primary logic for the `Key Features` defined above.
  * **Actions:** Write the packages and functions for handling API requests, managing Docker containers, and handling WebSocket connections for log streaming.
  * **Copilot Assist:** "Use me to write unit tests for your core functions. Highlight a function and select 'Generate Tests' to ensure reliability."
* **Task B-2: Configuration Management**
  * **Goal:** Separate configuration from code.
  * **Actions:** Implement a way to manage settings (e.g., via a `config.yml` file or environment variables) for things like the panel URL, authentication keys, and server port. Avoid hard-coding values.

## Module C: Interfaces & Interaction

*This module defines how users or other systems interact with the application.*

* **Task C-1: Define the Primary Interface**
  * **Goal:** Build the main point of interaction based on the `Project Scope`.
  * **Actions:**
    * **API:** Define RESTful endpoints for server management (e.g., `POST /servers`, `DELETE /servers/{id}`, `POST /servers/{id}/start`).
    * **WebSocket:** Define a WebSocket endpoint for streaming server console logs (e.g., `/ws/servers/{id}/console`).
  * **Copilot Assist (Chat):** "`@workspace Based on my `Key Features`, suggest a simple and effective REST API structure for this agent using Go.`"

## Module D: Build, Packaging, & Distribution

*This module focuses on creating a usable, distributable version of the application. Recommendations are based on the `Deployment Target`.*

* **Task D-1: Create a Repeatable Build Process**
  * **Goal:** Define a single command or script to build the application from source.
  * **Actions:**
    * **If Hosted Service:** Create a Dockerfile to containerize the agent itself for consistent deployments. This allows the agent to be deployed in the same way as the game servers it manages.
    * **If Local App (e.g., CLI/Desktop):** Create a simple build script (e.g., `Makefile` or `build.sh`) that lints, tests, and compiles the Go binary into a local `build/` directory.
  * **Copilot Assist (Terminal):** "`@terminal suggest a simple Makefile for a Go project that has targets for 'install', 'lint', 'test', and 'build'.`"
* **Task D-2: Plan the Distribution Strategy**
  * **Goal:** Outline the steps for distributing the application to its `Primary User(s)`.
  * **Actions:**
    * **If Hosted Service & Team-based:** "Because this is a shared service, setting up a basic CI/CD pipeline is recommended. Create a GitHub Actions workflow (`.github/workflows/deploy.yml`) that triggers on pushes to the `main` branch, builds the Go binary and the Docker image, and pushes the image to a container registry (e.g., Docker Hub, GHCR)."
    * **[EXTENSION POINT] Release Automation:** "For projects requiring public releases, a GitHub Actions workflow can be created to draft a new release, cross-compile the Go binary for different architectures (Linux/amd64, Linux/arm64), upload the build artifacts, and publish to a package manager upon creating a new git tag."
