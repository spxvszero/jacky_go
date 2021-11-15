
serverName="go_server"
latestURL="https://github.com/spxvszero/jacky_go/releases/latest/download/go_export_jacky_go_linux"

curDir="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
configFile="config.json"

systemdServiceDir="/usr/lib/systemd/system"
systemdServiceFile="/usr/lib/systemd/system/${serverName}.service"
firewallServiceXML="/usr/lib/firewalld/services/${serverName}.xml"

port=""
socksPort=""
downloadDir=""

function addSystemdService(){

	if ! command -v systemctl &> /dev/null
	then
	    echo "systemctl command could not be found."
	    exit
	fi

	serviceFile="[Unit]
	\nDescription=${serverName}
	\nAfter=network-online.target
	\n\n[Service]
	\nType=simple
	\nExecStart=${curDir}/${serverName} --config ${curDir}/${configFile}
	\n\n[Install]
	\nWantedBy=multi-user.target
	\n"

	if [[ -e ${systemdServiceDir} ]]; then
	else
		mkdir ${systemdServiceDir}
	fi

	echo -e ${serviceFile} > ${systemdServiceFile}

	systemctl daemon-reload
	systemctl enable ${serverName}

	systemctl start ${serverName}
}

function removeSystemdService(){

	if ! command -v systemctl &> /dev/null
	then
	    echo "systemctl command could not be found."
	    exit
	fi

	systemctl stop ${serverName}
	systemctl disable ${serverName}

	if [[ -e ${systemdServiceFile} ]]; then
		#statements
		rm -f ${systemdServiceFile}
	fi

	systemctl daemon-reload
}


function addFirewallService(){


	if ! command -v firewall-cmd &> /dev/null
	then
	    echo "firewalld command could not be found."
	    exit
	fi


	firewallService="<?xml version=\"1.0\" encoding=\"utf-8\"?>
\n<service>
\n  <short>${serverName}</short>
\n  <description>This server is made for custom services, looking more with site : https://github.com/spxvszero/jacky_go</description>
\n  <port protocol=\"tcp\" port=\"${port}\"/>
\n  <port protocol=\"tcp\" port=\"${socksPort}\"/>
\n</service>"
	
	echo -e ${firewallService} > ${firewallServiceXML}

	firewall-cmd --reload

	firewall-cmd --add-service=${serverName} --permanent
	firewall-cmd --reload
}

function removeFirewallService(){

	if ! command -v firewall-cmd &> /dev/null
	then
	    echo "firewalld command could not be found."
	    exit
	fi

	firewall-cmd --remove-service=${serverName} --permanent

	if [[ -e ${firewallServiceXML} ]]; then
		#statements
		rm -f ${firewallServiceXML}
	fi

	firewall-cmd --reload

}

function installServer(){
	#check if exist
	if [[ -e ${serverName} ]]; then
		#statements
		echo "${serverName} File exist, skip download..."
	else

		if command -v curl &> /dev/null;then
		    curl -o ${serverName} -L ${latestURL}
		    
		elif command -v wget &> /dev/null; then
			wget -O ${serverName} ${latestURL}

		else
			echo "Download Failed! Try Download yourself : ${latestURL}"
			echo "And retry this script."
			exit
		fi

		
	fi

	chmod +x ${serverName}
	./${serverName} --generate config.json

	#change port in config.json
	if [[ -n ${port} ]]; then
		#statements
		sed -i "s/\"port\": 8900/\"port\": ${port}/" config.json
	else
		port="8900"
	fi

	#socks5 port
	if [[ -n ${socksPort} ]]; then
		#statements
		sed -i "s/\"addr\": \"127.0.0.1:7777\"/\"addr\": \"${socksPort}\"/" config.json
		socksPort=`echo ${socksPort}|awk -F ":" '{print $2}'`
	else
		socksPort="7777"
	fi
	
	#download dir
	if [[ -n ${downloadDir} ]]; then
		#statements
		sed -i "s|\"download_dir_path\": \"/root/download\"|\"download_dir_path\": \"${downloadDir}\"|" config.json
		sed -i "s|\"dir_path\": \"/root/download\"|\"dir_path\": \"${downloadDir}\"|" config.json
		sed -i "s|\"save_dir_path\": \"/root/download\"|\"save_dir_path\": \"${downloadDir}\"|" config.json

	else
		downloadDir="/root/download"
	fi

	mkdir ${downloadDir}

	addSystemdService
	addFirewallService

	echo "Finished!"
}

function updateServer(){
	#stop server
	systemctl stop ${serverName}

	#download && update
	if command -v curl &> /dev/null;then
	    curl -o ${serverName} -L ${latestURL}
	    
	elif command -v wget &> /dev/null; then
		wget -O ${serverName} ${latestURL}

	else
		echo "Download Failed! Try Download yourself : ${latestURL}"
		echo "And retry this script."
		exit
	fi

	#restart server
	systemctl start ${serverName}
}

function uninstallServer(){
	removeFirewallService
	removeSystemdService
}

function welcomeInfo(){
	echo ""
	echo "** ** ** ** ** ** ** ** ** ** ** **"
	echo ""
	echo "Welcome jacky_go server Auto Setup Script! "
	echo "Tell me what you want to do :"
	echo ""
	echo "	1. Install"
	echo "	2. Update"
	echo "	3. Uninstall"
	echo ""
	printf "I want to : "
	read input


	if [[ ${input} == "1" ]]; then

		echo ""
		echo ""
		echo "Wanna change port ? Default listen on 8900"
		printf "Change Port (empty for default): "
		read port

		echo ""
		echo ""
		echo "Wanna open socks ? Default listen on 127.0.0.1:7777"
		printf "Change socks port (empty for default): "
		read socksPort

		echo ""
		echo ""
		echo "Wanna change download dir ? Default is /root/download"
		printf "Change download dir (empty for default): "
		read downloadDir

		echo "OK , Ready For Install."

		#statement
		installServer
	elif [[ ${input} == "2" ]]; then
		#statements
		updateServer
	elif [[ ${input} == "3" ]]; then
		uninstallServer
	else
		echo "Not Funny, Bye!"
	fi
}



welcomeInfo
