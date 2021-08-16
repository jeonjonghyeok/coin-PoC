package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jeonjonghyeok/coin/blockchain"
	"github.com/jeonjonghyeok/coin/p2p"
	"github.com/jeonjonghyeok/coin/utils"
	"github.com/jeonjonghyeok/coin/wallet"
)

var port string

type url string

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type urlDiscription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type addTxPayload struct {
	To     string
	Amount int
}

type myWalletResponse struct {
	Address string `json:"address"`
}

type addPeerPayload struct {
	Address string `json:"address"`
	Port    string `json:"port"`
}

func documentation(w http.ResponseWriter, r *http.Request) {
	data := []urlDiscription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         url("/status"),
			Method:      "GET",
			Description: "See the Status of the Blockchain",
		},
		{
			URL:         url("/blocks"),
			Method:      "GET",
			Description: "See All Block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add A Block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks{hash}"),
			Method:      "GET",
			Description: "See A Block",
		},
		{
			URL:         url("/balance{address}"),
			Method:      "GET",
			Description: "Get TxOuts for an Address",
		},
		{
			URL:         url("/transactions"),
			Method:      "POST",
			Description: "add transaction",
		},
		{
			URL:         url("/address"),
			Method:      "GET",
			Description: "get address",
		},
		{
			URL:         url("/ws"),
			Method:      "GET",
			Description: "Upgrade to WebSockets",
		},
	}
	json.NewEncoder(w).Encode(data)
}

func blocks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		utils.HandleErr(json.NewEncoder(w).Encode(blockchain.Blocks(blockchain.Blockchain())))
	case "POST":
		newBlock := blockchain.Blockchain().AddBlock()
		p2p.BroadcastNewBlock(newBlock)
		w.WriteHeader(http.StatusCreated)
	}

}
func block(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(w)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{fmt.Sprint(err)})

	} else {
		encoder.Encode(block)

	}
}

func status(w http.ResponseWriter, r *http.Request) {
	blockchain.Status(blockchain.Blockchain(), w)
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		next.ServeHTTP(w, r)
	})
}

func balance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	total := r.URL.Query().Get("total")
	switch total {
	case "true":
		amount := blockchain.BalanceByAddress(address, blockchain.Blockchain())
		json.NewEncoder(w).Encode(balanceResponse{address, amount})
	default:
		utils.HandleErr(json.NewEncoder(w).Encode(blockchain.UTxOutsByAddress(address, blockchain.Blockchain())))
	}
}
func mempool(w http.ResponseWriter, r *http.Request) {
	utils.HandleErr(json.NewEncoder(w).Encode(blockchain.Mempool.Txs))
}
func transactions(w http.ResponseWriter, r *http.Request) {
	var payload addTxPayload
	utils.HandleErr(json.NewDecoder(r.Body).Decode(&payload))
	err := blockchain.Mempool.AddTx(payload.To, payload.Amount)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{err.Error()})
		return
	}
	w.WriteHeader(http.StatusCreated)
}
func myWallet(w http.ResponseWriter, r *http.Request) {
	address := wallet.Wallet().Address
	utils.HandleErr(json.NewEncoder(w).Encode(myWalletResponse{address}))
}

func peers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var payload addPeerPayload
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&payload))
		p2p.AddPeer(payload.Address, payload.Port, port)
		w.WriteHeader(http.StatusOK)

	case "GET":
		utils.HandleErr(json.NewEncoder(w).Encode(p2p.AllPeers(&p2p.Peers)))
	}
}
func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)
	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware, loggerMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status).Methods("GET", "POST")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/wallet", myWallet).Methods("GET")
	router.HandleFunc("/transactions", transactions).Methods("POST")
	router.HandleFunc("/ws", p2p.Upgrade).Methods("GET")
	router.HandleFunc("/peers", peers).Methods("GET", "POST")
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe("localhost"+port, router))
}
