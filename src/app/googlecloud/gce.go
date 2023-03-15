package googlecloud

import (
	"context"
	"fmt"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
	"github.com/pottava/gpu-node-manager/src/app/util"
)

func CreateVM(ctx context.Context, name, email, menu string) error {
	imgClient, err := compute.NewImagesRESTClient(ctx, clientOption())
	if err != nil {
		return err
	}
	defer imgClient.Close()

	image, err := imgClient.GetFromFamily(ctx, &computepb.GetFromFamilyImageRequest{
		Project: "debian-cloud",
		Family:  "debian-10",
	})
	if err != nil {
		return err
	}

	client, err := compute.NewInstancesRESTClient(ctx, clientOption())
	if err != nil {
		return err
	}
	defer client.Close()

	req := &computepb.InsertInstanceRequest{
		Project: util.ProjectID(),
		Zone:    util.Zone,
		InstanceResource: &computepb.Instance{
				Name: proto.String(name),
				Disks: []*computepb.AttachedDisk{{
						InitializeParams: &computepb.AttachedDiskInitializeParams{
								DiskSizeGb:  proto.Int64(300),
								SourceImage: image.SelfLink,
								DiskType:    proto.String(fmt.Sprintf("zones/%s/diskTypes/pd-standard", util.Zone)),
						},
						AutoDelete: proto.Bool(true),
						Boot:       proto.Bool(true),
						Type:       proto.String(computepb.AttachedDisk_PERSISTENT.String()),
				}},
				NetworkInterfaces: []*computepb.NetworkInterface{{
						Name: proto.String("global/networks/default"),
				}},
		},
	}
	switch menu {
	case "cpu-01": // Intel 2 vCPU
		req.InstanceResource.MachineType = proto.String(fmt.Sprintf("zones/%s/machineTypes/%s", util.Zone, "n2-standard-2"))

	case "t4-01": // NVIDIA T4 1 基 + Intel 2 vCPU
		req.InstanceResource.MachineType = proto.String(fmt.Sprintf("zones/%s/machineTypes/%s", util.Zone, "n1-standard-2"))
		// runtime.VirtualMachine.VirtualMachineConfig.AcceleratorConfig = &notebookspb.RuntimeAcceleratorConfig{
		// 	Type:      notebookspb.RuntimeAcceleratorConfig_NVIDIA_TESLA_T4,
		// 	CoreCount: 1,
		// }
	}
	op, err := client.Insert(ctx, req)
	if err != nil {
		return err
	}
	if err = op.Wait(ctx); err != nil {
		return err
	}
	return nil
}

func StartVM(ctx context.Context, name string) error {
	client, err := compute.NewInstancesRESTClient(ctx, clientOption())
	if err != nil {
		return err
	}
	defer client.Close()

	req := &computepb.StartInstanceRequest{
		Project:  util.ProjectID(),
		Zone:     util.Zone,
		Instance: name,
	}
	op, err := client.Start(ctx, req)
	if err != nil {
		return err
	}
	if err = op.Wait(ctx); err != nil {
		return err
	}
	return nil
}

func StopVM(ctx context.Context, name string) error {
	client, err := compute.NewInstancesRESTClient(ctx, clientOption())
	if err != nil {
		return err
	}
	defer client.Close()

	req := &computepb.StopInstanceRequest{
		Project:  util.ProjectID(),
		Zone:     util.Zone,
		Instance: name,
	}
	op, err := client.Stop(ctx, req)
	if err != nil {
		return err
	}
	if err = op.Wait(ctx); err != nil {
		return err
	}
	return nil
}

func DeleteVM(ctx context.Context, name string) error {
	client, err := compute.NewInstancesRESTClient(ctx, clientOption())
	if err != nil {
		return err
	}
	defer client.Close()

	req := &computepb.DeleteInstanceRequest{
		Project:  util.ProjectID(),
		Zone:     util.Zone,
		Instance: name,
	}
	op, err := client.Delete(ctx, req)
	if err != nil {
		return err
	}
	if err = op.Wait(ctx); err != nil {
		return err
	}
	return nil
}
