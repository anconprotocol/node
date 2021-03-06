package handler

import (
	"bytes"
	"compress/gzip"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/0xPolygon/polygon-sdk/crypto"
	"github.com/0xPolygon/polygon-sdk/helper/keccak"
	"github.com/buger/jsonparser"
	"github.com/cosmos/iavl"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gin-gonic/gin"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"

	"github.com/mr-tron/base58/base58"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	dbm "github.com/tendermint/tm-db"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/anconprotocol/sdk"
	"github.com/anconprotocol/sdk/proofsignature"
)


type Commit struct {
	LastHash []byte
	Height   int64
}
type ProofHandler struct {
	*sdk.AnconSyncContext
	db dbm.DB

	LastCommit *Commit
	api        proofsignature.IavlProofAPI
	proofs     proofsignature.IavlProofService
	RootKey    string
	privateKey *ecdsa.PrivateKey
}

func (h *ProofHandler) Commit() (int64, string, error) {

	// p := fmt.Sprintf("%s/%s/hash", "/anconprotocol", h.RootKey)

	// parent, err := h.proofs.Hash(&emptypb.Empty{})
	// parentHash, _ := jsonparser.GetString(parent, "lastHash")
	// h.proofs.Set([]byte(p), []byte(parentHash))

	v, err := h.proofs.SaveVersion(&emptypb.Empty{})

	hash, err := jsonparser.GetString(v, "root_hash")
	version, err := jsonparser.GetInt(v, "version")
	lastHash := []byte(hash)
	blockNumber := cast.ToInt64(version)

	h.LastCommit = &Commit{Height: blockNumber, LastHash: lastHash}
	return cast.ToInt64(version), hash, err
}
func (h *ProofHandler) GetProofService() *proofsignature.IavlProofService {
	return &h.proofs

}

func (h *ProofHandler) GetProofAPI() *proofsignature.IavlProofAPI {
	return &h.api

}
func NewProofHandler(ctx *sdk.AnconSyncContext, privateKeyPath string) *ProofHandler {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	privateKey, err := crypto.GenerateOrReadPrivateKey(privateKeyPath)
	if err != nil {
//		panic(err)
	}

	// os.OpenFile(,,privateKeyPath)

	folder := filepath.Join(userHomeDir, dbPath)
	db, err := dbm.NewGoLevelDB(dbName, folder)
	if err != nil {
		panic(err)
	}
	proofs, _ := proofsignature.NewIavlAPI(ctx.Store, ctx.Exchange, db, 2000, 0)
	return &ProofHandler{AnconSyncContext: ctx, db: db, proofs: *proofs.Service, api: *proofs, privateKey: privateKey}

}
func (h *ProofHandler) VerifyGenesis(moniker string, key string) ([]byte, error) {

	version := 0
	tree, err := iavl.NewMutableTree(h.db, int(2000))
	if err != nil {
		return nil, errors.Wrap(err, "unable to create iavl tree")
	}
	if _, err = tree.LoadVersion(int64(version)); err != nil {
		return nil, errors.Wrapf(err, "unable to load version %d", version)
	}
	key = fmt.Sprintf("%s%s", moniker, key)

	_, v, err := tree.GetWithProof([]byte(key))
	if err != nil && v != nil {
		return nil, errors.Wrap(err, "Unable to get with proof")
	}

	bz := tree.Hash()
	err = v.Verify(bz)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get rawkey")
	}

	return bz, nil
}

