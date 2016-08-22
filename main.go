package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/google/go-github/github"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
)

func main() {
	var token string
	var command string
	var organization string

	app := cli.NewApp()
	app.Name = "do-with-gh-user-keys"
  app.Version = "0.1"
	app.Usage = "Runs a program for each users ssh key in an organization"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "token, t",
			Value:       "<unset>",
			Usage:       "A GitHub personal access token to be used for authentication",
			Destination: &token,
			EnvVar:      "DOW_GH_TOKEN",
		},
		cli.StringFlag{
			Name:        "command, c",
			Value:       "",
			Usage:       "Command that will reveive the piped key(s)",
			Destination: &command,
		},
		cli.BoolFlag{
			Name:  "userarg, u",
			Usage: "Add github username as last argument to command",
		},
		cli.StringSliceFlag{
			Name:  "args, a",
			Usage: "Args to command that will reveive the piped key(s)",
		},
		cli.StringFlag{
			Name:        "organization, o",
			Value:       "",
			Usage:       "List member keys in this organization",
			Destination: &organization,
		},
	}

	app.Action = func(c *cli.Context) error {
		if !c.IsSet("command") {
			return cli.NewExitError("Please specify command", 1)
		}
		if !c.IsSet("token") {
			return cli.NewExitError("Please specify token", 1)
		}
		if !c.IsSet("organization") {
			return cli.NewExitError("Please specify organization", 1)
		}

		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(oauth2.NoContext, ts)

		client := github.NewClient(tc)
		opt := &github.ListMembersOptions{
			ListOptions: github.ListOptions{PerPage: 10},
		}
		var allUsers []github.User
		for {
			users, resp, err := client.Organizations.ListMembers(organization, opt)
			if err != nil {
				log.Fatal(err)
			}
			allUsers = append(allUsers, users...)
			if resp.NextPage == 0 {
				break
			}
			opt.ListOptions.Page = resp.NextPage
		}
		for _, u := range allUsers {
			keys, _, err := client.Users.ListKeys(*u.Login, nil)
			if err != nil {
				log.Fatal(err)
			}
			cargs := make([]string, 1)
			if c.Bool("userarg") {
				cargs = append(cargs, c.StringSlice("args")...)
        cargs = append(cargs, *u.Login)
			} else {
				cargs = append(cargs, c.StringSlice("args")...)
			}
			for _, key := range keys {
				cmd := exec.Command(command, cargs...)
				stdin, err := cmd.StdinPipe()
				stdin.Write([]byte(*key.Key))
				stdin.Close()
				if err != nil {
					log.Fatal(err)
				}
				data, err := cmd.CombinedOutput()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(string(data))
			}
		}
		return nil
	}
	app.Run(os.Args)
}
