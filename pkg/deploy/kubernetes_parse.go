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

package deploy

import (
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
)

func ParseKubernetesYAML(content []byte, into runtime.Object) (runtime.Object, error) {
	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	codecFactory := serializer.NewCodecFactory(scheme)
	deserializer := codecFactory.UniversalDeserializer()

	stringContent := string(content)
	entities := strings.Split(stringContent, "---")

	if len(entities) > 1 {
		for _, entity := range entities {
			decoded, _, _ := deserializer.Decode([]byte(entity), nil, nil)
			if decoded.GetObjectKind().GroupVersionKind() == into.GetObjectKind().GroupVersionKind() {
				content = []byte(entity)
				break
			} else {
				continue
			}
		}
	}

	obj, _, err := deserializer.Decode(content, nil, into)

	return obj, err
}
