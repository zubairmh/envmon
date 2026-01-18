package cli

import "fmt"

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
