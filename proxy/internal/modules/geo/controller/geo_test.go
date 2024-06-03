package controller

// var (
// 	testLat = "12.123"
// 	testLng = "12.123"
// )

// func TestGeo_Geocode(t *testing.T) {
// 	geoMock := mocks.NewGeorer(t)
// 	geoMock.On("GeoCode", service.GeoCodeIn{Lat: testLat, Lng: testLng}).Return(service.GeoCodeOut{Lat: testLat, Lng: testLng})
// 	geo := &GeoController{
// 		geo:       geoMock,
// 		Responder: &responder.Respond{},
// 	}
// 	server := httptest.NewServer(http.HandlerFunc(geo.Geocode))
// 	r := GeocodeRequest{
// 		Lat: testLat,
// 		Lng: testLng,
// 	}

// 	body, _ := json.Marshal(r)

// 	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	client := http.DefaultClient

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer resp.Body.Close()

// 	assert.Equal(t, http.StatusOK, resp.StatusCode)

// }

// func TestGeo_Geocode_BadRequest(t *testing.T) {
// 	logger := logs.NewLogger(config.NewAppConf(), os.Stdout)
// 	respond := responder.NewResponder(logger)
// 	geo := GeoController{
// 		Responder: respond,
// 	}

// 	req := map[string]interface{}{"lat": 123}
// 	reqJSON, _ := json.Marshal(req)

// 	s := httptest.NewServer(http.HandlerFunc(geo.Geocode))
// 	defer s.Close()

// 	resp, err := http.Post(s.URL, "application/json", bytes.NewBuffer(reqJSON))
// 	if err != nil {
// 		t.Fatal("ошибка при выполнении тестового запроса:", err.Error())
// 	}
// 	defer resp.Body.Close()

// 	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
// }

// func TestGeo_Search(t *testing.T) {
// 	logger := logs.NewLogger(config.NewAppConf(), os.Stdout)
// 	respond := responder.NewResponder(logger)
// 	mockGeo := mocks.NewGeorer(t)
// 	mockGeo.On("SearchAddresses", service.SearchAddressesIn{Query: "test"}).Return(service.SearchAddressesOut{Address: models.Address{Lat: testLat, Lon: testLng}})
// 	geo := &GeoController{
// 		geo:       mockGeo,
// 		Responder: respond,
// 	}
// 	server := httptest.NewServer(http.HandlerFunc(geo.Search))
// 	r := SearchRequest{
// 		Query: "test",
// 	}

// 	body, _ := json.Marshal(r)

// 	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	client := http.DefaultClient

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer resp.Body.Close()

// 	assert.Equal(t, http.StatusOK, resp.StatusCode)
// }

// func TestGeo_Search_Error(t *testing.T) {
// 	logger := logs.NewLogger(config.NewAppConf(), os.Stdout)
// 	respond := responder.NewResponder(logger)

// 	geoMock := mocks.NewGeorer(t)

// 	geoMock.On("SearchAddresses", service.SearchAddressesIn{Query: "BadQuery"}).Return(service.SearchAddressesOut{Err: errors.New("error")})

// 	geo := &GeoController{
// 		geo:       geoMock,
// 		Responder: respond,
// 	}

// 	searchReq := SearchRequest{
// 		Query: "BadQuery",
// 	}

// 	reqBody, _ := json.Marshal(searchReq)

// 	s := httptest.NewServer(http.HandlerFunc(geo.Search))
// 	defer s.Close()

// 	resp, err := http.Post(s.URL, "application/json", bytes.NewBuffer(reqBody))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer resp.Body.Close()

// 	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
// }
