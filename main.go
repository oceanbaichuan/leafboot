package main

import (
	"github.com/hudgit2019/leafboot/appmain"
	"github.com/hudgit2019/leafboot/gameboot"
)

func main() {
	gameboot.StartGame(&appmain.GuessLogic{}, &appmain.GuessRobotLogic{})
}
