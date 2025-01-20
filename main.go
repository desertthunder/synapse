// module Synapse takes quotes and notes and posts them to an account on BlueSky
// on a set schedule, via the BlueSky API (AT Protocol). Notifications are sent
// to a Discord Server.
//
// ---
//
// # Bot
//
// Login to BlueSky
// Store the token in the database -> Refresh the token every 15-30 minutes
//
// Child Routine
// Auth with Discord -> Store token -> Notify that server is starting
//
// Child Routine
// Check for tasks that need to be executed -> Execute task in chronological order
// Child SubRoutine
// Notify on Discord
//
// Child Routine
// Check for notes that need tasks -> Create tasks
//
// Exit -> Handle signal, notify on Discord
//
// ---
//
// # Importer
//
// Read file, create in memory objects/structs -> Create batches for importing
//
// Each batch is a go routine
// If the transaction fails, split the batch
//
// Upon completion of _every_ transaction, notify on Discord
//
// ---
//
// # Helpers
//
// Logger: default level (stdout) should be Debug
// A level >= Info goes to database
//
// Dashboard
// HTMX powered
// Reads logs table and updates a log in real time
package main

var logger = DefaultLogger()

func main() {
	Login()
	SetupDb(false)
}
