package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// DeploymentReconciler reconciles a Deployment object
type DeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const (
	ecosStrategyAnnotation = "ecos/strategy"
	deployNameLabel        = "nkzren/deploy-name"
)

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;update;patch
func (r *DeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var deploy appsv1.Deployment
	if err := r.Get(ctx, req.NamespacedName, &deploy); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		logger.Error(err, "unable to fetch Deployment")
		return ctrl.Result{}, err
	}

	// currently only accepts "default" strategy
	hasEcosStrategy := deploy.Annotations[ecosStrategyAnnotation] != ""

	if hasEcosStrategy {
		logger.Info("adding affinities")
		var podSpec = &deploy.Spec.Template.Spec
		var newAffinity = setAffinity(podSpec.DeepCopy())
		podSpec.Affinity = newAffinity
	}

	if err := r.Update(ctx, &deploy); err != nil {
		if apierrors.IsConflict(err) {
			return ctrl.Result{Requeue: true}, nil
		}
		if apierrors.IsNotFound(err) {
			return ctrl.Result{Requeue: true}, nil
		}
		logger.Error(err, "unable to update deploy")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func setAffinity(podSpec *corev1.PodSpec) *corev1.Affinity {
	if podSpec.Affinity == nil {
		podSpec.Affinity = &corev1.Affinity{}
	}
	setNodeAffinity(podSpec.Affinity)
	return podSpec.Affinity
}

func setNodeAffinity(affinity *corev1.Affinity) {
	if affinity.NodeAffinity == nil {
		affinity.NodeAffinity = &corev1.NodeAffinity{}
	}
	setTerms(affinity.NodeAffinity)
}

func setTerms(nodeAffinity *corev1.NodeAffinity) {
	if nodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution == nil {
		nodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution = []corev1.PreferredSchedulingTerm{
			buildTerm(3, "NotIn", "bad"),
			buildTerm(2, "In", "good"),
			buildTerm(2, "NotIn", "neutral"),
		}
	} else {
		var t = nodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution
		t = append(t, buildTerm(3, "NotIn", "bad"))
		t = append(t, buildTerm(2, "In", "bad"))
		t = append(t, buildTerm(2, "NotIn", "bad"))
	}
}

func buildTerm(weight int, operator string, value string) corev1.PreferredSchedulingTerm {
	exp := []corev1.NodeSelectorRequirement{{
		Key:      "ecos",
		Operator: corev1.NodeSelectorOperator(operator),
		Values:   []string{value},
	}}
	return corev1.PreferredSchedulingTerm{
		Weight: int32(weight),
		Preference: corev1.NodeSelectorTerm{
			MatchExpressions: exp,
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		Complete(r)
}
