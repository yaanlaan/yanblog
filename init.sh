#!/bin/bash
# init.sh
# 用于初始化配置文件，将 config_template.yaml 复制为 config.yaml（Linux 版）

# 后端配置模板复制
if [ -f config/config_template.yaml ]; then
    cp config/config_template.yaml config/config.yaml
    echo "config/config_template.yaml 已复制为 config/config.yaml"
else
    echo "未找到 config/config_template.yaml，未进行复制。"
fi

# 前端配置模板复制
if [ -f web/frontend/public/config_template.yaml ]; then
    cp web/frontend/public/config_template.yaml web/frontend/public/config.yaml
    echo "web/frontend/public/config_template.yaml 已复制为 web/frontend/public/config.yaml"
else
    echo "未找到 web/frontend/public/config_template.yaml，未进行复制。"
fi