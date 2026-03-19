package inits

import (
	"fmt"
	"gospacex-tz/srv/basic/config"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
)

func NacosInit() {
	// 确保 GlobalConf 已初始化
	if config.GlobalConfig == nil {
		config.GlobalConfig = &config.AppConfig{}
	}

	// 构建 Nacos 客户端配置
	clientConfig := constant.ClientConfig{
		NamespaceId:         config.GlobalConfig.Nacos.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
	}

	// 构建 Nacos 服务器配置
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: config.GlobalConfig.Nacos.Addr,
			Port:   uint64(config.GlobalConfig.Nacos.Port),
		},
	}

	// 创建 Nacos 配置客户端
	nacosClient, err := clients.CreateConfigClient(map[string]interface{}{
		"clientConfig":  clientConfig,
		"serverConfigs": serverConfigs,
	})
	if err != nil {
		fmt.Printf("创建 Nacos 客户端失败: %v\n", err)
		return
	}

	// 从 Nacos 获取配置内容
	configContent, err := nacosClient.GetConfig(vo.ConfigParam{
		DataId: config.GlobalConfig.Nacos.DataID,
		Group:  config.GlobalConfig.Nacos.Group,
	})
	if err != nil {
		fmt.Printf("从 Nacos 获取配置失败: %v\n", err)
		return
	}

	// 使用 Viper 解析 YAML 配置并更新全局配置
	viper.SetConfigType("yaml")
	err = viper.ReadConfig(strings.NewReader(configContent))
	if err != nil {
		fmt.Printf("解析 Nacos 配置失败: %v\n", err)
		return
	}

	// 将配置解析到全局配置对象
	err = viper.Unmarshal(config.GlobalConfig)
	if err != nil {
		fmt.Printf("反序列化 Nacos 配置失败: %v\n", err)
		return
	}

	fmt.Println("Nacos 配置读取成功")
}
