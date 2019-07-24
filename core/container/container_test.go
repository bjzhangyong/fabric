/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package container_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/hyperledger/fabric/core/common/ccprovider"
	"github.com/hyperledger/fabric/core/container"
	"github.com/hyperledger/fabric/core/container/ccintf"
	"github.com/hyperledger/fabric/core/container/mock"
	"github.com/pkg/errors"
)

var _ = Describe("Container", func() {
	Describe("Router", func() {
		var (
			fakeVM *mock.VM
			router *container.Router
		)

		BeforeEach(func() {
			fakeVM = &mock.VM{}
			router = &container.Router{
				DockerVM: fakeVM,
			}
		})

		Describe("Start", func() {
			BeforeEach(func() {
				fakeVM.StartReturns(errors.New("fake-start-error"))
			})

			It("passes through to the docker impl", func() {
				err := router.Start(
					ccintf.CCID("start:name"),
					"fake-ccType",
					&ccintf.PeerConnection{
						Address: "peer-address",
						TLSConfig: &ccintf.TLSConfig{
							ClientKey:  []byte("key"),
							ClientCert: []byte("cert"),
							RootCert:   []byte("root"),
						},
					},
				)

				Expect(err).To(MatchError("fake-start-error"))
				Expect(fakeVM.StartCallCount()).To(Equal(1))
				ccid, args, peerConnection := fakeVM.StartArgsForCall(0)
				Expect(ccid).To(Equal(ccintf.CCID("start:name")))
				Expect(args).To(Equal("fake-ccType"))
				Expect(peerConnection).To(Equal(&ccintf.PeerConnection{
					Address: "peer-address",
					TLSConfig: &ccintf.TLSConfig{
						ClientKey:  []byte("key"),
						ClientCert: []byte("cert"),
						RootCert:   []byte("root"),
					},
				}))
			})
		})

		Describe("Stop", func() {
			BeforeEach(func() {
				fakeVM.StopReturns(errors.New("Boo"))
			})

			It("passes through to the docker impl", func() {
				err := router.Stop(ccintf.CCID("stop:name"))
				Expect(err).To(MatchError("Boo"))
				Expect(fakeVM.StopCallCount()).To(Equal(1))
				Expect(fakeVM.StopArgsForCall(0)).To(Equal(ccintf.CCID("stop:name")))
			})
		})

		Describe("Build", func() {
			BeforeEach(func() {
				fakeVM.BuildReturns(errors.New("fake-build-error"))
			})

			It("passes through to the docker impl", func() {
				err := router.Build(
					&ccprovider.ChaincodeContainerInfo{
						PackageID: "stop:name",
						Type:      "type",
						Path:      "path",
						Name:      "name",
						Version:   "version",
					},
					[]byte("code-bytes"),
				)
				Expect(err).To(MatchError("fake-build-error"))
				Expect(fakeVM.BuildCallCount()).To(Equal(1))
				ccci, codePackage := fakeVM.BuildArgsForCall(0)
				Expect(ccci).To(Equal(&ccprovider.ChaincodeContainerInfo{
					PackageID: "stop:name",
					Type:      "type",
					Path:      "path",
					Name:      "name",
					Version:   "version",
				}))
				Expect(codePackage).To(Equal([]byte("code-bytes")))
			})
		})

		Describe("Wait", func() {
			BeforeEach(func() {
				fakeVM.WaitReturns(7, errors.New("fake-build-error"))
			})

			It("passes through to the docker impl", func() {
				res, err := router.Wait(
					ccintf.CCID("stop:name"),
				)
				Expect(res).To(Equal(7))
				Expect(err).To(MatchError("fake-build-error"))
				Expect(fakeVM.WaitCallCount()).To(Equal(1))
				ccid := fakeVM.WaitArgsForCall(0)
				Expect(ccid).To(Equal(ccintf.CCID("stop:name")))
			})
		})
	})
})
