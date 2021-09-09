package main

import (
	"encoding/json"
	"fmt"
	"github.com/hristo-ganekov-sumup/genGraph/internal/sg"
	"github.com/hristo-ganekov-sumup/genGraph/internal/tfstate"
	"os"
)

var sgMap map[string]string

func main() {
	sgMap = make(map[string]string)
	state, err := tfstate.ParseTerraformStateFile("live.tfstate")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, resource := range state.Resources {
		if resource.Type == "aws_security_group" && resource.Mode == "managed" && resource.Name == "autosg" {
			//Populating the SG Map
			for _, instance := range resource.Instances {
				awsSg := &sg.AwsSecurityGroup{}
				_ = json.Unmarshal(instance.AttributesRaw, awsSg)
				sgMap[awsSg.Id] = awsSg.Name
			}

			//Printing out the rules
			for _, instance := range resource.Instances {
				awsSg := &sg.AwsSecurityGroup{}
				_ = json.Unmarshal(instance.AttributesRaw, awsSg)
				fmt.Printf("%s\n", awsSg.Name)

				for _, egress := range awsSg.Egress {
					if len(egress.SecurityGroups) > 0 {
						var tempSgs []string
						for _,sg := range egress.SecurityGroups {
							if val, ok := sgMap[sg]; ok{
								tempSgs = append(tempSgs, val )
							} else {
								tempSgs = append(tempSgs, sg )
							}
						}
						fmt.Printf("\t->%s [%d-%d] #%s \n", tempSgs, egress.FromPort, egress.ToPort, egress.Description)
					}
					if len(egress.CidrBlocks) > 0 {
						fmt.Printf("\t->%s [%d-%d] #%s \n", egress.CidrBlocks, egress.FromPort, egress.ToPort, egress.Description)
					}
				}

				for _, ingress := range awsSg.Ingress {
					if len(ingress.SecurityGroups) > 0 {
						var tempSgs []string
						for _,sg := range ingress.SecurityGroups {
							if val, ok := sgMap[sg]; ok {
								tempSgs = append(tempSgs, val)
							} else {
								tempSgs = append(tempSgs, sg)
							}
						}
						fmt.Printf("\t<-%s [%d-%d] #%s \n", tempSgs, ingress.FromPort, ingress.ToPort, ingress.Description)
					}
					if len(ingress.CidrBlocks) > 0 {
						fmt.Printf("\t<-%s [%d-%d] #%s \n", ingress.CidrBlocks, ingress.FromPort, ingress.ToPort, ingress.Description)
					}
				}

			}
		}
	}
}
