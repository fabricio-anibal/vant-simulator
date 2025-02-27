package processor

import (
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"vantsimulator/internal/commands/handler"
	"vantsimulator/internal/models"
	"vantsimulator/internal/processor/vant_util"
)

const (
	LIMIT_BOX = 100.0
)

var mu sync.Mutex

func ProcessSim() {
	root := "./data"

	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %v: %v\n", root, err)
		return
	}

	for _, file := range files {
		fmt.Println("====================================================================== PROCESSING FILE " + file + " ======================================================================")
		resultsRandom, networkRandom := processRandom(file)
		resultsCentroide, networkCentroide := processCentroide(file)

		fileResultRandom, _ := os.Create("./results/random/" + strings.Split(strings.Split(file, "/")[1], ".")[0] + ".txt")
		fileResultCentroide, _ := os.Create("./results/centroide/" + strings.Split(strings.Split(file, "/")[1], ".")[0] + ".txt")

		fileResultRandom.WriteString(networkRandom.ToString())

		fileResultRandom.WriteString(resultsRandom[0].ToString() + "\n")
		fileResultRandom.WriteString(resultsRandom[1].ToString() + "\n")
		fileResultRandom.WriteString(resultsRandom[2].ToString() + "\n")

		fileResultCentroide.WriteString(networkCentroide.ToString())

		fileResultCentroide.WriteString(resultsCentroide[0].ToString() + "\n")
		fileResultCentroide.WriteString(resultsCentroide[1].ToString() + "\n")
		fileResultCentroide.WriteString(resultsCentroide[2].ToString() + "\n")

		fileResultRandom.Close()
		fileResultCentroide.Close()

		fmt.Println("====================================================================== END OF FILE " + file + " ======================================================================")
	}
}

func testNetwork(vantsNetwork *models.Graph) (models.Stats, models.Stats, models.Stats) {
	results := make(map[int]float64)
	var wg sync.WaitGroup

	message := "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."

	// Broadcast para cada vant
	for _, vant := range vantsNetwork.Nodes {
		wg.Add(1)

		// Cópia de vant e vantsNetwork para evitar concorrência
		vantCopy := deepCopyVANT(vant)
		vantsNetworkCopy := deepCopyGraph(vantsNetwork)

		// Passando cópias para a goroutine
		go processBroadcast(*vantCopy, *vantsNetworkCopy, message, results, &wg)
	}

	wg.Wait()

	// Calculando as estatísticas de resultados
	min := models.Stats{
		Value: 0.0,
		Name:  "Minimo",
		Id:    0,
	}
	max := models.Stats{
		Value: 0.0,
		Name:  "Maximo",
		Id:    0,
	}
	avg := models.Stats{
		Value: 0.0,
		Name:  "Media",
		Id:    0,
	}
	count := 0
	for id, duration := range results {
		if min.Value == 0.0 || duration < min.Value {
			min.Value = duration
			min.Id = id
		}
		if max.Value == 0.0 || duration > max.Value {
			max.Value = duration
			max.Id = id
		}
		avg.Value += duration
		count++
	}
	avg.Value /= float64(count)

	return min, max, avg
}

func processBroadcast(vant models.VANT, vantsNetwork models.Graph, message string, results map[int]float64, wg *sync.WaitGroup) {
	defer wg.Done()
	//fmt.Printf("Starting broadcast[%d] \n", vant.ID)

	startBroadcast := time.Now()

	// Passando cópias para SendBroadcast (não mais ponteiros)
	vant_util.SendBroadcast(&vantsNetwork, &vant, message)

	endBroadcast := time.Now()

	duration := endBroadcast.Sub(startBroadcast)

	mu.Lock()
	results[vant.ID] = duration.Seconds()
	mu.Unlock()

	//vant_util.CleanBuffer(vantsNetwork)

	//fmt.Printf("Broadcast finished[%d]: %f\n", vant.ID, duration.Seconds())
}

