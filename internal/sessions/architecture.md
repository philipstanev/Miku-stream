# Overview

## General
- Session IDs are UUIDv4
- Session info contains user ID and role
- Sessions are accessed through interfaces 
- There can be many implementations but currently only in-memory sessions are implemented because the project is designed to run on one node

# Structure
- **abstract.go** - session interface and session data model
- **inMemorySession.go** - Implementation of sessions that uses a map and a mutex