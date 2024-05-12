package plan

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/filters"
	dockerclient "github.com/docker/docker/client"
)

func formatSize(size uint64) string {
	units := []string{"B", "KB", "MB", "GB"}
	v := float64(size)
	for i := range units {
		if v < 1024 {
			if i == 0 {
				return fmt.Sprintf(`%.0f%s`, v, units[i])
			}
			return fmt.Sprintf(`%.2f%s`, v, units[i])
		}
		v = v / 1024
	}
	return fmt.Sprintf(`%.2fGB`, v)
}

type PruneContainersExecutable struct {
	ApiClient *dockerclient.Client
}

func (exec *PruneContainersExecutable) Exec(ctx context.Context) error {
	logger.WithContext(ctx).Infof(`Pruning containers...`)
	report, e := exec.ApiClient.ContainersPrune(ctx, filters.NewArgs())
	if e != nil {
		return e
	}
	logger.WithContext(ctx).Infof(`Deleted %d containers and reclaimed %s space`, len(report.ContainersDeleted), formatSize(report.SpaceReclaimed))
	return nil
}

func (exec *PruneContainersExecutable) String() string {
	return `prune containers`
}

type PruneVolumesExecutable struct {
	ApiClient *dockerclient.Client
}

func (exec *PruneVolumesExecutable) Exec(ctx context.Context) error {
	logger.WithContext(ctx).Infof(`Pruning containers...`)
	report, e := exec.ApiClient.VolumesPrune(ctx, filters.NewArgs(filters.Arg("label!", "devenv.persist")))
	if e != nil {
		return e
	}
	logger.WithContext(ctx).Infof(`Deleted %d volumes and reclaimed %s space`, len(report.VolumesDeleted), formatSize(report.SpaceReclaimed))
	return nil
}

func (exec *PruneVolumesExecutable) String() string {
	return `prune volumes`
}

type PruneImagesExecutable struct {
	ApiClient *dockerclient.Client
}

func (exec *PruneImagesExecutable) Exec(ctx context.Context) error {
	logger.WithContext(ctx).Infof(`Pruning containers...`)
	report, e := exec.ApiClient.ImagesPrune(ctx, filters.NewArgs())
	if e != nil {
		return e
	}
	logger.WithContext(ctx).Infof(`Deleted %d images and reclaimed %s space`, len(report.ImagesDeleted), formatSize(report.SpaceReclaimed))
	return nil
}

func (exec *PruneImagesExecutable) String() string {
	return `prune images`
}
