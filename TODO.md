# TODO

This project should implement the following commands:

1. Pulse - command `pulse` schedules and sends posts to BlueSky on a set schedule,
  and logs/emits a heartbeat every 5-10s
2. Import - command `import` takes different types of inputs and stores them in a
  local SQLite database
3. Serve - command `serve` starts the monitoring server

## Objectives

Create a collection of concurrent web services built with the standard library.

Maintain minimal dependencies, and limit to [`<github.com/mattn/go-sqlite3>`](https://github.com/mattn/go-sqlite3)

MVP posts to BSky at least once a day.

## Database

- [ ] define post database model/table
  - post state
    - Drafted posts
    - Published posts - should have link & BSky ID
- [ ] define quote/note database model/table
  - "source" field
- [ ] create relationship between quote/note and post
- [ ] create relationship between task and post
- [ ] repository types
  - What should be abstract and what should be concrete?

## Import

- [ ] Define import command based on contents of [cli module](./cli.go)

## Pulse

- [ ] a "heartrate" command line flag to set the interval (in seconds)
- [ ] decide on and implement a format for the heartbeat method
- [ ] query the tasks table at _some_ level of precision

## System

- [ ] add caller information to log output
- [ ] self-hosting [research](TODO#Research)
  - [ ] containerize?

## Server/Dashboard

- [ ] Scaffold HTTP server
- [ ] Endpoints
  - [ ] Recently imported notes
  - [ ] Recently drafted posts
  - [ ] Health check
- [ ] templating & HTMX
  - [ ] vendor dependencies

## Research

What can SQLite Extensions do in terms of JSON storage?

How can I self-host a bot? Which machine is best served for it?

## Open Questions

What is the best way to check the health of the system?
