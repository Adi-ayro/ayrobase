package main

func main() {
	// Start local and global servers
	go startLocalServer()
	startGlobalServer()
}
