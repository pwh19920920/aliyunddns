package aliyun

import (
	"aliyunddns/config"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/kirinlabs/HttpRequest"
	"github.com/sirupsen/logrus"
)

func ExecuteUpdateAliYunDnsIp() {
	// 获取当前外网ip
	newIp, err := getIp()
	if err != nil {
		logrus.Error("请求外网ip失败", err.Error())
		return
	}

	// 对ip进行新旧检查
	oldIp, err := getOldIp()
	if err != nil {
		logrus.Error("获取旧Ip操作失败", err.Error())
		return
	}

	// ip比较一致
	if oldIp == newIp {
		logrus.WithFields(logrus.Fields{
			"oldIp": oldIp,
			"newIp": newIp,
		}).Info("域名解析Ip信息一致")
		return
	}

	// 更新ip
	//success, err := updateIp(newIp)
	//if success {
	//	logrus.Error("更新域名解析ip信息失败", err.Error())
	//	return
	//}

	logrus.WithFields(logrus.Fields{
		"oldIp": oldIp,
		"newIp": newIp,
	}).Info("更新域名解析ip信息成功")
}

// 检查新老ip是否一致
func getOldIp() (oldIp string, err error) {
	// 请求DescribeDomainRecordInfo
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", config.SystemConfig.AliYunConf.AliKeyId, config.SystemConfig.AliYunConf.AliKeySecret)

	request := alidns.CreateDescribeDomainRecordInfoRequest()
	request.Scheme = "https"
	request.RecordId = config.SystemConfig.AliYunConf.RecordId

	response, err := client.DescribeDomainRecordInfo(request)
	if err != nil {
		return "", err
	}

	// 返回比较
	return response.Value, nil
}

// 更新外网新ip
func updateIp(newIp string) (success bool, err error) {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", config.SystemConfig.AliYunConf.AliKeyId, config.SystemConfig.AliYunConf.AliKeySecret)

	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"

	request.RecordId = config.SystemConfig.AliYunConf.RecordId
	request.RR = config.SystemConfig.AliYunConf.RecordRr
	request.Type = config.SystemConfig.AliYunConf.RecordType
	request.TTL = requests.NewInteger(config.SystemConfig.AliYunConf.RecordTtl)
	request.Value = newIp

	_, err = client.UpdateDomainRecord(request)
	if err != nil {
		return false, err
	}

	return true, nil
}

// 获取外网ip
func getIp() (string, error) {
	req := HttpRequest.NewRequest()
	res, err := req.Get("http://ip-api.com/json")

	m := make(map[string]interface{}, 14)
	err = res.Json(&m)
	if err != nil {
		return "", err
	}

	return m["query"].(string), nil
}
