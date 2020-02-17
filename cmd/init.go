package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"
	"github.com/urfave/cli"
	"zerodha.tech/kubekutr/models"
	"zerodha.tech/kubekutr/utils"
)

// InitProject initializes git repo and copies a sample config
func (hub *Hub) InitProject(config models.Config) cli.Command {
	return cli.Command{
		Name:    "init",
		Aliases: []string{"i"},
		Usage:   "Initializes an empty git repo with a kubekutr config file.",
		Action:  hub.init,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "default, d",
				Usage: "Use the default config file",
			},
		},
	}
}

func (hub *Hub) init(cliCtx *cli.Context) error {
	// Initialize git repository
	cmd := exec.Command("git", "init")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error while initializing git repo: %v", err)
	}
	var configFile []byte
	if cliCtx.Bool("default") {
		// Copy sample config to local directory
		configFile, err = hub.Fs.Read("templates/config.sample.yml")
		if err != nil {
			return fmt.Errorf("error while copying sample config: %v", err)
		}
		// create kubekutr config
		f, err := os.Create("kubekutr.yml")
		if err != nil {
			return fmt.Errorf("error while creating sample config: %v", err)
		}
		_, err = f.Write(configFile)
		if err != nil {
			return fmt.Errorf("error while copying sample config: %v", err)
		}
	} else {
		workloadsLen := 0
		// workloads := []models.Workload{}
		resources := []models.Resource{}
		survey.AskOne(&survey.Input{
			Message: "How many workloads do you want to deploy?",
			Help:    "Workloads represent different application names.",
			Default: "1",
		}, &workloadsLen)
		// Iterate for all workloads
		for i := 0; i < workloadsLen; i++ {
			var (
				wd             = models.Workload{}
				deploymentsLen = 0
			)
			survey.AskOne(&survey.Input{
				Message: "What's the name of the application?",
			}, &wd.Name, survey.WithValidator(survey.Required))
			survey.AskOne(&survey.Input{
				Message: "How many deployments do you want to configure?",
				Help:    "Deployments represent different components of your application.",
				Default: "1",
			}, &deploymentsLen)
			for j := 0; j < deploymentsLen; j++ {
				var (
					dep           = models.Deployment{}
					containersLen = 0
				)
				survey.AskOne(&survey.Input{
					Message: "What's the name of deployment?",
					Help:    "Name of the deployment to be configured.",
					Default: "1",
				}, &dep.Name, survey.WithValidator(survey.Required))
				survey.AskOne(&survey.Input{
					Message: "How many containers do you want to configure?",
					Help:    "Containers form a part of a Pod. Typically each Pod is assosicated with one container.",
					Default: "1",
				}, &containersLen)
				for k := 0; k < containersLen; k++ {
					var ctr = models.Container{}
					// perform the questions
					err := survey.Ask(containerQs, &ctr)
					if err != nil {
						return fmt.Errorf("fake error: %v", err.Error())
					}
					dep.Containers = append(dep.Containers, ctr)
				}
				// wd.Deployments = append(wd.Deployments, dep)
				resources = append(resources, dep)
			}
			// create resources from the prompt config
			utils.CreateGitopsDirectory("/tmp/", wd.Name)
			err := prepareResources(resources, "/tmp/", wd.Name, hub.Fs)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
