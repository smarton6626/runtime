// Copyright (c) 2018 HyperHQ Inc.
//
// SPDX-License-Identifier: Apache-2.0
//

package containerdshim

import (
	"github.com/containerd/cgroups"
	"github.com/containerd/typeurl"

	google_protobuf "github.com/gogo/protobuf/types"
	vc "github.com/kata-containers/runtime/virtcontainers"
)

func marshalMetrics(s *service, containerID string) (*google_protobuf.Any, error) {
	stats, err := s.sandbox.StatsContainer(containerID)
	if err != nil {
		return nil, err
	}

	metrics := statsToMetrics(stats.CgroupStats)

	data, err := typeurl.MarshalAny(metrics)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func statsToMetrics(cgStats *vc.CgroupStats) *cgroups.Metrics {
	var hugetlb []*cgroups.HugetlbStat
	for pageSize, v := range cgStats.HugetlbStats {
		hugetlb = append(
			hugetlb,
			&cgroups.HugetlbStat{
				Usage:    v.Usage,
				Max:      v.MaxUsage,
				Failcnt:  v.Failcnt,
				Pagesize: pageSize,
			})
	}

	var perCPU []uint64
	perCPU = append(perCPU, cgStats.CPUStats.CPUUsage.PercpuUsage...)

	metrics := &cgroups.Metrics{
		Hugetlb: hugetlb,
		Pids: &cgroups.PidsStat{
			Current: cgStats.PidsStats.Current,
			Limit:   cgStats.PidsStats.Limit,
		},
		CPU: &cgroups.CPUStat{
			Usage: &cgroups.CPUUsage{
				Total:  cgStats.CPUStats.CPUUsage.TotalUsage,
				PerCPU: perCPU,
			},
		},
		Memory: &cgroups.MemoryStat{
			Cache: cgStats.MemoryStats.Cache,
			Usage: &cgroups.MemoryEntry{
				Limit: cgStats.MemoryStats.Usage.Limit,
				Usage: cgStats.MemoryStats.Usage.Usage,
			},
		},
	}

	return metrics
}
