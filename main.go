package main

import (
	"fmt"

	"github.com/georgebent/go-restorer/pkg/core"
	"github.com/georgebent/go-restorer/pkg/file_manager"
	"github.com/georgebent/go-restorer/pkg/io_manager"
	"github.com/georgebent/go-restorer/pkg/runner"
)

func main() {
	fmt.Println("Hello")

	options := map[string]string{}
	options["1"] = "Save"
	options["2"] = "Restore"
	options["3"] = "List"

	chosen := io_manager.Ask("Choose action", options)

	switch chosen {
	case "1":
		err := runner.Save()
		if err != nil {
			fmt.Println("Got error: ", err)
		}
	case "2":
		err := runner.Restore()
		if err != nil {
			fmt.Println("Got error: ", err)
		}
	case "3":
		backups := core.GetEnv("BACKUP_DIR")
		archives, err := file_manager.ListBackups(backups)
		if err != nil {
			fmt.Println("Got error: ", err)
		}

		for _, archive := range archives {
			println(archive)
		}
	}
}
