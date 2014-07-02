package mtconfig

import (
    "../zk"
    "io/ioutil"
    "net/http"
    "strings"
    "time"
)

type MtConfigServer struct {
    zklist []string
    cache_config map[string]map[string]string   //path-keys-values
    zookeeper *zk.Conn
}


func NewMtConfigServer(config_url string) (*MtConfigServer, error) {
    if(config_url == "") {
        config_url = "http://xxx.xxx.com/api/zkserverlist"
    }

    //try to get zklist by http get request
    res, err := http.Get(config_url)
    if(err != nil) {
        return nil, err
    }

    bodyByte, _ := ioutil.ReadAll(res.Body)

    res.Body.Close()

    ipList := strings.TrimSpace(string(bodyByte))
    ipArray := strings.Split(ipList, ",")

    zkCli, _, err := zk.Connect(ipArray, time.Second* 10)

    if(err != nil) {
        return nil, err
    }


    return &MtConfigServer{
                zklist: ipArray,
                zookeeper: zkCli,
                cache_config: make(map[string]map[string]string),
            }, err
}

func (p *MtConfigServer) Get(domain string, node string, key string) (string, error) {
    path := "/config/" + domain + "/" + node
    value, err := p.findInCache(path, key)

    if(value == "") {
        //find it in zk
        err = p.getFromZk(path)
        if(err != nil){
            return "", err
        }
    }else {
        return value, err
    }

    value, err = p.findInCache(path, key)

    return value, err
}

func (p *MtConfigServer) findInCache(path_ string, key_ string) (string, error) {
    //try to find config from cache_config
    for path, mapItem := range p.cache_config {
        if(path == path_){
            for key, value := range mapItem {
                if(key == key_){
                    return value, nil
                }
            }
        }
    }

    return "", nil
}

func (p *MtConfigServer) deleteFromCache(path string) {
    delete(p.cache_config, path)
}

func (p* MtConfigServer) getFromZk(path string) error {
    //first to delete the path cache
    p.deleteFromCache(path)

    p.zookeeper.SetCallBackFunc(p.onCallBack)
    //get data from zk
    data, _, _, err := p.zookeeper.GetW(path)
    
    if(err != nil) {
        return err
    }
    
    //parse the data //
    trimData := strings.TrimSpace(string(data))
    dataArry := strings.Split(trimData, "\n\n")
    
    item := make(map[string]string)

    for _, data := range dataArry {
        subDataArray := strings.Split(data, "=")
        
        item[strings.TrimSpace(string(subDataArray[0]))] = string(strings.TrimSpace(subDataArray[1]))
    }

    p.cache_config[path] = item

    return nil
}


func (p *MtConfigServer) onCallBack(path string) {
    p.deleteFromCache(path)
}

