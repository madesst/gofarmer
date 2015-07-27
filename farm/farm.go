package farm

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/awsutil"
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

//Move instances methods to instance.go
func Instances(c *cli.Context) {
	instances, fc := instancesByNameArg(c)

	t := clitable.New([]string{
		"ID",
		"Type",
		"State",
		"IP",
		"AMI",
		"Spot Request ID",
	})

	for _, inst := range instances {
		t.AddRow(map[string]interface{}{
			"ID":              *inst.InstanceID,
			"Type":            *inst.InstanceType,
			"State":           *inst.State.Name,
			"IP":              *inst.PrivateIPAddress,
			"AMI":             *inst.ImageID,
			"Spot Request ID": inst.SpotInstanceRequestID,
		})
	}
	t.Print()
}

func Start(c *cli.Context) {
	instances, fc := instancesByNameArg(c)
	for _, inst := range instances {
		svc := ec2.New(&aws.Config{Region: fc.Region})
		if *inst.State.Name == "running" {
			fmt.Println(fmt.Sprintln("	>>>", *inst.InstanceID, "already running"))
			continue
		}

		params := &ec2.StartInstancesInput{
			InstanceIDs: []*string{
				aws.String(*inst.InstanceID),
			},
		}
		if res, resp := checkResponse(svc.StartInstances(params)); !res {
			fmt.Println(resp)
			continue
		}

		fmt.Println(fmt.Sprintln("	>>>", *inst.InstanceID, "started successfully"))
	}

	fc.Status = 1
	cleanupLocks(inst, fc)
}

func Stop(c *cli.Context) {
	instances, fc := instancesByNameArg(c)
	for _, inst := range instances {
		svc := ec2.New(&aws.Config{Region: fc.Region})
		if *inst.State.Name == "stopped" {
			fmt.Println(fmt.Sprintln("	>>>", *inst.InstanceID, "already stopped"))
			continue
		}

		params := &ec2.StopInstancesInput{
			InstanceIDs: []*string{
				aws.String(*inst.InstanceID),
			},
		}
		if res, resp := checkResponse(svc.StopInstances(params)); !res {
			fmt.Println(fmt.Sprintln("	>>>", resp))
			continue
		}

		fmt.Println(fmt.Sprintln("	>>>", *inst.InstanceID, "stopped successfully"))
		saveInstanceLock(inst)
	}

	fc.Status = 0
	cleanupLocks(inst, fc)
}

func List(c *cli.Context) {
	farmConfigs := config.GetFarms()

	t := clitable.New([]string{
		"Name",
		"Created At",
		"Region",
		"Status",
		"AMI",
		"Last Updated At",
		"Quotas",
	})

	for _, f := range farmConfigs {
		t.AddRow(map[string]interface{}{
			"Name":       f.Name,
			"Created At": time.Unix(f.CreatedAt, 0),
			"Region":     f.Region,
			"Status":     statuses[f.Status],
			"AMI":        f.AMI,
			"Last Updated At": 0
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

	config.CreateFarm(fc)
}

func instancesByNameArg(c *cli.Context) ([]*ec2.Instance, config.FarmConfig) {
	name := c.Args().First()

	if name == "" {
		panic("This command requires a farm name argument")
	}

	return describeInstances(name)
}

func describeInstances(name string) ([]*ec2.Instance, config.FarmConfig) {
	fc := config.GetFarm(name)
	if fc == nil {
		panic("Undefined farm \"" + name + "\"")
	}
	svc := ec2.New(&aws.Config{Region: fc.Region})

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:gofarmer"),
				Values: []*string{
					aws.String(fc.Name),
				},
			},
		},
	}
	resp, err := svc.DescribeInstances(params)
	checkResponse(resp, err)

	if resp.Reservations[0] == nil {
		return []*ec2.Instance{}, *fc
	}

	for _, inst := range resp.Reservations[0].Instances {
		saveInstanceLock(inst)
	}
	cleanupLocks(resp.Reservations[0].Instances, *fc)

	return resp.Reservations[0].Instances, *fc
}

func saveInstanceLock(inst *ec2.Instance) {
	/**
	* 1. Prepare instance status struct
	* 2. Call config.SaveInstanceLock
	**/
}

func cleanupLocks(instances []*ec2.Instance, config.FarmConfig) {
	//cleanupLocks(resp.Reservations[0].Instances)
	config.SaveFarmConfig(fc)
}

func checkResponse(resp interface{}, err error) (bool, string) {
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			return false, fmt.Sprintln(awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
		}
	}
	// Pretty-print the response data.
	rawResponse := awsutil.StringValue(resp)
	return true, rawResponse
}
