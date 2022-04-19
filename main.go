package main

import (
	"bufio"
	"fmt"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
	"time"
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
}

func main() {
	log.Info("Creating new cron...")

	c := cron.New()

	// TODO parse our crontab file and load it here. Be careful that here we have second accuracy, but our crontab only has minute accuracy
	c.AddFunc("*/5 * * * *", func() { runShellCommand("date") })

	// Start cron with one scheduled job
	log.Info("Starting cron...")
	c.Start()
	//printCronEntries(c.Entries())
	//time.Sleep(2 * time.Minute)

	for {
		// Endless loop to keep the app running...
		time.Sleep(time.Minute)
	}

	// Funcs may also be added to a running Cron
	//log.Info("Add new job to a running cron")
	//entryID2, _ := c.AddFunc("*/2 * * * *", func() { log.Info("[Job 2]Every two minutes job\n") })
	//printCronEntries(c.Entries())
	//time.Sleep(5 * time.Minute)

	//Remove Job2 and add new Job2 that run every 1 minute
	//log.Info("Remove Job2 and add new Job2 with schedule run every minute")
	//c.Remove(entryID2)
	//c.AddFunc("*/1 * * * *", func() { log.Info("[Job 2]Every one minute job\n") })
	//time.Sleep(5 * time.Minute)

}

func runShellCommand(bashCmd string) {
	bashCmdFirstPart := "php"
	args := "-v"

	cmd := exec.Command(bashCmdFirstPart, strings.Split(args, " ")...)

	// TODO What about std?
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()

	cmd.Start()

	// TODO is this the best way to pass the stderr output from the bash command out of this Go app?
	scanner := bufio.NewScanner(stderr)
	//scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}

	// TODO is this the best way to pass the stdout output from the bash command out of this Go app?
	// TODO ideally, we want to differentiate stdout vs stderr output!
	scanner2 := bufio.NewScanner(stdout)
	//scanner2.Split(bufio.ScanWords)
	for scanner2.Scan() {
		m2 := scanner2.Text()
		fmt.Println(m2)
	}
	cmd.Wait()
}

//
//func printCronEntries(cronEntries []*cron.Entry) {
//	log.Infof("Cron Info: %+v\n", cronEntries)
//}
