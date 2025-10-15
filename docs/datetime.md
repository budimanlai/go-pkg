# DateTime Types Documentation

Paket types menyediakan custom time types untuk memastikan konsistensi dalam penanganan waktu, khususnya untuk JSON serialization/deserialization.

## UTCTime

`UTCTime` adalah custom type yang mengimplementasikan `time.Time` dengan behavior khusus untuk JSON:

- **MarshalJSON**: Selalu mengkonversi waktu ke UTC dan format RFC3339 dengan akhiran 'Z'
- **UnmarshalJSON**: Mem-parse string JSON dalam format RFC3339 dan menyimpannya sebagai UTC
- **String()**: Mengembalikan representasi string dalam format UTC RFC3339

### Mengapa UTCTime?

- Memastikan semua waktu disimpan dan dikirim dalam UTC
- Konsistensi format JSON (selalu dengan 'Z' untuk UTC)
- Mencegah masalah timezone dalam API

### Instalasi

Pastikan Anda memiliki dependensi standar Go (sudah termasuk dalam Go):

```go
import "time"
```

### Penggunaan

#### Deklarasi dan Inisialisasi

```go
import (
    "github.com/budimanlai/go-pkg/types"
    "time"
)

// Membuat UTCTime dari time.Time
now := time.Now()
utcTime := types.UTCTime(now)

// Atau langsung dari time.Date
specificTime := time.Date(2025, 10, 15, 12, 30, 45, 0, time.UTC)
utcTime := types.UTCTime(specificTime)
```

#### JSON Marshal

```go
import (
    "encoding/json"
    "github.com/budimanlai/go-pkg/types"
    "time"
)

type User struct {
    Name      string          `json:"name"`
    CreatedAt types.UTCTime   `json:"created_at"`
}

user := User{
    Name:      "John Doe",
    CreatedAt: types.UTCTime(time.Now()),
}

data, err := json.Marshal(user)
if err != nil {
    panic(err)
}

// Output: {"name":"John Doe","created_at":"2025-10-15T12:30:45Z"}
fmt.Println(string(data))
```

#### JSON Unmarshal

```go
jsonStr := `{"name":"John Doe","created_at":"2025-10-15T12:30:45Z"}`

var user User
err := json.Unmarshal([]byte(jsonStr), &user)
if err != nil {
    panic(err)
}

// user.CreatedAt sekarang berisi waktu dalam UTC
fmt.Println(user.CreatedAt) // 2025-10-15T12:30:45Z
```

#### String Representation

```go
utcTime := types.UTCTime(time.Now())
fmt.Println(utcTime.String()) // Output: 2025-10-15T12:30:45Z
```

### Behavior Khusus

#### Selalu UTC

UTCTime selalu mengkonversi waktu ke UTC, terlepas dari timezone asli:

```go
// Waktu dalam timezone Jakarta
loc, _ := time.LoadLocation("Asia/Jakarta")
jakartaTime := time.Date(2025, 10, 15, 19, 30, 45, 0, loc) // 19:30 WIB
utcTime := types.UTCTime(jakartaTime)

// Saat marshal, akan menjadi UTC (12:30)
data, _ := json.Marshal(utcTime)
// Output: "2025-10-15T12:30:45Z"
```

#### Format RFC3339 dengan Z

Format selalu menggunakan RFC3339 dan diakhiri dengan 'Z' untuk menunjukkan UTC:

```go
utcTime := types.UTCTime(time.Date(2025, 10, 15, 12, 30, 45, 0, time.UTC))
jsonBytes, _ := utcTime.MarshalJSON()
// string(jsonBytes) = "2025-10-15T12:30:45Z"
```

### Contoh Lengkap dalam Struct

```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/budimanlai/go-pkg/types"
    "time"
)

type Event struct {
    ID        int             `json:"id"`
    Title     string          `json:"title"`
    StartTime types.UTCTime   `json:"start_time"`
    EndTime   types.UTCTime   `json:"end_time"`
    CreatedAt types.UTCTime   `json:"created_at"`
}

func main() {
    // Membuat event
    event := Event{
        ID:        1,
        Title:     "Meeting",
        StartTime: types.UTCTime(time.Date(2025, 10, 15, 10, 0, 0, 0, time.UTC)),
        EndTime:   types.UTCTime(time.Date(2025, 10, 15, 11, 0, 0, 0, time.UTC)),
        CreatedAt: types.UTCTime(time.Now()),
    }

    // Marshal ke JSON
    data, err := json.Marshal(event)
    if err != nil {
        panic(err)
    }

    fmt.Println("JSON:", string(data))

    // Unmarshal dari JSON
    jsonStr := `{
        "id": 2,
        "title": "Workshop",
        "start_time": "2025-10-16T14:00:00Z",
        "end_time": "2025-10-16T16:00:00Z",
        "created_at": "2025-10-15T08:00:00Z"
    }`

    var newEvent Event
    err = json.Unmarshal([]byte(jsonStr), &newEvent)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Event: %+v\n", newEvent)
    fmt.Printf("Start Time: %s\n", newEvent.StartTime.String())
}
```

### Testing

Jalankan unit tests dengan:

```bash
go test ./types
```

Tests mencakup:
- MarshalJSON dengan format UTC yang benar
- UnmarshalJSON dari string JSON
- Method String() dengan format yang benar
- Round-trip marshal/unmarshal
- Konversi timezone ke UTC

### Catatan

- UTCTime mengimplementasikan `json.Marshaler` dan `json.Unmarshaler`
- Waktu selalu disimpan dalam UTC secara internal
- Format JSON selalu RFC3339 dengan akhiran 'Z'
- Nanoseconds tidak disertakan dalam format JSON (batasan RFC3339)
- Compatible dengan `time.Time` - bisa dikonversi dengan `time.Time(utcTime)`