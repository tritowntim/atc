package worker_test

import (
	"errors"
	"time"

	"code.cloudfoundry.org/clock/fakeclock"
	gfakes "code.cloudfoundry.org/garden/gardenfakes"

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/dbng/dbngfakes"
	. "github.com/concourse/atc/worker"
	wfakes "github.com/concourse/atc/worker/workerfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Worker", func() {
	var (
		logger                       *lagertest.TestLogger
		fakeVolumeClient             *wfakes.FakeVolumeClient
		fakeImageFactory             *wfakes.FakeImageFactory
		fakeGardenWorkerDB           *wfakes.FakeGardenWorkerDB
		fakeWorkerProvider           *wfakes.FakeWorkerProvider
		fakeClock                    *fakeclock.FakeClock
		fakeDBResourceCacheFactory   *dbngfakes.FakeResourceCacheFactory
		fakeResourceConfigFactory    *dbngfakes.FakeResourceConfigFactory
		fakeContainerProviderFactory *wfakes.FakeContainerProviderFactory
		fakeContainerProvider        *wfakes.FakeContainerProvider
		activeContainers             int
		resourceTypes                []atc.WorkerResourceType
		platform                     string
		tags                         atc.Tags
		teamID                       int
		workerName                   string
		workerStartTime              int64
		workerUptime                 uint64
		gardenWorker                 Worker
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("test")
		fakeVolumeClient = new(wfakes.FakeVolumeClient)
		fakeImageFactory = new(wfakes.FakeImageFactory)
		fakeGardenWorkerDB = new(wfakes.FakeGardenWorkerDB)
		fakeClock = fakeclock.NewFakeClock(time.Unix(123, 456))
		activeContainers = 42
		resourceTypes = []atc.WorkerResourceType{
			{
				Type:    "some-resource",
				Image:   "some-resource-image",
				Version: "some-version",
			},
		}
		platform = "some-platform"
		tags = atc.Tags{"some", "tags"}
		teamID = 17
		workerName = "some-worker"
		workerStartTime = fakeClock.Now().Unix()
		workerUptime = 0

		fakeDBResourceCacheFactory = new(dbngfakes.FakeResourceCacheFactory)
		fakeResourceConfigFactory = new(dbngfakes.FakeResourceConfigFactory)
		fakeWorkerProvider = new(wfakes.FakeWorkerProvider)
		fakeContainerProvider = new(wfakes.FakeContainerProvider)
		fakeContainerProviderFactory = new(wfakes.FakeContainerProviderFactory)
		fakeContainerProviderFactory.ContainerProviderForReturns(fakeContainerProvider)
	})

	JustBeforeEach(func() {
		gardenWorker = NewGardenWorker(
			fakeContainerProviderFactory,
			fakeVolumeClient,
			fakeGardenWorkerDB,
			fakeWorkerProvider,
			fakeClock,
			activeContainers,
			resourceTypes,
			platform,
			tags,
			teamID,
			workerName,
			workerStartTime,
		)

		fakeClock.IncrementBySeconds(workerUptime)
	})

	Describe("FindContainerByHandle", func() {
		var (
			handle            string
			foundContainer    Container
			existingContainer *wfakes.FakeContainer
			found             bool
			checkErr          error
		)

		BeforeEach(func() {
			handle = "we98lsv"
			existingContainer = new(wfakes.FakeContainer)
			fakeContainerProvider.FindContainerByHandleReturns(existingContainer, true, nil)
		})

		JustBeforeEach(func() {
			foundContainer, found, checkErr = gardenWorker.FindContainerByHandle(logger, handle, 42)
		})

		It("calls the container provider", func() {
			Expect(fakeContainerProviderFactory.ContainerProviderForCallCount()).To(Equal(1))

			Expect(fakeContainerProvider.FindContainerByHandleCallCount()).To(Equal(1))

			Expect(foundContainer).To(Equal(existingContainer))
			Expect(checkErr).ToNot(HaveOccurred())
			Expect(found).To(BeTrue())
		})

	})

	Describe("FindOrCreateBuildContainer", func() {
		var container Container
		var createErr error
		var imageSpec ImageSpec

		JustBeforeEach(func() {
			container, createErr = gardenWorker.FindOrCreateBuildContainer(
				logger,
				nil,
				nil,
				Identifier{},
				Metadata{},
				ContainerSpec{
					ImageSpec: imageSpec,
				},
				atc.VersionedResourceTypes{},
				map[string]string{},
			)
		})

		It("delegates container creation to the container provider", func() {
			Expect(fakeContainerProvider.FindOrCreateBuildContainerCallCount()).To(Equal(1))
		})
	})

	Describe("FindContainerForIdentifier", func() {
		var (
			id Identifier

			foundContainer Container
			found          bool
			lookupErr      error
		)

		BeforeEach(func() {
			id = Identifier{
				ResourceID: 1234,
			}
		})

		JustBeforeEach(func() {
			foundContainer, found, lookupErr = gardenWorker.FindContainerForIdentifier(logger, id)
		})

		Context("when the container can be found", func() {
			var (
				fakeContainer       *gfakes.FakeContainer
				fakeSavedContainer  db.SavedContainer
				fakeWorkerContainer *wfakes.FakeContainer
			)

			BeforeEach(func() {
				fakeContainer = new(gfakes.FakeContainer)
				fakeContainer.HandleReturns("provider-handle")

				fakeWorkerContainer = new(wfakes.FakeContainer)
				fakeWorkerContainer.HandleReturns("provider-handle")

				fakeSavedContainer = db.SavedContainer{
					Container: db.Container{
						ContainerIdentifier: db.ContainerIdentifier{
							ResourceID: 1234,
						},
						ContainerMetadata: db.ContainerMetadata{
							Handle:     "provider-handle",
							WorkerName: "some-worker",
						},
					},
				}

				fakeWorkerProvider.FindContainerForIdentifierReturns(fakeSavedContainer, true, nil)

				fakeContainerProvider.FindContainerByHandleReturns(fakeWorkerContainer, true, nil)

				fakeDBTeam := new(dbngfakes.FakeTeam)
				fakeDBTeam.FindContainerByHandleReturns(new(dbngfakes.FakeCreatedContainer), true, nil)
				fakeGardenWorkerDB.GetContainerReturns(fakeSavedContainer, true, nil)
			})

			It("succeeds", func() {
				Expect(lookupErr).NotTo(HaveOccurred())
			})

			It("looks for containers with matching properties via the container provider", func() {
				Expect(fakeWorkerProvider.FindContainerForIdentifierCallCount()).To(Equal(1))
				Expect(fakeWorkerProvider.FindContainerForIdentifierArgsForCall(0)).To(Equal(id))

				Expect(fakeContainerProvider.FindContainerByHandleCallCount()).To(Equal(1))
				_, lookupHandle, _ := fakeContainerProvider.FindContainerByHandleArgsForCall(0)
				Expect(lookupHandle).To(Equal("provider-handle"))
			})
		})

		Context("when looking up the container fails", func() {
			disaster := errors.New("nope")

			BeforeEach(func() {
				containerToReturn := db.SavedContainer{
					Container: db.Container{
						ContainerMetadata: db.ContainerMetadata{
							Handle: "handle",
						},
					},
				}

				fakeWorkerProvider.FindContainerForIdentifierReturns(containerToReturn, true, nil)
				fakeContainerProvider.FindContainerByHandleReturns(nil, false, disaster)
			})

			It("returns the error", func() {
				Expect(lookupErr).To(Equal(disaster))
			})
		})
	})

	Describe("Satisfying", func() {
		var (
			spec WorkerSpec

			satisfyingWorker Worker
			satisfyingErr    error

			customTypes atc.VersionedResourceTypes
		)

		BeforeEach(func() {
			spec = WorkerSpec{
				Tags:   []string{"some", "tags"},
				TeamID: teamID,
			}

			customTypes = atc.VersionedResourceTypes{
				{
					ResourceType: atc.ResourceType{
						Name:   "custom-type-b",
						Type:   "custom-type-a",
						Source: atc.Source{"some": "source"},
					},
					Version: atc.Version{"some": "version"},
				},
				{
					ResourceType: atc.ResourceType{
						Name:   "custom-type-a",
						Type:   "some-resource",
						Source: atc.Source{"some": "source"},
					},
					Version: atc.Version{"some": "version"},
				},
				{
					ResourceType: atc.ResourceType{
						Name:   "custom-type-c",
						Type:   "custom-type-b",
						Source: atc.Source{"some": "source"},
					},
					Version: atc.Version{"some": "version"},
				},
				{
					ResourceType: atc.ResourceType{
						Name:   "custom-type-d",
						Type:   "custom-type-b",
						Source: atc.Source{"some": "source"},
					},
					Version: atc.Version{"some": "version"},
				},
				{
					ResourceType: atc.ResourceType{
						Name:   "unknown-custom-type",
						Type:   "unknown-base-type",
						Source: atc.Source{"some": "source"},
					},
					Version: atc.Version{"some": "version"},
				},
			}
		})

		JustBeforeEach(func() {
			satisfyingWorker, satisfyingErr = gardenWorker.Satisfying(spec, customTypes)
		})

		Context("when the platform is compatible", func() {
			BeforeEach(func() {
				spec.Platform = "some-platform"
			})

			Context("when no tags are specified", func() {
				BeforeEach(func() {
					spec.Tags = nil
				})

				It("returns ErrIncompatiblePlatform", func() {
					Expect(satisfyingErr).To(Equal(ErrMismatchedTags))
				})
			})

			Context("when the worker has no tags", func() {
				BeforeEach(func() {
					tags = []string{}
					spec.Tags = []string{}
				})

				It("returns the worker", func() {
					Expect(satisfyingWorker).To(Equal(gardenWorker))
				})

				It("returns no error", func() {
					Expect(satisfyingErr).NotTo(HaveOccurred())
				})
			})

			Context("when all of the requested tags are present", func() {
				BeforeEach(func() {
					spec.Tags = []string{"some", "tags"}
				})

				It("returns the worker", func() {
					Expect(satisfyingWorker).To(Equal(gardenWorker))
				})

				It("returns no error", func() {
					Expect(satisfyingErr).NotTo(HaveOccurred())
				})
			})

			Context("when some of the requested tags are present", func() {
				BeforeEach(func() {
					spec.Tags = []string{"some"}
				})

				It("returns the worker", func() {
					Expect(satisfyingWorker).To(Equal(gardenWorker))
				})

				It("returns no error", func() {
					Expect(satisfyingErr).NotTo(HaveOccurred())
				})
			})

			Context("when any of the requested tags are not present", func() {
				BeforeEach(func() {
					spec.Tags = []string{"bogus", "tags"}
				})

				It("returns ErrMismatchedTags", func() {
					Expect(satisfyingErr).To(Equal(ErrMismatchedTags))
				})
			})
		})

		Context("when the platform is incompatible", func() {
			BeforeEach(func() {
				spec.Platform = "some-bogus-platform"
			})

			It("returns ErrIncompatiblePlatform", func() {
				Expect(satisfyingErr).To(Equal(ErrIncompatiblePlatform))
			})
		})

		Context("when the resource type is supported by the worker", func() {
			BeforeEach(func() {
				spec.ResourceType = "some-resource"
			})

			Context("when all of the requested tags are present", func() {
				BeforeEach(func() {
					spec.Tags = []string{"some", "tags"}
				})

				It("returns the worker", func() {
					Expect(satisfyingWorker).To(Equal(gardenWorker))
				})

				It("returns no error", func() {
					Expect(satisfyingErr).NotTo(HaveOccurred())
				})
			})

			Context("when some of the requested tags are present", func() {
				BeforeEach(func() {
					spec.Tags = []string{"some"}
				})

				It("returns the worker", func() {
					Expect(satisfyingWorker).To(Equal(gardenWorker))
				})

				It("returns no error", func() {
					Expect(satisfyingErr).NotTo(HaveOccurred())
				})
			})

			Context("when any of the requested tags are not present", func() {
				BeforeEach(func() {
					spec.Tags = []string{"bogus", "tags"}
				})

				It("returns ErrMismatchedTags", func() {
					Expect(satisfyingErr).To(Equal(ErrMismatchedTags))
				})
			})
		})

		Context("when the resource type is a custom type supported by the worker", func() {
			BeforeEach(func() {
				spec.ResourceType = "custom-type-c"
			})

			It("returns the worker", func() {
				Expect(satisfyingWorker).To(Equal(gardenWorker))
			})

			It("returns no error", func() {
				Expect(satisfyingErr).NotTo(HaveOccurred())
			})
		})

		Context("when the resource type is a custom type that overrides one supported by the worker", func() {
			BeforeEach(func() {
				customTypes = append(customTypes, atc.VersionedResourceType{
					ResourceType: atc.ResourceType{
						Name:   "some-resource",
						Type:   "some-resource",
						Source: atc.Source{"some": "source"},
					},
					Version: atc.Version{"some": "version"},
				})

				spec.ResourceType = "some-resource"
			})

			It("returns the worker", func() {
				Expect(satisfyingWorker).To(Equal(gardenWorker))
			})

			It("returns no error", func() {
				Expect(satisfyingErr).NotTo(HaveOccurred())
			})
		})

		Context("when the resource type is a custom type that results in a circular dependency", func() {
			BeforeEach(func() {
				customTypes = append(customTypes, atc.VersionedResourceType{
					ResourceType: atc.ResourceType{
						Name:   "circle-a",
						Type:   "circle-b",
						Source: atc.Source{"some": "source"},
					},
					Version: atc.Version{"some": "version"},
				}, atc.VersionedResourceType{
					ResourceType: atc.ResourceType{
						Name:   "circle-b",
						Type:   "circle-c",
						Source: atc.Source{"some": "source"},
					},
					Version: atc.Version{"some": "version"},
				}, atc.VersionedResourceType{
					ResourceType: atc.ResourceType{
						Name:   "circle-c",
						Type:   "circle-a",
						Source: atc.Source{"some": "source"},
					},
					Version: atc.Version{"some": "version"},
				})

				spec.ResourceType = "circle-a"
			})

			It("returns ErrUnsupportedResourceType", func() {
				Expect(satisfyingErr).To(Equal(ErrUnsupportedResourceType))
			})
		})

		Context("when the resource type is a custom type not supported by the worker", func() {
			BeforeEach(func() {
				spec.ResourceType = "unknown-custom-type"
			})

			It("returns ErrUnsupportedResourceType", func() {
				Expect(satisfyingErr).To(Equal(ErrUnsupportedResourceType))
			})
		})

		Context("when the type is not supported by the worker", func() {
			BeforeEach(func() {
				spec.ResourceType = "some-other-resource"
			})

			It("returns ErrUnsupportedResourceType", func() {
				Expect(satisfyingErr).To(Equal(ErrUnsupportedResourceType))
			})
		})

		Context("when spec specifies team", func() {
			BeforeEach(func() {
				teamID = 123
				spec.TeamID = teamID
			})

			Context("when worker belongs to same team", func() {
				It("returns the worker", func() {
					Expect(satisfyingWorker).To(Equal(gardenWorker))
				})

				It("returns no error", func() {
					Expect(satisfyingErr).NotTo(HaveOccurred())
				})
			})

			Context("when worker belongs to different team", func() {
				BeforeEach(func() {
					teamID = 777
				})

				It("returns ErrTeamMismatch", func() {
					Expect(satisfyingErr).To(Equal(ErrTeamMismatch))
				})
			})

			Context("when worker does not belong to any team", func() {
				It("returns the worker", func() {
					Expect(satisfyingWorker).To(Equal(gardenWorker))
				})

				It("returns no error", func() {
					Expect(satisfyingErr).NotTo(HaveOccurred())
				})
			})
		})

		Context("when spec does not specify a team", func() {
			Context("when worker belongs to no team", func() {
				BeforeEach(func() {
					teamID = 0
				})

				It("returns the worker", func() {
					Expect(satisfyingWorker).To(Equal(gardenWorker))
				})

				It("returns no error", func() {
					Expect(satisfyingErr).NotTo(HaveOccurred())
				})
			})

			Context("when worker belongs to any team", func() {
				BeforeEach(func() {
					teamID = 555
				})

				It("returns ErrTeamMismatch", func() {
					Expect(satisfyingErr).To(Equal(ErrTeamMismatch))
				})
			})
		})
	})
})
