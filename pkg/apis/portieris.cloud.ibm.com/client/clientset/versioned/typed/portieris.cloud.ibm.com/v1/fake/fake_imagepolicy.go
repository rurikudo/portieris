/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	portieriscloudibmcomv1 "github.com/rurikudo/portieris/pkg/apis/portieris.cloud.ibm.com/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeImagePolicies implements ImagePolicyInterface
type FakeImagePolicies struct {
	Fake *FakePortierisV1
	ns   string
}

var imagepoliciesResource = schema.GroupVersionResource{Group: "portieris.cloud.ibm.com", Version: "v1", Resource: "imagepolicies"}

var imagepoliciesKind = schema.GroupVersionKind{Group: "portieris.cloud.ibm.com", Version: "v1", Kind: "ImagePolicy"}

// Get takes name of the imagePolicy, and returns the corresponding imagePolicy object, and an error if there is any.
func (c *FakeImagePolicies) Get(name string, options v1.GetOptions) (result *portieriscloudibmcomv1.ImagePolicy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(imagepoliciesResource, c.ns, name), &portieriscloudibmcomv1.ImagePolicy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*portieriscloudibmcomv1.ImagePolicy), err
}

// List takes label and field selectors, and returns the list of ImagePolicies that match those selectors.
func (c *FakeImagePolicies) List(opts v1.ListOptions) (result *portieriscloudibmcomv1.ImagePolicyList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(imagepoliciesResource, imagepoliciesKind, c.ns, opts), &portieriscloudibmcomv1.ImagePolicyList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &portieriscloudibmcomv1.ImagePolicyList{ListMeta: obj.(*portieriscloudibmcomv1.ImagePolicyList).ListMeta}
	for _, item := range obj.(*portieriscloudibmcomv1.ImagePolicyList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested imagePolicies.
func (c *FakeImagePolicies) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(imagepoliciesResource, c.ns, opts))

}

// Create takes the representation of a imagePolicy and creates it.  Returns the server's representation of the imagePolicy, and an error, if there is any.
func (c *FakeImagePolicies) Create(imagePolicy *portieriscloudibmcomv1.ImagePolicy) (result *portieriscloudibmcomv1.ImagePolicy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(imagepoliciesResource, c.ns, imagePolicy), &portieriscloudibmcomv1.ImagePolicy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*portieriscloudibmcomv1.ImagePolicy), err
}

// Update takes the representation of a imagePolicy and updates it. Returns the server's representation of the imagePolicy, and an error, if there is any.
func (c *FakeImagePolicies) Update(imagePolicy *portieriscloudibmcomv1.ImagePolicy) (result *portieriscloudibmcomv1.ImagePolicy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(imagepoliciesResource, c.ns, imagePolicy), &portieriscloudibmcomv1.ImagePolicy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*portieriscloudibmcomv1.ImagePolicy), err
}

// Delete takes name of the imagePolicy and deletes it. Returns an error if one occurs.
func (c *FakeImagePolicies) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(imagepoliciesResource, c.ns, name), &portieriscloudibmcomv1.ImagePolicy{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeImagePolicies) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(imagepoliciesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &portieriscloudibmcomv1.ImagePolicyList{})
	return err
}

// Patch applies the patch and returns the patched imagePolicy.
func (c *FakeImagePolicies) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *portieriscloudibmcomv1.ImagePolicy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(imagepoliciesResource, c.ns, name, pt, data, subresources...), &portieriscloudibmcomv1.ImagePolicy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*portieriscloudibmcomv1.ImagePolicy), err
}
