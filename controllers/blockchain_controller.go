package controllers

import (
	"Blockchain/models"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"io"
	"net/http"
	"time"
)

// Blockchain controller model
type BlockchainController struct {
	Blockchain models.Blockchain
}

// Creates a new blockchain and returns are new blockchain controller
func NewBlockchainController() *BlockchainController {

	// Initialize new blockchain
	var blockchain models.Blockchain
	// Get current time
	t := time.Now()
	// Create Genesis Block
	blockchain.GenesisBlock = models.Block{
		Index:     0,
		Timestamp: t.String(),
		Data:      "",
		Hash:      "",
		PrevHash:  "",
	}

	// Append Genesis Block to Blockchain
	blockchain.Blocks = append(blockchain.Blocks, blockchain.GenesisBlock)
	// Return new blockchain controller
	return &BlockchainController{
		Blockchain: blockchain,
	}
}

func calculateHash(block models.Block) string {
	// record consists of the block index, timestamp, BPM, and previous hash
	record := string(block.Index) + block.Timestamp + block.Data + block.PrevHash
	// generate new SHA256 hash
	hash := sha256.New()
	// write the record as a slice of bytees
	hash.Write([]byte(record))
	// create the SHA256 checksum of the data
	hashed := hash.Sum(nil)
	// return the hexadecimal encoding of the hash
	return hex.EncodeToString(hashed)
}

// Generates and returns new block
func generateBlock(prevBlock models.Block, data string) models.Block {
	// Get current time
	t := time.Now()

	// Create new block
	newBlock := models.Block{
		Index:     prevBlock.Index + 1,
		Timestamp: t.String(),
		Data:      data,
		PrevHash:  prevBlock.Hash,
	}
	// Generate new hash from new block
	newBlock.Hash = calculateHash(newBlock)

	// Return new block
	return newBlock
}

// Checks if new block is valid
// Checks new blocks index, hash, and checks against previous block hash
func isBlockValid(newBlock, prevBlock models.Block) bool {
	// Check the index of the new block against the next block
	if prevBlock.Index+1 != newBlock.Index {
		return false
	}
	// Check if previous hash of the new block matches the previous block hash
	if prevBlock.Hash != newBlock.PrevHash {
		return false
	}
	// Check if the hash of the new block is valid
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

// Handler for displaying blockchain
func (bc *BlockchainController) HandleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(bc.Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}



// Handler for writing to the blockchain
func (bc *BlockchainController) HandleWriteBlock(w http.ResponseWriter, r *http.Request) {
	// Create new message data
	var m models.Message

	// Decode response data
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	// Generate new block using response data
	newBlock := generateBlock(bc.Blockchain.Blocks[len(bc.Blockchain.Blocks)-1], m.Data)

	// Ensure new block is valid
	if isBlockValid(newBlock, bc.Blockchain.Blocks[len(bc.Blockchain.Blocks)-1]) {
		// Create new blockchain
		var newBlockchain models.Blockchain
		// Append blocks to new blockchain
		newBlockchain.Blocks = append(bc.Blockchain.Blocks, newBlock)
		// Replace current blockchain if block length of new blockchain is greater
		bc.replaceChain(newBlockchain)
		spew.Dump(bc.Blockchain)
	}

	// Display JSON response
	respondWithJSON(w, r, http.StatusCreated, newBlock)

}

func (bc *BlockchainController) replaceChain(blockchain models.Blockchain) {
	// Check if new blockchain length is greater than current blockchain length
	if len(blockchain.Blocks) > len(bc.Blockchain.Blocks) {
		bc.Blockchain = blockchain
	}
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	// Format JSON response data
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	// Write response data
	w.WriteHeader(code)
	w.Write(response)
}
