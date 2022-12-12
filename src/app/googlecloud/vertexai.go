package googlecloud

import (
	"context"
	"fmt"

	notebooks "cloud.google.com/go/notebooks/apiv1"
	notebookspb "cloud.google.com/go/notebooks/apiv1/notebookspb"
	"github.com/pottava/gpu-node-manager/src/app/util"
)

func CreateNotebook(ctx context.Context, name, email, menu string) error {
	client, err := notebooks.NewNotebookClient(ctx, clientOption())
	if err != nil {
		return err
	}
	defer client.Close()

	// @see https://cloud.google.com/deep-learning-vm/docs/images
	// $ gcloud compute images list --project deeplearning-platform-release
	req := &notebookspb.CreateInstanceRequest{
		Parent:     fmt.Sprintf("projects/%s/locations/%s", util.ProjectID(), util.Zone),
		InstanceId: name,
		Instance: &notebookspb.Instance{
			Environment: &notebookspb.Instance_VmImage{
				VmImage: &notebookspb.VmImage{
					Project: "deeplearning-platform-release",
					Image: &notebookspb.VmImage_ImageFamily{
						ImageFamily: "tf-ent-2-10-cu113",
					},
				},
			},
			InstanceOwners:   []string{name},
			InstallGpuDriver: true,
			DataDiskType:     notebookspb.Instance_PD_SSD,
			DataDiskSizeGb:   300,
		},
	}
	switch menu {
	case "t4-01": // NVIDIA T4 1 基 + Intel 2 vCPU
		req.Instance.MachineType = "n1-standard-2"
		req.Instance.AcceleratorConfig = &notebookspb.Instance_AcceleratorConfig{
			Type:      notebookspb.Instance_NVIDIA_TESLA_T4,
			CoreCount: 1,
		}
	case "t4-02": //NVIDIA T4 1 基 + Intel 4 vCPU
		req.Instance.MachineType = "n1-highmem-4"
		req.Instance.AcceleratorConfig = &notebookspb.Instance_AcceleratorConfig{
			Type:      notebookspb.Instance_NVIDIA_TESLA_T4,
			CoreCount: 1,
		}
	case "a100-01": // NVIDIA A100 1 基 + Intel 12 vCPU
		req.Instance.MachineType = "a2-highgpu-1g"
		req.Instance.AcceleratorConfig = &notebookspb.Instance_AcceleratorConfig{
			Type:      notebookspb.Instance_NVIDIA_TESLA_A100,
			CoreCount: 1,
		}
	case "a100-02": // NVIDIA A100 4 基 + Intel 24 vCPU
		req.Instance.MachineType = "a2-highgpu-4g"
		req.Instance.AcceleratorConfig = &notebookspb.Instance_AcceleratorConfig{
			Type:      notebookspb.Instance_NVIDIA_TESLA_A100,
			CoreCount: 4,
		}
	}
	_, err = client.CreateInstance(ctx, req)
	return err
}

func CreateManagedNotebook(ctx context.Context, name, email, menu string) error {
	client, err := notebooks.NewManagedNotebookClient(ctx, clientOption())
	if err != nil {
		return err
	}
	defer client.Close()

	runtime := &notebookspb.Runtime_VirtualMachine{
		VirtualMachine: &notebookspb.VirtualMachine{
			VirtualMachineConfig: &notebookspb.VirtualMachineConfig{
				DataDisk: &notebookspb.LocalDisk{
					InitializeParams: &notebookspb.LocalDiskInitializeParams{
						DiskType:   notebookspb.LocalDiskInitializeParams_PD_SSD,
						DiskSizeGb: 300,
					},
				},
			},
		},
	}
	req := &notebookspb.CreateRuntimeRequest{
		Parent:    fmt.Sprintf("projects/%s/locations/%s", util.ProjectID(), util.Location),
		RuntimeId: name,
		Runtime: &notebookspb.Runtime{
			RuntimeType: runtime,
			SoftwareConfig: &notebookspb.RuntimeSoftwareConfig{
				InstallGpuDriver:    true,
				IdleShutdownTimeout: 30,
			},
			AccessConfig: &notebookspb.RuntimeAccessConfig{
				AccessType:   notebookspb.RuntimeAccessConfig_SINGLE_USER,
				RuntimeOwner: email,
			},
		},
	}
	switch menu {
	case "t4-01": // NVIDIA T4 1 基 + Intel 2 vCPU
		runtime.VirtualMachine.VirtualMachineConfig.MachineType = "n1-standard-2"
		runtime.VirtualMachine.VirtualMachineConfig.AcceleratorConfig = &notebookspb.RuntimeAcceleratorConfig{
			Type:      notebookspb.RuntimeAcceleratorConfig_NVIDIA_TESLA_T4,
			CoreCount: 1,
		}
	case "t4-02": //NVIDIA T4 1 基 + Intel 4 vCPU
		runtime.VirtualMachine.VirtualMachineConfig.MachineType = "n1-highmem-4"
		runtime.VirtualMachine.VirtualMachineConfig.AcceleratorConfig = &notebookspb.RuntimeAcceleratorConfig{
			Type:      notebookspb.RuntimeAcceleratorConfig_NVIDIA_TESLA_T4,
			CoreCount: 1,
		}
	case "a100-01": // NVIDIA A100 1 基 + Intel 12 vCPU
		runtime.VirtualMachine.VirtualMachineConfig.MachineType = "a2-highgpu-1g"
		runtime.VirtualMachine.VirtualMachineConfig.AcceleratorConfig = &notebookspb.RuntimeAcceleratorConfig{
			Type:      notebookspb.RuntimeAcceleratorConfig_NVIDIA_TESLA_A100,
			CoreCount: 1,
		}
	}
	req.Runtime.RuntimeType = runtime
	_, err = client.CreateRuntime(ctx, req)
	return err
}

func DescribeManagedNotebook(ctx context.Context, name string) (*notebookspb.Runtime, error) {
	client, err := notebooks.NewManagedNotebookClient(ctx, clientOption())
	if err != nil {
		return nil, err
	}
	defer client.Close()

	req := &notebookspb.GetRuntimeRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/runtimes/%s", util.ProjectID(), util.Location, name),
	}
	return client.GetRuntime(ctx, req)
}

func StartManagedNotebook(ctx context.Context, name string) error {
	client, err := notebooks.NewManagedNotebookClient(ctx, clientOption())
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.StartRuntime(ctx, &notebookspb.StartRuntimeRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/runtimes/%s", util.ProjectID(), util.Location, name),
	})
	return err
}

func StopManagedNotebook(ctx context.Context, name string) error {
	client, err := notebooks.NewManagedNotebookClient(ctx, clientOption())
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.StopRuntime(ctx, &notebookspb.StopRuntimeRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/runtimes/%s", util.ProjectID(), util.Location, name),
	})
	return err
}

func DeleteManagedNotebook(ctx context.Context, name string) error {
	client, err := notebooks.NewManagedNotebookClient(ctx, clientOption())
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.DeleteRuntime(ctx, &notebookspb.DeleteRuntimeRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/runtimes/%s", util.ProjectID(), util.Location, name),
	})
	return err
}