func processRandom(filePath string) ([]models.Stats, *models.Graph) {
	vants, err := handler.Read(filePath)

	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	fmt.Println("====================================RANDOM====================================")

	vantsNetwork := vant_util.BuildGraphNetwork(vants)

	vantsNetwork.AddProperty("AvgTransmitionRate", vant_util.AvgTransmitionRate(vantsNetwork))

	//vantsNetwork.PrintGraph()

	fmt.Println("VANT ADICIONADO")

	vant_util.AddVant(vantsNetwork, vant_util.RandomVant(vantsNetwork, LIMIT_BOX))

	vantsNetwork.PrintGraph()

	min, max, avg := testNetwork(vantsNetwork)

	fmt.Println("RANDOM")
	fmt.Printf("Minimo [%d]: %f\n", min.Id, min.Value)
	fmt.Printf("Maximo [%d]: %f\n", max.Id, max.Value)
	fmt.Printf("Media: %f\n", avg.Value)

	result := []models.Stats{min, max, avg}

	return result, vantsNetwork
}

func processCentroide(filePath string) ([]models.Stats, *models.Graph) {
	vants, err := handler.Read(filePath)

	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	fmt.Println("====================================CENTROIDE====================================")

	vantsNetwork := vant_util.BuildGraphNetwork(vants)

	vantsNetwork.AddProperty("AvgTransmitionRate", vant_util.AvgTransmitionRate(vantsNetwork))

	//vantsNetwork.PrintGraph()

	fmt.Println("VANT ADICIONADO")

	vant_util.AddVant(vantsNetwork, vant_util.CentroideVant(vantsNetwork))

	vantsNetwork.PrintGraph()

	min, max, avg := testNetwork(vantsNetwork)

	fmt.Println("CENTROIDE")
	fmt.Printf("Minimo [%d]: %f\n", min.Id, min.Value)
	fmt.Printf("Maximo [%d]: %f\n", max.Id, max.Value)
	fmt.Printf("Media: %f\n", avg.Value)

	result := []models.Stats{min, max, avg}

	return result, vantsNetwork
}

func deepCopyGraph(original *models.Graph) *models.Graph {
	// Criar uma nova instância de Graph
	copyGraph := &models.Graph{
		Edges:      make(map[int][]models.Edge),
		Properties: make(map[string]interface{}),
	}

	// Copiar os nós (deep copy dos VANTs)
	for _, vant := range original.Nodes {
		vantCopy := *vant // Criar uma cópia do VANT (não apenas o ponteiro)
		copyGraph.Nodes = append(copyGraph.Nodes, &vantCopy)
	}

	// Copiar as arestas (Edges)
	for key, edges := range original.Edges {
		copyGraph.Edges[key] = append([]models.Edge(nil), edges...)
	}

	// Copiar as propriedades
	for key, value := range original.Properties {
		copyGraph.Properties[key] = value
	}

	return copyGraph
}

func deepCopyVANT(original *models.VANT) *models.VANT {
	// Criar uma cópia do VANT
	copyVANT := &models.VANT{
		ID: original.ID,
		X:  original.X,
		Y:  original.Y,
		Z:  original.Z,
	}

	// Criar uma cópia do mapa MessagesBuffer (deep copy)
	copyVANT.MessagesBuffer = make(map[string][]int)
	for key, value := range original.MessagesBuffer {
		copyVANT.MessagesBuffer[key] = append([]int(nil), value...) // Copiar o slice
	}

	return copyVANT
}

func ProcessGenerator(qtdVants int, qtdFiles int) {
	/*err := os.RemoveAll("./data")
	if err != nil {
		return
	}

	os.Mkdir("./data", os.ModePerm)*/

	for i := 1; i <= qtdFiles; i++ {
		filePath := fmt.Sprintf("./data/scenario_%d_%d.csv", qtdVants, i)

		if _, err := os.Stat(filePath); err == nil {
			os.Remove(filePath)
		}
		file, _ := os.Create(filePath)
		defer file.Close()
		file.WriteString("id,x,y,z\n")
		for j := 1; j <= qtdVants; j++ {
			file.WriteString(fmt.Sprintf("%d,%.2f,%.2f,%.2f\n", j, rand.Float64()*100, rand.Float64()*100, rand.Float64()*100))
		}
	}
}
