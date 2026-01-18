package cli

import (
	"fmt"

	"envmon/internal/env"
)

func SwitchDeployment(deployment string) error {
	content, err := env.Process(deployment)
	if err != nil {
		return err
	}

	if err := env.Write(content); err != nil {
		return fmt.Errorf("could not write .env file: %w", err)
	}

	fmt.Printf("Switched to [%s] deployment\n", deployment)
	return nil
}

func ShowConfigs() error {
	configs, err := env.GetConfigs()
	if err != nil {
		return err
	}

	if len(configs) == 0 {
		fmt.Println("No deployment configurations detected")
		return nil
	}

	fmt.Println("Detected deployment configurations:")
	for _, cfg := range configs {
		fmt.Printf("  - %s\n", cfg)
	}
	return nil
}

func ShowCurrentDeployment() error {
	deployment, err := env.GetCurrentDeployment()
	if err != nil {
		return err
	}

	if deployment == "" {
		fmt.Println("not active")
	} else {
		fmt.Printf("Current deployment: %s\n", deployment)
	}
	return nil
}

func ShowHelp() {
	help := `envmon - Environment file deployment manager

Usage:
  envmon <deployment>    Switch to specified deployment configuration
  envmon                 Show current active deployment
  envmon configs         List all detected deployment configurations  
  envmon help            Show this help message

How it works:
  envmon reads the .env file in your current working directory and looks
  for deployment flags in comments. It uncomments lines matching the 
  specified flag while commenting out lines for other flags.

Flag Syntax:
  # [staging]            Single deployment flag
  # [staging/dev]        Multiple flags (active for both)
  # [~prod]              Inverse (active for all EXCEPT prod)
  # [~(staging/dev)]     Inverse multiple (active except staging and dev)

Example:
  # [staging] MongoDB config
  #MONGO_URI=staging-uri
  
  # [dev] MongoDB config  
  MONGO_URI=dev-uri

  Running 'envmon staging' will uncomment staging lines and comment dev lines.
`
	fmt.Print(help)
}
