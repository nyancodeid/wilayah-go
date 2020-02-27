package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

type ArrayObject []map[string]interface{}
type Object map[string]interface{}
type WilayahObject struct {
	DataProvinsi  Provinsi
	DataKabupaten Kabupaten
	DataKecamatan Kecamatan
	DataKelurahan KelurahanElement
}
type WriteChanData struct {
	Data []byte
	Path string
}

func main() {
	defer elapsed("running")()
	os.RemoveAll("./dist")
	makeDir("./dist/find")
	makeDir("./dist/provinsi")

	provinsi := loadProvinsi()
	kabupaten := loadKabupaten()
	kecamatan := loadKecamatan()
	kelurahan := loadKelurahan()

	generateProvinsi(provinsi)
	generateKabupaten(kabupaten, provinsi)
	generateKecamatan(kecamatan, kabupaten)
	generateKelurahan(provinsi, kabupaten, kecamatan, kelurahan)
	generateWilayah(provinsi, kabupaten, kecamatan, kelurahan)
}

func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", what, time.Since(start))
	}
}
func generateProvinsi(datas Provinsi) {
	dirAllJson := "./dist/provinsi.json"
	dataProvinsi := []map[string]interface{}{}

	for i := range datas {
		data := datas[i]
		dataProvinsi = append(dataProvinsi, map[string]interface{}{
			"id":   data.ID,
			"name": data.Provinsi,
			"bsni": data.Bsni,
		})

		dirPath := "./dist/provinsi/" + strconv.FormatInt(data.ID, 10)
		dirJson := "./dist/provinsi/" + strconv.FormatInt(data.ID, 10) + ".json"

		jsonData, _ := json.Marshal(map[string]interface{}{
			"id":   data.ID,
			"name": data.Provinsi,
			"bsni": data.Bsni,
		})

		go makeDir(dirPath)
		go writeFile(dirJson, jsonData)
	}

	provinsiDataJson, _ := json.Marshal(dataProvinsi)
	go writeFile(dirAllJson, provinsiDataJson)

	rateLog()
}
func generateKabupaten(datas Kabupaten, related Provinsi) {
	masterKabupaten := gubrak.From(datas).Map(func(each KabupatenElement) map[string]interface{} {
		return map[string]interface{}{
			"id":      each.ID,
			"pid":     each.ProvId,
			"name":    each.KabKota,
			"ibukota": each.Ibukota,
			"bsni":    each.Bsni,
		}
	}).GroupBy(func(each map[string]interface{}) int64 {
		return each["pid"].(int64)
	}).Result().(map[int64][]map[string]interface{})

	masterIbukota := gubrak.From(datas).Map(func(each KabupatenElement) map[string]interface{} {
		return map[string]interface{}{
			"id":   each.ID,
			"pid":  each.ProvId,
			"name": each.Ibukota,
			"bsni": each.Bsni,
		}
	}).GroupBy(func(each map[string]interface{}) int64 {
		return each["pid"].(int64)
	}).Result().(map[int64][]map[string]interface{})

	// Looping create directory for kabupaten
	for i := range datas {
		data := datas[i]

		dirPath := "./dist/provinsi/kabupaten/" + strconv.FormatInt(data.ID, 10)
		go makeDir(dirPath)
	}
	// Looping create json for kabupaten and ibukota
	for i := range related {
		data := related[i]

		jsonPathKabupaten := fmt.Sprintf("./dist/provinsi/%s/kabupaten.json",
			strconv.FormatInt(data.ID, 10))
		jsonPathIbukota := fmt.Sprintf("./dist/provinsi/%s/ibukota.json",
			strconv.FormatInt(data.ID, 10))

		jsonDataKabupaten := masterKabupaten[data.ID]
		jsonDataIbukota := masterIbukota[data.ID]

		jsonDataByteKab, _ := json.Marshal(jsonDataKabupaten)
		jsonDataByteIbu, _ := json.Marshal(jsonDataIbukota)

		go writeFile(jsonPathKabupaten, jsonDataByteKab)
		go writeFile(jsonPathIbukota, jsonDataByteIbu)
	}

	rateLog()
}
func generateKecamatan(datas Kecamatan, related Kabupaten) {
	masterKecamatan := gubrak.From(datas).Map(func(each KecamatanElement) map[string]interface{} {
		return map[string]interface{}{
			"id":    each.ID,
			"kabid": each.KabkotId,
			"name":  each.Kec,
		}
	}).GroupBy(func(each map[string]interface{}) int64 {
		return each["kabid"].(int64)
	}).Result().(map[int64][]map[string]interface{})

	for i := range datas {
		data := datas[i]
		dirPath := "./dist/provinsi/kabupaten/kecamatan/" + strconv.FormatInt(data.ID, 10)

		go makeDir(dirPath)
	}
	for i := range related {
		data := related[i]

		jsonPath := fmt.Sprintf("./dist/provinsi/kabupaten/%s/kecamatan.json",
			strconv.FormatInt(data.ID, 10))

		jsonDataByte, _ := json.Marshal(masterKecamatan[data.ID])

		go writeFile(jsonPath, jsonDataByte)
	}

	rateLog()
}
func generateKelurahan(provinsi Provinsi, kabupaten Kabupaten, kecamatan Kecamatan, datas Kelurahan) {
	masterKelurahan := gubrak.From(datas).Map(func(each KelurahanElement) map[string]interface{} {
		hash := makeHash("kec_" + strconv.FormatInt(each.ID, 10))

		return map[string]interface{}{
			"id":    each.ID,
			"kecid": each.KecId,
			"name":  each.Kelurahan,
			"hash":  hash,
		}
	}).GroupBy(func(each map[string]interface{}) int64 {
		return each["kecid"].(int64)
	}).Result().(map[int64][]map[string]interface{})

	for i := range datas {
		data := datas[i]
		hash := makeHash("kec_" + strconv.FormatInt(data.ID, 10))

		jsonPathHash := fmt.Sprintf("./dist/find/%s.json",
			hash)

		// Get related ID
		_kecamatan := gubrak.From(kecamatan).Find(func(each KecamatanElement) bool {
			return each.ID == data.KecId
		}).Result().(KecamatanElement)
		_kabupaten := gubrak.From(kabupaten).Find(func(each KabupatenElement) bool {
			return each.ID == _kecamatan.KabkotId
		}).Result().(KabupatenElement)
		_provinsi := gubrak.From(provinsi).Find(func(each ProvinsiElement) bool {
			return each.ID == _kabupaten.ProvId
		}).Result().(ProvinsiElement)

		jsonDataHash := map[string]interface{}{
			"id":    data.ID,
			"kecid": data.KecId,
			"name":  data.Kelurahan,
			"hash":  hash,
			"related": map[string]interface{}{
				"kecamatan": map[string]interface{}{
					"id":   data.KecId,
					"name": _kecamatan.Kec,
				},
				"kabupaten": map[string]interface{}{
					"id":      _kabupaten.ID,
					"name":    _kabupaten.KabKota,
					"ibukota": _kabupaten.Ibukota,
				},
				"provinsi": map[string]interface{}{
					"id":   _provinsi.ID,
					"name": _provinsi.Provinsi,
				},
			},
		}

		jsonDataHashByte, _ := json.Marshal(jsonDataHash)

		go writeFile(jsonPathHash, jsonDataHashByte)
	}
	for i := range kecamatan {
		data := kecamatan[i]

		jsonPath := fmt.Sprintf("./dist/provinsi/kabupaten/kecamatan/%s/kelurahan.json",
			strconv.FormatInt(data.ID, 10))

		jsonDataByte, _ := json.Marshal(masterKelurahan[data.ID])

		go writeFile(jsonPath, jsonDataByte)
	}

	rateLog()
}

