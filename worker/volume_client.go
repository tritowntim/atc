package worker

import (
	"errors"
	"time"

	"code.cloudfoundry.org/clock"
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/dbng"
	"github.com/concourse/baggageclaim"
)

const creatingVolumeRetryDelay = 1 * time.Second

//go:generate counterfeiter . VolumeClient

type VolumeClient interface {
	CreateVolumeForResourceCache(
		lager.Logger,
		VolumeSpec,
		*dbng.UsedResourceCache,
	) (Volume, error)
	FindOrCreateVolumeForContainer(
		lager.Logger,
		VolumeSpec,
		dbng.CreatingContainer,
		int,
		string,
	) (Volume, error)
	FindOrCreateVolumeForBaseResourceType(
		lager.Logger,
		VolumeSpec,
		int,
		string,
	) (Volume, error)
	FindInitializedVolumeForResourceCache(
		lager.Logger,
		*dbng.UsedResourceCache,
	) (Volume, bool, error)
	LookupVolume(lager.Logger, string) (Volume, bool, error)
}

var ErrVolumeExpiredImmediately = errors.New("volume expired immediately after saving")
var ErrCreatedVolumeNotFound = errors.New("failed-to-find-created-volume-in-baggageclaim")
var ErrBaseResourceTypeNotFound = errors.New("base-resource-type-not-found")

type volumeClient struct {
	baggageclaimClient              baggageclaim.Client
	db                              GardenWorkerDB
	dbVolumeFactory                 dbng.VolumeFactory
	dbWorkerBaseResourceTypeFactory dbng.WorkerBaseResourceTypeFactory
	clock                           clock.Clock
	dbWorker                        dbng.Worker
}

func NewVolumeClient(
	baggageclaimClient baggageclaim.Client,
	db GardenWorkerDB,
	dbVolumeFactory dbng.VolumeFactory,
	dbWorkerBaseResourceTypeFactory dbng.WorkerBaseResourceTypeFactory,
	clock clock.Clock,
	dbWorker dbng.Worker,
) VolumeClient {
	return &volumeClient{
		baggageclaimClient:              baggageclaimClient,
		db:                              db,
		dbVolumeFactory:                 dbVolumeFactory,
		dbWorkerBaseResourceTypeFactory: dbWorkerBaseResourceTypeFactory,
		clock:    clock,
		dbWorker: dbWorker,
	}
}

func (c *volumeClient) FindOrCreateVolumeForContainer(
	logger lager.Logger,
	volumeSpec VolumeSpec,
	container dbng.CreatingContainer,
	teamID int,
	mountPath string,
) (Volume, error) {
	return c.findOrCreateVolume(
		logger,
		volumeSpec,
		func() (dbng.CreatingVolume, dbng.CreatedVolume, error) {
			return c.dbVolumeFactory.FindContainerVolume(teamID, c.dbWorker, container, mountPath)
		},
		func() (dbng.CreatingVolume, error) {
			var v dbng.CreatingVolume
			var err error
			if cowStrategy, ok := volumeSpec.Strategy.(ContainerRootFSStrategy); ok {
				v, err = cowStrategy.Parent.CreateChildForContainer(container, mountPath)
			} else {
				v, err = c.dbVolumeFactory.CreateContainerVolume(teamID, c.dbWorker, container, mountPath)
			}
			if err != nil {
				return nil, err
			}

			logger.Debug("created-volume-for-container", lager.Data{"handle": v.Handle()})
			return v, nil
		},
	)
}

func (c *volumeClient) FindOrCreateVolumeForBaseResourceType(
	logger lager.Logger,
	volumeSpec VolumeSpec,
	teamID int,
	resourceTypeName string,
) (Volume, error) {
	workerBaseResourceType, found, err := c.dbWorkerBaseResourceTypeFactory.Find(resourceTypeName, c.dbWorker)
	if err != nil {
		return nil, err
	}
	if !found {
		logger.Error("base-resource-type-not-found", ErrBaseResourceTypeNotFound, lager.Data{"resource-type-name": resourceTypeName})
		return nil, ErrBaseResourceTypeNotFound
	}

	return c.findOrCreateVolume(
		logger,
		volumeSpec,
		func() (dbng.CreatingVolume, dbng.CreatedVolume, error) {
			return c.dbVolumeFactory.FindBaseResourceTypeVolume(teamID, workerBaseResourceType)
		},
		func() (dbng.CreatingVolume, error) {
			v, err := c.dbVolumeFactory.CreateBaseResourceTypeVolume(teamID, workerBaseResourceType)
			if err != nil {
				return nil, err
			}

			logger.Debug("created-volume-for-base-resource-type", lager.Data{"handle": v.Handle()})
			return v, nil
		},
	)
}

