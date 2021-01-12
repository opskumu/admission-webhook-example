package main

import (
	"context"
	"encoding/json"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type podMutate struct {
	Client  client.Client
	decoder *admission.Decoder
}

func (p *podMutate) Handle(ctx context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}
	podMutateLog := log.WithName("podMutate")

	err := p.decoder.Decode(req, pod)
	if err != nil {
		podMutateLog.Error(err, "failed decoder pod")
		return admission.Errored(http.StatusBadRequest, err)
	}

	podDNSConfig := []corev1.PodDNSConfigOption{}
	ndotsValue := "2"
	ndotsOpt := corev1.PodDNSConfigOption{
		Name:  "ndots",
		Value: &ndotsValue,
	}
	podDNSConfig = append(podDNSConfig, ndotsOpt)
	timeoutValue := "1"
	timeoutOpt := corev1.PodDNSConfigOption{
		Name:  "timeout",
		Value: &timeoutValue,
	}
	podDNSConfig = append(podDNSConfig, timeoutOpt)
	reopenOpt := corev1.PodDNSConfigOption{
		Name: "single-request-reopen",
	}
	podDNSConfig = append(podDNSConfig, reopenOpt)

	if pod.Spec.DNSConfig == nil {
		pod.Spec.DNSConfig = &corev1.PodDNSConfig{
			Options: podDNSConfig,
		}
	} else {
		if len(pod.Spec.DNSConfig.Options) == 0 {
			pod.Spec.DNSConfig.Options = podDNSConfig
		}
	}

	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		podMutateLog.Error(err, "failed marshal pod")
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

// podMutate implements admission.DecoderInjector.
// A decoder will be automatically injected.

// InjectDecoder injects the decoder.
func (p *podMutate) InjectDecoder(d *admission.Decoder) error {
	p.decoder = d
	return nil
}
