module github.com/hugot/whapp-deltachat

go 1.13

require (
	github.com/Rhymen/go-whatsapp v0.1.0
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/hugot/go-deltachat v0.0.0-20200103100028-e93f1f3d8b97
	github.com/mdp/qrterminal v1.0.1
	github.com/skip2/go-qrcode v0.0.0-20191027152451-9434209cb086
	go.etcd.io/bbolt v1.3.3
	golang.org/x/crypto v0.0.0-20191227163750-53104e6ec876 // indirect
	gopkg.in/yaml.v2 v2.2.2
)

replace github.com/hugot/go-deltachat v0.0.0-20200103100028-e93f1f3d8b97 => ../go-deltachat
