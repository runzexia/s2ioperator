package validating

import (
	"strings"

	"github.com/docker/distribution/reference"
	api "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	"github.com/kubesphere/s2ioperator/pkg/errors"
)

// ValidateConfig returns a list of error from validation.
func ValidateConfig(config *api.S2iConfig, fromTemplate bool) []error {
	allErrs := make([]error, 0)
	if len(config.SourceURL) == 0 {
		allErrs = append(allErrs, errors.NewFieldRequired("sourceUrl"))
	}
	if !fromTemplate && len(config.BuilderImage) == 0 {
		allErrs = append(allErrs, errors.NewFieldRequired("builderImage"))
	}
	switch config.BuilderPullPolicy {
	case api.PullNever, api.PullAlways, api.PullIfNotPresent:
	default:
		allErrs = append(allErrs, errors.NewFieldInvalidValue("builderPullPolicy"))
	}
	if config.DockerNetworkMode != "" && !validateDockerNetworkMode(config.DockerNetworkMode) {
		allErrs = append(allErrs, errors.NewFieldInvalidValue("dockerNetworkMode"))
	}
	if config.Labels != nil {
		for k := range config.Labels {
			if len(k) == 0 {
				allErrs = append(allErrs, errors.NewFieldInvalidValue("labels"))
			}
		}
	}
	if config.BuilderImage != "" {
		if err := validateDockerReference(config.BuilderImage); err != nil {
			allErrs = append(allErrs, errors.NewFieldInvalidValueWithReason("builderImage", err.Error()))
		}
	}
	return allErrs
}

func validateDockerReference(ref string) error {
	_, err := reference.Parse(ref)
	return err
}

// validateDockerNetworkMode checks wether the network mode conforms to the docker remote API specification (v1.19)
// Supported values are: bridge, host, container:<name|id>, and netns:/proc/<pid>/ns/net
func validateDockerNetworkMode(mode api.DockerNetworkMode) bool {
	switch mode {
	case api.DockerNetworkModeBridge, api.DockerNetworkModeHost:
		return true
	}
	if strings.HasPrefix(string(mode), api.DockerNetworkModeContainerPrefix) {
		return true
	}
	if strings.HasPrefix(string(mode), api.DockerNetworkModeNetworkNamespacePrefix) {
		return true
	}
	return false
}
