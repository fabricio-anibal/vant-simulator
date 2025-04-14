package vant_util

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand/v2"
	"time"
	"vantsimulator/internal/models"
)

const (
	LIMIT_BOX = 100.0
)

func maxConnectionDistance(numDrones int) float64 {
	volume := LIMIT_BOX * LIMIT_BOX * LIMIT_BOX
	c := 1.4
	distance := c * math.Cbrt(volume/float64(numDrones))
	return distance
	//return 9.0
}

func BuildGraphNetwork(vants []models.VANT) *models.Graph {
	var graph models.Graph

	graph.Edges = make(map[int][]models.Edge)
	graph.Properties = make(map[string]interface{})

	for _, vant := range vants {
		graph.Nodes = append(graph.Nodes, &vant)
		for _, otherVant := range vants {
			if vant.ID != otherVant.ID {
				dist := distance(vant, otherVant)
				if dist < maxConnectionDistance(len(vants)) {
					graph.Edges[vant.ID] = append(graph.Edges[vant.ID], models.Edge{
						To:              otherVant,
						Weight:          dist,
						TransmitionRate: transmitionRate(vant, otherVant),
						Id:              fmt.Sprintf("%d-%d", vant.ID, otherVant.ID),
					})
				}
			}
		}
	}
	return &graph
}

func AddVant(graph *models.Graph, vant models.VANT) {
	graph.Nodes = append(graph.Nodes, &vant)
	for _, otherVant := range graph.Nodes {
		if vant.ID != otherVant.ID {
			dist := distance(vant, *otherVant)
			if dist < maxConnectionDistance(len(graph.Nodes)) {
				graph.Edges[vant.ID] = append(graph.Edges[vant.ID], models.Edge{
					To:              *otherVant,
					Weight:          dist,
					TransmitionRate: transmitionRate(vant, *otherVant),
					Id:              fmt.Sprintf("%d-%d", vant.ID, otherVant.ID),
				})
				graph.Edges[otherVant.ID] = append(graph.Edges[otherVant.ID], models.Edge{
					To:              vant,
					Weight:          dist,
					TransmitionRate: transmitionRate(vant, *otherVant),
					Id:              fmt.Sprintf("%d-%d", vant.ID, otherVant.ID),
				})

			}
		}
	}
}

func CentroideVant(graph *models.Graph) models.VANT {
	var sumX float64
	var sumY float64
	var sumZ float64
	var count float64

	for _, vant := range graph.Nodes {
		sumX += vant.X
		sumY += vant.Y
		sumZ += vant.Z
		count++
	}

	return models.VANT{
		ID: int(count) + 1,
		X:  sumX / count,
		Y:  sumY / count,
		Z:  sumZ / count,
	}
}

func RandomVant(graph *models.Graph, limitBox int) models.VANT {
	var vant models.VANT
	vant.X = float64(limitBox) * rand.Float64()
	vant.Y = float64(limitBox) * rand.Float64()
	vant.Z = float64(limitBox) * rand.Float64()
	vant.ID = len(graph.Nodes) + 1
	return vant
}

func distance(vant1 models.VANT, vant2 models.VANT) float64 {
	x1, y1, z1 := vant1.X, vant1.Y, vant1.Z
	x2, y2, z2 := vant2.X, vant2.Y, vant2.Z
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2) + math.Pow(z2-z1, 2))
}

func pathLoss(vant1 models.VANT, vant2 models.VANT) float64 {
	const Alpha = 2.03
	const D0 = 1.0
	const PathLostD0 = 46.6

	dist := distance(vant1, vant2)

	pathLost := PathLostD0 + (10 * Alpha * math.Log10(dist/D0))

	return pathLost
}

func transmitionRate(vant1 models.VANT, vant2 models.VANT) float64 {
	const B = 40.0
	const Ptx = 20.0
	const T = 290.0
	k := 1.38 * math.Pow(10, -23)

	pl := pathLoss(vant1, vant2)
	prx := Ptx - pl

	prxWatts := math.Pow(10, (prx-30)/10)

	pn := k * T * B

	c := B * math.Log2(1+(prxWatts/pn))

	return c
}

