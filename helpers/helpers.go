package helpers

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/xmp-er/Auto_Clicker_Go/models"
)

func Click_on_interval(x int, y int, interval int, sigShutDown chan os.Signal, coordinate_Shift chan models.Coordinates, TempWait chan models.TempWait, ctx context.Context, ctxCancel context.CancelFunc) {

	isFollowMouse := true
	wg := sync.WaitGroup{}
	for {
		select {
		case <-ctx.Done():
			return
		case <-sigShutDown: //shutdown signal from system
			sigShutDown <- os.Kill
			return
		case val := <-TempWait:
			if !val.IsKill {
				time.Sleep(time.Duration(val.SleepVal) * time.Second)
			}
			if val.IsKill {
				go func() {
					time.Sleep(time.Duration(val.SleepVal) * time.Second)
					ctxCancel()
					sigShutDown <- os.Kill
				}()
			}
		case val := <-coordinate_Shift:
			if val.IsFollowMouse {
				isFollowMouse = val.IsFollowMouse
			}
			if val.Interval != -1 {
				interval = val.Interval
			}
		default:
			wg.Add(1)
			go func() {
				defer wg.Done()
				if !isFollowMouse {
					robotgo.Move(x, y) //move the mouse to the given position
				}
				robotgo.Click("left")
				time.Sleep(time.Duration(interval) * time.Second)
			}()
			wg.Wait()
		}
	}
}

func ConvertToInt(str string) int {
	val, _ := strconv.Atoi(str)
	return val
}

func GetTimeValue(str string) int {
	switch str {
	case "min":
		return 60
	case "hrs":
		return 3600
	case "days":
		return 86400
	}
	return 1
}
