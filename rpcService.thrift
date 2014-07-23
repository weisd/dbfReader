namespace go dbfReader.rpc
namespace php dbfReader.rpc

service RpcService {
    // 取记录列表
    list<map<string,string>> GetRecords(1:string filePath);
}
