package consul

import (
	"NULL/consul-gin/pkg/setting"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

//启动consul客户端
func StartConsul() {
	fmt.Println("start consul client。。。")
	node := fmt.Sprintf("-node=%s", setting.ServerSetting.NodeName)
	bind := fmt.Sprintf("-bind=%s", setting.ServerSetting.LocalIP)
	arg := []string{"agent", "-config-dir=consul.d", node, bind}
	cmd := exec.Command("consul.exe", arg...)
	b, err := cmd.CombinedOutput()
	if strings.Contains(string(b), "Only one usage of each socket address ") {
		fmt.Println("there has one consul client running ...")
		return
	}
	if err != nil {
		log.Fatalf("start consul Error:%v\n,consul log:%v\n", err, string(b))
	}
}

type RegService struct {
	ID                string
	Name              string
	Tags              []string
	Address           string
	Port              int
	Checks            []Checks
	EnableTagOverride bool
}
type Checks struct {
	Http     string
	Interval string
}

//注册服务
func RegisterService() {
	srv := fmt.Sprintf("webmeeting-%s", setting.ServerSetting.NodeName)
	rs := RegService{
		ID:                srv,
		Name:              srv,
		Tags:              []string{"webmeeting"},
		Address:           setting.ServerSetting.LocalIP,
		Port:              setting.ServerSetting.HttpPort,
		EnableTagOverride: false,
		Checks: []Checks{{
			Http: "http://" + setting.ServerSetting.LocalIP + ":" +
				strconv.Itoa(setting.ServerSetting.HttpPort) + "/health",
			Interval: "10s"},
		},
	}
	data, err := json.Marshal(&rs)
	if err != nil {
		log.Fatalf("Marshal json err:%v", err)
	}
	//log.Println(string(data))
	_, _, errs := gorequest.New().
		Put("http://127.0.0.1:9500/v1/agent/service/register").
		Type(gorequest.TypeJSON).Send(string(data)).End()
	if len(errs) != 0 {
		log.Printf("reg service err:%v", errs[0])
		if strings.Contains(errs[0].Error(), "the target machine actively refused it.") {
			time.Sleep(3)
			RegisterService()
		}
	}
}

//注销服务
func DeregisterService(srvID string) error {
	uri := fmt.Sprintf("http://127.0.0.1:9500/v1/agent/service/deregister/%s", srvID)
	_, _, errs := gorequest.New().Put(uri).End()
	if len(errs) != 0 {
		log.Printf("deregister service err:%v", errs[0])
		return errs[0]
	}
	return nil
}

//获取所有健康服务IP
func GetServiceIPs() []string {
	srvIps := make([]string, 0)
	srvs := GetHealthServices()
	for _, srv := range srvs {
		ips := CatalogService(srv)
		srvIps = append(srvIps, ips...)
	}
	return srvIps
}

//获取KV清单
func GetKVs() ([]map[string]string, error) {
	resp, body, errs := gorequest.New().
		Get("http://127.0.0.1:9500/v1/kv/?recurse").End()
	if len(errs) != 0 {
		log.Printf("GetKVs err:%v", errs[0])
		return nil, errs[0]
	}
	if resp.StatusCode == 404 {
		return nil, nil
	}
	res := make([]interface{}, 0)
	kvs := make([]map[string]string, 0)
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		log.Printf("keys Unmarshal json err:%v", err)
		return nil, err
	} else {
		for _, d := range res {
			key := d.(map[string]interface{})["Key"].(string)
			val := d.(map[string]interface{})["Value"].(string)
			decodeVal, _ := base64.StdEncoding.DecodeString(val)
			kvs = append(kvs, map[string]string{"Key": key, "Value": string(decodeVal)})
		}
	}
	return kvs, nil
}

