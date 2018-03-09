package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

type HealthCheck struct {
	App      HealthCheckItem `json:app`
	Database HealthCheckItem `json:database`
	Env      HealthCheckItem `json:env`
	Redis    HealthCheckItem `json:redis`
	Site     HealthCheckItem `json:site`
}

type HealthCheckItem struct {
	Message string  `json:message`
	Success bool    `json:success`
	Time    float32 `json:time`
}

func main() {
	fmt.Println("Checking...\n")

	resp, err := http.Get("https://api.lendesk.com/_hc/all.json")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	var healthCheck HealthCheck
	json.Unmarshal(contents, &healthCheck)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"Service", "Message", "Success"})

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)

	table.SetColumnColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{},
		tablewriter.Colors{},
	)

	services := map[string]HealthCheckItem{
		"App":      healthCheck.App,
		"Database": healthCheck.Database,
		"Env":      healthCheck.Env,
		"Redis":    healthCheck.Redis,
		"Site":     healthCheck.Site,
	}

	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	for name, service := range services {
		successStr := green(strconv.FormatBool(service.Success))
		if service.Success == false {
			successStr = red(strconv.FormatBool(service.Success))
		}

		row := []string{name, service.Message, successStr}
		table.Append(row)
	}

	table.Render()
}
