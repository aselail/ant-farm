package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Room struct {
	name string //Room Stucter
	x, y int    // coordinates
}

type AntFarm struct {
	numAnts    int                 //Number of ants
	rooms      map[string]Room     //map of rooms
	tunnels    map[string][]string //map of rooms connected
	start, end string
}

func parseInput(file string) (AntFarm, error) {
	f, err := os.Open(file)
	if err != nil {
		return AntFarm{}, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f) //scan to read the file line by line
	var numAnts int
	var rooms = make(map[string]Room)
	var tunnels = make(map[string][]string)
	var start, end string
	phase := 0 // 0 for ants, 1 for rooms, 2 for tunnels

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			if line == "##start" {
				phase = 1
				start = ""
			} else if line == "##end" {
				phase = 1
				end = ""
			}
			continue
		}

		if phase == 0 {
			fmt.Sscanf(line, "%d", &numAnts)
			phase = 1
		} else if phase == 1 {
			parts := strings.Fields(line)
			if len(parts) == 3 {
				x, y := 0, 0
				fmt.Sscanf(parts[1], "%d", &x)
				fmt.Sscanf(parts[2], "%d", &y)
				room := Room{name: parts[0], x: x, y: y}
				rooms[parts[0]] = room
				if start == "" {
					start = parts[0]
				} else if end == "" {
					end = parts[0]
				}
			} else {
				phase = 2
			}
		}
		if phase == 2 {
			parts := strings.Split(line, "-")
			if len(parts) == 2 {
				tunnels[parts[0]] = append(tunnels[parts[0]], parts[1])
				tunnels[parts[1]] = append(tunnels[parts[1]], parts[0])
			} else {
				return AntFarm{}, errors.New("ERROR: invalid data format")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return AntFarm{}, err
	}
	if start == "" || end == "" {
		return AntFarm{}, errors.New("ERROR: invalid data format, no start or end room found")
	}

	return AntFarm{
		numAnts: numAnts,
		rooms:   rooms,
		tunnels: tunnels,
		start:   start,
		end:     end,
	}, nil
}

func bfs(farm AntFarm) []string {
	// BFS to find the shortest path
	var queue [][]string
	var visited = make(map[string]bool)
	queue = append(queue, []string{farm.start})
	visited[farm.start] = true

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		room := path[len(path)-1]

		if room == farm.end {
			return path
		}

		for _, neighbor := range farm.tunnels[room] {
			if !visited[neighbor] {
				newPath := make([]string, len(path))
				copy(newPath, path)
				newPath = append(newPath, neighbor)
				queue = append(queue, newPath)
				visited[neighbor] = true
			}
		}
	}

	return nil
}

func simulateAnts(farm AntFarm, path []string) {
	antPositions := make([]string, farm.numAnts)
	for i := 0; i < farm.numAnts; i++ {
		antPositions[i] = farm.start
	}

	moves := []string{}
	for len(moves) < farm.numAnts {
		var move []string
		for i := range antPositions {
			if antPositions[i] != farm.end {
				nextRoom := ""
				for j := 0; j < len(path)-1; j++ {
					if path[j] == antPositions[i] {
						nextRoom = path[j+1]
						break
					}
				}
				if nextRoom != "" {
					antPositions[i] = nextRoom
					move = append(move, fmt.Sprintf("L%d-%s", i+1, nextRoom))
				}
			}
		}
		if len(move) > 0 {
			fmt.Println(strings.Join(move, " "))
			moves = append(moves, move...)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <input_file>")
		return
	}

	inputFile := os.Args[1]
	farm, err := parseInput(inputFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(farm.numAnts)
	for _, room := range farm.rooms {
		fmt.Printf("%s %d %d\n", room.name, room.x, room.y)
	}
	for room, neighbors := range farm.tunnels {
		for _, neighbor := range neighbors {
			fmt.Printf("%s-%s\n", room, neighbor)
		}
	}

	path := bfs(farm)
	if path == nil {
		fmt.Println("ERROR: no path from start to end")
		return
	}

	simulateAnts(farm, path)
}
