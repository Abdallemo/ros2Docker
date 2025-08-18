/*
Copyright © 2025 Abdallemo <learn3038it@gmail.com>
*/
package main

import (
	"log"

	"github.com/Abdallemo/ros2Docker/cmd"
	"github.com/Abdallemo/ros2Docker/internals/docker"
	"github.com/Abdallemo/ros2Docker/internals/utils"
	"github.com/docker/docker/client"
)

func main() {

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal("unable to connect to docker", err)
	}
	dckr := docker.NewDocker(cli)
	_, err = utils.NewConfig("ros2docker")
	if err != nil {
		log.Fatal("unable to setup Program Configurations", err)
	}

	cmd.Execute(dckr)

}
