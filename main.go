package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/robfig/cron"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s cron_spec cmd [ args ... ]\n", os.Args[0])
	}

	args := os.Args[1:]
	spec := args[0]
	program := args[1]
	programArgs := args[2:]

	c := cron.New()
	err := c.AddFunc(spec, func() {
		cmd := exec.Command(program, programArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Println(err)
		}
	})

	if err != nil {
		log.Fatalln("failed to add func")
	}

	c.Start()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-sig:
		log.Print("received signal, shutting down...")
		c.Stop()
	}
}
