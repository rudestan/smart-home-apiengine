package main

import "smh-apiengine/pkg/rmqproc"

func main() {
	rmqproc.ListenHttpAndPublish()
}
