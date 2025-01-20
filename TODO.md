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

See [entry point](./main.go) for more thorough overview

## Research

What can SQLite Extensions do in terms of JSON storage?

How can I self-host a bot? Which machine is best served for it?

## Open Questions

What is the best way to check the health of the system?
