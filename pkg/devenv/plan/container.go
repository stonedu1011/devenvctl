package plan

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/stonedu1011/devenvctl/pkg/utils/tmplutils"
	"io"
	"os"
	"strings"
	"text/template"
	"time"
)

type ContainerResolver func(name string, c *types.Container) string

func ExactNameContainerResolver() ContainerResolver {
	return func(name string, c *types.Container) string {
		for _, n := range c.Names {
			if n == name {
				return n
			}
		}
		return ""
	}
}

func ComposeContainerResolver(profileName string) ContainerResolver {
	return func(name string, c *types.Container) string {
		// first try to use labels
		if p, ok := c.Labels["com.docker.compose.project"]; ok && p == profileName {
			if s, ok := c.Labels["com.docker.compose.service"]; ok && s == name {
				return c.ID
			}
		}
		// fallback to generated name format: <profile>-<service>-#
		contains := fmt.Sprintf(`%s-%s-`, profileName, name)
		for _, n := range c.Names {
			if strings.Contains(n, contains) {
				return n
			}
		}
		return ""
	}
}

func NewContainerMonitorExecutable(client *client.Client, opts ...func(exec *ContainerMonitorExecutable)) *ContainerMonitorExecutable {
	exec := &ContainerMonitorExecutable{
		ApiClient: client,
		Resolver:  ExactNameContainerResolver(),
		Desc:      "containers",
	}
	for _, fn := range opts {
		fn(exec)
	}
	return exec
}

var containerColorPool = []string{
	"cyan", "yellow", "magenta", "green", "blue", "red",
	"cyan_b", "yellow_b", "magenta_b", "green_b", "blue_b", "red_b",
}

type containerEvent struct {
	Container string
	Entry     interface{}
}

type ContainerMonitorExecutable struct {
	ApiClient *client.Client
	Names    []string
	Resolver ContainerResolver
	Desc     string
	tmpls     map[string]*template.Template
}

func (exec *ContainerMonitorExecutable) Exec(ctx context.Context) error {
	ctx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	// prepare
	mapping, e := exec.resolveContainerNames(ctx)
	if e != nil {
		return fmt.Errorf(`unable to find containers: %v`, e)
	}
	exec.prepareTemplates()

	// start monitor all containers
	ch := make(chan containerEvent, 1)
	defer close(ch)
	for _, name := range exec.Names {
		cName, ok := mapping[name]
		if !ok {
			return fmt.Errorf(`unable to find container for [%s]`, name)
		}
		go exec.monitor(ctx, name, cName, ch)
	}

	// wait for all containers to finish
	logger.Infof(`Waiting for %s to finish ...`, exec.Desc)
	finished := map[string]error{}
	for len(finished) < len(exec.Names) {
		select {
		case <-ctx.Done():
			return context.DeadlineExceeded
		case evt := <-ch:
			switch v := evt.Entry.(type) {
			case error:
				finished[evt.Container] = v
				if errors.Is(v, io.EOF) {
					exec.printEvent(containerEvent{Container: evt.Container, Entry: "exited with code 0"})
				} else {
					exec.printEvent(containerEvent{Container: evt.Container, Entry: fmt.Sprintf("error: %v", e)})
				}
			case string:
				exec.printEvent(evt)
			}
		}
	}

	// check error
	errored := make([]string, 0, len(finished))
	for k, v := range finished {
		if !errors.Is(v, io.EOF) {
			errored = append(errored, k)
		}
	}
	if len(errored) != 0 {
		return fmt.Errorf(`containers [%s] didn't finish without error'`, strings.Join(errored, ", "))
	}
	return nil
}

func (exec *ContainerMonitorExecutable) String() string {
	if len(exec.Desc) == 0 {
		exec.Desc = "containers"
	}
	switch {
	case len(exec.Names) == 1:
		return fmt.Sprintf("%s: %s", exec.Desc, exec.Names[0])
	case len(exec.Names) == 0:
		return "no-op"
	default:
		return fmt.Sprintf("%s: \n    %s", exec.Desc, strings.Join(exec.Names, "\n    "))
	}
}

func (exec *ContainerMonitorExecutable) WithTimeout(timeout time.Duration) Executable {
	return &TimeoutExecutableWrapper{
		Timeout:  timeout,
		Delegate: exec,
	}
}

func (exec *ContainerMonitorExecutable) prepareTemplates() {
	exec.tmpls = map[string]*template.Template{}
	maxWidth := 0
	for _, c := range exec.Names {
		if len(c) > maxWidth {
			maxWidth = len(c)
		}
	}
	format := `{{pad %d .Container | %s}} | {{.Entry}}`
	for i, c := range exec.Names {
		color := containerColorPool[i%len(containerColorPool)]
		tmplText := fmt.Sprintf(format, maxWidth+4, color)
		if tmpl, e := tmplutils.Parse(tmplText); e == nil {
			exec.tmpls[c] = tmpl
		}

	}
}

func (exec *ContainerMonitorExecutable) printEvent(evt containerEvent) {
	if !strings.HasSuffix(evt.Entry.(string), "\n") {
		evt.Entry = evt.Entry.(string) + "\n"
	}
	if tmpl, ok := exec.tmpls[evt.Container]; ok {
		if e := tmpl.Execute(os.Stdout, evt); e == nil {
			return
		}
	}
	fmt.Printf(`[%s] | %s`, evt.Container, evt.Entry)
}

func (exec *ContainerMonitorExecutable) resolveContainerNames(ctx context.Context) (map[string]string, error) {
	containers, e := exec.ApiClient.ContainerList(ctx, container.ListOptions{})
	if e != nil {
		return nil, e
	}
	idMapping := map[string]string{}
	for i := range containers {
		for _, name := range exec.Names {
			if cn := exec.Resolver(name, &containers[i]); len(cn) != 0 {
				idMapping[name] = cn
				break
			}
		}
	}
	return idMapping, nil
}

func (exec *ContainerMonitorExecutable) monitor(ctx context.Context, name string, cName string, ch chan containerEvent) {
	reader, e := exec.ApiClient.ContainerLogs(ctx, cName, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Details:    true,
	})
	if e != nil {
		ch <- containerEvent{Container: name, Entry: e}
		return
	}
	defer func() { _ = reader.Close() }()
	header := make([]byte, 8)
LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		default:
		}
		switch _, e := reader.Read(header); {
		case e != nil:
			ch <- containerEvent{Container: name, Entry: e}
			break LOOP
		}
		size := binary.BigEndian.Uint32(header[4:])
		buf := make([]byte, size)
		if _, e = reader.Read(buf); e != nil {
			ch <- containerEvent{Container: name, Entry: e}
			break LOOP
		}
		ch <- containerEvent{Container: name, Entry: string(buf)}
	}
	return
}
