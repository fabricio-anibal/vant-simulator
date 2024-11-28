package processor

import (
	"fmt"
	"time"
	"vantsimulator/internal/commands/handler"
	"vantsimulator/internal/processor/vant_util"
)

func Process() {
	vants, err := handler.Read("./data/vants.csv")

	if err != nil {
		fmt.Println(err)
		return
	}

	vantsNetwork := vant_util.BuildGraphNetwork(vants)

	//fmt.Println(vantsNetwork)

	vantsNetwork.AddProperty("AvgTransmitionRate", vant_util.AvgTransmitionRate(vantsNetwork))

	vantsNetwork.PrintGraph()

	//fmt.Println(vantsNetwork.GetNeighbors(&vants[0]))

	//vant_util.SendMessage(vantsNetwork, &vants[0], &vants[1], "Hello, World!")

	//fmt.Println(&vants[1].MessagesBuffer)

	startBroadcast := time.Now()

	vant_util.SendBroadcast(vantsNetwork, vantsNetwork.GetVantByID(1), "Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!")

	endBroadcast := time.Now()

	duration := endBroadcast.Sub(startBroadcast)

	fmt.Println("Duration:", duration)

	//fmt.Println(vant_util.AvgTransmitionRate(vantsNetwork))
}
