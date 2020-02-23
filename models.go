package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

type Provinsi []ProvinsiElement
type Kecamatan []KecamatanElement
type Kabupaten []KabupatenElement
type Kelurahan []KelurahanElement

func UnmarshalProvinsi(data []byte) (Provinsi, error) {
	var r Provinsi
	err := json.Unmarshal(data, &r)
	return r, err
}
func (r *Provinsi) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
type ProvinsiElement struct {
	ID       int64  `json:"id"`
	Provinsi string `json:"provinsi"`
	Bsni     string `json:"p_bsni"`
}

func UnmarshalKabupaten(data []byte) (Kabupaten, error) {
	var r Kabupaten
	err := json.Unmarshal(data, &r)
	return r, err
}
func (r *Kabupaten) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
type KabupatenElement struct {
	ID      int64  `json:"id"`      
	ProvId  int64  `json:"prov_id"` 
	KabKota string `json:"kab_kota"`
	Ibukota string `json:"ibukota"` 
	Bsni   string `json:"k_bsni"`  
}

func UnmarshalKecamatan(data []byte) (Kecamatan, error) {
	var r Kecamatan
	err := json.Unmarshal(data, &r)
	return r, err
}
func (r *Kecamatan) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
type KecamatanElement struct {
	ID       int64  `json:"id"`       
	KabkotId int64  `json:"kabkot_id"`
	Kec      string `json:"kec"`      
}

func UnmarshalKelurahan(data []byte) (Kelurahan, error) {
	var r Kelurahan
	err := json.Unmarshal(data, &r)
	return r, err
}
func (r *Kelurahan) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
type KelurahanElement struct {
	ID    int64  `json:"id"`    
	KecId int64  `json:"kec_id"`
	Kelurahan  string `json:"kelu"`  
	Pos   int64  `json:"pos"`   
}


func loadProvinsi() Provinsi {
	filePath, _ := filepath.Abs("./dirty/provinsi.json")
	file, _ := ioutil.ReadFile(filePath)

	provinsi, _ := UnmarshalProvinsi(file)

	return provinsi
}
func loadKabupaten() Kabupaten {
	filePath, _ := filepath.Abs("./dirty/kabupaten.json")
	file, _ := ioutil.ReadFile(filePath)

	kabupaten, _ := UnmarshalKabupaten(file)

	return kabupaten
}
func loadKecamatan() Kecamatan {
	filePath, _ := filepath.Abs("./dirty/kecamatan.json")
	file, _ := ioutil.ReadFile(filePath)

	kecamatan, _ := UnmarshalKecamatan(file)

	return kecamatan
}
func loadKelurahan() Kelurahan {
	filePath, _ := filepath.Abs("./dirty/kelurahan.json")
	file, _ := ioutil.ReadFile(filePath)

	kelurahan, _ := UnmarshalKelurahan(file)

	return kelurahan
}
