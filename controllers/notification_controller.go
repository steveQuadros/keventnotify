package controllers

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NotificationReconciler reconciles a Notification object
type NotificationReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Notifier Notifier
}

// +kubebuilder:rbac:groups=events.quad.com,resources=notifications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=events.quad.com,resources=notifications/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=,resources=events/status,verbs=get;list

func (r *NotificationReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("notification", req.NamespacedName)

	// your logic here
	event := &v1.Event{}
	if err := r.Get(ctx, req.NamespacedName, event); err != nil {
		if errors.IsNotFound(err) {
			log.Info("event not found")
			return ctrl.Result{}, nil
		}
		log.Error(err, "fetching event")
		return ctrl.Result{}, err
	}

	if err := r.Notifier.Notify(*event); err != nil {
		log.Error(err, "notifying")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *NotificationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Event{}).
		Complete(r)
}

type Notifier interface {
	Notify(event v1.Event) error
}

type PrintNotifier struct{}

func (n *PrintNotifier) Notify(event v1.Event) error {
	fmt.Println(event)
	return nil
}
