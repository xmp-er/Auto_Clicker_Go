package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-vgo/robotgo"
	"github.com/xmp-er/Auto_Clicker_Go/helpers"
	"github.com/xmp-er/Auto_Clicker_Go/models"
	"github.com/xmp-er/Auto_Clicker_Go/validators"
)

func main() {
	//taking the positions from the command line
	var x_position int
	var y_position int
	var interval int

	x_default, y_default := robotgo.Location() //get the default position of the mouse

	flag.IntVar(&x_position, "x", x_default, "x-coordinate of the mouse")
	flag.IntVar(&y_position, "y", y_default, "y-coordinate of the mouse")
	flag.IntVar(&interval, "t", 5, "time interval for the mouse to click")

	flag.Parse()

	fmt.Println("stop: stop the program")
	fmt.Println("stop after <time> <unit>: stop the program after a given time interval")
	fmt.Println("sec <time>: change the time interval on which clicks are performed")
	fmt.Println("timer <time> <unit>(sec/min/hrs/days): set a timer for the clicks to start again")

	sigShutDown := make(chan os.Signal, 1)
	signal.Notify(sigShutDown, os.Interrupt, syscall.SIGTERM)
	co_ordinates := make(chan models.Coordinates, 1)
	tempWait := make(chan models.TempWait, 1)

	scanner := bufio.NewScanner(os.Stdin)
	ctx, ctxCancel := context.WithCancel(context.Background())

	go func() {
		helpers.Click_on_interval(x_position, y_position, interval, sigShutDown, co_ordinates, tempWait, ctx, ctxCancel)
	}()

	go func() {
		for {
			select {
			case <-sigShutDown:
				sigShutDown <- os.Kill
				return
			default:
				var input string = ""
				if scanner.Scan() {
					input = scanner.Text()
				}
				str := strings.Split(input, " ")
				switch {
				case str[0] == "stop":
					if len(str) == 1 {
						sigShutDown <- os.Kill
						return
					}
					if len(str) == 4 {
						if str[1] == "after" && validators.IsInt(str[2]) && validators.IsTimeUnit(str[3]) {
							val := helpers.ConvertToInt(str[2])
							switch str[3] {
							case "sec":
								tempWait <- models.TempWait{SleepVal: val, IsKill: true}
							case "min":
								tempWait <- models.TempWait{SleepVal: val * 60, IsKill: true}
							case "hrs":
								tempWait <- models.TempWait{SleepVal: val * 3600, IsKill: true}
							case "days":
								tempWait <- models.TempWait{SleepVal: val * 86400, IsKill: true}
							}
						}
						continue
					}
					fmt.Println("Please enter valid input")
				case str[0] == "timer":
					if !validators.IsArgs(str, 3) || !validators.IsInt(str[1]) {
						fmt.Println("Please enter valid input")
					}
					val := helpers.ConvertToInt(str[1])
					switch str[2] {
					case "sec":
						tempWait <- models.TempWait{SleepVal: val, IsKill: false}
					case "min":
						tempWait <- models.TempWait{SleepVal: val * 60, IsKill: false}
					case "hrs":
						tempWait <- models.TempWait{SleepVal: val * 3600, IsKill: false}
					case "days":
						tempWait <- models.TempWait{SleepVal: val * 86400, IsKill: false}
					}
				case str[0] == "sec":
					if !validators.IsArgs(str, 2) || !validators.IsInt(str[1]) {
						fmt.Println("Please enter valid input")
						continue
					} else {
						co_ordinates <- models.Coordinates{IsFollowMouse: false, Interval: helpers.ConvertToInt(str[1])}
					}
				default:
					select {
					case <-sigShutDown:
						sigShutDown <- os.Kill
						return
					default:
						fmt.Println("Please enter valid input")
					}
				}
			}
		}
	}()
	<-sigShutDown
	fmt.Println("Closing the program")

}
