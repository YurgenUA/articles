package shared

import (
	"context"
	"fmt"
	"strconv"

	quotav1alpha1 "github.com/YurgenUA/k8s-client-quota/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const API_KEY_ANNOTATION = "quota.operator.k8s.yfenyuk.io/api-key"

func ReadOrInitQuotaMap(ctx context.Context, r client.Client, quota *quotav1alpha1.ClientQuota) (map[string]int, error) {
	cmName := "clientquota-" + quota.Name
	cm := &corev1.ConfigMap{}
	err := r.Get(ctx, client.ObjectKey{Name: cmName, Namespace: quota.Namespace}, cm)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create a new ConfigMap with initial quota from spec
			data := make(map[string]string)
			for _, client := range quota.Spec.Clients {
				data[client.Name] = strconv.Itoa(client.QuotaMinutes)
			}

			newCM := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      cmName,
					Namespace: quota.Namespace,
					Labels: map[string]string{
						"app": "clientquota",
					},
				},
				Data: data,
			}

			if err := r.Create(ctx, newCM); err != nil {
				return nil, err
			}
			quotaMap := make(map[string]int)
			for k, v := range data {
				quotaMap[k], _ = strconv.Atoi(v)
			}
			return quotaMap, nil
		}
		return nil, err // real error
	}

	// Existing CM found — convert Data to map[string]int
	quotaMap := make(map[string]int)
	for k, v := range cm.Data {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid quota value for client %s: %v", k, err)
		}
		quotaMap[k] = i
	}

	return quotaMap, nil
}
