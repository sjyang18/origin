package client

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	kclient "github.com/GoogleCloudPlatform/kubernetes/pkg/client"

	"github.com/openshift/origin/pkg/api/latest"
	"github.com/openshift/origin/pkg/version"
)

// Interface exposes methods on OpenShift resources.
type Interface interface {
	BuildsNamespacer
	BuildConfigsNamespacer
	ImagesNamespacer
	ImageRepositoriesNamespacer
	ImageRepositoryMappingsNamespacer
	ImageRepositoryTagsNamespacer
	DeploymentsNamespacer
	DeploymentConfigsNamespacer
	RoutesNamespacer
	UsersInterface
	UserIdentityMappingsInterface
	ProjectsInterface
	PoliciesNamespacer
	RolesNamespacer
	RoleBindingsNamespacer
	PolicyBindingsNamespacer
	ResourceAccessReviewsNamespacer
	RootResourceAccessReviews
	SubjectAccessReviewsNamespacer
	TemplatesNamespacer
}

func (c *Client) Builds(namespace string) BuildInterface {
	return newBuilds(c, namespace)
}

func (c *Client) BuildConfigs(namespace string) BuildConfigInterface {
	return newBuildConfigs(c, namespace)
}

func (c *Client) Images(namespace string) ImageInterface {
	return newImages(c, namespace)
}

// ImageRepositories provides a REST client for ImageRepository
func (c *Client) ImageRepositories(namespace string) ImageRepositoryInterface {
	return newImageRepositories(c, namespace)
}

// ImageRepositoryMappings provides a REST client for ImageRepositoryMapping
func (c *Client) ImageRepositoryMappings(namespace string) ImageRepositoryMappingInterface {
	return newImageRepositoryMappings(c, namespace)
}

// ImageRepositoryTags provides a REST client for ImageRepositoryTag
func (c *Client) ImageRepositoryTags(namespace string) ImageRepositoryTagInterface {
	return newImageRepositoryTags(c, namespace)
}

// Deployments provides a REST client for Deployment
func (c *Client) Deployments(namespace string) DeploymentInterface {
	return newDeployments(c, namespace)
}

// DeploymentConfigs provides a REST client for DeploymentConfig
func (c *Client) DeploymentConfigs(namespace string) DeploymentConfigInterface {
	return newDeploymentConfigs(c, namespace)
}

// Routes provides a REST client for Route
func (c *Client) Routes(namespace string) RouteInterface {
	return newRoutes(c, namespace)
}

// Users provides a REST client for User
func (c *Client) Users() UserInterface {
	return newUsers(c)
}

// UserIdentityMappings provides a REST client for UserIdentityMapping
func (c *Client) UserIdentityMappings() UserIdentityMappingInterface {
	return newUserIdentityMappings(c)
}

// Projects provides a REST client for Projects
func (c *Client) Projects() ProjectInterface {
	return newProjects(c)
}

// TemplateConfigs provides a REST client for TemplateConfig
func (c *Client) TemplateConfigs(namespace string) TemplateConfigInterface {
	return newTemplateConfigs(c, namespace)
}

// TemplateConfigs provides a REST client for TemplateConfig
func (c *Client) Templates(namespace string) TemplateInterface {
	return newTemplates(c, namespace)
}

func (c *Client) Policies(namespace string) PolicyInterface {
	return newPolicies(c, namespace)
}

func (c *Client) PolicyBindings(namespace string) PolicyBindingInterface {
	return newPolicyBindings(c, namespace)
}

func (c *Client) Roles(namespace string) RoleInterface {
	return newRoles(c, namespace)
}

func (c *Client) RoleBindings(namespace string) RoleBindingInterface {
	return newRoleBindings(c, namespace)
}

func (c *Client) ResourceAccessReviews(namespace string) ResourceAccessReviewInterface {
	return newResourceAccessReviews(c, namespace)
}

func (c *Client) RootResourceAccessReviews() ResourceAccessReviewInterface {
	return newRootResourceAccessReviews(c)
}

func (c *Client) SubjectAccessReviews(namespace string) SubjectAccessReviewInterface {
	return newSubjectAccessReviews(c, namespace)
}

// Client is an OpenShift client object
type Client struct {
	*kclient.RESTClient
}

// New creates an OpenShift client for the given config. This client works with builds, deployments,
// templates, routes, and images. It allows operations such as list, get, update and delete on these
// objects. An error is returned if the provided configuration is not valid.
func New(c *kclient.Config) (*Client, error) {
	config := *c
	if err := SetOpenShiftDefaults(&config); err != nil {
		return nil, err
	}
	client, err := kclient.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &Client{client}, nil
}

func SetOpenShiftDefaults(config *kclient.Config) error {
	if config.Prefix == "" {
		config.Prefix = "/osapi"
	}
	if len(config.UserAgent) == 0 {
		config.UserAgent = DefaultOpenShiftUserAgent()
	}
	if config.Version == "" {
		// Clients default to the preferred code API version
		// TODO: implement version negotiation (highest version supported by server)
		config.Version = latest.Version
	}
	version := config.Version
	versionInterfaces, err := latest.InterfacesFor(version)
	if err != nil {
		return fmt.Errorf("API version '%s' is not recognized (valid values: %s)", version, strings.Join(latest.Versions, ", "))
	}
	if config.Codec == nil {
		config.Codec = versionInterfaces.Codec
	}
	config.LegacyBehavior = (config.Version == "v1beta1")
	return nil
}

// NewOrDie creates an OpenShift client and panics if the provided API version is not recognized.
func NewOrDie(c *kclient.Config) *Client {
	client, err := New(c)
	if err != nil {
		panic(err)
	}
	return client
}

// DefaultOpenShiftUserAgent returns the default user agent that clients can use.
func DefaultOpenShiftUserAgent() string {
	commit := version.Get().GitCommit
	if len(commit) > 7 {
		commit = commit[:7]
	}
	if len(commit) == 0 {
		commit = "unknown"
	}
	version := version.Get().GitVersion
	seg := strings.SplitN(version, "-", 2)
	version = seg[0]
	return fmt.Sprintf("%s/%s (%s/%s) openshift/%s", path.Base(os.Args[0]), version, runtime.GOOS, runtime.GOARCH, commit)
}
