###	云函数部署

使用CustomRuntime进行部署， bootstrap 文件在 scripts/bootstrap 中已给出。 在部署前，请在本地完成登录，并将 config.yml ， device.json ， bootstrap 和 QN 一起打包。

在触发器中创建一个 API 网关触发器，并启用继承响应， 创建完成后即可通过api网关访问 QN (建议配置 AccessToken)。

    scripts/bootstrap 中使用的工作路径为 /tmp, 这个目录最大能容下 500M 文件, 如需长期使用， 请挂载文件存储(CFS).