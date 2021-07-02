/*
 * MIT License
 *
 * Copyright (c) 2020 schulterklopfer/__escapee__
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILIT * Y, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package test_helpers

import (
  "context"
  "github.com/docker/docker/api/types"
  "github.com/docker/docker/api/types/container"
  "github.com/docker/docker/client"
  "github.com/docker/docker/pkg/jsonmessage"
  "github.com/docker/docker/pkg/term"
  "github.com/docker/go-connections/nat"
  "golang.org/x/net/nettest"
  "net"
  "os"
)

const PostgresContainerName = "cyphernode_admin_tests"
const PostgresImageName = "docker.io/library/postgres:13-alpine"
const PostgresUser = "cnadmin"
const PostgresPassword = "cnadmin"
const PostgresDb = "cnadmin"


func StartupPostgres() (string, error) {

  containerConfig := &container.Config{
    Image: PostgresImageName,
    ExposedPorts: nat.PortSet{
      "5432/tcp": struct{}{},
    },
    Env: []string{
      "POSTGRES_PASSWORD="+PostgresPassword,
      "POSTGRES_USER="+PostgresUser,
      "POSTGRES_DB="+PostgresDb,
    },
  }

  psqlPort, err := nat.NewPort("tcp", "5432")
  if err != nil {
    return "", err
  }

  portBindings := nat.PortMap{
    psqlPort: []nat.PortBinding{{
      HostIP:   "0.0.0.0",
      HostPort: "5432",
    }},
  }

  hostConfig := &container.HostConfig{
    AutoRemove: true,
    PortBindings: portBindings,
  }

  id, err := RunInContainer( containerConfig, hostConfig, PostgresContainerName)

  if err != nil {
    return "", err
  }

  return id, nil

}


func RunInContainer(containerConfig *container.Config, hostConfig *container.HostConfig, containerName string) (string, error) {

  cli, err := client.NewEnvClient()
  if err != nil {
    return "", err
  }

  containers, err := cli.ContainerList( context.Background(), types.ContainerListOptions{} )

  for _, container := range containers {
    for _, name := range container.Names {
      if name == "/"+containerName {
        return container.ID, nil
      }
    }
  }

  reader, err := cli.ImagePull( context.Background(), containerConfig.Image, types.ImagePullOptions{})

  if err != nil {
    return "", err
  }

  defer reader.Close()

  termFd, isTerm := term.GetFdInfo(os.Stderr)
  err = jsonmessage.DisplayJSONMessagesStream(reader, os.Stderr, termFd, isTerm, nil)

  cont, err := cli.ContainerCreate(
    context.Background(),
    containerConfig,
    hostConfig,
    nil,
    nil,
    containerName)

  if err != nil {
    return "", err
  }

  err = cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})

  if err != nil {
    return "", err
  }

  return cont.ID, nil
}

func CleanupContainer(containerID string ) error {
  cli, err := client.NewEnvClient()
  if err != nil {
    return err
  }

  err = cli.ContainerStop(context.Background(), containerID, nil)
  if err != nil {
    return err
  }

  err = cli.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{
    RemoveVolumes: true,
    RemoveLinks: true,
    Force: true,
  } )
  if err != nil {
    return err
  }

  return nil
}

func GetLocalIP() *net.IP {
  iface, err := nettest.RoutedInterface("ip4", net.FlagUp | net.FlagBroadcast )

  if err != nil {
    return nil
  }

  addrs, err := iface.Addrs()
  // handle err
  if err != nil {
    return nil
  }
  for _, addr := range addrs {
    switch v := addr.(type) {
    case *net.IPNet:
      a := addr.(*net.IPNet)
      _, bits := a.Mask.Size()
      if bits != 32 {
        // no ipv4
        continue
      }
      return &v.IP
    }
    // process IP address
  }
  return nil
}
