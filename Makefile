
all: pump_sk.go imitator_sk.go
	go build uniset-go-example.go controller.go imitator.go pump_sk.go imitator_sk.go

pump_sk.go: pump.src.xml
	./uniset2-codegen-go -l ./ -n Pump pump.src.xml

imitator_sk.go: imitator.src.xml
	./uniset2-codegen-go -l ./ -n Imitator imitator.src.xml
