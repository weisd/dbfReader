package dbfReader

import(
    "crypto/tls"
    "fmt"
    "os"
	"git.apache.org/thrift.git/lib/go/thrift"
    "github.com/weisd/dbfReader"
    "github.com/weisd/dbfReader/gen-go/dbfReader/rpc"
    
)

// 实现thrift中的接口

type RpcServiceHandle struct {
}

func (this *RpcServiceHandle) GetRecords(path string) (list []map[string]string, err error){

    list = []map[string]string{}
    fileName := "/Users/weisd/gocode/src/github.com/weisd/dbfReader/show2003.dbf"
    fp, err := os.Open(fileName)
    if err != nil {
        return
    }
    defer fp.Close()
    list = dbfReader.GetRecords(fp)
    return
}


// 启动服务

func RunServer(transportFactory thrift.TTransportFactory, protocolFactory thrift.TProtocolFactory, addr string, secure bool) error {
	var transport thrift.TServerTransport
	var err error
	if secure {
		cfg := new(tls.Config)
		if cert, err := tls.LoadX509KeyPair("server.crt", "server.key"); err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
		} else {
			return err
		}
		transport, err = thrift.NewTSSLServerSocket(addr, cfg)
	} else {
		transport, err = thrift.NewTServerSocket(addr)
	}
	
	if err != nil {
		return err
	}
	fmt.Printf("%T\n", transport)
	handler := &RpcServiceHandle{}
	processor := rpc.NewRpcServiceProcessor(handler)
	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)

	fmt.Println("Starting the simple server... on ", addr)
	return server.Serve()
}
