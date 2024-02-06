package docker

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

func docker() {
	fmt.Println("Super docker runner")

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	// pull image
	// reader, err := cli.ImagePull(context.Background(), "docker.io/library/alpine", types.ImagePullOptions{})
	// if err != nil { panic(err) }

	// build image
	// res, err := cli.ImageBuild(context.Background(), nil, types.ImageBuildOptions{
	// 	Dockerfile: filepath.Base("C:\\Users\\kajos\\git\\go\\docker\\javascript\\Dockerfile"),
	// 	Tags: []string{"test-js2"},
	// 	// Remove: true,
	// 	// ForceRemove: true,
	// })
	// if err != nil { panic(err) }
	// fmt.Println("Image build output:", res.OSType)
	// defer res.Body.Close()

	jsContainer, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image: "test-py",
		// Volumes: map[string]struct{}{
		// 	"/app/main.js:C:\\Users\\kajos\\git\\go\\docker\\javascript\\main.js": {},
		// },
		// Cmd: []string{"node","-e", "const b = 5 + 5; console.log(\"yooo: \" + a);"},
		// Env: []string{"CODE=const a = 5 + 5; console.log('yooo: ' + a);"},
		Env: []string{"CODE=b=5+5;print('yooo: ', a)"},
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(context.Background(), jsContainer.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	// get container output
	reader, err := cli.ContainerLogs(context.Background(), jsContainer.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	buf := new(strings.Builder)
	bufError := new(strings.Builder)

	// demultiplex output with github.com/docker/docker/pkg/stdcopy.StdCopy
	// _, err = stdcopy.StdCopy(os.Stdout, os.Stderr, reader) // to host stdout
	_, err = stdcopy.StdCopy(buf, bufError, reader) // to buffer
	if err != nil {
		log.Fatal(err)
	}

	// demultiplex output with io.Copy
	// _, err = io.Copy(os.Stdout, reader)
	// if err != nil && err != io.EOF { log.Fatal(err) }

	// n, err := io.Copy(buf, reader)
	// fmt.Println("N", n)
	fmt.Println("Buffer output:\n--------\n\t", buf.String())
	fmt.Println("BufferError output:\n--------\n\t", bufError.String())

	// stopping and removing container
	if err := cli.ContainerStop(context.Background(), jsContainer.ID, container.StopOptions{}); err != nil {
		panic(err)
	}
	// if err := cli.ContainerKill(context.Background(), jsContainer.ID, "SIGKILL"); err != nil { panic(err) }
	if err := cli.ContainerRemove(context.Background(), jsContainer.ID, types.ContainerRemoveOptions{}); err != nil {
		panic(err)
	}

	// listing containers
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}

}
