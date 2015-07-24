package farm

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/codegangsta/cli"
	"github.com/crackcomm/go-clitable"
	"github.com/gofarmer/utils/config"
	"time"
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

func Instances(c *cli.Context) {
	name := c.Args().First()

	if name == "" {
		fmt.Println("This command requires a farm name argument")
		return
	}

	fc := config.GetFarm(name)
	fmt.Println(fc.Region)
	svc := ec2.New(&aws.Config{Region: fc.Region})

	// Sample
	resp, err := svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("> Number of reservation sets: ", len(resp.Reservations))
	for idx, res := range resp.Reservations {
		fmt.Println("  > Number of instances: ", len(res.Instances))
		for _, inst := range resp.Reservations[idx].Instances {
			fmt.Println("    - Instance ID: ", *inst.InstanceID)
		}
	}
}

func List(c *cli.Context) {
	farmConfigs := config.GetFarms()

	t := clitable.New([]string{
		"Name",
		"Created At",
		"Region",
		"Status",
		"AMI",
		"Quotas",
	})

	for _, f := range farmConfigs {
		t.AddRow(map[string]interface{}{
			"Name":       f.Name,
			"Created At": time.Unix(f.CreatedAt, 0),
			"Region":     f.Region,
			"Status":     statuses[f.Status],
			"AMI":        f.AMI,
			"Quotas":     f.Quotas.Merge().String(),
		})
	}
	t.Print()
}

func Create(c *cli.Context) {
	name := c.Args().First()
	ami := c.Args().Get(1)
	region := c.Args().Get(2)

	if name == "" {
		fmt.Println("This command requires a farm name argument")
		return
	}

	existFarm := config.GetFarm(name)
	if existFarm != nil {
		fmt.Println("Farm with this name is already exist")
	}

	if region == "" {
		region = config.GetGlobal().DefaultRegion
	}

	fc := config.FarmConfig{
		Name:      name,
		Status:    0,
		Region:    region,
		CreatedAt: int64(time.Now().Unix()),
		AMI:       ami,
		Quotas:    farmQuotas,
	}

	config.CreateFarm(name, fc)
	/*
		1. Check and prepare internal dirs
		2. Check and read global config
		3. Check cli input if auth info does not exist in global config
		4. Create new dir with name from input
		5. Save typical farm config in new dir from step 4
	*/
}