//获取value
func GetKV(key string) string {
	uri := fmt.Sprintf("http://127.0.0.1:9500/v1/kv/%s", key)
	resp, body, errs := gorequest.New().Get(uri).End()
	if len(errs) != 0 {
		log.Printf("GetKV err:%v", errs[0])
	}
	if resp.StatusCode == 404 {
		return ""
	}
	res := make([]interface{}, 0)
	vals := make([]string, 0)
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		log.Printf("key[%s] Unmarshal json err:%v", key, err)
	} else {
		for _, d := range res {
			val := d.(map[string]interface{})["Value"].(string)
			decodeVal, _ := base64.StdEncoding.DecodeString(val)
			vals = append(vals, string(decodeVal))
		}
	}
	return vals[0]
}

//创建&更新KV
func PutKV(key, val string) string {
	uri := fmt.Sprintf("http://127.0.0.1:9500/v1/kv/%s", key)
	_, body, errs := gorequest.New().Put(uri).Type(gorequest.TypeText).Send(val).End()
	if len(errs) != 0 {
		log.Printf("PutKV err:%v", errs[0])
		return fmt.Sprintf("false:%v", errs[0])
	}
	return body
}

//删除KV
func DeleteKV(key string) string {
	uri := fmt.Sprintf("http://127.0.0.1:9500/v1/kv/%s", key)
	_, body, errs := gorequest.New().Delete(uri).End()
	if len(errs) != 0 {
		log.Printf("PutKV err:%v", errs[0])
		return fmt.Sprintf("false:%v", errs[0])
	}
	return body
}

//获取所有节点IP
func GetAgentIPs() ([]string, error) {
	_, body, errs := gorequest.New().
		Get("http://127.0.0.1:9500/v1/catalog/nodes").End()
	if len(errs) != 0 {
		log.Printf("Get AgentIPs err:%v", errs[0])
		return nil, errs[0]
	}
	res := make([]interface{}, 0)
	ips := make([]string, 0)
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		log.Printf("Get AgentIPs Unmarshal json err:%v", err)
		return nil, err
	} else {
		for _, d := range res {
			ip := d.(map[string]interface{})["Address"].(string)
			ips = append(ips, ip)
		}
	}
	return ips, nil
}

//查看所有已注册服务
func CatalogServices() []string {
	_, body, errs := gorequest.New().
		Get("http://127.0.0.1:9500/v1/catalog/services").End()
	if len(errs) != 0 {
		log.Printf("Agent Checks err:%v", errs[0])
	}
	res := make(map[string]interface{}, 0)
	srvs := make([]string, 0)
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		log.Printf("Agent Checks Unmarshal json err:%v", err)
	} else {
		for k := range res {
			if k == "consul" {
				continue
			}
			srvs = append(srvs, k)
		}
	}
	return srvs
}

//获取所有健康的服务
func GetHealthServices() []string {
	_, body, errs := gorequest.New().
		Get("http://127.0.0.1:9500/v1/health/state/passing").End()
	if len(errs) != 0 {
		log.Printf("Get Health Services err:%v", errs[0])
	}
	res := make([]interface{}, 0)
	srvs := make([]string, 0)
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		log.Printf("Get Health Services Unmarshal json err:%v", err)
	} else {
		for _, r := range res {
			if r.(map[string]interface{})["CheckID"].(string) == "serfHealth" {
				continue
			}
			srvName := r.(map[string]interface{})["ServiceName"].(string)
			srvs = append(srvs, srvName)
		}
	}
	return srvs
}

//查看服务详情
func CatalogService(srvName string) []string {
	uri := fmt.Sprintf("http://127.0.0.1:9500/v1/catalog/service/%s", srvName)
	_, body, errs := gorequest.New().Get(uri).End()
	if len(errs) != 0 {
		log.Printf("Agent Checks err:%v", errs[0])
	}
	res := make([]interface{}, 0)
	ips := make([]string, 0)
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		log.Printf("Catalog Service[%s] Unmarshal json err:%v", srvName, err)
	} else {
		for _, d := range res {
			//ServiceName := d.(map[string]interface{})["ServiceName"].(string)
			ip := d.(map[string]interface{})["Address"].(string)
			ips = append(ips, ip)
		}
	}
	return ips
}
