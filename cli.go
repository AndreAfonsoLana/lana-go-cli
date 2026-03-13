package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: You need to provide the project name.")
		fmt.Println("Correct use: go run cli.go <name-new-project> <name-git>")
		os.Exit(1)

	}

	projectName := os.Args[1]
	fmt.Printf("Starting the project creation: %s ...\n", projectName)

	directories := []string{
		filepath.Join(projectName, "cmd", projectName),
		filepath.Join(projectName, "api"),
		filepath.Join(projectName, "configs"),
		filepath.Join(projectName, "internal"),
		filepath.Join(projectName, "pkg"),
		filepath.Join(projectName, "test"),
	}

	moduleName := projectName

	if len(os.Args) >= 3 {
		gitName := os.Args[2]
		fmt.Printf("📁 Git name detected: %s\n", gitName)
		moduleName = fmt.Sprintf("github.com/%s/%s", gitName, projectName)
	} else {
		fmt.Println("📁 Git name not provided. The module will only use the project name.")
	}

	fmt.Printf("Iniciando a criação do projeto: %s...\n", projectName)
	for _, dir := range directories {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Printf("Create error %s: %v\n", dir, err)
			os.Exit(1)
		}
		fmt.Printf("📁 Created: %s/\n", dir)
	}

	mainFilePath := filepath.Join(projectName, "cmd", projectName, "main.go")

	templateCode := `package main
	import "fmt"

	func main() {
		fmt.Println("Welcome ` + projectName + `!")
	}
	`
	err := os.WriteFile(mainFilePath, []byte(templateCode), 0644)
	if err != nil {
		fmt.Printf("Error creating the main.go file: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("\n✅ Project structure successfully created.!")

	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Dir = projectName
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error running go mod init: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Project Go successfully initialized!")

	mainPath := filepath.Join("cmd", projectName, "main.go")
	cmd = exec.Command("go", "run", mainPath)

	cmd.Dir = projectName
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error running : %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Project Go Executing!")
}
