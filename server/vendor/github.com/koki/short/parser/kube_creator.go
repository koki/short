package parser

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	admissionregistrationv1alpha1 "k8s.io/api/admissionregistration/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	authenticationv1 "k8s.io/api/authentication/v1"
	authenticationv1beta1 "k8s.io/api/authentication/v1beta1"
	authorizationv1 "k8s.io/api/authorization/v1"
	authorizationv1beta1 "k8s.io/api/authorization/v1beta1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	autoscalingv2beta1 "k8s.io/api/autoscaling/v2beta1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	batchv2alpha1 "k8s.io/api/batch/v2alpha1"
	certificatesv1beta1 "k8s.io/api/certificates/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	rbacv1alpha1 "k8s.io/api/rbac/v1alpha1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	schedulingv1alpha1 "k8s.io/api/scheduling/v1alpha1"
	settingsv1alpha1 "k8s.io/api/settings/v1alpha1"
	storagev1 "k8s.io/api/storage/v1"
	storagev1beta1 "k8s.io/api/storage/v1beta1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	// Scheme knows about all kubernetes types
	creator        = runtime.NewScheme()
	codecs         = serializer.NewCodecFactory(creator)
	parameterCodec = runtime.NewParameterCodec(creator)
)

func init() {
	v1.AddToGroupVersion(creator, schema.GroupVersion{Version: "v1"})
	AddToScheme(creator)
}

func AddToScheme(scheme *runtime.Scheme) {
	admissionregistrationv1alpha1.AddToScheme(scheme)
	appsv1beta1.AddToScheme(scheme)
	appsv1beta2.AddToScheme(scheme)
	appsv1.AddToScheme(scheme)
	authenticationv1.AddToScheme(scheme)
	authenticationv1beta1.AddToScheme(scheme)
	authorizationv1.AddToScheme(scheme)
	authorizationv1beta1.AddToScheme(scheme)
	autoscalingv1.AddToScheme(scheme)
	autoscalingv2beta1.AddToScheme(scheme)
	batchv1.AddToScheme(scheme)
	batchv1beta1.AddToScheme(scheme)
	batchv2alpha1.AddToScheme(scheme)
	certificatesv1beta1.AddToScheme(scheme)
	corev1.AddToScheme(scheme)
	extensionsv1beta1.AddToScheme(scheme)
	networkingv1.AddToScheme(scheme)
	policyv1beta1.AddToScheme(scheme)
	rbacv1.AddToScheme(scheme)
	rbacv1beta1.AddToScheme(scheme)
	rbacv1alpha1.AddToScheme(scheme)
	schedulingv1alpha1.AddToScheme(scheme)
	settingsv1alpha1.AddToScheme(scheme)
	storagev1beta1.AddToScheme(scheme)
	storagev1.AddToScheme(scheme)

}
