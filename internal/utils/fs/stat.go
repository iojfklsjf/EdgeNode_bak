// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package fsutils

import (
	"github.com/TeaOSLab/EdgeNode/internal/utils/fasttime"
	"golang.org/x/sys/unix"
	"sync"
)

// Stat device contains the path
func Stat(path string) (*StatResult, error) {
	var stat = &unix.Statfs_t{}
	err := unix.Statfs(path, stat)
	if err != nil {
		return nil, err
	}
	return NewStatResult(stat), nil
}

var locker = &sync.RWMutex{}
var cacheMap = map[string]*StatResult{} // path => StatResult

const cacheLife = 3 // seconds

// StatCache stat device with cache
func StatCache(path string) (*StatResult, error) {
	locker.RLock()
	stat, ok := cacheMap[path]
	if ok && stat.updatedAt >= fasttime.Now().Unix()-cacheLife {
		locker.RUnlock()
		return stat, nil
	}
	locker.RUnlock()

	locker.Lock()
	defer locker.Unlock()

	stat, err := Stat(path)
	if err != nil {
		return nil, err
	}

	cacheMap[path] = stat
	return stat, nil
}

type StatResult struct {
	rawStat   *unix.Statfs_t
	blockSize uint64

	updatedAt int64
}

func NewStatResult(rawStat *unix.Statfs_t) *StatResult {
	var blockSize = rawStat.Bsize
	if blockSize < 0 {
		blockSize = 0
	}

	return &StatResult{
		rawStat:   rawStat,
		blockSize: uint64(blockSize),
		updatedAt: fasttime.Now().Unix(),
	}
}

func (this *StatResult) AvailableSize() uint64 {
	return this.rawStat.Bavail * this.blockSize
}

func (this *StatResult) TotalSize() uint64 {
	return this.rawStat.Blocks * this.blockSize
}

func (this *StatResult) UsedSize() uint64 {
	if this.rawStat.Bavail <= this.rawStat.Blocks {
		return (this.rawStat.Blocks - this.rawStat.Bavail) * this.blockSize
	}
	return 0
}
