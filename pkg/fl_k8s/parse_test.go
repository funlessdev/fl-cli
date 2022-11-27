// Copyright 2022 Giuseppe De Palma, Matteo Trentin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fl_k8s

import (
	"testing"

	apiCoreV1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiRbacV1 "k8s.io/api/rbac/v1"
)

func TestParse(t *testing.T) {
	t.Run("should return an error when given malformed YAML", func(t *testing.T) {
		bad_yaml := `
tag:
  key: "value"
- index:
  key: "value"`
		_, err := ParseKubernetesYAML([]byte(bad_yaml), nil)
		require.Error(t, err)
	})

	t.Run("should return an error when given empty YAML", func(t *testing.T) {
		bad_yaml := ""
		_, err := ParseKubernetesYAML([]byte(bad_yaml), nil)
		require.Error(t, err)
	})

	t.Run("should return the correct object when given a single-entity document", func(t *testing.T) {
		yaml := `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: fl-svc-account
  namespace: fl`
		obj, err := ParseKubernetesYAML([]byte(yaml), &apiCoreV1.ServiceAccount{})

		expected := &apiCoreV1.ServiceAccount{
			TypeMeta: v1.TypeMeta{
				Kind:       "ServiceAccount",
				APIVersion: "v1",
			},
			ObjectMeta: v1.ObjectMeta{
				Name:         "fl-svc-account",
				GenerateName: "",
				Namespace:    "fl",
				SelfLink:     "",
			},
			AutomountServiceAccountToken: (*bool)(nil),
		}

		require.NoError(t, err)
		assert.Equal(t, expected, obj)
	})

	t.Run("should return the first object of the correct type when given a multi-entity document", func(t *testing.T) {
		yaml := `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: fl-svc-account
  namespace: fl

---

apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: fl-role
  namespace: fl
rules:
  - apiGroups:
    - ""
    resources:
      - pods
    verbs:
      - list
      - get
      - watch`
		obj1, err1 := ParseKubernetesYAML([]byte(yaml), &apiCoreV1.ServiceAccount{TypeMeta: v1.TypeMeta{Kind: "ServiceAccount", APIVersion: "v1"}})
		obj2, err2 := ParseKubernetesYAML([]byte(yaml), &apiRbacV1.Role{TypeMeta: v1.TypeMeta{Kind: "Role", APIVersion: "rbac.authorization.k8s.io/v1"}})

		expected1 := &apiCoreV1.ServiceAccount{
			TypeMeta: v1.TypeMeta{
				Kind:       "ServiceAccount",
				APIVersion: "v1",
			},
			ObjectMeta: v1.ObjectMeta{
				Name:         "fl-svc-account",
				GenerateName: "",
				Namespace:    "fl",
				SelfLink:     "",
			},
			AutomountServiceAccountToken: (*bool)(nil),
		}

		expected2 := &apiRbacV1.Role{
			TypeMeta: v1.TypeMeta{
				Kind:       "Role",
				APIVersion: "rbac.authorization.k8s.io/v1",
			},
			ObjectMeta: v1.ObjectMeta{
				Name:         "fl-role",
				GenerateName: "",
				Namespace:    "fl",
				SelfLink:     "",
			},
			Rules: []apiRbacV1.PolicyRule{
				{
					Verbs: []string{
						"list",
						"get",
						"watch",
					},
					APIGroups:       []string{""},
					Resources:       []string{"pods"},
					ResourceNames:   []string(nil),
					NonResourceURLs: []string(nil),
				},
			},
		}

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.Equal(t, expected1, obj1)
		assert.Equal(t, expected2, obj2)
	})
}
