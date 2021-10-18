package assignment02IBC

import (
  "strconv"
  "crypto/sha256"
  "fmt"
)

const miningReward = 100
const rootUser = "Satoshi"

type BlockData struct {
  Title    string
  Sender   string
  Receiver string
  Amount   int
}

type Block struct {
  Data        []BlockData
  PrevPointer *Block
  PrevHash    string
  CurrentHash string
}

func CalculateBalance(userName string, chainHead *Block) int {
  balance := 0
  for chainHead != nil {
    for i := 0; i < len(chainHead.Data); i++ {
      if (chainHead.Data[i].Sender == userName) {
        balance -= chainHead.Data[i].Amount
      } else if (chainHead.Data[i].Receiver == userName) {
        balance += chainHead.Data[i].Amount
      }
    }
    chainHead = chainHead.PrevPointer
	}
  return balance
}

func CalculateHash(inputBlock *Block) string {
  dataString := ""
  for i := 0; i < len(inputBlock.Data); i++ {
    dataString += (inputBlock.Data[i].Title + inputBlock.Data[i].Sender +
      inputBlock.Data[i].Receiver + strconv.Itoa(inputBlock.Data[i].Amount))
  }
  return fmt.Sprintf("%x", sha256.Sum256([]byte(dataString)))
}

func VerifyTransaction(transaction *BlockData, chainHead *Block) bool {
  balance := CalculateBalance(transaction.Sender, chainHead);
  if (transaction.Title != "Coinbase" && balance < transaction.Amount) {
    fmt.Println("\nInvalid Transaction: Amount exceeds balance!");
    return false;
  }
  return true;
}

func InsertBlock(blockData []BlockData, chainHead *Block) *Block {

  // Validating Transactions in Chain
  for i :=0; i < len(blockData); i++ {
    if (!VerifyTransaction(&blockData[i], chainHead)) {
      return chainHead;
    }
  }
  // Validating Transactions in BlockData
  for i :=0; i < len(blockData); i++ {
    // Adding a temporary block without updating chainHead
    var tempBlock* Block = new(Block)
    var tempBlockData []BlockData
    for j := 0; j < i; j++ {
      tempBlockData = append(tempBlockData, blockData[j])
    }
    tempBlock.Data = tempBlockData

    if chainHead == nil {
  		tempBlock.PrevPointer = nil
  		tempBlock.PrevHash = "g3n3515"
  	} else {
      tempBlock.PrevPointer = chainHead
  		tempBlock.PrevHash = chainHead.CurrentHash
  	}
  	tempBlock.CurrentHash = CalculateHash(tempBlock)
    for j := 0; j < i; j++ {
      if (!VerifyTransaction(&tempBlockData[j], tempBlock)) {
        return chainHead;
      }
    }
  }

  var newBlock *Block = new(Block)
	newBlock.Data = blockData

  if chainHead == nil {
    newBlock.PrevPointer = nil
		newBlock.PrevHash = "g3n3515"
	} else {
    newBlock.PrevPointer = chainHead
    newBlock.PrevHash = chainHead.CurrentHash
	}
	newBlock.CurrentHash = CalculateHash(newBlock)

  fmt.Println("\nAlert: Block inserted successfully!");
	return newBlock
}


func ListBlocks(chainHead *Block) {
  for chainHead != nil {
    fmt.Println("\nBlock - ")
    fmt.Println("Data:        ", chainHead.Data)
    fmt.Println("PrevPointer: ", &chainHead.PrevPointer)
    fmt.Println("PrevHash:    ", chainHead.PrevHash)
    fmt.Println("CurrentHash: ", chainHead.CurrentHash)
    chainHead = chainHead.PrevPointer
  }
}

func VerifyChain(chainHead *Block) {
  for chainHead != nil {
    if chainHead.CurrentHash != CalculateHash(chainHead) {
      fmt.Println("\nInvalid Chain: Current Hash of block is tampered!")
			return
		}
    if (chainHead.PrevPointer != nil && chainHead.PrevPointer.CurrentHash != chainHead.PrevHash) {
      fmt.Println("\nInvalid Chain: CurrentHash of Previous block doesn't match the value of PrevHash in the current block!")
      fmt.Println(chainHead.Data)
			return
    }
    chainHead = chainHead.PrevPointer
	}
  fmt.Println("\nAlert! Blockchain is valid!")
}

func PremineChain(chainHead *Block, numBlocks int) *Block {
  blockData := []BlockData{{Title: "Coinbase", Sender: "System", Receiver: "Satoshi", Amount: 100}}
  for i := 0; i < numBlocks; i++ {
    chainHead = InsertBlock(blockData, chainHead)
  }
  return chainHead
}
