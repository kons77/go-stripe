### JSON `omitempty` patameter

`omitempty` means that if the field is empty (""), it will not be included in the JSON.

```
type jsonResponce struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"` 
	Content string `json:"content,omitempty"`
	ID      int    `json:"id,omitempty"`
}
```



## Way to read JSON from HTTP requets

```go
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576 // 1Mb is going to allow

	// protect from large request body
    r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes)) 

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	// ensure that only one JSON object is in the body
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only have a single JSON value")
	}
	return nil
}

```

And because we've actually received our data interface as a reference to a variable, we're actually just changing a pointer value. And this is a really clean way to read any kind of JSON we get as a request body. Assuming that that request body has only a single JSON value.  

> Переменная `data` передается **по указателю** (`interface{}` передается по ссылке), поэтому изменения, внесенные в нее, будут доступны вызывающему коду. Этот способ позволяет **гибко** читать JSON любых форматов (структурированные данные, массивы, мапы), так как `data` может быть любого типа. Ограничение на **один JSON-объект** в запросе предотвращает ошибки при чтении **потенциально неправильных данных**.

