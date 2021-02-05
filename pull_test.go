/*
   Copyright The containerd Authors.

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

package containerd

import (
	"context"
	"strings"
	"testing"

	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/platforms"
	_ "github.com/containerd/containerd/runtime"
)

func TestImagePullPerf(t *testing.T) {
	// Max length of namespace should be 76
	namespace := strings.Repeat("n", 76)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = namespaces.WithNamespace(ctx, namespace)

	client, err := newClient(t, address)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	image, err := client.Pull(ctx, testImage,
		WithPlatformMatcher(platforms.Default()),
		WithPullUnpack,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := client.imageStore.Delete(ctx, image.Name(), images.SynchronousDelete()); err != nil {
		t.Errorf("cannot delete image %q, %v", image.Name(), err)
	}
}