func InitGenesis(hostName string, moniker string, cidlink datamodel.Link, priv *ecdsa.PrivateKey) (string, string, error) {
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
		return " ", " ", errors.Wrap(err, "unable to create iavl tree")
	}
	if _, err = tree.LoadVersion(int64(version)); err != nil {
		return " ", " ", errors.Wrapf(err, "unable to load version %d", version)
	}

	// cidlink := sdk.CreateCidLink(signed[0:32])

	key := fmt.Sprintf("%s%s", moniker, cidlink.String())
	value := fmt.Sprintf(
		`{
		data: "%s",
		signature: "%s",
		}`,
		hostName, cidlink.String(),
	)

	tree.Set([]byte(key), []byte(value))

	_, _, err = tree.SaveVersion()

	if err != nil {
		return " ", " ", errors.Wrap(err, "Unable to commit")
	}

	_, proof, err := tree.GetWithProof([]byte(key))
	if err != nil {
		return " ", " ", errors.Wrap(err, "Unable to get with proof")
	}
	var proofData []byte
	proofData, err = proof.ToProto().Marshal()
	// err = proto.Unmarshal(proofData, )
	if err != nil {
		return " ", " ", errors.Wrap(err, "Unable to marshal")
	}

	hash := tree.Hash()
	rawKey, err := crypto.MarshalPrivateKey(priv)

	if err != nil {
		return " ", " ", errors.Wrap(err, "Unable to get rawkey")
	}
	stringRawKey := hexutil.Encode(rawKey)

	message := fmt.Sprintf(
		`Ancon protocol node initialize with: 
		*Sep256k1 private key: %s
		*Genesis value: %s
		*Genesis key: %s
		*Proof: %s
		*Last header hash: %s
		`, stringRawKey, hex.EncodeToString([]byte(value)), cidlink.String(), hex.EncodeToString(proofData), hex.EncodeToString(hash),
	)

	return message, key, nil
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
	sig := c.Query("sig")

	if sig == "true" {
		var digest []byte
		// priv, err := crypto.GenerateKey()
		keccak.Keccak256(digest, []byte(lastHash))
		signed, err := dagctx.PrivateKey.Sign(rand.Reader, digest, nil) //priv.Sign(rand.Reader, digest, nil)
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("Sig query Error %v", err).Error(),
			})
			return
		}
		c.JSON(201, gin.H{
			"lastHash":  lastHash,
			"signature": signed,
		})
		return
	}

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
// @Router /v0/proofs [post]
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
// @Router /v0/proofs/get/{path} [get]
func (dagctx *ProofHandler) Read(c *gin.Context) {

	key, _ := c.Params.Get("key")

	if key == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing key").Error(),
		})
		return
	}
	height, _ := c.GetQuery("height")

	if height == "" {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("missing height").Error(),
		})
		return
	}
	version := cast.ToInt64(height)
	internalKey, _ := base64.StdEncoding.DecodeString(key)
	data, err := dagctx.proofs.GetCommitmentProof([]byte(internalKey), version)

	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("decode Error %v", err).Error(),
		})
		return
	}

	exportAs, _ := c.GetQuery("export")
	if exportAs == "qr" {
		qrc, err := qrcode.New(string(data))
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("decode Error %v", err).Error(),
			})
			return
		}

		bg := c.Query("bgcolor")
		if bg == "" {
			bg = "#ffffff"
		} else {
			bg = "#" + bg
		}
		fg := c.Query("fgcolor")
		if fg == "" {
			fg = "#000000"
		} else {
			fg = "#" + fg
		}
		buf := &bytes.Buffer{}
		buf2 := &bytes.Buffer{}
		wr := gzip.NewWriter(buf)

		w := standard.NewWithWriter(wr,
			standard.WithBuiltinImageEncoder(standard.PNG_FORMAT),
			standard.WithBgColorRGBHex(bg),
			standard.WithFgColorRGBHex(fg),
		)
		qrc.Save(w)
		w.Close()
		rdr, err := gzip.NewReader(buf)
		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("error %v", err).Error(),
			})
			return
		}

		data, err := io.ReadAll(rdr)
		buf2.Write(data)
		defer rdr.Close()

		if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Errorf("error %v", err).Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"qr": base64.StdEncoding.EncodeToString(buf2.Bytes()),
		})
	} else {
		c.JSON(200, data)
	}
}
