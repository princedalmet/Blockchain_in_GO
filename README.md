# Blockchain_in_GO

![blockchain-3019120_640](https://github.com/user-attachments/assets/bf621c19-659f-4b3c-9567-a0196c7ecc3a)

Block Structure: A block contains data, a timestamp, a hash of the previous block and its own hash

Hashing: cryptographic Hashing SHA-256 to link block

Mechanism:Proof of work

Storing the Blockchain in memmory or a database

Networking:Nodes communicate to synchrnize the Blockchain

Networking peer to peer communication to sync blocks between nodes

proof of Stake

Transaction model


func main() {

    node := Node{Address: "localhost:3000"}
    
    node.Start()
    
    node.ConnectToPeer("localhost:3001")
    
    bc := InitBlockchain("my-address")
    
    tx := CreateTransaction(from, to, amount, privKey)
    
    node.Broadcast(Message{Type: "tx", Data: tx.Serialize()})
    
    stakingManager.AddStake("my-address", 100)
    
    if stakingManager.SelectValidator() == "my-address" {
    
    newBlock := CreateBlock(bc.Blocks[len(bc.Blocks)-1].Hash, txs)
    
    bc.AddBlock(newBlock)
    
    node.Broadcast(Message{Type: "block", Data: newBlock.Serialize()})
  
  } }
