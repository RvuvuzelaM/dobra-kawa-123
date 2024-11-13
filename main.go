package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sync"
	"time"
)

type Cowboy struct {
	Name   string `json:"name"`
	Health int    `json:"health"`
	Damage int    `json:"damage"`
}

type SafeCowboySlice struct {
	Cowboys []Cowboy
	mu      sync.Mutex
}

func (s *SafeCowboySlice) TakeDamage(idx int, amt int) {
	s.mu.Lock()
	s.Cowboys[idx].Health -= amt
	s.mu.Unlock()
}

func (s *SafeCowboySlice) GetHealth(idx int) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Cowboys[idx].Health
}

func (s *SafeCowboySlice) GetDamage(idx int) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Cowboys[idx].Damage
}

func (s *SafeCowboySlice) GetName(idx int) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Cowboys[idx].Name
}

func main() {
	jsonFile, err := os.Open("input.json")
	if err != nil {
		panic(fmt.Sprintf("failed to read %v", err))
	}
	defer jsonFile.Close()

	jsonCowboys, err := io.ReadAll(jsonFile)
	if err != nil {
		panic(fmt.Sprintf("failed to read cowboys from file %v", err))
	}

	var cowboys []Cowboy

	if err := json.Unmarshal(jsonCowboys, &cowboys); err != nil {
		panic(fmt.Sprintf("failed to unmarshal cowboys %v", err))
	}

	safeCowboys := &SafeCowboySlice{
		Cowboys: cowboys,
	}

	fmt.Println("Fight!!!")
	wg := sync.WaitGroup{}
	for i := range cowboys {
		wg.Add(1)
		go worker(i, safeCowboys, &wg)
	}
	wg.Wait()
}

func worker(idx int, safeCowboys *SafeCowboySlice, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		cowboys := safeCowboys.Cowboys
		name := safeCowboys.GetName(idx)
		health := safeCowboys.GetHealth(idx)

		if health <= 0 {
			fmt.Printf("Cowboy %s is dead\n", name)
			return
		}

		availableIdxs := []int{}
		for i := range cowboys {
			if safeCowboys.GetHealth(i) > 0 {
				availableIdxs = append(availableIdxs, i)
			}
		}

		// TOOD: extract to other function
		if len(availableIdxs) == 1 {
			fmt.Printf("%s wins\n", safeCowboys.GetName(availableIdxs[0]))
			return
		}

		// TODO: dirty hack
		newAvailableIdxs := []int{}
		for _, availIdx := range availableIdxs {
			if availIdx == idx {
				continue
			}
			newAvailableIdxs = append(newAvailableIdxs, availIdx)
		}

		damage := safeCowboys.GetDamage(idx)
		randomIdx := newAvailableIdxs[rand.Intn(len(newAvailableIdxs))]
		safeCowboys.TakeDamage(randomIdx, damage)
		fmt.Printf("%s hit %s on %d hits\n", name, safeCowboys.GetName(randomIdx), damage)
		time.Sleep(1 * time.Second)
	}
}
