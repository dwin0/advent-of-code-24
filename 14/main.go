package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
)

type Velocity struct {
	tilesPerSecondX int
	tilesPerSecondY int
}

type Position struct {
	x int
	y int
}

type Robot struct {
	position Position
	velocity Velocity
}

type Config struct {
	inputPath       string
	tilesXDirection int
	tilesYDirection int
}

// Easter Bunny Headquater
type EBHQ struct {
	robots []Robot
}

func main() {
	// config := Config{
	// 	inputPath:       "./input-small.txt",
	// 	tilesXDirection: 11,
	// 	tilesYDirection: 7,
	// }
	config := Config{ // Part 2 does not work with the small input
		inputPath:       "./input.txt",
		tilesXDirection: 101,
		tilesYDirection: 103,
	}

	ebhq := EBHQ{
		robots: []Robot{},
	}

	err := ebhq.readListOfRobots(config)
	if err != nil {
		log.Fatal(err)
	}

	// Part 1
	// seconds := 100
	// ebhq.elapseSecond(config, seconds)
	// safetyFactor := ebhq.calculateSafetyFactor(config)
	// fmt.Println(safetyFactor)

	// Part 2
	seconds := 0
	for {
		amountOfCloseRobots := ebhq.robotsNextToEachOther()

		// random threshold
		if amountOfCloseRobots > 350 {
			fmt.Println(ebhq.createRobotsMap(config))
			break
		}
		seconds++
		ebhq.elapseSecond(config, 1)
	}

	fmt.Println("seconds:", seconds)
}

// Part 2
// https://www.reddit.com/r/adventofcode/comments/1hee2nr/day_14_part_2_losing_my_mind_need_help/
func (e *EBHQ) robotsNextToEachOther() int {
	sum := 0
	for _, robot1 := range e.robots {
		for _, robot2 := range e.robots {
			if robot1 == robot2 {
				continue
			}

			if math.Abs(float64(robot1.position.x-robot2.position.x))+math.Abs(float64(robot1.position.y-robot2.position.y)) == 1 {
				sum++
			}
		}
	}

	return sum
}

// Part 1
// func (e *EBHQ) calculateSafetyFactor(config Config) int {
// 	robotsInQuadrants := []int{0, 0, 0, 0}

// 	for _, robot := range e.robots {
// 		if robot.position.x < config.tilesXDirection/2 && robot.position.y < config.tilesYDirection/2 {
// 			robotsInQuadrants[0]++
// 		} else if robot.position.x > config.tilesXDirection/2 && robot.position.y < config.tilesYDirection/2 {
// 			robotsInQuadrants[1]++
// 		} else if robot.position.x < config.tilesXDirection/2 && robot.position.y > config.tilesYDirection/2 {
// 			robotsInQuadrants[2]++
// 		} else if robot.position.x > config.tilesXDirection/2 && robot.position.y > config.tilesYDirection/2 {
// 			robotsInQuadrants[3]++
// 		}
// 	}

// 	return robotsInQuadrants[0] * robotsInQuadrants[1] * robotsInQuadrants[2] * robotsInQuadrants[3]
// }

// for debugging
func (e *EBHQ) createRobotsMap(config Config) [][]string {
	robotsMap := make([][]int, config.tilesYDirection)
	for i := range robotsMap {
		robotsMap[i] = make([]int, config.tilesXDirection)
	}

	for _, robot := range e.robots {
		robotsMap[robot.position.y][robot.position.x]++
	}

	robotsMapString := make([][]string, len(robotsMap))
	for i := range robotsMap {
		robotsMapString[i] = make([]string, len(robotsMap[i]))
		for j := range robotsMap[i] {
			if robotsMap[i][j] > 0 {
				robotsMapString[i][j] = "#" //strconv.Itoa(robotsMap[i][j])
			} else {
				robotsMapString[i][j] = "."
			}
		}
	}

	return robotsMapString
}

func (e *EBHQ) elapseSecond(config Config, seconds int) {
	for i := range e.robots {
		e.robots[i].position.x = (e.robots[i].position.x + seconds*(e.robots[i].velocity.tilesPerSecondX+config.tilesXDirection)) % config.tilesXDirection
		e.robots[i].position.y = (e.robots[i].position.y + seconds*(e.robots[i].velocity.tilesPerSecondY+config.tilesYDirection)) % config.tilesYDirection
	}
}

func (e *EBHQ) readListOfRobots(config Config) error {
	file, err := os.Open(config.inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)

	for sc.Scan() {
		line := sc.Text()
		px, py, vx, vy, err := parseInputLine(line)

		if err != nil {
			return err
		}

		e.robots = append(e.robots, Robot{
			position: Position{
				x: px,
				y: py,
			},
			velocity: Velocity{
				tilesPerSecondX: vx,
				tilesPerSecondY: vy,
			},
		})
	}

	if err := sc.Err(); err != nil {
		return err
	}

	return nil
}

func parseInputLine(line string) (px, py, vx, vy int, err error) {
	regex, regexErr := regexp.Compile(`p=(-?\d*),(-?\d*) v=(-?\d*),(-?\d*)`)
	if regexErr != nil {
		return 0, 0, 0, 0, regexErr
	}

	stringValues := regex.FindStringSubmatch(line)

	// index 0 is the full string match (to ignore)
	px, conversionErr := strconv.Atoi(stringValues[1])
	if conversionErr != nil {
		return 0, 0, 0, 0, conversionErr
	}
	py, conversionErr = strconv.Atoi(stringValues[2])
	if conversionErr != nil {
		return 0, 0, 0, 0, conversionErr
	}
	vx, conversionErr = strconv.Atoi(stringValues[3])
	if conversionErr != nil {
		return 0, 0, 0, 0, conversionErr
	}
	vy, conversionErr = strconv.Atoi(stringValues[4])
	if conversionErr != nil {
		return 0, 0, 0, 0, conversionErr
	}

	return px, py, vx, vy, nil
}
