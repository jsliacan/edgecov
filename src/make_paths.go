package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// vertices numbered 0,...,n-1
// each slice graph[i] is the list of adjacent vertices to vertex i
var graph = [][]int{}

func main() {

	err := loadGraph("../dat/graph.txt")
	if err != nil {
		fmt.Errorf("Could not load the graph.\n")
	}
	for vi := range graph {
		fmt.Println(graph[vi])
	}
}

// load graph from a file
// line i is space delim list of vertices adjacent to vertex i
func loadGraph(filename string) error {

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Couldn't open file %s.\n", filename)
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		ln, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		line := string(ln)
		if len(line) == 0 {
			graph = append(graph, nil)
			continue
		}
		edges_str := strings.Split(line, ",")
		edges := []int{}

		for i := range edges_str {
			e, _ := strconv.Atoi(edges_str[i])
			edges = append(edges, e)
		}
		graph = append(graph, edges)
	}

	return nil
}
