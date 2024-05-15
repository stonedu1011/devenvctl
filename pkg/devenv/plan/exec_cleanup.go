package plan

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	dockerclient "github.com/docker/docker/client"
	"github.com/stonedu1011/devenvctl/pkg/utils"
)

type pruneExecutable struct{}

func (pruneExecutable) formatSize(size uint64) string {
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

func (exec pruneExecutable) printVerboseReport(header string, values[]string, reclaimedSize uint64) {
	logger.Debugf("")
	if len(values) == 0 {
		fmt.Printf("%s: NONE\n", header)
	} else {
		fmt.Printf("%s:\n", header )
		for _, v := range values {
			fmt.Printf("    %s\n", v)
		}
	}
	fmt.Printf("Total reclaimed space: %s\n", exec.formatSize(reclaimedSize))
}

type PruneContainersExecutable struct {
	pruneExecutable
	ApiClient *dockerclient.Client
}

func (exec *PruneContainersExecutable) Exec(ctx context.Context, opts ExecOption) error {
	if opts.DryRun {
		fmt.Printf("- %v\n", exec)
		return nil
	}
	logger.WithContext(ctx).Infof(`Pruning containers...`)
	report, e := exec.ApiClient.ContainersPrune(ctx, filters.NewArgs())
	if e != nil {
		return e
	}

	if opts.Verbose {
		exec.printVerboseReport("Deleted Containers", report.ContainersDeleted, report.SpaceReclaimed)
	} else {
		logger.WithContext(ctx).Infof(`Deleted %d containers and reclaimed %s space`, len(report.ContainersDeleted), exec.formatSize(report.SpaceReclaimed))
	}
	return nil
}

func (exec *PruneContainersExecutable) String() string {
	return `prune containers`
}

type PruneVolumesExecutable struct {
	pruneExecutable
	ApiClient *dockerclient.Client
}

func (exec *PruneVolumesExecutable) Exec(ctx context.Context, opts ExecOption) error {
	if opts.DryRun {
		fmt.Printf("- %v\n", exec)
		return nil
	}
	logger.WithContext(ctx).Infof(`Pruning volumes...`)
	report, e := exec.ApiClient.VolumesPrune(ctx, filters.NewArgs(filters.Arg("label!", "devenv.persist")))
	if e != nil {
		return e
	}

	if opts.Verbose {
		exec.printVerboseReport("Deleted Volumes", report.VolumesDeleted, report.SpaceReclaimed)
	} else {
		logger.WithContext(ctx).Infof(`Deleted %d volumes and reclaimed %s space`, len(report.VolumesDeleted), exec.formatSize(report.SpaceReclaimed))
	}
	return nil
}

func (exec *PruneVolumesExecutable) String() string {
	return `prune volumes`
}

type PruneImagesExecutable struct {
	pruneExecutable
	ApiClient *dockerclient.Client
}

func (exec *PruneImagesExecutable) Exec(ctx context.Context, opts ExecOption) error {
	if opts.DryRun {
		fmt.Printf("- %v\n", exec)
		return nil
	}
	logger.WithContext(ctx).Infof(`Pruning images...`)
	report, e := exec.ApiClient.ImagesPrune(ctx, filters.NewArgs())
	if e != nil {
		return e
	}
	if opts.Verbose {
		values := utils.ConvertSlice(report.ImagesDeleted, func(resp image.DeleteResponse) string { return resp.Deleted })
		exec.printVerboseReport("Deleted Images", values, report.SpaceReclaimed)
		values = utils.ConvertSlice(report.ImagesDeleted, func(resp image.DeleteResponse) string { return resp.Untagged })
		exec.printVerboseReport("Untagged Images", values, report.SpaceReclaimed)
	} else {
		logger.WithContext(ctx).Infof(`Deleted %d images and reclaimed %s space`, len(report.ImagesDeleted), exec.formatSize(report.SpaceReclaimed))
	}
	return nil
}

func (exec *PruneImagesExecutable) String() string {
	return `prune images`
}
