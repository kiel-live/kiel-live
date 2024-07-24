package main

type KielLiveRPC struct{}

func (t *KielLiveRPC) Hello(name string) (string, error) {
	return "Hello, " + name, nil
}
