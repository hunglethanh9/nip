package cl

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"

	"bufio"

	"strings"

	"github.com/mswift42/nip/tv"
	"github.com/urfave/cli"
)

func extractIndex(c *cli.Context) (int, error) {
	if len(c.Args()) < 1 {
		fmt.Println("Please enter valid index number.")
	}
	ind := c.Args().Get(0)
	index, err := strconv.ParseInt(ind, 10, 64)
	if err != nil {
		return 0, err
	}
	return int(index), nil
}

// InitCli loads the ProgrammeDB into memory
// and sets up the command line commands.
func InitCli() *cli.App {
	dbpath := tv.DBPath()
	db, err := tv.RestoreProgrammeDB(dbpath + tv.NipDB)
	if err != nil {
		panic(err)
	}
	app := cli.NewApp()
	app.Setup()
	app.Name = "nip"
	app.Version = "0.0.1"
	app.Copyright = "2018 (c) Martin Haesler"
	app.Usage = "search for iplayer tv programmes."
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
COPYRIGHT:
   {{.Copyright}}
VERSION:
   {{.Version}}
`
	app.Commands = []cli.Command{
		{
			Name:     "list",
			Aliases:  []string{"l"},
			Usage:    "list all available categories.",
			HelpName: "list",
			Action: func(c *cli.Context) error {
				fmt.Println(db.ListAvailableCategories())
				return nil
			},
		},
		{
			Name:      "category",
			Aliases:   []string{"c"},
			Usage:     "list all programmes for a category.",
			HelpName:  "category",
			ArgsUsage: "[index]",
			Action: func(c *cli.Context) error {
				fmt.Println(db.ListCategory(c.Args().Get(0)))
				return nil
			},
		},
		{
			Name:      "search",
			Aliases:   []string{"s"},
			Usage:     "search for a programme.",
			HelpName:  "search",
			ArgsUsage: "[searchterm]",
			Action: func(c *cli.Context) error {
				fmt.Println(db.FindTitle(c.Args().Get(0)))
				return nil
			},
		},
		{
			Name:      "show",
			Aliases:   []string{"sh"},
			Usage:     "open a programme's homepage.",
			HelpName:  "show",
			ArgsUsage: "[index]",
			Action: func(c *cli.Context) error {
				index, err := extractIndex(c)
				if err != nil {
					fmt.Println(err)
					return nil
				}
				url, err := db.FindURL(index)
				if err != nil {
					fmt.Println(err)
				}
				switch runtime.GOOS {
				case "linux":
					err = exec.Command("xdg-open", url).Start()
				case "darwin":
					err = exec.Command("open", url).Start()
				case "windows":
					err = exec.Command("cmd", "/c", "start", url).Start()
				default:
					fmt.Println("Unsupported platform.")
				}
				if err != nil {
					fmt.Println(err)
				}
				return nil
			},
		},
		{
			Name:      "synopsis",
			Aliases:   []string{"syn"},
			Usage:     "print programme's synopsis",
			HelpName:  "synopsis",
			ArgsUsage: "[index]",
			Action: func(c *cli.Context) error {
				index, err := extractIndex(c)
				if err != nil {
					fmt.Println(err)
					return nil
				}
				prog, err := db.FindProgramme(index)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(prog.String() + "\n" + prog.Synopsis)
				}
				return nil
			},
		},
		{
			Name:      "links",
			Aliases:   []string{"lnk"},
			Usage:     "show related links for a programme",
			HelpName:  "links",
			ArgsUsage: "[index]",
			Action: func(c *cli.Context) error {
				index, err := extractIndex(c)
				if err != nil {
					fmt.Println("Please enter valid index number.")
				}
				rl, err := db.FindRelatedLinks(index)
				if err != nil {
					fmt.Println(err)
				}
				if len(rl) == 0 {
					fmt.Println("Sorry, no related links were found.")
				} else {
					for _, i := range rl {
						fmt.Println(i.Title, " : ", i.URL)
					}
				}
				return nil
			},
		},
		{
			Name:      "download",
			Aliases:   []string{"g", "d", "get"},
			Usage:     "use youtube-dl to download programme with index n",
			UsageText: "download programmes with index n. If no format is specified,\nbest available format is used.",
			HelpName:  "download",
			ArgsUsage: "[index] [format]",
			Action: func(c *cli.Context) error {
				var format string
				if len(c.Args()) == 2 {
					format = c.Args().Get(1)
				}
				ind, err := extractIndex(c)
				if err != nil {
					fmt.Println("Please enter valid index number.")
					return nil
				}
				prog, err := db.FindProgramme(ind)
				if err != nil {
					fmt.Println(err)
					return nil
				}
				fmt.Println("Downloading Programme \n", prog.String())
				u := tv.BBCPrefix + prog.URL
				var cmd *exec.Cmd
				switch runtime.GOOS {
				case "windows":
					if format != "" {
						cmd = exec.Command("cmd", "/c", "youtube-dl -f "+format+" "+u)
					} else {
						cmd = exec.Command("cmd", "/c", "youtube-dl -f best "+u)
					}
				default:
					if format != "" {
						cmd = exec.Command("/bin/sh", "-c", "youtube-dl -f "+format+" "+u)
					} else {
						cmd = exec.Command("/bin/sh", "-c", "youtube-dl -f best "+u)
					}
				}
				outpipe, err := cmd.StdoutPipe()
				if err != nil {
					fmt.Println(err)
				}
				err = cmd.Start()
				if err != nil {
					fmt.Println(err)
				}
				scanner := bufio.NewScanner(outpipe)
				scanner.Split(bufio.ScanRunes)
				var target string
				for scanner.Scan() {
					fmt.Print(scanner.Text())
					target += scanner.Text()
				}
				err = cmd.Wait()
				if err != nil {
					fmt.Println(err)
				}
				split := strings.Split(target, "\n")
				for _, i := range split {
					if strings.Contains(i, "Destination:") {
						path := tv.DBPath()
						db.MarkSaved(path + i[24:])
					}
				}
				return nil
			},
		},
		{
			Name:      "formats",
			Aliases:   []string{"f"},
			Usage:     "list youtube-dl formats for programme with index n",
			HelpName:  "formats",
			ArgsUsage: "[index]",
			Action: func(c *cli.Context) error {
				ind, err := extractIndex(c)
				if err != nil {
					fmt.Println("Please enter valid index number.")
					return nil
				}
				prog, err := db.FindProgramme(ind)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("Listing Formats for Programme \n", prog.String())
				u := tv.BBCPrefix + prog.URL
				var cmd *exec.Cmd
				switch runtime.GOOS {
				case "windows":
					cmd = exec.Command("cmd", "/c", "youtube-dl -F "+u)
				default:
					cmd = exec.Command("/bin/sh", "-c", "youtube-dl -F "+u)
				}
				if err != nil {
					fmt.Println(err)
				}
				outpipe, err := cmd.StdoutPipe()
				if err != nil {
					fmt.Println(err)
				}
				err = cmd.Start()
				if err != nil {
					fmt.Println(err)
				}
				scanner := bufio.NewScanner(outpipe)
				for scanner.Scan() {
					fmt.Println(scanner.Text())
				}
				err = cmd.Wait()
				if err != nil {
					fmt.Println(err)
				}
				return nil
			},
		},
		{
			Name:      "refresh",
			Aliases:   []string{"r"},
			Usage:     "refresh programme db",
			HelpName:  "refresh",
			ArgsUsage: " ",
			Action: func(c *cli.Context) error {
				path := tv.DBPath()
				tv.RefreshDB(path + tv.NipDB)
				return nil
			},
		},
	}
	return app

}
