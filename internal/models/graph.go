package models

import (
	"fmt"
)

type Graph struct {
	Nodes      []*VANT
	Edges      map[int][]Edge
	Properties map[string]interface{}
}

func (g *Graph) PrintGraph() {
	for _, node := range g.Nodes {
		fmt.Printf("Node %d\n", node.ID)
		for _, edge := range g.Edges[node.ID] {
			fmt.Printf("  -> Node %d, Weight: %f, TransmitionRate: %f\n", edge.To.ID, edge.Weight, edge.TransmitionRate)
		}
	}
}

func (g *Graph) AddProperty(key string, value interface{}) {
	g.Properties[key] = value
}

func (g *Graph) GetProperty(key string) (interface{}, bool) {
	value, ok := g.Properties[key]
	return value, ok
}

func (g *Graph) GetNeighbors(node *VANT) []*VANT {
	var neighbors []*VANT

	for _, edge := range g.Edges[node.ID] {
		//fmt.Println("Aresta", edge)
		neighbors = append(neighbors, &edge.To)
	}

	return neighbors
}

func (g *Graph) GetEdge(origem *VANT, destino *VANT) *Edge {
	for _, edge := range g.Edges[origem.ID] {
		if edge.To.ID == destino.ID {
			return &edge
		}
	}

	return nil
}

func (g *Graph) GetVantByID(id int) *VANT {
	for _, node := range g.Nodes {
		if node.ID == id {
			//fmt.Printf("Node: %p\n", &node)
			return node
		}
	}

	return nil
}