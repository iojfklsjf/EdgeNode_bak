package iplibrary

import "github.com/TeaOSLab/EdgeNode/internal/utils"

type IPItem struct {
	Id        int64
	IPFrom    uint32
	IPTo      uint32
	ExpiredAt int64
}

func (this *IPItem) Contains(ip uint32) bool {
	if this.IPTo == 0 {
		if this.IPFrom != ip {
			return false
		}
	} else {
		if this.IPFrom > ip || this.IPTo < ip {
			return false
		}
	}
	if this.ExpiredAt > 0 && this.ExpiredAt < utils.UnixTime() {
		return false
	}
	return true
}