func (c *volumeClient) CreateVolumeForResourceCache(
	logger lager.Logger,
	volumeSpec VolumeSpec,
	usedResourceCache *dbng.UsedResourceCache,
) (Volume, error) {
	return c.findOrCreateVolume(
		logger,
		volumeSpec,
		func() (dbng.CreatingVolume, dbng.CreatedVolume, error) {
			return nil, nil, nil
		},
		func() (dbng.CreatingVolume, error) {
			v, err := c.dbVolumeFactory.CreateResourceCacheVolume(c.dbWorker, usedResourceCache)
			if err != nil {
				return nil, err
			}

			logger.Debug("created-volume-for-resource-cache", lager.Data{"handle": v.Handle()})
			return v, nil
		},
	)
}

func (c *volumeClient) FindInitializedVolumeForResourceCache(
	logger lager.Logger,
	usedResourceCache *dbng.UsedResourceCache,
) (Volume, bool, error) {
	dbVolume, found, err := c.dbVolumeFactory.FindResourceCacheInitializedVolume(c.dbWorker, usedResourceCache)
	if err != nil {
		logger.Error("failed-to-lookup-initialized-volume-in-db", err)
		return nil, false, err
	}

	if !found {
		return nil, false, nil
	}

	bcVolume, found, err := c.baggageclaimClient.LookupVolume(logger, dbVolume.Handle())
	if err != nil {
		logger.Error("failed-to-lookup-volume-in-bc", err)
		return nil, false, err
	}

	if !found {
		return nil, false, nil
	}

	return NewVolume(bcVolume, dbVolume), true, nil
}

func (c *volumeClient) LookupVolume(logger lager.Logger, handle string) (Volume, bool, error) {
	dbVolume, found, err := c.dbVolumeFactory.FindCreatedVolume(handle)
	if err != nil {
		logger.Error("failed-to-lookup-volume-in-db", err)
		return nil, false, err
	}

	if !found {
		return nil, false, nil
	}

	bcVolume, found, err := c.baggageclaimClient.LookupVolume(logger, handle)
	if err != nil {
		logger.Error("failed-to-lookup-volume-in-bc", err)
		return nil, false, err
	}

	if !found {
		return nil, false, nil
	}

	return NewVolume(bcVolume, dbVolume), true, nil
}

func (c *volumeClient) findOrCreateVolume(
	logger lager.Logger,
	volumeSpec VolumeSpec,
	findVolumeFunc func() (dbng.CreatingVolume, dbng.CreatedVolume, error),
	createVolumeFunc func() (dbng.CreatingVolume, error),
) (Volume, error) {
	var bcVolume baggageclaim.Volume
	var bcVolumeFound bool

	creatingVolume, createdVolume, err := findVolumeFunc()
	if err != nil {
		logger.Error("failed-to-find-volume-in-db", err)
		return nil, err
	}

	if createdVolume != nil {
		bcVolume, bcVolumeFound, err = c.baggageclaimClient.LookupVolume(
			logger.Session("lookup-volume"),
			createdVolume.Handle(),
		)
		if err != nil {
			logger.Error("failed-to-lookup-volume-in-baggageclaim", err)
			return nil, err
		}

		if !bcVolumeFound {
			logger.Error("failed-to-lookup-volume-in-baggageclaim", ErrCreatedVolumeNotFound, lager.Data{"handle": createdVolume.Handle()})
			return nil, ErrCreatedVolumeNotFound
		}

		return NewVolume(bcVolume, createdVolume), nil
	} else {
		if creatingVolume != nil {
			bcVolume, bcVolumeFound, err = c.baggageclaimClient.LookupVolume(
				logger.Session("create-volume"),
				creatingVolume.Handle(),
			)
			if err != nil {
				logger.Error("failed-to-lookup-volume-in-baggageclaim", err)
				return nil, err
			}
		} else {
			creatingVolume, err = createVolumeFunc()
			if err != nil {
				logger.Error("failed-to-create-volume-in-db", err)
				return nil, err
			}
		}

		lock, acquired, err := c.db.AcquireVolumeCreatingLock(logger, creatingVolume.ID())
		if err != nil {
			logger.Error("failed-to-acquire-volume-creating-lock", err)
			return nil, err
		}

		if !acquired {
			c.clock.Sleep(creatingVolumeRetryDelay)
			return c.findOrCreateVolume(logger, volumeSpec, findVolumeFunc, createVolumeFunc)
		}

		defer lock.Release()

		if !bcVolumeFound {
			logger.Debug("creating-volume-in-baggageclaim", lager.Data{"handle": creatingVolume.Handle()})
			bcVolume, err = c.baggageclaimClient.CreateVolume(
				logger.Session("create-volume"),
				creatingVolume.Handle(),
				volumeSpec.baggageclaimVolumeSpec(),
			)
			if err != nil {
				logger.Error("failed-to-create-volume-in-baggageclaim", err)
				return nil, err
			}
		}

		createdVolume, err = creatingVolume.Created()
		if err != nil {
			logger.Error("failed-to-initialize-volume", err)
			return nil, err
		}
	}

	return NewVolume(bcVolume, createdVolume), nil
}
