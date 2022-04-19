package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	// By default, the logger outputs to stderr, but this may not be ideal...
	log.SetOutput(os.Stdout)
}

var (
	fileArg = flag.String("file", "", "file-name")
)

func main() {
	flag.Parse()

	c := cron.New()

	// Read crontab file line by line
	file, err := os.Open(*fileArg)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	secondAccuracyRegex := regexp.MustCompile(`^(.+?)\s+(.+?)\s+(.+?)\s+(.+?)\s+(.+?)\s+(.+?)\s+(.+?)\s+(.+)$`)
	minuteAccuracyRegex := regexp.MustCompile(`^(.+?)\s+(.+?)\s+(.+?)\s+(.+?)\s+(.+?)\s+(.+?)\s+(.+)$`)

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " ")

		if !strings.HasPrefix(line, "#") {
			timing := ""
			cmdToRun := ""

			if secondAccuracyRegex.MatchString(line) {
				resSec := secondAccuracyRegex.FindAllStringSubmatch(line, -1)
				for i := range resSec {
					//fmt.Printf("year: %s, month: %s, day: %s\n", res[i][1], res[i][2], res[i][3])
					timing = resSec[i][1] + " " + resSec[i][2] + " " + resSec[i][3] + " " + resSec[i][4] + " " + resSec[i][5] + " " + resSec[i][6]
					cmdToRun = resSec[i][8]
				}
			} else if minuteAccuracyRegex.MatchString(line) {
				res := minuteAccuracyRegex.FindAllStringSubmatch(line, -1)
				for i := range res {
					//fmt.Printf("year: %s, month: %s, day: %s\n", res[i][1], res[i][2], res[i][3])
					timing = "0 " + res[i][1] + " " + res[i][2] + " " + res[i][3] + " " + res[i][4] + " " + res[i][5]
					cmdToRun = res[i][7]
				}
			}

			if len(timing) > 0 && len(cmdToRun) > 0 {
				// We don't need to change the stdout and stderr streams like cron does
				cmdToRun = strings.ReplaceAll(cmdToRun, " > /proc/1/fd/1 2>/proc/1/fd/2", "")

				c.AddFunc(timing, func() { runShellCommand(cmdToRun) })
			}

		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Start cron with one scheduled job
	log.Info("Starting go-cron with " + strconv.Itoa(len(c.Entries())) + " jobs...")
	c.Start()
	//printCronEntries(c.Entries())
	//time.Sleep(2 * time.Minute)

	for {
		// Endless loop to keep the app running...
		time.Sleep(time.Minute)
	}

	// Funcs can also be added while running
	// log.Info("Add new job to a running cron")
	// entryID2, _ := c.AddFunc("*/2 * * * *", func() { log.Info("[Job 2]Every two minutes job\n") })
	// showCronEntries(c.Entries())

	// Or funcs can be removed while running
	// log.Info("Remove Job2 and add new Job2 with schedule run every minute")
	// c.Remove(entryID2)
	// c.AddFunc("*/1 * * * *", func() { log.Info("[Job 2]Every one minute job\n") })
}

func runShellCommand(bashCmd string) {
	//fmt.Fprintf(os.Stdout, "Running: "+bashCmd+"\n")
	parts := strings.Split(bashCmd, " ")
	bashCmdFirstPart := parts[0]
	args := parts[1:]

	// TODO ENV vars are not interpreted. Example: "echo $HOME" will output this exactly.
	cmd := exec.Command(bashCmdFirstPart, args...)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	cmd.Start()

	stdoutScanner := bufio.NewScanner(stdout)
	stderrScanner := bufio.NewScanner(stderr)

	for stdoutScanner.Scan() {
		stdoutMessage := stdoutScanner.Text()
		fmt.Fprintf(os.Stdout, stdoutMessage+"\n")
	}

	for stderrScanner.Scan() {
		stderrMessage := stderrScanner.Text()
		fmt.Fprintf(os.Stderr, stderrMessage+"\n")
	}

	cmd.Wait()
}

//
//func printCronEntries(cronEntries []*cron.Entry) {
//	log.Infof("Cron Info: %+v\n", cronEntries)
//}
