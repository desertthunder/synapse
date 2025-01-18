package main

var logger = DefaultLogger()

func main() {
	w := NewWorker(3, 3, logger)
	w.Run()
}
