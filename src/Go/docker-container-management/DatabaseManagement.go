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

// JSON for the user
type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
	Team         Team   `json:"team"`
}

type Team struct {
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	Organization Organization `json:"organization"`
}

type Organization struct {
	Name        string `json:"name"`
	Description string `json:"description"`
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

// DockerContainer bucket name
const DockerContainerBucketName = "DockerContainers"

// Host machines bucket name
const HostMachineBucketName = "HostMachines"

// Users bucket name
const UserBucketName = "Users"

// Teams bucket name
const TeamBucketName = "Teams"

// Organization bucket name
const OrganizationBucketName = "Organizations"

// InitDB starts and creates the DB and the buckets
func InitDB(filepath string) (*DBService, error) {
	db, err := bolt.Open(filepath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("❗ could not open the database: %v", err)
	}

	// Creating the buckets
	err = db.Update(func(tx *bolt.Tx) error {
		// Defining the buckets that are going to be created
		buckets := []string{DockerContainerBucketName, HostMachineBucketName, UserBucketName, TeamBucketName, OrganizationBucketName}

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

// Saving the information of the DockerContainer
func (s *DBService) SaveDockerContainer(d DockerContainer) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DockerContainerBucketName))

		buf, err := json.Marshal(d)
		if err != nil {
			return err
		}
		return b.Put([]byte(d.ID), buf)
	})
}

// Saving the information of the HostMachines
func (s *DBService) SaveHostMachines(h HostMachines) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(HostMachineBucketName))

		buf, err := json.Marshal(h)
		if err != nil {
			return err
		}
		return b.Put([]byte(h.ID), buf)
	})
}

func (s DBService) GetDockerContainerByID(id string) (*DockerContainer, error) {
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

func (s DBService) GetHostMachineByID(id string) (*HostMachines, error) {
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

func (s DBService) GetAllDockerContainers() ([]DockerContainer, error) {
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

func (s DBService) GetAllHostMachines() ([]HostMachines, error) {
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

// Saving the information of the User
func (s *DBService) SaveUser(u User) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UserBucketName))

		buf, err := json.Marshal(u)
		if err != nil {
			return err
		}
		return b.Put([]byte(u.Username), buf)
	})
}

// GetUserByUsername retrieves a user by their username
func (s DBService) GetUserByUsername(username string) (*User, error) {
	var u User

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UserBucketName))

		data := b.Get([]byte(username))

		if data == nil {
			return fmt.Errorf("❗ user with username : %s, was not found", username)
		}

		return json.Unmarshal(data, &u)
	})

	if err != nil {
		return nil, err
	}

	return &u, nil
}

// GetAllUsers retrieves all users from the database
func (s DBService) GetAllUsers() ([]User, error) {
	var users []User

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UserBucketName))

		return b.ForEach(func(k, v []byte) error {
			var u User
			if err := json.Unmarshal(v, &u); err != nil {
				return err
			}
			users = append(users, u)
			return nil
		})
	})

	return users, err
}

// Saving the information of the Team
func (s *DBService) SaveTeam(t Team) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TeamBucketName))

		buf, err := json.Marshal(t)
		if err != nil {
			return err
		}
		return b.Put([]byte(t.Name), buf)
	})
}

// GetTeamByName retrieves a team by its name
func (s DBService) GetTeamByName(name string) (*Team, error) {
	var t Team

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TeamBucketName))

		data := b.Get([]byte(name))

		if data == nil {
			return fmt.Errorf("❗ team with name : %s, was not found", name)
		}

		return json.Unmarshal(data, &t)
	})

	if err != nil {
		return nil, err
	}

	return &t, nil
}

// GetAllTeams retrieves all teams from the database
func (s DBService) GetAllTeams() ([]Team, error) {
	var teams []Team

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TeamBucketName))

		return b.ForEach(func(k, v []byte) error {
			var t Team
			if err := json.Unmarshal(v, &t); err != nil {
				return err
			}
			teams = append(teams, t)
			return nil
		})
	})

	return teams, err
}

// Saving the information of the Organization
func (s *DBService) SaveOrganization(o Organization) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(OrganizationBucketName))

		buf, err := json.Marshal(o)
		if err != nil {
			return err
		}
		return b.Put([]byte(o.Name), buf)
	})
}

// GetOrganizationByName retrieves an organization by its name
func (s DBService) GetOrganizationByName(name string) (*Organization, error) {
	var o Organization

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(OrganizationBucketName))

		data := b.Get([]byte(name))

		if data == nil {
			return fmt.Errorf("❗ organization with name : %s, was not found", name)
		}

		return json.Unmarshal(data, &o)
	})

	if err != nil {
		return nil, err
	}

	return &o, nil
}

// GetAllOrganizations retrieves all organizations from the database
func (s DBService) GetAllOrganizations() ([]Organization, error) {
	var organizations []Organization

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(OrganizationBucketName))

		return b.ForEach(func(k, v []byte) error {
			var o Organization
			if err := json.Unmarshal(v, &o); err != nil {
				return err
			}
			organizations = append(organizations, o)
			return nil
		})
	})

	return organizations, err
}
