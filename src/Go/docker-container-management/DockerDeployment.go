package dockercontainermanagement

// This file contains fucntions that can be called to
// start a docker container on the host machine

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

// The DockerDeployment struct is a json that will be recived from the UI to deploy a docker Container
type DockerDeployment struct {
	TYPE                 string `json:"type"`
	NAME                 string `json:"name"`
	IMAGENAME            string `json:"imageName"`
	HOSTMACHINE          string `json:"hostMachine"`
	RESTARTPOLICY        bool   `json:"restartPolicy"`
	PORTMAPPING          string `json:"portMapping"`
	ENVIRONMENTVARIABLES string `json:"environementVariables"`
	VOLUMEMOUNTS         string `json:"volumeMounts"`
}

// DeployContainerBackground deploys a container as background process
// on the host machine much like docker run -d container would

func DeployContainerBackground(dockerContainer DockerContainer, s *DBService) error {

	ctx := context.Background()

	// Create Docker API client
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		log.Printf("❗ Error creating the Docker API client : %v", err)
		return err
	}
	defer cli.Close()

	// Pull the image
	out, err := cli.ImagePull(ctx, dockerContainer.Name, image.PullOptions{})
	if err != nil {
		log.Printf("❗ Error pylling the docker image : %v", err)
		return err
	}
	defer out.Close()
	io.Copy(os.Stdout, out)

	// Create container
	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{Image: dockerContainer.Name},
		nil, nil, nil,
		dockerContainer.Container, // container name
	)
	if err != nil {
		log.Printf("❗ Error creating the docker container : %v", err)
		return err
	}

	// Start container
	if err := cli.ContainerStart(
		ctx,
		resp.ID,
		container.StartOptions{},
	); err != nil {
		log.Printf("❗ Error starting the docker container : %v", err)
		return err
	}

	// Insert the container inside the Bolt database
	s.SaveActiveDockerContainer(dockerContainer)

	log.Println("🧰 Started container:", resp.ID)

	return nil
}
