/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	gcpcompute "google.golang.org/api/compute/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cloudbuilderv1alpha1 "cloudbuilder/api/v1alpha1"
	oauth2 "golang.org/x/oauth2/google"

	"github.com/oracle/oci-go-sdk/v38/common/auth"
)

// ComputeReconciler reconciles a Compute object
type ComputeReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cloudbuilder.example.com,resources=computes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cloudbuilder.example.com,resources=computes/status,verbs=get;update;patch

func (r *ComputeReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("compute", req.NamespacedName)

	log.Info("This is not a perfect code and is only for the purpose of the blog")
	compute := &cloudbuilderv1alpha1.Compute{}
	err := r.Get(ctx, req.NamespacedName, compute)
	if err != nil {
		return ctrl.Result{}, err
	}
	computeName := compute.Spec.ComputeName
	log.Info("Compute resource record created in cluster. Creating compute on cloud platform with name: " + computeName)
	switch compute.Spec.CloudProviderName {
	case "oci":
		if error := r.createOCICompute(ctx, log, compute); error != nil {
			log.Error(err, "Compute create failed for: "+computeName)
			return ctrl.Result{}, error
		}
	case "gcp":
		if error := r.createGCPCompute(ctx, log, compute); error != nil {
			log.Error(err, "Compute create failed for: "+computeName)
			return ctrl.Result{}, error
		}
	case "aws":
		if error := r.createAWSCompute(ctx, log, compute); error != nil {
			log.Error(err, "Compute create failed for: "+computeName)
			return ctrl.Result{}, error
		}
	case "azure":
		if error := r.createAzureCompute(ctx, log, compute); error != nil {
			log.Error(err, "Compute create failed for: "+computeName)
			return ctrl.Result{}, error
		}
	default:
		log.Info("Unknown cloud provider")
	}
	return ctrl.Result{}, nil
}

func (r *ComputeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cloudbuilderv1alpha1.Compute{}).
		Complete(r)
}

func (r *ComputeReconciler) createOCICompute(ctx context.Context, log logr.Logger, compute *cloudbuilderv1alpha1.Compute) error {
	log.Info("Create compute on OCI")
	computeName := compute.Spec.ComputeName
	osImage := compute.Spec.OSImage
	shape := compute.Spec.Shape
	region := compute.Spec.Region
	zone := compute.Spec.Zone
	network := compute.Spec.Network
	subnet := compute.Spec.Subnet
	log.Info("Processing compute resource to create compute resource: " + computeName)
	log.Info("Image: " + osImage + " | Shape: " + shape + " | Network: " + network + " | Subnet: " + subnet + " | Region: " + region + " | Zone: " + zone)
	auth.InstancePrincipalConfigurationProvider()
	// code to add an instance goes here
	return nil
}

func (r *ComputeReconciler) createGCPCompute(ctx context.Context, log logr.Logger, compute *cloudbuilderv1alpha1.Compute) error {
	log.Info("Create compute on GCP")
	computeName := compute.Spec.ComputeName
	osImage := compute.Spec.OSImage
	shape := compute.Spec.Shape
	region := compute.Spec.Region
	zone := compute.Spec.Zone
	network := compute.Spec.Network
	subnet := compute.Spec.Subnet
	log.Info("Processing compute resource to create compute resource: " + computeName)
	log.Info("Image: " + osImage + " | Shape: " + shape + " | Network: " + network + " | Subnet: " + subnet + " | Region: " + region + " | Zone: " + zone)
	log.Info("Fetching default credentials")
	credentials, err := oauth2.FindDefaultCredentials(ctx, gcpcompute.ComputeScope)
	if err != nil {
		return err
	}
	log.Info("Fetching Project ID")
	projectID := credentials.ProjectID
	log.Info("Project ID: " + projectID)
	service, err := gcpcompute.NewService(ctx)
	if err != nil {
		log.Error(err, "Error while creating new service")
		return err
	}
	computeInstance := &gcpcompute.Instance{
		Name:        computeName,
		Description: "Instance created by Kubernetes operator",
		MachineType: "https://www.googleapis.com/compute/v1/projects/" + projectID + "/zones/" + zone + "/machineTypes/" + shape,
		Disks: []*gcpcompute.AttachedDisk{
			{
				AutoDelete: true,
				Boot:       true,
				Type:       "PERSISTENT",
				InitializeParams: &gcpcompute.AttachedDiskInitializeParams{
					DiskName:    "demo-disk",
					SourceImage: "https://www.googleapis.com/compute/v1/projects/centos-cloud/global/images/" + osImage,
				},
			},
		},
		NetworkInterfaces: []*gcpcompute.NetworkInterface{
			{
				Network:    "https://www.googleapis.com/compute/v1/projects/" + projectID + "/global/networks/" + network,
				Subnetwork: "https://www.googleapis.com/compute/v1/projects/" + projectID + "/regions/" + region + "/subnetworks/" + subnet,
			},
		},
		ServiceAccounts: []*gcpcompute.ServiceAccount{
			{
				Email: "default",
				Scopes: []string{
					gcpcompute.ComputeScope,
				},
			},
		},
	}
	operation, err := service.Instances.Insert(projectID, zone, computeInstance).Do()
	if err != nil {
		log.Error(err, "Error occurred during compute create")
	}
	log.Info("Successfully completed instance create " + operation.Header.Get("Etag"))
	return nil
}

func (r *ComputeReconciler) createAWSCompute(ctx context.Context, log logr.Logger, compute *cloudbuilderv1alpha1.Compute) error {
	log.Info("Create compute on AWS")
	return nil
}

func (r *ComputeReconciler) createAzureCompute(ctx context.Context, log logr.Logger, compute *cloudbuilderv1alpha1.Compute) error {
	log.Info("Create compute on Azure")
	return nil
}
