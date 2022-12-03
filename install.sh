#!/bin/bash

Green_font="\033[32m"
Yellow_font="\033[33m"
Red_font="\033[31m"
Font_suffix="\033[0m"
Info="${Green_font}[Info]${Font_suffix}"
Error="${Red_font}[Error]${Font_suffix}"
Tips="${Green_font}[Tips]${Font_suffix}"
Temp_path="/var/tmp/client"

checkRootPermit() {
    [[ $EUID -ne 0 ]] && echo -e "${Error} 请使用sudo/root权限运行本脚本" && exit 1
}

checkSystemArch() {
    arch=$(uname -m)
    if [[ $arch == "x86_64" ]]; then
        archParam="amd64"
    fi

    if [[ $arch == "aarch64" ]]; then
        archParam="arm64"
    fi
}

checkSystemDistribution() {
    case "$OSTYPE" in
    linux*)
        osDistribution="linux"
        conf_dir="/usr/local/traceClient/"
        downPath="/usr/local/traceClient/client"
        ;;
    *)
        echo "unknown: $OSTYPE"
        exit 1
        ;;
    esac
}

getNodeInfo() {
    res=$(curl -s "https://$1/api/node/token/$2")
    array=(${res//|/ })
    if [[ ${array[0]} == "" ]]; then
        echo -e "${Error} 获取节点信息失败，请确认参数是否复制正确"
        exit 1
    fi
    echo "创建用户：${array[0]}"
    echo "用户权限：${array[2]}"
    echo "节点 IP：${array[1]}"
    read -p "请确认您想要对接的节点信息是否正确 (y/n) [y]:" input
    if [[ $input == "n" ]]; then
        exit 1
    fi
}

installWgetPackage() {
    echo -e "${Info} wget 正在安装中..."
    # try apt
    apt -h &>/dev/null
    if [ $? -eq 0 ]; then
        # 先更新一下数据源，有些机器数据源比较老可能会404
        apt update -y &>/dev/null
        apt install wget -y &>/dev/null
    fi

    # try yum
    yum -h &>/dev/null
    if [ $? -eq 0 ]; then
        yum install wget -y &>/dev/null
    fi

    # try dnf
    dnf -h &>/dev/null
    if [ $? -eq 0 ]; then
        dnf install wget -y &>/dev/null
    fi

    # try pacman
    pacman -h &>/dev/null
    if [ $? -eq 0 ]; then
        pacman -Sy
        pacman -S wget
    fi

}

getLocation() {
    countryCode=$(curl -s "http://ip-api.com/line/?fields=countryCode")
}

checkWgetPackage() {
    wget -h &>/dev/null
    if [ $? -ne 0 ]; then
        installWgetPackage
    fi
}

downloadBinrayFile() {
    if [ ! -d $conf_dir ]; then
        mkdir $conf_dir
    fi
    if [ ! -e ${downPath} ]; then
        latestURL="https://leo.moe/traceClient/traceClient_${osDistribution}_${archParam}"
        echo $latestURL
        echo -e "${Info} 正在下载二进制文件..."
        wget -O ${Temp_path} ${latestURL} &>/dev/null
        if [ $? -eq 0 ]; then
            changeMode
            mv ${Temp_path} ${downPath}
            echo -e "${Info} Done!"
            # 初始化配置文件
            		echo -e "backcallurl: https://$1
token: $2" > $conf_dir/config.yaml
            # 创建 Systemctl
            echo -e "[Unit]
Description=Trace Client Service
After=network.target
[Service]
Type=simple
Restart=always
WoringDirectory=${conf_dir}
ExecStart=${downPath}
[Install]
WantedBy=multi-user.target" >/etc/systemd/system/traceClient.service
            systemctl daemon-reload
            systemctl start traceClient.service
            systemctl enable traceClient.service
        else
            echo -e "${Error} 下载失败，请检查您的网络是否正常"
            exit 1
        fi
    fi
}

changeMode() {
    chmod +x ${Temp_path} &>/dev/null
}

checkRootPermit
checkSystemArch
checkSystemDistribution
getNodeInfo $1 $2
checkWgetPackage
downloadBinrayFile $1 $2
