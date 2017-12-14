package docker

import (
    "github.com/fsouza/go-dockerclient"
    "bytes"
    "strings"
    "strconv"
)

type ContainerWrapper struct {
    Container *docker.Container
}

type Client struct {
    ContainerName string
    DockerClient  *docker.Client
}

func CreateClient(containerName string) (*Client, error) {
    dockerClient, err := docker.NewClient("unix:///var/run/docker.sock")
    if err != nil {
        return nil, err
    }
    client := Client{
        ContainerName: containerName,
        DockerClient: dockerClient,
    }
    return &client, nil;
}

func (client *Client) GetLogs(numLines int) ([]string, error) {
    var buffer bytes.Buffer

    options := docker.LogsOptions{
        Container: client.ContainerName,
        Tail: strconv.Itoa(numLines),
        RawTerminal: false,
        OutputStream: &buffer,
        Stdout: true,
        Timestamps: true,
    }
    client.DockerClient.Logs(options)
    return strings.Split(buffer.String(), "\n"), nil
}

func (client *Client) InspectContainer() (*ContainerWrapper, error) {
    container, err := client.DockerClient.InspectContainer(client.ContainerName)
    return &ContainerWrapper{Container: container}, err
}

func (client *Client) Restart(timeout uint) (error) {
    return client.DockerClient.RestartContainer(client.ContainerName, timeout)
}
