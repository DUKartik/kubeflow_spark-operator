/*
Copyright The Kubeflow Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util_test

import (
	"github.com/kubeflow/spark-operator/v2/pkg/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

var _ = Describe("GetContainerByNameOrFirst", func() {
	It("returns nil when input container list is empty", func() {
		Expect(util.GetContainerByNameOrFirst(nil, "spark")).To(BeNil())
	})

	It("returns the named container when it exists", func() {
		containers := []corev1.Container{
			{Name: "sidecar", Image: "busybox"},
			{Name: "spark", Image: "apache/spark"},
		}

		container := util.GetContainerByNameOrFirst(containers, "spark")

		Expect(container).NotTo(BeNil())
		Expect(container.Name).To(Equal("spark"))
		Expect(container.Image).To(Equal("apache/spark"))
	})

	It("returns the first container when the named container is absent", func() {
		containers := []corev1.Container{
			{Name: "first", Image: "first-image"},
			{Name: "second", Image: "second-image"},
		}

		container := util.GetContainerByNameOrFirst(containers, "spark")

		Expect(container).NotTo(BeNil())
		Expect(container.Name).To(Equal("first"))
		Expect(container.Image).To(Equal("first-image"))
	})

	It("returns a pointer to the original slice element", func() {
		containers := []corev1.Container{
			{Name: "first"},
			{Name: "spark"},
		}

		container := util.GetContainerByNameOrFirst(containers, "spark")
		Expect(container).NotTo(BeNil())

		container.Image = "updated-image"

		Expect(containers[1].Image).To(Equal("updated-image"))
	})
})

var _ = Describe("SetContainerCPUResources", func() {
	var (
		container *corev1.Container
	)

	BeforeEach(func() {
		container = &corev1.Container{}
	})

	It("sets the CPU request when only a request is given", func() {
		request := resource.MustParse("500m")

		returned := util.SetContainerCPUResources(container, &request, nil)

		Expect(returned.Resources.Requests).To(HaveKey(corev1.ResourceCPU))
		Expect(returned.Resources.Requests.Cpu().MilliValue()).To(Equal(int64(500)))
		Expect(returned.Resources.Limits).To(BeNil())
	})

	It("sets the CPU limit when only a limit is given", func() {
		limit := resource.MustParse("1")

		returned := util.SetContainerCPUResources(container, nil, &limit)

		Expect(returned.Resources.Limits).To(HaveKey(corev1.ResourceCPU))
		Expect(returned.Resources.Limits.Cpu().MilliValue()).To(Equal(int64(1000)))
		Expect(returned.Resources.Requests).To(BeNil())
	})

	It("sets both the CPU request and limit when both are given", func() {
		request := resource.MustParse("1.5")
		limit := resource.MustParse("2.5")

		returned := util.SetContainerCPUResources(container, &request, &limit)

		Expect(returned.Resources.Requests.Cpu().MilliValue()).To(Equal(int64(1500)))
		Expect(returned.Resources.Limits.Cpu().MilliValue()).To(Equal(int64(2500)))
	})

	It("creates no resource maps when both are nil", func() {
		returned := util.SetContainerCPUResources(container, nil, nil)

		Expect(returned.Resources.Requests).To(BeNil())
		Expect(returned.Resources.Limits).To(BeNil())
	})

	It("preserves other resource keys", func() {
		container = &corev1.Container{
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceMemory: resource.MustParse("1Gi"),
				},
				Limits: corev1.ResourceList{
					corev1.ResourceMemory: resource.MustParse("1Gi"),
				},
			},
		}
		request := resource.MustParse("500m")

		returned := util.SetContainerCPUResources(container, &request, nil)

		Expect(returned.Resources.Requests).To(HaveKey(corev1.ResourceMemory))
		Expect(returned.Resources.Requests.Memory().Value()).To(Equal(int64(1) << 30))
		Expect(returned.Resources.Requests.Cpu().MilliValue()).To(Equal(int64(500)))
		Expect(returned.Resources.Limits.Memory().Value()).To(Equal(int64(1) << 30))
	})

	It("returns the container it was given", func() {
		request := resource.MustParse("500m")

		returned := util.SetContainerCPUResources(container, &request, nil)

		Expect(returned).To(BeIdenticalTo(container))
	})
})
