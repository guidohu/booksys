package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"server/mynautique"
)

// ErrorResponse is an error reply from the API.
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse is a successful reply from the API.
type SuccessResponse struct {
	Success   bool   `json:"success"`
	AuthUntil string `json:"auth_until"`
}

func printJSON(v any) {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling JSON: %v\n", err)
		return
	}
	fmt.Println(string(bytes))
}

func printError(msg string) {
	printJSON(ErrorResponse{Error: msg})
}

func main() {
	userFlag := flag.String("user", "", "myNautique username/email")
	passwordFlag := flag.String("password", "", "myNautique password")
	idFlag := flag.String("id", "", "Boat ID for GetBoatTelemetry")
	apiKeyFlag := flag.String("apikey", "", "Firebase Auth API Key (optional)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <function>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Functions:\n  Login\n  GetFleet\n  GetBoatTelemetry\n\nFlags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *userFlag == "" {
		printError("username is required, specify with --user")
		flag.Usage()
		os.Exit(1)
	}

	if *passwordFlag == "" {
		printError("password is required, specify with --password")
		flag.Usage()
		os.Exit(1)
	}

	if flag.NArg() < 1 {
		printError("function to test is required (choose from: Login, GetFleet, GetBoatTelemetry)")
		flag.Usage()
		os.Exit(1)
	}

	function := flag.Arg(0)

	apiKey := *apiKeyFlag
	if apiKey == "" {
		apiKey = os.Getenv("BOOKSYS_MYNAUTIQUE_API_KEY")
	}
	if apiKey == "" {
		apiKey = "AIzaSyAb8S3Owvo8k-gI8eK_DEztFOn0FcZFRxw"
	}

	client := mynautique.NewClient(&mynautique.Options{
		User:       *userFlag,
		Password:   *passwordFlag,
		AuthAPIKey: apiKey,
	})

	switch function {
	case "Login":
		err := client.Login()
		if err != nil {
			printError(fmt.Sprintf("Login failed: %v", err))
			os.Exit(1)
		}
		printJSON(SuccessResponse{
			Success:   true,
			AuthUntil: client.AuthUntil.String(),
		})

	case "GetFleet":
		err := client.GetFleet()
		if err != nil {
			printError(fmt.Sprintf("GetFleet failed: %v", err))
			os.Exit(1)
		}
		printJSON(client.Fleet)

	case "GetBoatTelemetry":
		if *idFlag == "" {
			printError("boat ID is required for GetBoatTelemetry, specify with --id")
			os.Exit(1)
		}
		boatID, err := strconv.ParseInt(*idFlag, 10, 64)
		if err != nil {
			printError(fmt.Sprintf("invalid boat ID %q: must be an integer", *idFlag))
			os.Exit(1)
		}
		telemetry, err := client.GetBoatTelemetry(boatID)
		if err != nil {
			printError(fmt.Sprintf("GetBoatTelemetry failed: %v", err))
			os.Exit(1)
		}
		printJSON(telemetry)

	default:
		printError(fmt.Sprintf("unknown function %q. Allowed functions: Login, GetFleet, GetBoatTelemetry", function))
		flag.Usage()
		os.Exit(1)
	}
}
