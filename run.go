package main

import (
	"context"
	"fmt"
	"os"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/linux/runctypes"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/sirupsen/logrus"
)

func main() {
	argsWithProg := os.Args
	runContainer(argsWithProg[1:])
}

func runContainer(args []string) {
	if len(args) < 2 {
		logrus.Infof("./run image id runtime arg...")
		return
	}
	var (
		imageref = args[0]
		id       = args[1]
		runtime  = args[2]
		arg      = args[3:]
	)

	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		logrus.Errorf("fail to create containerd %s", err)
		return
	}
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), "default")

	image, err := client.GetImage(ctx, imageref)
	if err != nil {
		logrus.Errorf("fail to get image %s %s", imageref, err)
		return
	}

	specOpts := []oci.SpecOpts{
		oci.WithImageConfig(image),
	}

	if len(arg) > 0 {
		specOpts = append(specOpts, oci.WithProcessArgs(arg...))
	}

	opts := []containerd.NewContainerOpts{
		containerd.WithImage(image),
		containerd.WithNewSnapshot(fmt.Sprintf("%s-snapshot", id), image),
		containerd.WithRuntime("io.containerd.runtime.v1.linux", &runctypes.RuncOptions{
			Runtime: runtime,
		}),
		containerd.WithNewSpec(specOpts...),
	}
	container, err := client.NewContainer(ctx, id, opts...)
	if err != nil {
		logrus.Errorf("fail to create container %s", err)
		return
	}

	logrus.Infof("start new task")
	task, err := container.NewTask(ctx, cio.Stdio)
	if err != nil {
		logrus.Errorf("fail to create task %s", err)
		return
	}
	defer task.Delete(ctx)

	logrus.Infof("start wait task")

	exitStatusC, err := task.Wait(ctx)
	if err != nil {
		logrus.Warningf("fail to call task wait %s", err)
	}

	if err := task.Start(ctx); err != nil {
		logrus.Errorf("fail to start task %s", err)
		return
	}

	logrus.Infof("start task successful")
	// wait for the process to fully exit and print out the exit status

	status := <-exitStatusC
	code, _, err := status.Result()
	if err != nil {
		logrus.Errorf("fail to get exit result %s", err)
		return
	}

	logrus.Infof("task exit with status %v", code)
}