func generateWilayah(provinsi Provinsi, kabupaten Kabupaten, kecamatan Kecamatan, kelurahan Kelurahan) {
	var dataWilayah ArrayObject
	var wg sync.WaitGroup

	queue := make(chan Object, 1)

	// Create our data and send it into the queue.
	wg.Add(len(kelurahan))

	for i := range kelurahan {
		data := kelurahan[i]
		var wilayah WilayahObject

		wilayah.DataProvinsi = provinsi
		wilayah.DataKabupaten = kabupaten
		wilayah.DataKecamatan = kecamatan
		wilayah.DataKelurahan = data

		go func(data WilayahObject) {
			hash := makeHash("kec_" + strconv.FormatInt(data.DataKelurahan.ID, 10))
			_provinsi, _kabupaten, _kecamatan, _kelurahan := _generateDataWilayah(data)

			wilayah := fmt.Sprintf("%s, %s, %s, %s", _kelurahan.Kelurahan, _kecamatan.Kec, _kabupaten.KabKota, _provinsi.Provinsi)

			jsonData := Object{
				"id":   _kelurahan.ID,
				"name": wilayah,
				"hash": hash,
				"related": map[string]interface{}{
					"kecamatan": map[string]interface{}{
						"id":   _kelurahan.KecId,
						"name": _kecamatan.Kec,
					},
					"kabupaten": map[string]interface{}{
						"id":      _kabupaten.ID,
						"name":    _kabupaten.KabKota,
						"ibukota": _kabupaten.Ibukota,
					},
					"provinsi": map[string]interface{}{
						"id":   _provinsi.ID,
						"name": _provinsi.Provinsi,
					},
				},
			}

			queue <- jsonData
		}(wilayah)
	}

	go func() {
		// defer wg.Done() <- Never gets called since the 100 `Done()` calls are made above, resulting in the `Wait()` to continue on before this is executed
		for data := range queue {
			dataWilayah = append(dataWilayah, data)
			wg.Done() // ** move the `Done()` call here
		}
	}()

	wg.Wait()

	jsonDataByte, err := json.Marshal(dataWilayah)
	if err != nil {
		log.Fatalln(err)
	} else {
		writeFile("./dist/wilayah.json", jsonDataByte)
	}

	rateLog()
}

func rateLog() {
	fmt.Println("Writing operation: " + strconv.FormatInt(counter.Rate(), 10) + "s")
}
