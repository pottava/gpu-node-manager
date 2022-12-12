package googlecloud

import (
	"github.com/pottava/gpu-node-manager/src/app/util"
	"github.com/revel/revel"
	rm "google.golang.org/api/cloudresourcemanager/v1"
)

func BindRole(r *revel.Request, member, role string) error {
	svc, err := rm.NewService(r.Context())
	if err != nil {
		return err
	}
	policy, err := svc.Projects.GetIamPolicy(
		util.ProjectID(), new(rm.GetIamPolicyRequest)).Do()
	if err != nil {
		return err
	}
	var binding *rm.Binding
	for _, b := range policy.Bindings {
		if b.Role == role {
			binding = b
			break
		}
	}
	if binding != nil {
		binding.Members = append(binding.Members, member)
	} else {
		policy.Bindings = append(policy.Bindings, &rm.Binding{
			Role:    role,
			Members: []string{member},
		})
	}
	request := new(rm.SetIamPolicyRequest)
	request.Policy = policy
	_, err = svc.Projects.SetIamPolicy(util.ProjectID(), request).Do()
	return err
}

func RemoveMember(r *revel.Request, member, role string) error {
	svc, err := rm.NewService(r.Context())
	if err != nil {
		return err
	}
	policy, err := svc.Projects.GetIamPolicy(
		util.ProjectID(), new(rm.GetIamPolicyRequest)).Do()
	if err != nil {
		return err
	}
	var binding *rm.Binding
	var bindingIndex int
	for i, b := range policy.Bindings {
		if b.Role == role {
			binding = b
			bindingIndex = i
			break
		}
	}
	if len(binding.Members) == 1 {
		last := len(policy.Bindings) - 1
		policy.Bindings[bindingIndex] = policy.Bindings[last]
		policy.Bindings = policy.Bindings[:last]
	} else {
		var memberIndex int
		for i, mm := range binding.Members {
			if mm == member {
				memberIndex = i
			}
		}
		last := len(policy.Bindings[bindingIndex].Members) - 1
		binding.Members[memberIndex] = binding.Members[last]
		binding.Members = binding.Members[:last]
	}
	request := new(rm.SetIamPolicyRequest)
	request.Policy = policy
	svc.Projects.SetIamPolicy(util.ProjectID(), request).Do()
	return nil
}
