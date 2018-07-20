
all:
	go get -u github.com/xiaomingfuckeasylife/job/db
	go get github.com/astaxie/beego
	go get github.com/skip2/go-qrcode
	go get github.com/elastos/Elastos.ELA.Utility/common
	go build -o frontend main.go

format:
	go fmt ./...

clean:
	rm -rf *.8 *.o *.out *.6
