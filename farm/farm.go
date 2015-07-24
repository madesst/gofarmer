package farm

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/crackcomm/go-clitable"
	"github.com/gofarmer/utils/config"
)

var farmQuotas config.FarmQuotas = config.FarmQuotas{}
var statuses map[int]string = map[int]string{0: "Off", 1: "On"}

func CheckCredentialsConfig(c *cli.Context) error {
	farmQuotas.Quotas = config.Quotas{
		MaxAmount:    -1,
		MaxPrice:     -1,
		MinInstances: -1,
		MaxInstances: -1,
	}

	farmQuotas.FromGlobal = false
	if !c.IsSet("m-a") || !c.IsSet("m-p") || !c.IsSet("max-i") || !c.IsSet("min-i") {
		farmQuotas.FromGlobal = true
	}

	if c.IsSet("m-a") {
		farmQuotas.Quotas.MaxAmount = c.Float64("m-a")
	}
	if c.IsSet("m-p") {
		farmQuotas.Quotas.MaxPrice = c.Float64("m-p")
	}
	if c.IsSet("min-i") {
		farmQuotas.Quotas.MinInstances = c.Int("min-i")
	}
	if c.IsSet("max-i") {
		farmQuotas.Quotas.MaxInstances = c.Int("max-i")
	}

	return nil
}

func List(c *cli.Context) {
	farmConfigs := config.GetFarms()

	t := clitable.New([]string{
		"Name",
		"Created At",
		"Status",
		"AMI",
		"AWS Tag Name",
		"Quotas",
	})

	for _, f := range farmConfigs {
		t.AddRow(map[string]interface{}{
			"Name":         f.Name,
			"Created At":   f.Name,
			"Status":       statuses[f.Status],
			"AMI":          f.Name,
			"AWS Tag Name": f.AwsTagName,
			"Quotas":       f.Quotas.Merge().String(),
		})
	}
	t.Print()
}

func Create(c *cli.Context) {
	name := c.Args().First()

	if name == "" {
		fmt.Println("This command requires a farm name argument")
		return
	}

	existFarm := config.GetFarm(name)
	if existFarm != nil {
		fmt.Println("Farm with this name is already exist")
	}

	config.CreateFarm(name, farmQuotas)
	/*
		1. Check and prepare internal dirs
		2. Check and read global config
		3. Check cli input if auth info does not exist in global config
		4. Create new dir with name from input
		5. Save typical farm config in new dir from step 4
	*/
}
