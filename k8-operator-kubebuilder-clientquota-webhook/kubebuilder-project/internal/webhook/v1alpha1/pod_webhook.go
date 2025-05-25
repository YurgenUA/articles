package v1alpha1

import (
	"context"
	"encoding/json"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"github.com/YurgenUA/k8s-client-quota/internal/shared"
	quotav1alpha1 "github.com/YurgenUA/k8s-client-quota/api/v1alpha1"
)

type PodValidator struct {
	client.Client
	CfgMapNamespace string
	CfgMapName      string
}

// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch

func (v *PodValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	log := logf.FromContext(ctx)
	log.Info("Start handling...")
	log.Info("Request", "req", req)

	// Only process Pods
	if req.Kind.Kind != "Pod" {
		return admission.Allowed("Not a Pod")
	}

	// Allow the Pod to be created in others than 'playground' namespaces
	if req.Namespace != "playground" {
		return admission.Allowed("Pod is being created in non-guarded namespace")
	}

	apiKey, errResp := v.getApiKey(req)
	if errResp != nil {
		return *errResp
	}
	log.Info("API Key Annotation", "apiKey", apiKey)

	leftClientQuota, errResp := v.readQuotaForApiKey(ctx, apiKey)
	if errResp != nil {
		return *errResp
	}

	if leftClientQuota <= 0 {
		log.Info("Client quota exceeded", "apiKey", apiKey)
		return admission.Denied("Client quota exceeded for API Key: " + apiKey)
	}
	log.Info("Client still has quota", "apiKey", apiKey, "leftQuota", leftClientQuota)
	return admission.Allowed("Quota ok")
}

func (v *PodValidator) getApiKey(req admission.Request) (string, *admission.Response) {
	var pod corev1.Pod
	if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
		resp := admission.Errored(http.StatusBadRequest, err)
		return "", &resp
	}

	apiKey := pod.ObjectMeta.Annotations[shared.API_KEY_ANNOTATION]
	if apiKey == "" {
		resp := admission.Denied("Missing required annotation: " + shared.API_KEY_ANNOTATION)
		return "", &resp
	}
	return apiKey, nil
}

func (v *PodValidator) readQuotaForApiKey(ctx context.Context, apiKeyToFind string ) (int, *admission.Response) {
	log := logf.FromContext(ctx)
	var quota quotav1alpha1.ClientQuota
	if err := v.Get(ctx, client.ObjectKey{
		Namespace: "playground",
		Name:      "client-quota",
	}, &quota); err != nil {
		resp := admission.Errored(http.StatusInternalServerError, err)
		return 0, &resp
	}

	quotaMap, err := shared.ReadOrInitQuotaMap(ctx, v.Client, &quota)
	if err != nil {
		resp := admission.Errored(http.StatusInternalServerError, err)
		return 0, &resp
	}

	// find pod_api_key in quota.Spec.Clients
	for _, client := range quota.Spec.Clients {
		if client.APIKey == apiKeyToFind {
			log.Info("Found Client in quota", "client", client.Name, "leftQuota", quotaMap[client.Name])
			return quotaMap[client.Name], nil
		}
	}
	resp := admission.Denied("API Key not found in quota: " + apiKeyToFind)
	return 0, &resp
}