
serverName="go_server"
latestURL="https://github.com/spxvszero/jacky_go/releases/latest/download/go_export_jacky_go_linux"

curDir="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
configFile="config.json"

systemdServiceFile="/usr/lib/systemd/system/${serverName}.service"
firewallServiceXML="/usr/lib/firewalld/services/${serverName}.xml"

function addSystemdService(){
	serviceFile="[Unit]
	\nDescription=${serverName}
	\nAfter=network-online.target
	\n\n[Service]
	\nType=simple
	\nExecStart=${curDir}/${serverName} --config ${curDir}/${configFile}
	\n\n[Install]
	\nWantedBy=multi-user.target
	\n"

	echo -e ${serviceFile} > ${systemdServiceFile}

	systemctl daemon-reload
	systemctl enable ${serverName}

	systemctl start ${serverName}
}

function removeSystemdService(){
	systemctl stop ${serverName}
	systemctl disable ${serverName}

	if [[ -e ${systemdServiceFile} ]]; then
		#statements
		rm -f ${systemdServiceFile}
	fi

	systemctl daemon-reload
}


function addFirewallService(){
	firewallService="<?xml version=\"1.0\" encoding=\"utf-8\"?>
\n<service>
\n  <short>${serverName}</short>
\n  <description>This server is made for custom services, looking more with site : https://github.com/spxvszero/jacky_go</description>
\n  <port protocol=\"tcp\" port=\"8900\"/>
\n  <port protocol=\"tcp\" port=\"7777\"/>
\n</service>"
	
	echo -e ${firewallService} > ${firewallServiceXML}

	firewall-cmd --reload

	firewall-cmd --add-service=${serverName} --permanent
	firewall-cmd --reload
}

function removeFirewallService(){
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
		echo "jacky_go File exist, skip download..."
	else
		curl -o ${serverName} -L ${latestURL}
	fi

	chmod +x ${serverName}
	./${serverName} --generate config.json
	mkdir download

	addSystemdService
	addFirewallService

	echo "Finished!"
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
	echo "	2. Uninstall"
	echo ""
	printf "I want to : "
	read input
	if [[ ${input} == "1" ]]; then
		#statement
		installServer
	elif [[ ${input} == "2" ]]; then
		#statements
		uninstallServer
	else
		echo "Not Funny, Bye!"
	fi
}



welcomeInfo