func AvgTransmitionRate(graph *models.Graph) float64 {
	var sum float64
	var count float64

	for _, edges := range graph.Edges {
		for _, edge := range edges {
			sum += edge.TransmitionRate
			count++
		}
	}

	return sum / count
}

func SendMessage(graph *models.Graph, origem *models.VANT, destino *models.VANT, message string) {
	n := graph.GetNeighbors(origem)

	bits := stringToBits(message)

	//fmt.Println("Bits:", len(bits))

	for i := range n {
		neighbor := graph.GetVantByID(n[i].ID)
		//transmitionAvailable := graph.GetEdge(origem, neighbor).TransmitionAvailable
		if neighbor.ID == destino.ID {
			//fmt.Println("Enviando mensagem de", origem.ID, "para", destino.ID)
			messageEnd := false

			edge := graph.GetEdge(origem, neighbor)
			transmitionRate := edge.TransmitionRate

			rateLimiter := NewRateLimiter(int64(math.Round(transmitionRate)), time.Second)

			offset := 0

			for !messageEnd {
				/*availableData := rateLimiter.GetAvailableData()
				fmt.Println("Loop:", offset)
				fmt.Println("Bits:", len(bits))
				fmt.Println("Available:", availableData)*/

				allowed, availableData := rateLimiter.AllowSoft(int64(len(bits)) - int64(offset))

				//fmt.Println("Available:", availableData)

				if allowed {
					if availableData > 0 {
						transmitionAvailable := math.Min(float64(availableData), float64(len(bits)))
						//fmt.Println("transmitionAvailable:", transmitionAvailable)
						bitsToSend := []int{}
						if offset+int(transmitionAvailable) > len(bits) {
							bitsToSend = bits[offset:]
						} else {
							bitsToSend = bits[offset : offset+int(transmitionAvailable)]
						}
						//fmt.Println("bitsToSend:", bitsToSend)
						offset = offset + int(transmitionAvailable)
						//fmt.Println("bitsToSend:", bitsToSend)
						destino.ReceiveMessage(generateHash(message), bitsToSend)
					}

					if offset >= len(bits) {
						messageEnd = true
					}
				}
			}
			return
		}
	}
}

func stringToBits(s string) []int {
	var bits []int
	for _, char := range []byte(s) {
		for i := 7; i >= 0; i-- {
			bit := (char >> i) & 1
			bits = append(bits, int(bit))
		}
	}
	return bits
}

func SendBroadcast(graph *models.Graph, origem *models.VANT, message string) {
	n := graph.GetNeighbors(origem)

	//fmt.Println("Origem:", origem)
	//fmt.Println("Vizinhos:", n[0])

	continueBraodcast := false

	for i := range n {
		neighbor := graph.GetVantByID(n[i].ID)
		//fmt.Printf("%p\n", graph.GetVantByID(n[i].ID))
		//fmt.Println(&neighbor)
		//fmt.Println("Neigbhor buffer", neighbor.MessagesBuffer)
		if !neighbor.HasMessage(generateHash(message)) {
			continueBraodcast = true
			SendMessage(graph, origem, neighbor, message)
			//fmt.Println(neighbor.MessagesBuffer)
			//fmt.Println(&neighbor)
			//fmt.Println(graph.GetVantByID(neighbor.ID).MessagesBuffer)
			//fmt.Printf("%p\n", graph.GetVantByID(neighbor.ID))
		}
	}

	if continueBraodcast {
		for i := range n {
			neighbor := n[i]
			SendBroadcast(graph, neighbor, message)
		}
	}
}

func generateHash(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	hashedBytes := hash.Sum(nil)

	// Convertendo para string hexadecimal
	hashedString := hex.EncodeToString(hashedBytes)

	return hashedString
}

func CleanBuffer(graph *models.Graph) {
	for _, vant := range graph.Nodes {
		vant.CleanBuffer()
	}
}
