package models

import (
	"encoding/hex"
	"gorm.io/gorm"
	"time"

	"golang.org/x/crypto/ssh"
)

// Node contains information about a node.
type Node struct {
	Name       string            `json:"name" gorm:"<-:create;notnull;unique"`
	ID         string            `json:"id" gorm:"<-:create;primarykey"`
	IDShort    string            `json:"-" gorm:"<-:create;notnull;unique"`
	PrivateKey []byte            `json:"private_key" gorm:"<-:create;notnull"`
	PublicKey  []byte            `json:"public_key" gorm:"<-:create;notnull"`
	HostKey    []byte            `json:"host_key" gorm:"<-:create;notnull;unique"`
	IP         string            `json:"ip" gorm:"<-:create;notnull;unique"`
	Username   string            `json:"username" gorm:"<-:create;notnull;default:root"`
	Password   string            `json:"password,omitempty" gorm:"-"`
	Port       string            `json:"port" gorm:"<-:create;default:22"`
	Config     *ssh.ClientConfig `json:"-" gorm:"-"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// GetIDShort returns the first 8 bytes of the node ID.
func (n *Node) GetIDShort() string { x, _ := hex.DecodeString(n.ID); return hex.EncodeToString(x[:8]) }

func (n *Node) BeforeCreate(tx *gorm.DB) (err error) {
	n.IDShort = n.GetIDShort()
	return nil
}

func (n *Node) AfterFind(tx *gorm.DB) (err error) {
	return n.Setup()
}

func (n *Node) Setup() error {
	signer, err := ssh.ParsePrivateKey(n.PrivateKey)
	if err != nil {
		return err
	}

	hostKey, err := ssh.ParsePublicKey(n.HostKey)
	if err != nil {
		return err
	}

	n.Config = &ssh.ClientConfig{
		User: n.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	return nil
}
