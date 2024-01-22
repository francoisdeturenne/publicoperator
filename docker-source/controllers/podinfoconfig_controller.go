/*
Copyright 2023.

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
	"bytes"
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	tutov1 "podinfoOperator/api/v1"
)

const (
	REQUEUE_DURATION = 20 * time.Second
	UPDATING_TIME    = 30 * time.Second

	ENABLE  = "enable"
	DISABLE = "disable"
)

var (
	URL = GetEnv("URL", "http://localhost:31198") + "/"
)

// PodinfoConfigReconciler reconciles a PodinfoConfig object
type PodinfoConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tuto.config,resources=podinfoconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tuto.config,resources=podinfoconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tuto.config,resources=podinfoconfigs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PodinfoConfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *PodinfoConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.Log
	/*
	 Get the Custom Resource Config
	*/
	cr_config := &tutov1.PodinfoConfig{}

	logger.Info("- - - - - - - - - START - - - - - - - - -")
	logger.Info("- GET INTENT CONFIG - ")
	if err := r.Get(ctx, req.NamespacedName, cr_config); err != nil {
		err_msg := "Reconcile:: Unable to retrieve the requested CRD instance"
		logger.Error(err, err_msg)
		return ctrl.Result{RequeueAfter: REQUEUE_DURATION}, errors.New(err_msg)
	}
	logger.Info("Reconcile:: Retrieved CR config from cache", "config.spec", cr_config.Spec)
	intent_config := cr_config.Spec.Readyz
	var real_config string

	//GetDeviceConfig()
	logger.Info("- GET REAL CONFIG -")
	if err := r.GetDeviceConfig(&real_config); err != nil {
		err_msg := "Reconcile:: Unable to get device configuration"
		logger.Error(err, err_msg)
		return ctrl.Result{RequeueAfter: REQUEUE_DURATION}, errors.New(err_msg)
	}
	logger.Info("Reconcile:: GetDeviceConfig", "config.real.readys", real_config)

	//CompareConfigs()
	isDifferent := intent_config != real_config
	logger.Info("- COMPARE CONFIGs -", "isDifferent", isDifferent)

	//Configure()
	if isDifferent {
		logger.Info("- CONFIGURE -", "config.intent.readys", intent_config)
		err := r.Configure(&intent_config)
		logger.Info("- - - - - - - - - END - - - - - - - - -")
		if err != nil {
			logger.Error("config error detected", err)
			return ctrl.Result{RequeueAfter: REQUEUE_DURATION}, nil
		}
		return ctrl.Result{RequeueAfter: UPDATING_TIME}, nil
	}
	logger.Info("- - - - - - - - - END - - - - - - - - -")
	return ctrl.Result{RequeueAfter: REQUEUE_DURATION}, nil
}

func (r *PodinfoConfigReconciler) GetDeviceConfig(real_config *string) error {
	resp, err := http.Get(URL + "readyz")
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		*real_config = ENABLE
	} else {
		*real_config = DISABLE
	}
	return nil
}

func (r *PodinfoConfigReconciler) Configure(intent_config *string) error {
	logger := log.Log
	if *intent_config == ENABLE {
		logger.Info("Post", "intent", "enable")
		_, err := http.Post(URL+"readyz/enable", "application/json", bytes.NewBufferString(""))
		if err != nil {
			return err
		}
	} else {
		logger.Info("Post", "intent", "disable")
		_, err := http.Post(URL+"readyz/disable", "application/json", bytes.NewBufferString(""))
		if err != nil {
			return err
		}
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodinfoConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tutov1.PodinfoConfig{}).
		Complete(r)
}

/*
Get en environment variable (key) or a default value (fallback)
*/
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
