package structur

import "time"

type SliceUserRequest struct {
	Nam  string `json:"nam"`
	Unm  string `json:"unm"`
	Pass string `json:"pass"`
	Rlcd string `json:"rlcd"`
	Rlnm string `json:"rlnm"`
	Almt string `json:"almt"`
	Gdr  string `json:"gdr"`
	Pn   string `json:"pn"`
	Spcd string `json:"spcd"`
	Spnm string `json:"spnm"`
}

type SliceShopRequest struct {
	Spnm string `json:"spnm"`
	Almt string `json:"almt"`
	Spcd string `json:"spcd"`
}

type SliceProductRequest struct {
	Pnm   string `json:"pnm"`
	Pcd   string `json:"pcd"`
	Qty   int64  `json:"qty"`
	Price int64  `json:"price"`
	Catcd string `json:"catcd"`
	Catnm string `json:"catnm"`
	Spcd  string `json:"spcd"`
	Spnm  string `json:"spnm"`
	Crby  string `json:"crby"`
}

type SizeGetDataRequest struct {
	// user
	Nam string `json:"nam"`
	Unm string `json:"unm"`

	// shop
	Spcd string `json:"spcd"`
	Spnm string `json:"spnm"`

	// product
	Pnm   string `json:"pnm"`
	Pcd   string `json:"pcd"`
	Catcd string `json:"catcd"`

	// pages
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type ChangePasswordRequest struct {
	Unm     string `json:"unm"`
	Spcd    string `json:"spcd"`
	Pass    string `json:"pass"`
	NewPass string `json:"newPass"`
}

type CreateProductRequest struct {
	Data []SliceProductRequest `json:"data"`
}

type Token struct {
	Name string
	Role string
	Spcd string
	Exp  time.Time
}

type LoginRequest struct {
	Unm  string `json:"unm"`
	Pass string `json:"pass"`
}

type SliceTransactionRequest struct {
	Price       int64  `json:"price"`
	Total       int64  `json:"total"`
	Quantity    int64  `json:"qty"`
	ProductCode string `json:"pcd"`
	ShopCode    string `json:"spcd"`
}

type CreateTransactionRequest struct {
	Data []SliceTransactionRequest `json:"data"`
}
