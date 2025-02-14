// Block Structure: A block contains data, a timestamp, a hash of the previous block and its own hash
// Hashing: cryptographic Hashing SHA-256 to link block
// Mechanism:Proof of work
// Storing the Blockchain in memmory or a database
// Networking:Nodes communicate to synchrnize the Blockchain

//Defining the block structure
package main

import (
	"bytes"
	"crypto/sha256"
	//"encoding/hex"
	"fmt"
	"math/big"
	"time"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net"
	"sync"
)


type Block struct {
	Timestamp		int64
	Transactions	[]*Transaction
	Validator		string
	Data			[]byte
	PrevBlockHash	[]byte
	Hash			[]byte
	Nonce			int //proof of work

}

//Transaction o/p
type TXOutput struct {
	Value		int	//Amount
	PubKeyHash	[]byte	//public key
}


// Transaction Structure
type Transaction struct {

	ID		[]byte
	Vin		[]TXOutput
	Vout	[]TXOutput
}

type Node struct {
	Address 	string
	Server 		net.Listener
	Peers		[]string
	mu			sync.Mutex
}
//to Generate Block HAshes

func (b *Block) SetHash()	{
	data := bytes.Join(
		[][]byte{
			b.PrevBlockHash,
			b.Data,
			[]byte(string(b.Timestamp)),
			[]byte(string(b.Nonce)),
		},
		[]byte{},
	)
	hash := sha256.Sum256(data)
	b.Hash = hash[:]
}

// Proof of work

const targetBits = 16

func (b *Block) MineBlock() {
	target := big.NewInt(1)
	target.Lsh(target, uint(256 - targetBits))

	var hashInt big.Int
	nonce := 0

	for {

		b.Nonce = nonce
		b.SetHash()
		hashInt.SetBytes(b.Hash)

		if hashInt.Cmp(target) == -1  {// Hash < target
			break
		}
		nonce++
	}
}

//Blockchain Structure
//Blockchain is a linked list of blocks
type Blockchain struct {
	Blocks []*Block
}

//The first block in the chain

func NewGenesisBlock()	*Block {
	return &Block {
		Timestamp:			time.Now().Unix(),
		Data:				[]byte("Genesis Block"),
		PrevBlockHash:		[]byte{},
		Nonce:				0,
	}
}

//Initialize the Blockchain

func InitBlockchain() *Blockchain {
	genesis := NewGenesisBlock()
	genesis.MineBlock()
	return &Blockchain {Blocks: []*Block{genesis}}
}

//Adding new Block

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := &Block{
		Timestamp:		time.Now().Unix(),
		Data:			[]byte(data),
		PrevBlockHash:	prevBlock.Hash,
		Nonce:			0,
	}
	newBlock.MineBlock()
	bc.Blocks = append(bc.Blocks, newBlock)
}

// Checking the block chain

func (bc *Blockchain) IsValid()	bool {
	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		prevBlock := bc.Blocks[i-1]

		//previous hash
		if !bytes.Equal(currentBlock.PrevBlockHash, prevBlock.Hash) {
			return false
		}
		tempHash := currentBlock.Hash
		currentBlock.SetHash()
		if !bytes.Equal(tempHash, currentBlock.Hash) {
			return false
		}
	}
	return true
}


func (tx *Transaction) Serialize()	[]byte {
	data,	_ := json.Marshal(tx)
	return data
}

func NewCoinbaseTX(to	[]byte) *Transaction {
	txin := TXInput{[]byte{}, -1, nil, nil}
	txout := TXOutput{Value: 100, PubKeyHash: HashPubKey(to)}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.ID = tx.Hash()
	return &tx
}

func (tx *Transaction) Hash() []byte{
	hash := sha256.Sum256(tx.Serialize())
	return hash[:]
}

type StakingManager struct {
	Stakes map[string]int
}

func	(sm *StakingManager) AddStake(address string, amount int) {
	sm.Stakes[address] += amount
}

func (sm *StakingManager) selectvalidator() string {
	var validator string
	maxStake := -1

	for addr, stake := range sm.Stakes {
		if stake > maxStake {
			maxStake = stake
			validator = addr
		}
	}
	return validator	
}

func (b *Block) MineBlock(stake string)bool {
	if validateStaker(stake) {
		b.SetHash()
		return true
	}
	return false
}

func (n *Node) Start() error{
	ln, err	:= net.Listen("tcp", n.Address)
	if err != nil {
		return err
	}
	n.Server = ln 	

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				continue	
			}
			go n.HandleConnection(conn)
		}
	}()
	return nil
}

fun (n *Node) HandleConnection(conn net.Conn) {
	defer conn.Close()

	var msg Message
	decoder := gob.NewDecoder(conn)
	err := decoder.Decode(&msg)
	if err != nil {
		return
	}
	switch msg.Type {
	case "block":
		n.ProcessBlock(msg.Data)
	case "tx":
		n.ProcessTransaction(msg.Data)
	}
}	

func (n *Node) Broadcast(msg Message) {
	for _, peer := range n.Peers {
		conn, err := net.Dial("tcp",peer)
		if err != nil {
			continue
		}
		encoder := gob.NewEncoder(conn)
		encoder.encoder(msg)
		conn.Close()
	}	
}

type Message struct {
	Type string
	Data []byte
}

func main() {

	bc := InitBlockchain()
	bc.AddBlock("Block 1 Data")
	bc.AddBlock("Block 2 Data")

	if bc.IsValid() {

		fmt.Println("Blockchain is valid")
	}	else {
		fmt.Println("Blockchain is invalid!")
	}
}