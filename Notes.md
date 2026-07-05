## 1
- **go.work**, bir monorepo'da (tek depoda birden fazla Go modülü) çalışırken, Go'nun tüm modülleri aynı anda görmesinive tanımasını sağlayan bir dosyadır.




---

Bir örnek:


protoc --go_out=./çıktı --go-grpc_out=./çıktı proto/auth/v1/auth.proto

Bu komut şunu der:

protoc: Derleyiciyi çalıştır.

--go_out=./çıktı: .pb.go dosyalarını ./çıktı dizinine yaz.

--go-grpc_out=./çıktı: *_grpc.pb.go dosyalarını ./çıktı dizinine yaz.

proto/auth/v1/auth.proto: Derlenecek proto dosyası.

## 2
- `gen/go` altındaki dosyalar hata veriyordu çünkü bu dosyalar `google.golang.org/grpc` ve `google.golang.org/protobuf` paketlerini kullanıyordu, ama `gen/go/go.mod` içinde bu bağımlılıklar tanımlı değildi.
- Ayrıca `go.work` sadece `./gen/go` modülünü içeriyordu. Bu yüzden kök dizinde `go test ./...` çalıştırınca workspace hatası alınıyordu.

### Çözüm
- `go.work` içine kök modül de eklendi:
  - `./`
  - `./gen/go`
- `gen/go/go.mod` içine gerekli paketler eklendi:
  - `google.golang.org/grpc`
  - `google.golang.org/protobuf`

### Sonuç
- Kök dizinde `go test ./...` çalıştı.
- `gen/go` altındaki tüm paketler başarıyla derlendi.
