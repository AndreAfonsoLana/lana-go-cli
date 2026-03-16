package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type LanaConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

var configProjetCli = "lana-cli.json"
var version = "2026.1"

func main() {

	commands := []string{
		"new-project",
		"np",
		"new-entity",
		"ne",
	}

	if len(os.Args) < 2 {
		fmt.Println("Error: You need to provide a command.")
		fmt.Println("Available commands: new-project (np), new-entity (ne)")
		os.Exit(1)
	}

	command := strings.ToLower(os.Args[1])

	isValidCommand := false
	for _, c := range commands {
		if command == c {
			isValidCommand = true
			break
		}
	}

	if !isValidCommand {
		fmt.Printf("Error: '%s' is not a valid command.\n", command)
		fmt.Println("Available commands: new-project (np), new-entity (ne)")
		os.Exit(1)
	}

	// ==========================================
	// COMAND: NEW-PROJECT
	// ==========================================
	if command == "new-project" || command == "np" {
		if len(os.Args) < 3 {
			fmt.Println("Error: You need to provide the project name.")
			fmt.Println("Correct use: go run cli.go np <name-new-project> <name-git>")
			os.Exit(1)
		}

		projectName := os.Args[2]
		fmt.Printf("Starting the project creation: %s ...\n", projectName)

		directories := []string{
			filepath.Join(projectName, "cmd", projectName),
			filepath.Join(projectName, "api"),
			filepath.Join(projectName, "configs"),
			filepath.Join(projectName, "internal", "entity"),
			filepath.Join(projectName, "pkg"),
			filepath.Join(projectName, "test"),
		}

		moduleName := projectName

		if len(os.Args) >= 4 {
			gitName := os.Args[3]
			fmt.Printf("📁 Git name detected: %s\n", gitName)
			moduleName = fmt.Sprintf("github.com/%s/%s", gitName, projectName)
		} else {
			fmt.Println("📁 Git name not provided. The module will only use the project name.")
		}

		for _, dir := range directories {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				fmt.Printf("Create error %s: %v\n", dir, err)
				os.Exit(1)
			}
			fmt.Printf("📁 Created: %s/\n", dir)
		}

		lanaCliFilePath := filepath.Join(configProjetCli)
		templateFileProject := fmt.Sprintf(`
{
	"name" : "%s",
	"version" : "%s"
}`, projectName, version)

		err := os.WriteFile(lanaCliFilePath, []byte(templateFileProject), 0644)
		if err != nil {
			fmt.Printf("Error creating the lana-cli.json file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("📁 lana-cli.json created!")

		mainFilePath := filepath.Join(projectName, "cmd", projectName, "main.go")
		templateCode := `package main
import "fmt"

func main() {
	fmt.Println("Welcome ` + projectName + `!")
}
`
		err = os.WriteFile(mainFilePath, []byte(templateCode), 0644)
		if err != nil {
			fmt.Printf("Error creating the main.go file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("📁 Main created!")

		configsFilePath := filepath.Join(projectName, "configs", "config.go")
		templateCode = `package configs

var cfg *conf
type conf struct {
	DBDriver      string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	WebServerPort int
}

func LoadConfig(path string) (*conf, error) {
	// ...
	return cfg, nil
}
`
		err = os.WriteFile(configsFilePath, []byte(templateCode), 0644)
		if err != nil {
			fmt.Printf("Error creating the config.go file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("📁 Configs created!")

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
			fmt.Printf("Error running project: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Project Go Executing!")

		fmt.Printf("\n🎉 Project '%s' created successfully!\n", projectName)
		fmt.Println("👉 Next steps:")
		fmt.Printf("   cd %s\n", projectName)
		fmt.Println("   go mod tidy")
	}

	// ==========================================
	// COMAND: NEW-ENTITY
	// ==========================================
	if command == "new-entity" || command == "ne" {
		configsProject, err := ReadLanaCliConf()
		if err != nil {
			fmt.Println("❌ Error: lana-cli.json not found.")
			fmt.Println("👉 Make sure you are inside a project folder to create an entity!")
			os.Exit(1)
		}

		projectName := configsProject.Name

		if len(os.Args) < 3 {
			fmt.Println("Error: You need to provide the entity name.")
			fmt.Println("Correct use: go run cli.go ne <entity-name>")
			os.Exit(1)
		}

		entityName := os.Args[2]
		fmt.Printf("Starting the entity creation: %s ...\n", entityName)

		entityDir := filepath.Join(projectName, "internal", "entity")
		err = os.MkdirAll(entityDir, 0755)
		if err != nil {
			fmt.Printf("Error creating entity directory: %v\n", err)
			os.Exit(1)
		}

		entityFilePath := filepath.Join(entityDir, strings.ToLower(entityName)+".go")

		templateCode := fmt.Sprintf(`package entity
import (
	"errors"
	"time"
)

type %s struct {
	Name string %cjson:"name"%c
	CreatedAt time.Time %cjson:"created_at"%c
}

func NewProduct(name string, price int) (*%s, error) {

}
func (p *%s) Validate() error {

}
`, entityName, '`', '`', '`', '`', entityName, entityName)

		err = os.WriteFile(entityFilePath, []byte(templateCode), 0644)
		if err != nil {
			fmt.Printf("Error creating the entity file: %v\n", err)
			os.Exit(1)
		}

		entityFilePath = filepath.Join(entityDir, strings.ToLower(entityName)+"_test.go")

		templateCode = fmt.Sprintf(`package entity
import (
	"testing"

	"github.com/stretchr/testify/assert"
)
func TestNewUse(t *testing.T) {

}
	`)
		err = os.WriteFile(entityFilePath, []byte(templateCode), 0644)
		if err != nil {
			fmt.Printf("Error creating the entity file: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("📁 " + entityName + " created successfully!")

	}
}

func ReadLanaCliConf() (LanaConfig, error) {
	fileBytes, err := os.ReadFile(configProjetCli)
	if err != nil {
		fmt.Println("❌ Error: lana-cli.json not found.")
		fmt.Println("👉 Make sure you are inside a project folder to create an entity!")
		os.Exit(1)
		return LanaConfig{}, err
	}

	var config LanaConfig
	err = json.Unmarshal(fileBytes, &config)
	if err != nil {
		fmt.Printf("❌ Error reading JSON format: %v\n", err)
		os.Exit(1)
	}

	return config, nil
}
