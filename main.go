package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Load struct {
	id      int
	pickup  [2]float64
	dropoff [2]float64
}

type Driver struct {
	loads        []Load
	totalMinutes float64
}

const (
	maxMinuts     = 12.0 * 60.0
	initialTemp   = 10000.0
	coolingRate   = 0.999995
	maxIterations = 100000
	costPerDrivr  = 1000.0
)

var distanceMatix map[[2][2]float64]float64

func EuclideanDistance(p1, p2 [2]float64) float64 {
	if dist, exists := distanceMatix[[2][2]float64{p1, p2}]; exists {
		return dist
	}
	dist := math.Sqrt(math.Pow(p2[0]-p1[0], 2) + math.Pow(p2[1]-p1[1], 2))
	distanceMatix[[2][2]float64{p1, p2}] = dist
	distanceMatix[[2][2]float64{p2, p1}] = dist
	return dist
}

func TotalTime(driver Driver) float64 {
	var totalTime float64
	depot := [2]float64{0, 0}
	if len(driver.loads) > 0 {
		totalTime += EuclideanDistance(depot, driver.loads[0].pickup)
	}
	for i := 0; i < len(driver.loads); i++ {
		totalTime += EuclideanDistance(driver.loads[i].pickup, driver.loads[i].dropoff)
		if i+1 < len(driver.loads) {
			totalTime += EuclideanDistance(driver.loads[i].dropoff, driver.loads[i+1].pickup)
		}
	}
	if len(driver.loads) > 0 {
		totalTime += EuclideanDistance(driver.loads[len(driver.loads)-1].dropoff, depot)
	}
	return totalTime
}

func ParseInput(filePath string) ([]Load, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lodes []Load
	scanner := bufio.NewScanner(file)
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		id, _ := strconv.Atoi(fields[0])
		pickupCoords := strings.Trim(fields[1], "()")
		dropoffCoords := strings.Trim(fields[2], "()")
		pickupParts := strings.Split(pickupCoords, ",")
		dropoffParts := strings.Split(dropoffCoords, ",")

		pickupX, _ := strconv.ParseFloat(pickupParts[0], 64)
		pickupY, _ := strconv.ParseFloat(pickupParts[1], 64)
		dropoffX, _ := strconv.ParseFloat(dropoffParts[0], 64)
		dropoffY, _ := strconv.ParseFloat(dropoffParts[1], 64)

		lode := Load{id, [2]float64{pickupX, pickupY}, [2]float64{dropoffX, dropoffY}}
		lodes = append(lodes, lode)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lodes, nil
}

func AssignLoads(lodes []Load) []Driver {
	var drivers []Driver
	sort.Slice(lodes, func(i, j int) bool {
		return EuclideanDistance([2]float64{0, 0}, lodes[i].pickup) < EuclideanDistance([2]float64{0, 0}, lodes[j].pickup)
	})
	loadAssigned := make(map[int]bool)
	for _, lode := range lodes {
		assigned := false
		for i := range drivers {
			if loadAssigned[lode.id] {
				continue
			}
			drivers[i].loads = append(drivers[i].loads, lode)
			newTotalMinuts := TotalTime(drivers[i])
			if newTotalMinuts <= maxMinuts {
				drivers[i].totalMinutes = newTotalMinuts
				loadAssigned[lode.id] = true
				assigned = true
				break
			} else {
				drivers[i].loads = drivers[i].loads[:len(drivers[i].loads)-1]
			}
		}
		if !assigned {
			newDriver := Driver{loads: []Load{lode}, totalMinutes: TotalTime(Driver{loads: []Load{lode}})}
			drivers = append(drivers, newDriver)
			loadAssigned[lode.id] = true
		}
	}
	return drivers
}

func CopyDrivers(drivers []Driver) []Driver {
	newDrivers := make([]Driver, len(drivers))
	for i, driver := range drivers {
		newLoads := make([]Load, len(driver.loads))
		copy(newLoads, driver.loads)
		newDrivers[i] = Driver{loads: newLoads, totalMinutes: driver.totalMinutes}
	}
	return newDrivers
}

func CalculateTotalCost(drivers []Driver) float64 {
	var totalDrivenMinuts float64
	for _, driver := range drivers {
		totalDrivenMinuts += driver.totalMinutes
	}
	return costPerDrivr*float64(len(drivers)) + totalDrivenMinuts
}

func OptimizeSolution(drivers []Driver) []Driver {
	currentSolution := CopyDrivers(drivers)
	bestSolution := CopyDrivers(currentSolution)
	bestCost := CalculateTotalCost(bestSolution)
	rand.New(rand.NewSource(time.Now().UnixNano()))
	temp := initialTemp
	for temp > 1 {
		newSolution := CopyDrivers(currentSolution)
		driverId1 := rand.Intn(len(newSolution))
		driverId2 := rand.Intn(len(newSolution))
		if len(newSolution[driverId1].loads) == 0 || len(newSolution[driverId2].loads) == 0 {
			continue
		}
		lodeIdx1 := rand.Intn(len(newSolution[driverId1].loads))
		loadIdx2 := rand.Intn(len(newSolution[driverId2].loads))
		newSolution[driverId1].loads[lodeIdx1], newSolution[driverId2].loads[loadIdx2] = newSolution[driverId2].loads[loadIdx2], newSolution[driverId1].loads[lodeIdx1]
		newSolution[driverId1].totalMinutes = TotalTime(newSolution[driverId1])
		newSolution[driverId2].totalMinutes = TotalTime(newSolution[driverId2])
		if newSolution[driverId1].totalMinutes <= maxMinuts && newSolution[driverId2].totalMinutes <= maxMinuts {
			newCost := CalculateTotalCost(newSolution)
			if newCost < bestCost {
				bestSolution = CopyDrivers(newSolution)
				bestCost = newCost
			}
			if newCost < CalculateTotalCost(currentSolution) || math.Exp((CalculateTotalCost(currentSolution)-newCost)/temp) > rand.Float64() {
				currentSolution = CopyDrivers(newSolution)
			}
		}
		temp *= coolingRate
	}
	return bestSolution
}

func PrintSolution(drivers []Driver) {
	for _, driver := range drivers {
		fmt.Print("[")
		for i, lode := range driver.loads {
			if i > 0 {
				fmt.Print(",")
			}
			fmt.Print(lode.id)
		}
		fmt.Println("]")
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide the input file path")
		return
	}
	filePath := os.Args[1]
	lodes, err := ParseInput(filePath)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}
	distanceMatix = make(map[[2][2]float64]float64)
	drivers := AssignLoads(lodes)
	drivers = OptimizeSolution(drivers)
	PrintSolution(drivers)
}
