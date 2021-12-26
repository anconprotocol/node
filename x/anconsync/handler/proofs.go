package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/0xPolygon/polygon-sdk/crypto"
	"github.com/0xPolygon/polygon-sdk/helper/keccak"
	"github.com/buger/jsonparser"
	"github.com/cosmos/iavl"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gin-gonic/gin"
	"github.com/mr-tron/base58/base58"
	"github.com/pkg/errors"
	dbm "github.com/tendermint/tm-db"

	"github.com/anconprotocol/sdk"
	"github.com/anconprotocol/sdk/proofsignature"
)

var GENESISKEY = "/anconprotocol/"

type ProofHandler struct {
	*sdk.AnconSyncContext
	db dbm.DB

	api    proofsignature.IavlProofAPI
	proofs proofsignature.IavlProofService
}

func (h *ProofHandler) GetProofService() *proofsignature.IavlProofService {
	return &h.proofs

}

func (h *ProofHandler) GetProofAPI() *proofsignature.IavlProofAPI {
	return &h.api

}
func NewProofHandler(ctx *sdk.AnconSyncContext) *ProofHandler {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	folder := filepath.Join(userHomeDir, dbPath)
	db, err := dbm.NewGoLevelDB(dbName, folder)
	if err != nil {
		panic(err)
	}
	proofs, _ := proofsignature.NewIavlAPI(ctx.Store, ctx.Exchange, db, 2000, 0)
	return &ProofHandler{AnconSyncContext: ctx, db: db, proofs: *proofs.Service, api: *proofs}

}
func (h *ProofHandler) VerifyGenesis(root, key string) error {

	version := 0
	tree, err := iavl.NewMutableTree(h.db, int(2000))
	if err != nil {
		return errors.Wrap(err, "unable to create iavl tree")
	}
	if _, err = tree.LoadVersion(int64(version)); err != nil {
		return errors.Wrapf(err, "unable to load version %d", version)
	}

	_, v, err := tree.GetWithProof([]byte(key))
	if err != nil {
		return errors.Wrap(err, "Unable to get with proof")
	}

	bz, err := hex.DecodeString(root)
	err = v.Verify(bz)
	if err != nil {
		return errors.Wrap(err, "Unable to get rawkey")
	}

	return nil
}

func InitGenesis(hostName string) (string, error) {
	userHomeDir, err := os.UserHomeDir()
	version := 0
	if err != nil {
		panic(err)
	}

	folder := filepath.Join(userHomeDir, dbPath)
	db, err := dbm.NewGoLevelDB(dbName, folder)
	if err != nil {
		panic(err)
	}

	tree, err := iavl.NewMutableTree(db, int(2000))
	if err != nil {
		return " ", errors.Wrap(err, "unable to create iavl tree")
	}
	if _, err = tree.LoadVersion(int64(version)); err != nil {
		return " ", errors.Wrapf(err, "unable to load version %d", version)
	}

	// Set your own keypair
	priv, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	var digest []byte

	keccak.Keccak256(digest, []byte(hostName))
	signed, err := priv.Sign(rand.Reader, digest, nil)

	if err != nil {
		return " ", errors.Wrap(err, "Unable to sign")
	}

	cidlink := sdk.CreateCidLink(signed[0:32])

	key := fmt.Sprintf("%s%s", GENESISKEY, cidlink.String())
	value := fmt.Sprintf(
		`{
		data: "%s",
		signature: "%s",
		}`,
		hostName, signed,
	)

	tree.Set([]byte(key), []byte(value))

	_, _, err = tree.SaveVersion()

	if err != nil {
		return " ", errors.Wrap(err, "Unable to commit")
	}

	_, proof, err := tree.GetWithProof([]byte(key))
	if err != nil {
		return " ", errors.Wrap(err, "Unable to get with proof")
	}
	var proofData []byte
	proofData, err = proof.ToProto().Marshal()
	// err = proto.Unmarshal(proofData, )
	if err != nil {
		return " ", errors.Wrap(err, "Unable to marshal")
	}

	hash := tree.Hash()
	rawKey, err := crypto.MarshalPrivateKey(priv)

	if err != nil {
		return " ", errors.Wrap(err, "Unable to get rawkey")
	}
	stringRawKey := hexutil.Encode(rawKey)

	message := fmt.Sprintf(
		`Ancon protocol node initialize with: 
		*Sep256k1 private key: %s
		*Genesis value: %s
		*Genesis key: %s
		*Proof: %s
		*Last header hash: %s
		`, stringRawKey, hex.EncodeToString([]byte(value)), key, hex.EncodeToString(proofData), hex.EncodeToString(hash),
	)

	return message, nil
}

func GenerateKeys() (string, error) {
	// Set your own keypair
	priv, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}

	rawKey, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return " ", errors.Wrap(err, "Unable to get private key")
	}
	pub := crypto.MarshalPublicKey(&priv.PublicKey)
	if err != nil {
		return " ", errors.Wrap(err, "Unable to get public key")
	}
	stringRawKey := hexutil.Encode(rawKey)
	publicKeyBase58 := base58.Encode(pub)
	pubhex := hexutil.Encode(pub)
	message := fmt.Sprintf(
		`Generated keys
		Sep256k1 private key (hex): %s
		Sep256k1 public key (hex): %s
		Sep256k1 public key (base58): %s`,
		stringRawKey, pubhex, publicKeyBase58,
	)

	return message, nil
}

// @BasePath /v0
// Verify godoc
// @Summary Verifies an ics23 proofs
// @Schemes
// @Description Verifies an ics23 proof
// @Tags proofs
// @Accept json
// @Produce json
// @Success 201 {string} cid
// @Router /v0/proofs/verify [post]
func (dagctx *ProofHandler) ReadCurrentRootHash(c *gin.Context) {

	lastHash, err := dagctx.proofs.Hash(nil)

	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("lastHash Error %v", err).Error(),
		})
		return
	}
	c.JSON(201, gin.H{
		"lastHash": lastHash,
	})
}

// @BasePath /v0
// Create godoc
// @Summary Create
// @Schemes
// @Description Writes an ics23 proof
// @Tags proofs
// @Accept json
// @Produce json
// @Success 201 {string} cid
// @Router /v0/dagjson [post]
func (dagctx *ProofHandler) Create(c *gin.Context) {

	v, _ := c.GetRawData()

	key, _ := jsonparser.GetString(v, "key")

	if key == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing key").Error(),
		})
		return
	}

	value, _ := jsonparser.GetString(v, "value")

	if value == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing value").Error(),
		})
		return
	}

	message, err := dagctx.proofs.Set([]byte(key), []byte(value))

	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("decode Error %v", err).Error(),
		})
		return
	}

	c.JSON(201, message)
	//	impl.PushBlock(c.Request.Context(), dagctx.Exchange, dagctx.IPFSPeer, cid)
}

// @BasePath /v0
// Read godoc
// @Summary Reads an existing proof
// @Schemes
// @Description Returns JSON
// @Tags proofs
// @Accept json
// @Produce json
// @Success 200
// @Router /v0/dagjson/{cid}/{path} [get]
func (dagctx *ProofHandler) Read(c *gin.Context) {

	v, _ := c.GetRawData()

	key, _ := jsonparser.GetString(v, "key")

	if key == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing key").Error(),
		})
		return
	}

	data, err := dagctx.proofs.GetWithProof([]byte(key))

	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("decode Error %v", err).Error(),
		})
		return
	}
	c.JSON(200, data)
}
