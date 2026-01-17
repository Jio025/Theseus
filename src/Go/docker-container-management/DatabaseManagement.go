package dockercontainermanagement

import (
	"encoding/json"
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

// JSON for the host machines
type HostMachines struct {
	ID     string `json:"id"`
	IP     string `json:"ip"`
	Status string `json:"status"`
}

type PortBinding struct {
	Internal int `json:"internal"`
	External int `json:"external"`
}

// JSON for the dockercontainer running on the host
type DockerContainer struct {
	// ID is the PK in the database
	ID string `json:"id"`
	// Name is the image name
	Name string `json:"name"`
	// Container is the container name
	Container string `json:"container"`
	// Host machine that the container is running on
	HostMachine HostMachines `json:"hostmachine"`
	// Restart policy for the container
	RestartPolicy string `json:"restartpolicy"`
	// Ports of the docker container
	Port []PortBinding `json:"ports"`
	// Environment variables of a docker container
	EnvironmentVariable map[string]string `json:"environmentvariables"`
	// Volume mounts of the docker container
	VolumeMounts map[string]string `json:"volumemounts"`
	// shm_size of the docker container
	ShmSize string `json:"shmsize"`
	// Container Status eg: Active, Stopped, Restarting
	Status string `json:"status"`
}

// DB service wraps the BoltDB instance
type DBService struct {
	db *bolt.DB
}

// DockerContainer running on the host bucket name
const DockerContainerBucketName = "DockerContainers"

// Host machines bucket name
const HostMachineBucketName = "HostMachines"

// InitDB starts and creates the DB and the buckets
func InitDB(filepath string) (*DBService, error) {
	db, err := bolt.Open(filepath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("❗ could not open the database: %v", err)
	}

	// Creating the buckets
	err = db.Update(func(tx *bolt.Tx) error {
		// Defining the buckets that are going to be created
		buckets := []string{DockerContainerBucketName, HostMachineBucketName}

		for _, bucket := range buckets {
			if _, err := tx.CreateBucketIfNotExists([]byte(bucket)); err != nil {
				return fmt.Errorf("❗ error creating bucket %s: %s", bucket, err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("❗ error creating bucket in database: %s", err)
	}

	return &DBService{db: db}, nil
}

func (s *DBService) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

// Saving the information of the running DockerContainer
func (s *DBService) SaveActiveDockerContainer(d DockerContainer) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DockerContainerBucketName))

		buf, err := json.Marshal(d)
		if err != nil {
			return err
		}
		return b.Put([]byte(d.ID), buf)
	})
}

// Saving the information of the active HostMachines
func (s *DBService) SaveActiveHostMachines(h HostMachines) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(HostMachineBucketName))

		buf, err := json.Marshal(h)
		if err != nil {
			return err
		}
		return b.Put([]byte(h.ID), buf)
	})
}

func (s DBService) GetActiveDockerContainerByID(id string) (*DockerContainer, error) {
	var d DockerContainer

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DockerContainerBucketName))

		// Get the raw bytes
		data := b.Get([]byte(id))

		// if the data is nil, the key dosent exist
		if data == nil {
			return fmt.Errorf("❗ the Docker Container for id : %s, was not found", id)
		}

		// Deserializing the JSON back to Struct
		return json.Unmarshal(data, &d)
	})

	if err != nil {
		return nil, fmt.Errorf("❗ error reading the data for id : %s", id)
	}

	return &d, nil
}

func (s DBService) GetActiveHostMachineByID(id string) (*HostMachines, error) {
	var h HostMachines

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(HostMachineBucketName))

		// Get the raw bytes
		data := b.Get([]byte(id))

		// if the data is nil, the key dosent exist
		if data == nil {
			return fmt.Errorf("❗ the host machine for id : %s, was not found", id)
		}

		// Deserializing the JSON back to Struct
		return json.Unmarshal(data, &h)
	})

	if err != nil {
		return nil, fmt.Errorf("❗ error reading the data for id : %s", id)
	}

	return &h, nil
}

func (s DBService) GetAllActiveDockerContainers() ([]DockerContainer, error) {
	var dockerContainers []DockerContainer

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DockerContainerBucketName))

		// Iterates for all key-values in the bucket
		return b.ForEach(func(k, v []byte) error {
			var d DockerContainer
			if err := json.Unmarshal(v, &d); err != nil {
				return err
			}
			dockerContainers = append(dockerContainers, d)
			return nil
		})
	})

	return dockerContainers, err
}

func (s DBService) GetAllActiveHostMachines() ([]HostMachines, error) {
	var hostMachines []HostMachines

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(HostMachineBucketName))

		// Iterates for all key-values in the bucket
		return b.ForEach(func(k, v []byte) error {
			var h HostMachines
			if err := json.Unmarshal(v, &h); err != nil {
				return err
			}
			hostMachines = append(hostMachines, h)
			return nil
		})
	})

	return hostMachines, err
}
