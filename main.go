package main

import (
	"github.com/towski/Golang-AStar/utils"
)

func main() {
	var scene utils.Scene
	scene.InitScene(23, 70)
	scene.AddWalls(10)
	utils.InitAstar(&scene)

	for {
		utils.FindPath(&scene)
		//time.Sleep(50 * time.Millisecond)
        if utils.Result != 10 {
            break
        }
	}
}
