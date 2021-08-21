package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ImageContentSourcePolicy holds cluster-wide information about how to handle registry mirror rules.
// When multiple policies are defined, the outcome of the behavior is defined on each field.
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type ImageContentSourcePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
	// +kubebuilder:validation:Required
	// +required
	Spec ImageContentSourcePolicySpec `json:"spec"`
}

// ImageContentSourcePolicySpec is the specification of the ImageContentSourcePolicy CRD.
type ImageContentSourcePolicySpec struct {
	// repositoryMirrors allows images referenced in pods to be
	// pulled from alternative mirrored repository locations. The image pull specification
	// provided to the pod will be compared to the source locations described in RepositoryMirrors
	// and the image may be pulled down from any of the mirrors in the list instead of the
	// specified repository allowing administrators to choose a potentially faster mirror.
	// Image specifications can be referenced by tags or digests.
	//
	// Each “source” repository is treated independently; configurations for different “source”
	// repositories don’t interact.
	//
	// If the "mirrors" is not specified, the image will continue to be pulled from the specified
	// repository in the pull spec.
	//
	// When multiple policies are defined for the same “source” repository, the sets of defined
	// mirrors will be merged together, preserving the relative order of the mirrors, if possible.
	// For example, if policy A has mirrors `a, b, c` and policy B has mirrors `c, d, e`, the
	// mirrors will be used in the order `a, b, c, d, e`.  If the orders of mirror entries conflict
	// (e.g. `a, b` vs. `b, a`) the configuration is not rejected but the resulting order is unspecified.
	//
	// +optional
	// +listType=map
	// +listMapKey=source
	RepositoryMirrors []RepositoryMirrors `json:"repositoryMirrors"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ImageContentSourcePolicyList lists the items in the ImageContentSourcePolicy CRD.
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type ImageContentSourcePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ImageContentSourcePolicy `json:"items"`
}

// RepositoryMirrors holds cluster-wide information about how to handle mirrors in the registries config.
// Note: this is different from the RepositoryDigestMirrors that the mirrors only
// work when pulling the images that are referenced by their digests.
type RepositoryMirrors struct {
	// source is the repository that users refer to, e.g. in image pull specifications.
	// see https://github.com/containers/image/blob/main/docs/containers-registries.conf.5.md#choosing-a-registry-toml-table for format specifications
	// i.e. using the following format
	// host[:port]
	// host[:port]/namespace[/namespace…]
	// host[:port]/namespace[/namespace…]/repo
	// +kubebuilder:validation:Pattern=`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])(:[0-9]+)?(\/[^\/]+)*$`
	// +required
	Source string `json:"source"`
	// If true, the mirrors only work when pulling the images that are referenced by their digests.
	// Pulling images by tag can potentially yield different images, depending on which endpoint
	// we pull from. Forcing digest-pulls for mirrors avoids that issue.
	// +optional
	AllowMirrorByTags bool `json:"allowMirrorByTags,omitempty"`
	// mirrors is zero or many repositories that may also contain the same images.
	// A zero length mirrors or undefined mirrors field will not set a mirror, the image will continue to be pulled from the specified
	// repository in the pull spec.
	// For format of each mirror specifications, see https://github.com/containers/image/blob/main/docs/containers-registries.conf.5.md#choosing-a-registry-toml-table
	// i.e. using the following format
	// host[:port]
	// host[:port]/namespace[/namespace…]
	// host[:port]/namespace[/namespace…]/repo
	// an empty string is not a valid mirror
	// The order of mirrors in this list is treated as the user's desired priority, while source
	// is by default considered lower priority than all mirrors. Other cluster configuration,
	// including (but not limited to) other repositoryMirrors objects,
	// may impact the exact order mirrors are contacted in, or some mirrors may be contacted
	// in parallel, so this should be considered a preference rather than a guarantee of ordering.
	// +optional
	// +listType=set
	Mirrors []string `json:"mirrors,omitempty"`
}
