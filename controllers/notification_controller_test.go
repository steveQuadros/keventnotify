package controllers

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"k8s.io/client-go/kubernetes/scheme"
	"time"
)

var _ = Describe("[Notifications Controller]", func() {
	var ctx = context.Background()

	It("notifies of events", func() {
		decode := scheme.Codecs.UniversalDeserializer().Decode
		eventsYaml, err := ioutil.ReadFile("testdata/rancherevent.yaml")
		Expect(err).NotTo(HaveOccurred())
		obj, _, _ := decode(eventsYaml, nil, nil)
		fmt.Printf("%+v", obj)
		Expect(k8sClient.Create(ctx, obj)).NotTo(HaveOccurred())
		Eventually(func() bool {
			return false
		}, time.Second*5, time.Second*1).Should(BeTrue())
	})
})
