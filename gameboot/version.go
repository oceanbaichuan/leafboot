package gameboot

import (
	"github.com/name5566/leaf/log"
)

const LEAFBOOT_VERSION = "20.1.15.01"

func init() {
	log.Debug("LeafBoot:%v initialized", LEAFBOOT_VERSION)
}
