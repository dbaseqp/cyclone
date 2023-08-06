package main

import (
	// "errors"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	// "strings"
	"context"
	"net/url"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"

	"bruharmy/models"
)

var (
	availablePortGroups = &models.RWPortGroupMap{}
	vDSName             string
)

func LoadPortGroups() error {
	availablePortGroups = &models.RWPortGroupMap{
		Data: make(map[int]string),
	}

	vDSName = "MainDSW"

	ctx := context.Background()
	u, err := soap.ParseURL(tomlConf.VCenterURL)
	if err != nil {
		fmt.Printf("Error parsing vCenter URL: %s\n", err)
		os.Exit(1)
	}

	u.User = url.UserPassword(tomlConf.VCenterUsername, tomlConf.VCenterPassword)

	client, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		fmt.Printf("Error creating vSphere client: %s\n", err)
		os.Exit(1)
	}
	defer client.Logout(ctx)

	finder := find.NewFinder(client.Client, true)

	// Find datacenter in the vSphere environment
	dc, err := finder.Datacenter(ctx, tomlConf.Datacenter)
	if err != nil {
		fmt.Printf("Error finding datacenter: %s\n", err)
		os.Exit(1)
	}

	finder.SetDatacenter(dc)

	// Find all distributed virtual switches in the datacenter
	podNetworks, err := finder.NetworkList(ctx, "*_PodNetwork")

	if err != nil {
		fmt.Printf("Error listing networks: %s\n", err)
		os.Exit(1)
	}

	// Collect found DistributedVirtualPortgroup refs
	var refs []types.ManagedObjectReference
	for _, pgRef := range podNetworks {
		refs = append(refs, pgRef.Reference())
	}

	pc := property.DefaultCollector(client.Client)

	// Collect property from references list
	var pgs []mo.DistributedVirtualPortgroup
	err = pc.Retrieve(ctx, refs, []string{"name"}, &pgs)
	if err != nil {
		fmt.Printf("Error collecting references: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Port groups for vDS '%s':\n", vDSName)
	for _, pg := range pgs {
		r, _ := regexp.Compile("^\\d+")
		match := r.FindString(pg.Name)
		pgNumber, _ := strconv.Atoi(match)
		if pgNumber >= tomlConf.StartingPortGroup && pgNumber < tomlConf.EndingPortGroup {
			availablePortGroups.Data[pgNumber] = pg.Name
			fmt.Printf("%d\n", pgNumber)
		}
	}
	log.Printf("Found %d port groups within on-demand DistributedPortGroup range: %d - %d", len(availablePortGroups.Data), tomlConf.StartingPortGroup, tomlConf.EndingPortGroup)
	return nil
}

func TemplateGuestView() ([]string, error) {
	var templates []string

	ctx := context.Background()
	u, err := soap.ParseURL(tomlConf.VCenterURL)
	if err != nil {
		fmt.Printf("Error parsing vCenter URL: %s\n", err)
		os.Exit(1)
	}

	u.User = url.UserPassword(tomlConf.VCenterUsername, tomlConf.VCenterPassword)

	client, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		fmt.Printf("Error creating vSphere client: %s\n", err)
		os.Exit(1)
	}
	defer client.Logout(ctx)

	finder := find.NewFinder(client.Client, true)

	// Find datacenter in the vSphere environment
	dc, err := finder.Datacenter(ctx, tomlConf.Datacenter)
	if err != nil {
		fmt.Printf("Error finding datacenter: %s\n", err)
		os.Exit(1)
	}

	finder.SetDatacenter(dc)

	templateResourcePool, err := finder.ResourcePool(ctx, tomlConf.TemplateResourcePool)

	if err != nil {
		fmt.Printf("Error finding guest templates: %s\n", err)
	}

	var trp mo.ResourcePool
	err = templateResourcePool.Properties(ctx, templateResourcePool.Reference(), []string{"resourcePool"}, &trp)
	if err != nil {
		fmt.Printf("Error getting child resource pools: %s\n", err)
		os.Exit(1)
	}

	pc := property.DefaultCollector(client.Client)

	var rps []mo.ResourcePool
	err = pc.Retrieve(ctx, trp.ResourcePool, []string{"name"}, &rps)
	if err != nil {
		fmt.Printf("Error collecting references: %s\n", err)
		os.Exit(1)
	}

	for _, rp := range rps {
		templates = append(templates, rp.Name)
	}

	return templates, nil
}

func CloneOnDemand(data models.InvokeCloneOnDemandForm, username string) (string, error) {
	var nextAvailablePortGroup string
	availablePortGroups.Mu.Lock()
	for i := tomlConf.StartingPortGroup; i < tomlConf.EndingPortGroup; i++ {
		if _, exists := availablePortGroups.Data[i]; !exists {
			nextAvailablePortGroup = strconv.Itoa(i)
			availablePortGroups.Data[i] = fmt.Sprintf("%s_PodNetwork", nextAvailablePortGroup)
			break
		}
	}
	availablePortGroups.Mu.Unlock()
	cmd := exec.Command("powershell", ".\\pwsh\\cloneondemand.ps1", data.Template, username, nextAvailablePortGroup, tomlConf.TargetResourcePool, tomlConf.Domain, tomlConf.WanPortGroup)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return "", err
	}

	fmt.Println("bruh ", out.String(), stderr.String())

	return nextAvailablePortGroup, nil
}

func DeletePod(data models.DeletePodForm, username string) error {
	cmd := exec.Command("powershell", ".\\pwsh\\deletepod.ps1", username, data.Target)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	availablePortGroups.Mu.Lock()
	deleted_pg, _ := strconv.Atoi(strings.Split(data.Target, "_")[0])
	delete(availablePortGroups.Data, deleted_pg)
	availablePortGroups.Mu.Unlock()

	return nil
}
