package main

import "os"

func main() {
	configFile, err := os.Open("./config.yml")
	if err != nil {
		panic(err)
	}

	config, err := ParseConfig(configFile)
	if err != nil {
		panic(err)
	}

	server, err := NewServer(config.Host, config.Port, config)
	if err != nil {
		panic(err)
	}

	panic(server.Start())
}
