
serverName="go_server"
latestURL="https://github.com/spxvszero/jacky_go/releases/latest/download/go_export_jacky_go_linux_linux"

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
	\nWorkingDirectory=${curDir}
	\nExecStart=${curDir}/${serverName} --config ${curDir}/${configFile}
	\nRestart=on-failure
	\n\n[Install]
	\nWantedBy=multi-user.target
	\n"

	if [[ -e ${systemdServiceDir} ]]; then
		echo "SystemService Dir Exist."
	else
		echo "SystemService Dir Not Exist. Create One. ${systemdServiceDir}"
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

function iptablesRules(){
	#block scan 
	iptables -N block-scan
	iptables -A block-scan -p tcp —tcp-flags SYN,ACK,FIN,RST RST -m limit —limit 1/s -j RETURN
	iptables -A block-scan -j DROP

	#block bad port
	badport="135,136,137,138,139,445"
	iptables -A INPUT -p tcp -m multiport --dport $badport -j DROP
	iptables -A INPUT -p udp -m multiport --dport $badport -j DROP

	#block ddos
	iptables -A INPUT -p tcp --dport 80 -m limit --limit 20/minute --limit-burst 100 -j ACCEPT
	#portect more from ddos
	echo 1 > /proc/sys/net/ipv4/ip_forward
	echo 1 > /proc/sys/net/ipv4/tcp_syncookies
	echo 0 > /proc/sys/net/ipv4/conf/all/accept_redirects
	echo 0 > /proc/sys/net/ipv4/conf/all/accept_source_route
	echo 1 > /proc/sys/net/ipv4/conf/all/rp_filter
	echo 1 > /proc/sys/net/ipv4/conf/lo/rp_filter
	echo 1 > /proc/sys/net/ipv4/conf/lo/arp_ignore
	echo 2 > /proc/sys/net/ipv4/conf/lo/arp_announce
	echo 1 > /proc/sys/net/ipv4/conf/all/arp_ignore
	echo 2 > /proc/sys/net/ipv4/conf/all/arp_announce
	echo 0 > /proc/sys/net/ipv4/icmp_echo_ignore_all
	echo 1 > /proc/sys/net/ipv4/icmp_echo_ignore_broadcasts
	echo 30 > /proc/sys/net/ipv4/tcp_fin_timeout
	echo 1800 > /proc/sys/net/ipv4/tcp_keepalive_time
	echo 1 > /proc/sys/net/ipv4/tcp_window_scaling 
	echo 0 > /proc/sys/net/ipv4/tcp_sack
	echo 1280 > /proc/sys/net/ipv4/tcp_max_syn_backlog

	#block smtp
	iptables -A OUTPUT -p tcp --dport 25 -j DROP

	#only accept ip to connect mysql port
	iptables -A INPUT -p tcp -s 192.168.1.0/24 --dport 3306 -m state --state NEW,ESTABLISHED -j ACCEPT

	#block icmp(ping)
	#outgoing
	iptables -A OUTPUT -p icmp --icmp-type 8 -j DROP
	#incoming
	iptables -I INPUT -p icmp --icmp-type 8 -j DROP
}

function addLaunchServiceOnMac(){
workingDirectory="/Users"
filePath="/Users/go_server"
configFilePath="/Users/configFilePath"
serviceName="com.jacky.goserver"
plistFile="/Users/${serverName}.plist"
echo "
<?xml version=\"1.0\" encoding=\"UTF-8\"?>
<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">
<plist version=\"1.0\">
<dict>
	<key>Label</key>
	<string>${serviceName}</string>
	<key>ServiceDescription</key>
	<string>Go server From jacky</string>
	<key>RunAtLoad</key>
	<string>false</string>
	<key>WorkingDirectory</key>
	<string>${workingDirectory}</string>
	<key>MachServices</key>
    <dict>
         <key>${serviceName}</key>
         <true/>
    </dict>
    <key>Program</key>
    <string>${filePath}</string>
    <key>ProgramArguments</key>
        <array>
            <string>go_server</string>
            <string> --config ${configFilePath}</string>
        </array>
    <key>KeepAlive</key>
    <true/> 
</dict>
</plist>
" > plistFile
	
	launchctl load plistFile
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
