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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gcpcompute "google.golang.org/api/compute/v1"

	cloudbuilderv1alpha1 "cloudbuilder/api/v1alpha1"
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
        ctx = context.Background()
	log = r.Log.WithValues("compute", req.NamespacedName)

	log.Info("DEBUG MSG 1")
	compute := &cloudbuilderv1alpha1.Compute{}
	err := r.Get(ctx, req.NamespacedName, compute)
    if err != nil {
        return ctrl.Result{}, err
    }
    log.Info("Compute resource created")
    log.Info(compute.Spec.ComputeName)
    log.Info("DEBUG MSG 2")
    computeName := compute.Spec.ComputeName
    log.Info("DEBUG MSG 3")
    switch compute.Spec.ComputeName {
        case oci:
            if error := r.createOCICompute(ctx, log, &compute); error != nil {
                log.Error(err, "Compute create failed for: " + computeName)
                return ctrl.Result{}, error
            }
        case gcp:
            if error := r.createGCPCompute(ctx, log, &compute); error != nil {
                log.Error(err, "Compute create failed for: " + computeName)
                return ctrl.Result{}, error
            }
        case aws:
            if error := r.createAWSCompute(ctx, log, &compute); error != nil {
                log.Error(err, "Compute create failed for: " + computeName)
                return ctrl.Result{}, error
            }
        case azure:
            if error := r.createAzureCompute(ctx, log, &compute); error != nil {
                log.Error(err, "Compute create failed for: " + computeName)
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

func (r *ComputeReconciler)createOCICompute(ctx context.Context, log logr.Logger, compute *cloudbuilderv1alpha1.Compute){
    log.Info("DEBUG MSG OCI")
}

func (r *ComputeReconciler)createGCPCompute(ctx context.Context, log logr.Logger, compute *cloudbuilderv1alpha1.Compute) err {
    log.Info("DEBUG MSG GCP")
    computeName := compute.Spec.ComputeName
    osImage := compute.Spec.OSImage
    shape := compute.Spec.Shape
    region := compute.Spec.Region
    zone := compute.Spec.Zone
    log.Info("Processing request to create compute resource: " + computeName)

    credentials, err := google.FindDefaultCredentials(ctx,compute.ComputeScope)
    if err != nil {
        return err
    }
    log.Info("DEBUG MSG 10")
    projectID := credentials.ProjectID
    service, err := gcpcompute.NewService(ctx)
    if err != nil {
        log.Error("Error while creating new service")
        return err
    }
    log.Info("DEBUG MSG 10")
    computeInstance := &gcpcompute.Instance{
        Name:        computeName,
        Description: "Instance created by Kubernetes operator",
        MachineType: "https://www.googleapis.com/compute/v1/projects/" + projectID + "/zones/" + zone + "/machineTypes/e2-micro",
        Disks: []*gcpcompute.AttachedDisk{
            {
                AutoDelete: true,
                Boot:       true,
                Type:       "PERSISTENT",
                InitializeParams: &gcpcompute.AttachedDiskInitializeParams{
                    DiskName:    "demo-disk",
                    SourceImage: "https://www.googleapis.com/compute/v1/projects/centos-cloud/global/images/centos-8-v20210316",
                },
            },
        },
        NetworkInterfaces: []*gcpcompute.NetworkInterface{
            {
                AccessConfigs: []*gcpcompute.AccessConfig{
                    {
                        Type: "ONE_TO_ONE_NAT",
                        Name: "External NAT",
                    },
                },
                Network: "https://www.googleapis.com/compute/v1/projects/" + projectID + "/global/networks/default",
            },
        },
        ServiceAccounts: []*gcpcompute.ServiceAccount{
            {
                Email: "default",
                Scopes: []string{
                    gcpcompute.DevstorageFullControlScope,
                    gcpcompute.ComputeScope,
                },
            },
        },
    }
    return nil
}

func (r *ComputeReconciler)createAWSCompute(ctx context.Context, log logr.Logger, compute *cloudbuilderv1alpha1.Compute){
    log.Info("DEBUG MSG AWS")
}

func (r *ComputeReconciler)createAzureCompute(ctx context.Context, log logr.Logger, compute *cloudbuilderv1alpha1.Compute){
    log.Info("DEBUG MSG Azure")
}