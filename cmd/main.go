package main

import (
	"fmt"
	"log"
	"net/http"

	"iteasy.wrappedAnsible/internal/ansible"
	"iteasy.wrappedAnsible/internal/handlers"
	"iteasy.wrappedAnsible/internal/router"
)

func init() {

	fmt.Println("Welcome Ansible Test Init.")
	initJsonData := `{
			"type": "init",
			"name": "초기화테스터",
		  "options": {}
	  }
	`
	extendAnsible := ansible.GetAnsibleFromFactory([]byte(initJsonData))
	ansible.Excuter(extendAnsible)
}

func main() {
	fmt.Println("Welcome Ansible Test.")

	// 	outputJsonData := `
	// {
	//   "custom_stats": {},
	//   "global_custom_stats": {},
	//   "plays": [
	//     {
	//       "play": {
	//         "duration": {
	//           "end": "2024-06-26T06:23:15.519469Z",
	//           "start": "2024-06-26T06:23:14.135209Z"
	//         },
	//         "id": "5254004d-77d3-22f8-e962-000000000007",
	//         "name": "Example of using ansible_facts"
	//       },
	//       "tasks": [
	//         {
	//           "hosts": {
	//             "172.16.74.100": {
	//               "_ansible_no_log": false,
	//               "_ansible_verbose_override": true,
	//               "action": "gather_facts",
	//               "ansible_facts": {
	//                 "ansible_all_ipv4_addresses": [
	//                   "172.16.74.100",
	//                   "10.0.2.15"
	//                 ],
	//                 "ansible_all_ipv6_addresses": [
	//                   "fe80::a00:27ff:fe2a:23e4",
	//                   "fe80::5054:ff:fe4d:77d3"
	//                 ],
	//                 "ansible_apparmor": {
	//                   "status": "disabled"
	//                 },
	//                 "ansible_architecture": "x86_64",
	//                 "ansible_bios_date": "12/01/2006",
	//                 "ansible_bios_vendor": "innotek GmbH",
	//                 "ansible_bios_version": "VirtualBox",
	//                 "ansible_board_asset_tag": "NA",
	//                 "ansible_board_name": "VirtualBox",
	//                 "ansible_board_serial": "0",
	//                 "ansible_board_vendor": "Oracle Corporation",
	//                 "ansible_board_version": "1.2",
	//                 "ansible_chassis_asset_tag": "NA",
	//                 "ansible_chassis_serial": "NA",
	//                 "ansible_chassis_vendor": "Oracle Corporation",
	//                 "ansible_chassis_version": "NA",
	//                 "ansible_cmdline": {
	//                   "BOOT_IMAGE": "/boot/vmlinuz-3.10.0-1127.el7.x86_64",
	//                   "LANG": "en_US.UTF-8",
	//                   "biosdevname": "0",
	//                   "console": "ttyS0,115200n8",
	//                   "crashkernel": "auto",
	//                   "elevator": "noop",
	//                   "net.ifnames": "0",
	//                   "no_timer_check": true,
	//                   "ro": true,
	//                   "root": "UUID=1c419d6c-5064-4a2b-953c-05b2c67edb15"
	//                 },
	//                 "ansible_date_time": {
	//                   "date": "2024-06-17",
	//                   "day": "17",
	//                   "epoch": "1718606899",
	//                   "hour": "06",
	//                   "iso8601": "2024-06-17T06:48:19Z",
	//                   "iso8601_basic": "20240617T064819271628",
	//                   "iso8601_basic_short": "20240617T064819",
	//                   "iso8601_micro": "2024-06-17T06:48:19.271628Z",
	//                   "minute": "48",
	//                   "month": "06",
	//                   "second": "19",
	//                   "time": "06:48:19",
	//                   "tz": "UTC",
	//                   "tz_dst": "UTC",
	//                   "tz_offset": "+0000",
	//                   "weekday": "Monday",
	//                   "weekday_number": "1",
	//                   "weeknumber": "25",
	//                   "year": "2024"
	//                 },
	//                 "ansible_default_ipv4": {
	//                   "address": "10.0.2.15",
	//                   "alias": "eth0",
	//                   "broadcast": "10.0.2.255",
	//                   "gateway": "10.0.2.2",
	//                   "interface": "eth0",
	//                   "macaddress": "52:54:00:4d:77:d3",
	//                   "mtu": 1500,
	//                   "netmask": "255.255.255.0",
	//                   "network": "10.0.2.0",
	//                   "type": "ether"
	//                 },
	//                 "ansible_default_ipv6": {},
	//                 "ansible_device_links": {
	//                   "ids": {
	//                     "sda": [
	//                       "ata-VBOX_HARDDISK_VB425aa53d-5adf779d"
	//                     ],
	//                     "sda1": [
	//                       "ata-VBOX_HARDDISK_VB425aa53d-5adf779d-part1"
	//                     ]
	//                   },
	//                   "labels": {},
	//                   "masters": {},
	//                   "uuids": {
	//                     "sda1": [
	//                       "1c419d6c-5064-4a2b-953c-05b2c67edb15"
	//                     ]
	//                   }
	//                 },
	//                 "ansible_devices": {
	//                   "sda": {
	//                     "holders": [],
	//                     "host": "IDE interface: Intel Corporation 82371AB/EB/MB PIIX4 IDE (rev 01)",
	//                     "links": {
	//                       "ids": [
	//                         "ata-VBOX_HARDDISK_VB425aa53d-5adf779d"
	//                       ],
	//                       "labels": [],
	//                       "masters": [],
	//                       "uuids": []
	//                     },
	//                     "model": "VBOX HARDDISK",
	//                     "partitions": {
	//                       "sda1": {
	//                         "holders": [],
	//                         "links": {
	//                           "ids": [
	//                             "ata-VBOX_HARDDISK_VB425aa53d-5adf779d-part1"
	//                           ],
	//                           "labels": [],
	//                           "masters": [],
	//                           "uuids": [
	//                             "1c419d6c-5064-4a2b-953c-05b2c67edb15"
	//                           ]
	//                         },
	//                         "sectors": "83884032",
	//                         "sectorsize": 512,
	//                         "size": "40.00 GB",
	//                         "start": "2048",
	//                         "uuid": "1c419d6c-5064-4a2b-953c-05b2c67edb15"
	//                       }
	//                     },
	//                     "removable": "0",
	//                     "rotational": "1",
	//                     "sas_address": null,
	//                     "sas_device_handle": null,
	//                     "scheduler_mode": "noop",
	//                     "sectors": "83886080",
	//                     "sectorsize": "512",
	//                     "serial": "VB425aa53d",
	//                     "size": "40.00 GB",
	//                     "support_discard": "0",
	//                     "vendor": "ATA",
	//                     "virtual": 1
	//                   }
	//                 },
	//                 "ansible_distribution": "CentOS",
	//                 "ansible_distribution_file_parsed": true,
	//                 "ansible_distribution_file_path": "/etc/redhat-release",
	//                 "ansible_distribution_file_variety": "RedHat",
	//                 "ansible_distribution_major_version": "7",
	//                 "ansible_distribution_release": "Core",
	//                 "ansible_distribution_version": "7.8",
	//                 "ansible_dns": {
	//                   "nameservers": [
	//                     "10.0.2.3"
	//                   ]
	//                 },
	//                 "ansible_domain": "",
	//                 "ansible_effective_group_id": 0,
	//                 "ansible_effective_user_id": 0,
	//                 "ansible_env": {
	//                   "HOME": "/root",
	//                   "LANG": "C",
	//                   "LC_ALL": "C",
	//                   "LC_MESSAGES": "C",
	//                   "LESSOPEN": "||/usr/bin/lesspipe.sh %s",
	//                   "LOGNAME": "root",
	//                   "LS_COLORS": "rs=0:di=38;5;27:ln=38;5;51:mh=44;38;5;15:pi=40;38;5;11:so=38;5;13:do=38;5;5:bd=48;5;232;38;5;11:cd=48;5;232;38;5;3:or=48;5;232;38;5;9:mi=05;48;5;232;38;5;15:su=48;5;196;38;5;15:sg=48;5;11;38;5;16:ca=48;5;196;38;5;226:tw=48;5;10;38;5;16:ow=48;5;10;38;5;21:st=48;5;21;38;5;15:ex=38;5;34:*.tar=38;5;9:*.tgz=38;5;9:*.arc=38;5;9:*.arj=38;5;9:*.taz=38;5;9:*.lha=38;5;9:*.lz4=38;5;9:*.lzh=38;5;9:*.lzma=38;5;9:*.tlz=38;5;9:*.txz=38;5;9:*.tzo=38;5;9:*.t7z=38;5;9:*.zip=38;5;9:*.z=38;5;9:*.Z=38;5;9:*.dz=38;5;9:*.gz=38;5;9:*.lrz=38;5;9:*.lz=38;5;9:*.lzo=38;5;9:*.xz=38;5;9:*.bz2=38;5;9:*.bz=38;5;9:*.tbz=38;5;9:*.tbz2=38;5;9:*.tz=38;5;9:*.deb=38;5;9:*.rpm=38;5;9:*.jar=38;5;9:*.war=38;5;9:*.ear=38;5;9:*.sar=38;5;9:*.rar=38;5;9:*.alz=38;5;9:*.ace=38;5;9:*.zoo=38;5;9:*.cpio=38;5;9:*.7z=38;5;9:*.rz=38;5;9:*.cab=38;5;9:*.jpg=38;5;13:*.jpeg=38;5;13:*.gif=38;5;13:*.bmp=38;5;13:*.pbm=38;5;13:*.pgm=38;5;13:*.ppm=38;5;13:*.tga=38;5;13:*.xbm=38;5;13:*.xpm=38;5;13:*.tif=38;5;13:*.tiff=38;5;13:*.png=38;5;13:*.svg=38;5;13:*.svgz=38;5;13:*.mng=38;5;13:*.pcx=38;5;13:*.mov=38;5;13:*.mpg=38;5;13:*.mpeg=38;5;13:*.m2v=38;5;13:*.mkv=38;5;13:*.webm=38;5;13:*.ogm=38;5;13:*.mp4=38;5;13:*.m4v=38;5;13:*.mp4v=38;5;13:*.vob=38;5;13:*.qt=38;5;13:*.nuv=38;5;13:*.wmv=38;5;13:*.asf=38;5;13:*.rm=38;5;13:*.rmvb=38;5;13:*.flc=38;5;13:*.avi=38;5;13:*.fli=38;5;13:*.flv=38;5;13:*.gl=38;5;13:*.dl=38;5;13:*.xcf=38;5;13:*.xwd=38;5;13:*.yuv=38;5;13:*.cgm=38;5;13:*.emf=38;5;13:*.axv=38;5;13:*.anx=38;5;13:*.ogv=38;5;13:*.ogx=38;5;13:*.aac=38;5;45:*.au=38;5;45:*.flac=38;5;45:*.mid=38;5;45:*.midi=38;5;45:*.mka=38;5;45:*.mp3=38;5;45:*.mpc=38;5;45:*.ogg=38;5;45:*.ra=38;5;45:*.wav=38;5;45:*.axa=38;5;45:*.oga=38;5;45:*.spx=38;5;45:*.xspf=38;5;45:",
	//                   "MAIL": "/var/mail/root",
	//                   "PATH": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin",
	//                   "PWD": "/root",
	//                   "SELINUX_LEVEL_REQUESTED": "",
	//                   "SELINUX_ROLE_REQUESTED": "",
	//                   "SELINUX_USE_CURRENT_RANGE": "",
	//                   "SHELL": "/bin/bash",
	//                   "SHLVL": "2",
	//                   "SSH_CLIENT": "172.16.74.10 40598 22",
	//                   "SSH_CONNECTION": "172.16.74.10 40598 172.16.74.100 22",
	//                   "SSH_TTY": "/dev/pts/1",
	//                   "TERM": "xterm-256color",
	//                   "USER": "root",
	//                   "XDG_RUNTIME_DIR": "/run/user/0",
	//                   "XDG_SESSION_ID": "17",
	//                   "_": "/usr/bin/python"
	//                 },
	//                 "ansible_eth0": {
	//                   "active": true,
	//                   "device": "eth0",
	//                   "features": {
	//                     "busy_poll": "off [fixed]",
	//                     "fcoe_mtu": "off [fixed]",
	//                     "generic_receive_offload": "on",
	//                     "generic_segmentation_offload": "on",
	//                     "highdma": "off [fixed]",
	//                     "hw_tc_offload": "off [fixed]",
	//                     "l2_fwd_offload": "off [fixed]",
	//                     "large_receive_offload": "off [fixed]",
	//                     "loopback": "off [fixed]",
	//                     "netns_local": "off [fixed]",
	//                     "ntuple_filters": "off [fixed]",
	//                     "receive_hashing": "off [fixed]",
	//                     "rx_all": "off",
	//                     "rx_checksumming": "off",
	//                     "rx_fcs": "off",
	//                     "rx_gro_hw": "off [fixed]",
	//                     "rx_udp_tunnel_port_offload": "off [fixed]",
	//                     "rx_vlan_filter": "on [fixed]",
	//                     "rx_vlan_offload": "on",
	//                     "rx_vlan_stag_filter": "off [fixed]",
	//                     "rx_vlan_stag_hw_parse": "off [fixed]",
	//                     "scatter_gather": "on",
	//                     "tcp_segmentation_offload": "on",
	//                     "tx_checksum_fcoe_crc": "off [fixed]",
	//                     "tx_checksum_ip_generic": "on",
	//                     "tx_checksum_ipv4": "off [fixed]",
	//                     "tx_checksum_ipv6": "off [fixed]",
	//                     "tx_checksum_sctp": "off [fixed]",
	//                     "tx_checksumming": "on",
	//                     "tx_fcoe_segmentation": "off [fixed]",
	//                     "tx_gre_csum_segmentation": "off [fixed]",
	//                     "tx_gre_segmentation": "off [fixed]",
	//                     "tx_gso_partial": "off [fixed]",
	//                     "tx_gso_robust": "off [fixed]",
	//                     "tx_ipip_segmentation": "off [fixed]",
	//                     "tx_lockless": "off [fixed]",
	//                     "tx_nocache_copy": "off",
	//                     "tx_scatter_gather": "on",
	//                     "tx_scatter_gather_fraglist": "off [fixed]",
	//                     "tx_sctp_segmentation": "off [fixed]",
	//                     "tx_sit_segmentation": "off [fixed]",
	//                     "tx_tcp6_segmentation": "off [fixed]",
	//                     "tx_tcp_ecn_segmentation": "off [fixed]",
	//                     "tx_tcp_mangleid_segmentation": "off",
	//                     "tx_tcp_segmentation": "on",
	//                     "tx_udp_tnl_csum_segmentation": "off [fixed]",
	//                     "tx_udp_tnl_segmentation": "off [fixed]",
	//                     "tx_vlan_offload": "on [fixed]",
	//                     "tx_vlan_stag_hw_insert": "off [fixed]",
	//                     "udp_fragmentation_offload": "off [fixed]",
	//                     "vlan_challenged": "off [fixed]"
	//                   },
	//                   "hw_timestamp_filters": [],
	//                   "ipv4": {
	//                     "address": "10.0.2.15",
	//                     "broadcast": "10.0.2.255",
	//                     "netmask": "255.255.255.0",
	//                     "network": "10.0.2.0"
	//                   },
	//                   "ipv6": [
	//                     {
	//                       "address": "fe80::5054:ff:fe4d:77d3",
	//                       "prefix": "64",
	//                       "scope": "link"
	//                     }
	//                   ],
	//                   "macaddress": "52:54:00:4d:77:d3",
	//                   "module": "e1000",
	//                   "mtu": 1500,
	//                   "pciid": "0000:00:03.0",
	//                   "promisc": false,
	//                   "speed": 1000,
	//                   "timestamping": [
	//                     "tx_software",
	//                     "rx_software",
	//                     "software"
	//                   ],
	//                   "type": "ether"
	//                 },
	//                 "ansible_eth1": {
	//                   "active": true,
	//                   "device": "eth1",
	//                   "features": {
	//                     "busy_poll": "off [fixed]",
	//                     "fcoe_mtu": "off [fixed]",
	//                     "generic_receive_offload": "on",
	//                     "generic_segmentation_offload": "on",
	//                     "highdma": "off [fixed]",
	//                     "hw_tc_offload": "off [fixed]",
	//                     "l2_fwd_offload": "off [fixed]",
	//                     "large_receive_offload": "off [fixed]",
	//                     "loopback": "off [fixed]",
	//                     "netns_local": "off [fixed]",
	//                     "ntuple_filters": "off [fixed]",
	//                     "receive_hashing": "off [fixed]",
	//                     "rx_all": "off",
	//                     "rx_checksumming": "off",
	//                     "rx_fcs": "off",
	//                     "rx_gro_hw": "off [fixed]",
	//                     "rx_udp_tunnel_port_offload": "off [fixed]",
	//                     "rx_vlan_filter": "on [fixed]",
	//                     "rx_vlan_offload": "on",
	//                     "rx_vlan_stag_filter": "off [fixed]",
	//                     "rx_vlan_stag_hw_parse": "off [fixed]",
	//                     "scatter_gather": "on",
	//                     "tcp_segmentation_offload": "on",
	//                     "tx_checksum_fcoe_crc": "off [fixed]",
	//                     "tx_checksum_ip_generic": "on",
	//                     "tx_checksum_ipv4": "off [fixed]",
	//                     "tx_checksum_ipv6": "off [fixed]",
	//                     "tx_checksum_sctp": "off [fixed]",
	//                     "tx_checksumming": "on",
	//                     "tx_fcoe_segmentation": "off [fixed]",
	//                     "tx_gre_csum_segmentation": "off [fixed]",
	//                     "tx_gre_segmentation": "off [fixed]",
	//                     "tx_gso_partial": "off [fixed]",
	//                     "tx_gso_robust": "off [fixed]",
	//                     "tx_ipip_segmentation": "off [fixed]",
	//                     "tx_lockless": "off [fixed]",
	//                     "tx_nocache_copy": "off",
	//                     "tx_scatter_gather": "on",
	//                     "tx_scatter_gather_fraglist": "off [fixed]",
	//                     "tx_sctp_segmentation": "off [fixed]",
	//                     "tx_sit_segmentation": "off [fixed]",
	//                     "tx_tcp6_segmentation": "off [fixed]",
	//                     "tx_tcp_ecn_segmentation": "off [fixed]",
	//                     "tx_tcp_mangleid_segmentation": "off",
	//                     "tx_tcp_segmentation": "on",
	//                     "tx_udp_tnl_csum_segmentation": "off [fixed]",
	//                     "tx_udp_tnl_segmentation": "off [fixed]",
	//                     "tx_vlan_offload": "on [fixed]",
	//                     "tx_vlan_stag_hw_insert": "off [fixed]",
	//                     "udp_fragmentation_offload": "off [fixed]",
	//                     "vlan_challenged": "off [fixed]"
	//                   },
	//                   "hw_timestamp_filters": [],
	//                   "ipv4": {
	//                     "address": "172.16.74.100",
	//                     "broadcast": "172.16.74.255",
	//                     "netmask": "255.255.255.0",
	//                     "network": "172.16.74.0"
	//                   },
	//                   "ipv6": [
	//                     {
	//                       "address": "fe80::a00:27ff:fe2a:23e4",
	//                       "prefix": "64",
	//                       "scope": "link"
	//                     }
	//                   ],
	//                   "macaddress": "08:00:27:2a:23:e4",
	//                   "module": "e1000",
	//                   "mtu": 1500,
	//                   "pciid": "0000:00:08.0",
	//                   "promisc": false,
	//                   "speed": 1000,
	//                   "timestamping": [
	//                     "tx_software",
	//                     "rx_software",
	//                     "software"
	//                   ],
	//                   "type": "ether"
	//                 },
	//                 "ansible_fibre_channel_wwn": [],
	//                 "ansible_fips": false,
	//                 "ansible_form_factor": "Other",
	//                 "ansible_fqdn": "managed10",
	//                 "ansible_hostname": "managed10",
	//                 "ansible_hostnqn": "",
	//                 "ansible_interfaces": [
	//                   "lo",
	//                   "eth1",
	//                   "eth0"
	//                 ],
	//                 "ansible_is_chroot": false,
	//                 "ansible_iscsi_iqn": "",
	//                 "ansible_kernel": "3.10.0-1127.el7.x86_64",
	//                 "ansible_kernel_version": "#1 SMP Tue Mar 31 23:36:51 UTC 2020",
	//                 "ansible_lo": {
	//                   "active": true,
	//                   "device": "lo",
	//                   "features": {
	//                     "busy_poll": "off [fixed]",
	//                     "fcoe_mtu": "off [fixed]",
	//                     "generic_receive_offload": "on",
	//                     "generic_segmentation_offload": "on",
	//                     "highdma": "on [fixed]",
	//                     "hw_tc_offload": "off [fixed]",
	//                     "l2_fwd_offload": "off [fixed]",
	//                     "large_receive_offload": "off [fixed]",
	//                     "loopback": "on [fixed]",
	//                     "netns_local": "on [fixed]",
	//                     "ntuple_filters": "off [fixed]",
	//                     "receive_hashing": "off [fixed]",
	//                     "rx_all": "off [fixed]",
	//                     "rx_checksumming": "on [fixed]",
	//                     "rx_fcs": "off [fixed]",
	//                     "rx_gro_hw": "off [fixed]",
	//                     "rx_udp_tunnel_port_offload": "off [fixed]",
	//                     "rx_vlan_filter": "off [fixed]",
	//                     "rx_vlan_offload": "off [fixed]",
	//                     "rx_vlan_stag_filter": "off [fixed]",
	//                     "rx_vlan_stag_hw_parse": "off [fixed]",
	//                     "scatter_gather": "on",
	//                     "tcp_segmentation_offload": "on",
	//                     "tx_checksum_fcoe_crc": "off [fixed]",
	//                     "tx_checksum_ip_generic": "on [fixed]",
	//                     "tx_checksum_ipv4": "off [fixed]",
	//                     "tx_checksum_ipv6": "off [fixed]",
	//                     "tx_checksum_sctp": "on [fixed]",
	//                     "tx_checksumming": "on",
	//                     "tx_fcoe_segmentation": "off [fixed]",
	//                     "tx_gre_csum_segmentation": "off [fixed]",
	//                     "tx_gre_segmentation": "off [fixed]",
	//                     "tx_gso_partial": "off [fixed]",
	//                     "tx_gso_robust": "off [fixed]",
	//                     "tx_ipip_segmentation": "off [fixed]",
	//                     "tx_lockless": "on [fixed]",
	//                     "tx_nocache_copy": "off [fixed]",
	//                     "tx_scatter_gather": "on [fixed]",
	//                     "tx_scatter_gather_fraglist": "on [fixed]",
	//                     "tx_sctp_segmentation": "on",
	//                     "tx_sit_segmentation": "off [fixed]",
	//                     "tx_tcp6_segmentation": "on",
	//                     "tx_tcp_ecn_segmentation": "on",
	//                     "tx_tcp_mangleid_segmentation": "on",
	//                     "tx_tcp_segmentation": "on",
	//                     "tx_udp_tnl_csum_segmentation": "off [fixed]",
	//                     "tx_udp_tnl_segmentation": "off [fixed]",
	//                     "tx_vlan_offload": "off [fixed]",
	//                     "tx_vlan_stag_hw_insert": "off [fixed]",
	//                     "udp_fragmentation_offload": "on",
	//                     "vlan_challenged": "on [fixed]"
	//                   },
	//                   "hw_timestamp_filters": [],
	//                   "ipv4": {
	//                     "address": "127.0.0.1",
	//                     "broadcast": "",
	//                     "netmask": "255.0.0.0",
	//                     "network": "127.0.0.0"
	//                   },
	//                   "ipv6": [
	//                     {
	//                       "address": "::1",
	//                       "prefix": "128",
	//                       "scope": "host"
	//                     }
	//                   ],
	//                   "mtu": 65536,
	//                   "promisc": false,
	//                   "timestamping": [
	//                     "rx_software",
	//                     "software"
	//                   ],
	//                   "type": "loopback"
	//                 },
	//                 "ansible_local": {},
	//                 "ansible_lsb": {},
	//                 "ansible_machine": "x86_64",
	//                 "ansible_machine_id": "1559907b4611764f85b48cabb6c31e97",
	//                 "ansible_memfree_mb": 724,
	//                 "ansible_memory_mb": {
	//                   "nocache": {
	//                     "free": 855,
	//                     "used": 135
	//                   },
	//                   "real": {
	//                     "free": 724,
	//                     "total": 990,
	//                     "used": 266
	//                   },
	//                   "swap": {
	//                     "cached": 0,
	//                     "free": 2047,
	//                     "total": 2047,
	//                     "used": 0
	//                   }
	//                 },
	//                 "ansible_memtotal_mb": 990,
	//                 "ansible_mounts": [
	//                   {
	//                     "block_available": 9701554,
	//                     "block_size": 4096,
	//                     "block_total": 10480385,
	//                     "block_used": 778831,
	//                     "device": "/dev/sda1",
	//                     "fstype": "xfs",
	//                     "inode_available": 20943272,
	//                     "inode_total": 20971008,
	//                     "inode_used": 27736,
	//                     "mount": "/",
	//                     "options": "rw,seclabel,relatime,attr2,inode64,noquota",
	//                     "size_available": 39737565184,
	//                     "size_total": 42927656960,
	//                     "uuid": "1c419d6c-5064-4a2b-953c-05b2c67edb15"
	//                   }
	//                 ],
	//                 "ansible_nodename": "managed10",
	//                 "ansible_os_family": "RedHat",
	//                 "ansible_pkg_mgr": "yum",
	//                 "ansible_proc_cmdline": {
	//                   "BOOT_IMAGE": "/boot/vmlinuz-3.10.0-1127.el7.x86_64",
	//                   "LANG": "en_US.UTF-8",
	//                   "biosdevname": "0",
	//                   "console": [
	//                     "tty0",
	//                     "ttyS0,115200n8"
	//                   ],
	//                   "crashkernel": "auto",
	//                   "elevator": "noop",
	//                   "net.ifnames": "0",
	//                   "no_timer_check": true,
	//                   "ro": true,
	//                   "root": "UUID=1c419d6c-5064-4a2b-953c-05b2c67edb15"
	//                 },
	//                 "ansible_processor": [
	//                   "0",
	//                   "GenuineIntel",
	//                   "Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz"
	//                 ],
	//                 "ansible_processor_cores": 1,
	//                 "ansible_processor_count": 1,
	//                 "ansible_processor_nproc": 1,
	//                 "ansible_processor_threads_per_core": 1,
	//                 "ansible_processor_vcpus": 1,
	//                 "ansible_product_name": "VirtualBox",
	//                 "ansible_product_serial": "0",
	//                 "ansible_product_uuid": "1559907B-4611-764F-85B4-8CABB6C31E97",
	//                 "ansible_product_version": "1.2",
	//                 "ansible_python": {
	//                   "executable": "/usr/bin/python",
	//                   "has_sslcontext": true,
	//                   "type": "CPython",
	//                   "version": {
	//                     "major": 2,
	//                     "micro": 5,
	//                     "minor": 7,
	//                     "releaselevel": "final",
	//                     "serial": 0
	//                   },
	//                   "version_info": [
	//                     2,
	//                     7,
	//                     5,
	//                     "final",
	//                     0
	//                   ]
	//                 },
	//                 "ansible_python_version": "2.7.5",
	//                 "ansible_real_group_id": 0,
	//                 "ansible_real_user_id": 0,
	//                 "ansible_selinux": {
	//                   "config_mode": "enforcing",
	//                   "mode": "enforcing",
	//                   "policyvers": 31,
	//                   "status": "enabled",
	//                   "type": "targeted"
	//                 },
	//                 "ansible_selinux_python_present": true,
	//                 "ansible_service_mgr": "systemd",
	//                 "ansible_ssh_host_key_ecdsa_public": "AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBKI1Y36wVgubRWJZIq5Zk22TA/k9eICPxbqqhXymqMwwqet8KocoxI66F3qnFQ6TOam7JrThyrZoG5hckS47aR4=",
	//                 "ansible_ssh_host_key_ecdsa_public_keytype": "ecdsa-sha2-nistp256",
	//                 "ansible_ssh_host_key_ed25519_public": "AAAAC3NzaC1lZDI1NTE5AAAAIH5O9zutqVTryEKkCjqFsHfTADyOyriX1AxcolXS8SgE",
	//                 "ansible_ssh_host_key_ed25519_public_keytype": "ssh-ed25519",
	//                 "ansible_ssh_host_key_rsa_public": "AAAAB3NzaC1yc2EAAAADAQABAAABAQC+0MTYLYwXGQHgV8SR5HJLUigOwd5RKdbKEdmYKxKolP0NpowkeVaCoSfMFIsMV1a4mfIgRZQnc+KS/ikiypSkNT1+FW7btnxyF9VfRO7IRjTWC5X1GDMTY95mgWYSeCRmaeFacqgo65EAp7sWLwZ+WHWHAPVjZrkg94U+6ASPNAD6bIAx1YqzhEhCtfbI3JlpP3GyHM6Jdn4w+RH6bT42ivWPQLCGjmlNRx7TnB6EhSsMPC4/Ur1wiGqAXULosrfSjrEd55sKv0lBEr8Rd60NPe4R61eS1MM0GNlctbak92kf54Nfj6CzN96zlIciitIfzwUPs/ChceyqP/feqNDX",
	//                 "ansible_ssh_host_key_rsa_public_keytype": "ssh-rsa",
	//                 "ansible_swapfree_mb": 2047,
	//                 "ansible_swaptotal_mb": 2047,
	//                 "ansible_system": "Linux",
	//                 "ansible_system_capabilities": [
	//                   "cap_chown",
	//                   "cap_dac_override",
	//                   "cap_dac_read_search",
	//                   "cap_fowner",
	//                   "cap_fsetid",
	//                   "cap_kill",
	//                   "cap_setgid",
	//                   "cap_setuid",
	//                   "cap_setpcap",
	//                   "cap_linux_immutable",
	//                   "cap_net_bind_service",
	//                   "cap_net_broadcast",
	//                   "cap_net_admin",
	//                   "cap_net_raw",
	//                   "cap_ipc_lock",
	//                   "cap_ipc_owner",
	//                   "cap_sys_module",
	//                   "cap_sys_rawio",
	//                   "cap_sys_chroot",
	//                   "cap_sys_ptrace",
	//                   "cap_sys_pacct",
	//                   "cap_sys_admin",
	//                   "cap_sys_boot",
	//                   "cap_sys_nice",
	//                   "cap_sys_resource",
	//                   "cap_sys_time",
	//                   "cap_sys_tty_config",
	//                   "cap_mknod",
	//                   "cap_lease",
	//                   "cap_audit_write",
	//                   "cap_audit_control",
	//                   "cap_setfcap",
	//                   "cap_mac_override",
	//                   "cap_mac_admin",
	//                   "cap_syslog",
	//                   "35",
	//                   "36+ep"
	//                 ],
	//                 "ansible_system_capabilities_enforced": "True",
	//                 "ansible_system_vendor": "innotek GmbH",
	//                 "ansible_uptime_seconds": 4522,
	//                 "ansible_user_dir": "/root",
	//                 "ansible_user_gecos": "root",
	//                 "ansible_user_gid": 0,
	//                 "ansible_user_id": "root",
	//                 "ansible_user_shell": "/bin/bash",
	//                 "ansible_user_uid": 0,
	//                 "ansible_userspace_architecture": "x86_64",
	//                 "ansible_userspace_bits": "64",
	//                 "ansible_virtualization_role": "guest",
	//                 "ansible_virtualization_tech_guest": [
	//                   "virtualbox"
	//                 ],
	//                 "ansible_virtualization_tech_host": [],
	//                 "ansible_virtualization_type": "virtualbox",
	//                 "discovered_interpreter_python": "/usr/bin/python",
	//                 "gather_subset": [
	//                   "all"
	//                 ],
	//                 "module_setup": true
	//               },
	//               "changed": false,
	//               "deprecations": [],
	//               "warnings": []
	//             },
	//             "172.16.74.101": {
	//               "_ansible_no_log": false,
	//               "_ansible_verbose_override": true,
	//               "action": "gather_facts",
	//               "ansible_facts": {
	//                 "ansible_all_ipv4_addresses": [
	//                   "172.16.74.101",
	//                   "10.0.2.15"
	//                 ],
	//                 "ansible_all_ipv6_addresses": [
	//                   "fe80::a00:27ff:fef3:c25e",
	//                   "fe80::5054:ff:fe4d:77d3"
	//                 ],
	//                 "ansible_apparmor": {
	//                   "status": "disabled"
	//                 },
	//                 "ansible_architecture": "x86_64",
	//                 "ansible_bios_date": "12/01/2006",
	//                 "ansible_bios_vendor": "innotek GmbH",
	//                 "ansible_bios_version": "VirtualBox",
	//                 "ansible_board_asset_tag": "NA",
	//                 "ansible_board_name": "VirtualBox",
	//                 "ansible_board_serial": "0",
	//                 "ansible_board_vendor": "Oracle Corporation",
	//                 "ansible_board_version": "1.2",
	//                 "ansible_chassis_asset_tag": "NA",
	//                 "ansible_chassis_serial": "NA",
	//                 "ansible_chassis_vendor": "Oracle Corporation",
	//                 "ansible_chassis_version": "NA",
	//                 "ansible_cmdline": {
	//                   "BOOT_IMAGE": "/boot/vmlinuz-3.10.0-1127.el7.x86_64",
	//                   "LANG": "en_US.UTF-8",
	//                   "biosdevname": "0",
	//                   "console": "ttyS0,115200n8",
	//                   "crashkernel": "auto",
	//                   "elevator": "noop",
	//                   "net.ifnames": "0",
	//                   "no_timer_check": true,
	//                   "ro": true,
	//                   "root": "UUID=1c419d6c-5064-4a2b-953c-05b2c67edb15"
	//                 },
	//                 "ansible_date_time": {
	//                   "date": "2024-06-17",
	//                   "day": "17",
	//                   "epoch": "1718606913",
	//                   "hour": "06",
	//                   "iso8601": "2024-06-17T06:48:33Z",
	//                   "iso8601_basic": "20240617T064833957403",
	//                   "iso8601_basic_short": "20240617T064833",
	//                   "iso8601_micro": "2024-06-17T06:48:33.957403Z",
	//                   "minute": "48",
	//                   "month": "06",
	//                   "second": "33",
	//                   "time": "06:48:33",
	//                   "tz": "UTC",
	//                   "tz_dst": "UTC",
	//                   "tz_offset": "+0000",
	//                   "weekday": "Monday",
	//                   "weekday_number": "1",
	//                   "weeknumber": "25",
	//                   "year": "2024"
	//                 },
	//                 "ansible_default_ipv4": {
	//                   "address": "10.0.2.15",
	//                   "alias": "eth0",
	//                   "broadcast": "10.0.2.255",
	//                   "gateway": "10.0.2.2",
	//                   "interface": "eth0",
	//                   "macaddress": "52:54:00:4d:77:d3",
	//                   "mtu": 1500,
	//                   "netmask": "255.255.255.0",
	//                   "network": "10.0.2.0",
	//                   "type": "ether"
	//                 },
	//                 "ansible_default_ipv6": {},
	//                 "ansible_device_links": {
	//                   "ids": {
	//                     "sda": [
	//                       "ata-VBOX_HARDDISK_VBe40f73af-159d8e8c"
	//                     ],
	//                     "sda1": [
	//                       "ata-VBOX_HARDDISK_VBe40f73af-159d8e8c-part1"
	//                     ]
	//                   },
	//                   "labels": {},
	//                   "masters": {},
	//                   "uuids": {
	//                     "sda1": [
	//                       "1c419d6c-5064-4a2b-953c-05b2c67edb15"
	//                     ]
	//                   }
	//                 },
	//                 "ansible_devices": {
	//                   "sda": {
	//                     "holders": [],
	//                     "host": "IDE interface: Intel Corporation 82371AB/EB/MB PIIX4 IDE (rev 01)",
	//                     "links": {
	//                       "ids": [
	//                         "ata-VBOX_HARDDISK_VBe40f73af-159d8e8c"
	//                       ],
	//                       "labels": [],
	//                       "masters": [],
	//                       "uuids": []
	//                     },
	//                     "model": "VBOX HARDDISK",
	//                     "partitions": {
	//                       "sda1": {
	//                         "holders": [],
	//                         "links": {
	//                           "ids": [
	//                             "ata-VBOX_HARDDISK_VBe40f73af-159d8e8c-part1"
	//                           ],
	//                           "labels": [],
	//                           "masters": [],
	//                           "uuids": [
	//                             "1c419d6c-5064-4a2b-953c-05b2c67edb15"
	//                           ]
	//                         },
	//                         "sectors": "83884032",
	//                         "sectorsize": 512,
	//                         "size": "40.00 GB",
	//                         "start": "2048",
	//                         "uuid": "1c419d6c-5064-4a2b-953c-05b2c67edb15"
	//                       }
	//                     },
	//                     "removable": "0",
	//                     "rotational": "1",
	//                     "sas_address": null,
	//                     "sas_device_handle": null,
	//                     "scheduler_mode": "noop",
	//                     "sectors": "83886080",
	//                     "sectorsize": "512",
	//                     "serial": "VBe40f73af",
	//                     "size": "40.00 GB",
	//                     "support_discard": "0",
	//                     "vendor": "ATA",
	//                     "virtual": 1
	//                   }
	//                 },
	//                 "ansible_distribution": "CentOS",
	//                 "ansible_distribution_file_parsed": true,
	//                 "ansible_distribution_file_path": "/etc/redhat-release",
	//                 "ansible_distribution_file_variety": "RedHat",
	//                 "ansible_distribution_major_version": "7",
	//                 "ansible_distribution_release": "Core",
	//                 "ansible_distribution_version": "7.8",
	//                 "ansible_dns": {
	//                   "nameservers": [
	//                     "10.0.2.3"
	//                   ]
	//                 },
	//                 "ansible_domain": "",
	//                 "ansible_effective_group_id": 0,
	//                 "ansible_effective_user_id": 0,
	//                 "ansible_env": {
	//                   "HOME": "/root",
	//                   "LANG": "C",
	//                   "LC_ALL": "C",
	//                   "LC_MESSAGES": "C",
	//                   "LESSOPEN": "||/usr/bin/lesspipe.sh %s",
	//                   "LOGNAME": "root",
	//                   "LS_COLORS": "rs=0:di=38;5;27:ln=38;5;51:mh=44;38;5;15:pi=40;38;5;11:so=38;5;13:do=38;5;5:bd=48;5;232;38;5;11:cd=48;5;232;38;5;3:or=48;5;232;38;5;9:mi=05;48;5;232;38;5;15:su=48;5;196;38;5;15:sg=48;5;11;38;5;16:ca=48;5;196;38;5;226:tw=48;5;10;38;5;16:ow=48;5;10;38;5;21:st=48;5;21;38;5;15:ex=38;5;34:*.tar=38;5;9:*.tgz=38;5;9:*.arc=38;5;9:*.arj=38;5;9:*.taz=38;5;9:*.lha=38;5;9:*.lz4=38;5;9:*.lzh=38;5;9:*.lzma=38;5;9:*.tlz=38;5;9:*.txz=38;5;9:*.tzo=38;5;9:*.t7z=38;5;9:*.zip=38;5;9:*.z=38;5;9:*.Z=38;5;9:*.dz=38;5;9:*.gz=38;5;9:*.lrz=38;5;9:*.lz=38;5;9:*.lzo=38;5;9:*.xz=38;5;9:*.bz2=38;5;9:*.bz=38;5;9:*.tbz=38;5;9:*.tbz2=38;5;9:*.tz=38;5;9:*.deb=38;5;9:*.rpm=38;5;9:*.jar=38;5;9:*.war=38;5;9:*.ear=38;5;9:*.sar=38;5;9:*.rar=38;5;9:*.alz=38;5;9:*.ace=38;5;9:*.zoo=38;5;9:*.cpio=38;5;9:*.7z=38;5;9:*.rz=38;5;9:*.cab=38;5;9:*.jpg=38;5;13:*.jpeg=38;5;13:*.gif=38;5;13:*.bmp=38;5;13:*.pbm=38;5;13:*.pgm=38;5;13:*.ppm=38;5;13:*.tga=38;5;13:*.xbm=38;5;13:*.xpm=38;5;13:*.tif=38;5;13:*.tiff=38;5;13:*.png=38;5;13:*.svg=38;5;13:*.svgz=38;5;13:*.mng=38;5;13:*.pcx=38;5;13:*.mov=38;5;13:*.mpg=38;5;13:*.mpeg=38;5;13:*.m2v=38;5;13:*.mkv=38;5;13:*.webm=38;5;13:*.ogm=38;5;13:*.mp4=38;5;13:*.m4v=38;5;13:*.mp4v=38;5;13:*.vob=38;5;13:*.qt=38;5;13:*.nuv=38;5;13:*.wmv=38;5;13:*.asf=38;5;13:*.rm=38;5;13:*.rmvb=38;5;13:*.flc=38;5;13:*.avi=38;5;13:*.fli=38;5;13:*.flv=38;5;13:*.gl=38;5;13:*.dl=38;5;13:*.xcf=38;5;13:*.xwd=38;5;13:*.yuv=38;5;13:*.cgm=38;5;13:*.emf=38;5;13:*.axv=38;5;13:*.anx=38;5;13:*.ogv=38;5;13:*.ogx=38;5;13:*.aac=38;5;45:*.au=38;5;45:*.flac=38;5;45:*.mid=38;5;45:*.midi=38;5;45:*.mka=38;5;45:*.mp3=38;5;45:*.mpc=38;5;45:*.ogg=38;5;45:*.ra=38;5;45:*.wav=38;5;45:*.axa=38;5;45:*.oga=38;5;45:*.spx=38;5;45:*.xspf=38;5;45:",
	//                   "MAIL": "/var/mail/root",
	//                   "PATH": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin",
	//                   "PWD": "/root",
	//                   "SELINUX_LEVEL_REQUESTED": "",
	//                   "SELINUX_ROLE_REQUESTED": "",
	//                   "SELINUX_USE_CURRENT_RANGE": "",
	//                   "SHELL": "/bin/bash",
	//                   "SHLVL": "2",
	//                   "SSH_CLIENT": "172.16.74.10 37842 22",
	//                   "SSH_CONNECTION": "172.16.74.10 37842 172.16.74.101 22",
	//                   "SSH_TTY": "/dev/pts/1",
	//                   "TERM": "xterm-256color",
	//                   "USER": "root",
	//                   "XDG_RUNTIME_DIR": "/run/user/0",
	//                   "XDG_SESSION_ID": "17",
	//                   "_": "/usr/bin/python"
	//                 },
	//                 "ansible_eth0": {
	//                   "active": true,
	//                   "device": "eth0",
	//                   "features": {
	//                     "busy_poll": "off [fixed]",
	//                     "fcoe_mtu": "off [fixed]",
	//                     "generic_receive_offload": "on",
	//                     "generic_segmentation_offload": "on",
	//                     "highdma": "off [fixed]",
	//                     "hw_tc_offload": "off [fixed]",
	//                     "l2_fwd_offload": "off [fixed]",
	//                     "large_receive_offload": "off [fixed]",
	//                     "loopback": "off [fixed]",
	//                     "netns_local": "off [fixed]",
	//                     "ntuple_filters": "off [fixed]",
	//                     "receive_hashing": "off [fixed]",
	//                     "rx_all": "off",
	//                     "rx_checksumming": "off",
	//                     "rx_fcs": "off",
	//                     "rx_gro_hw": "off [fixed]",
	//                     "rx_udp_tunnel_port_offload": "off [fixed]",
	//                     "rx_vlan_filter": "on [fixed]",
	//                     "rx_vlan_offload": "on",
	//                     "rx_vlan_stag_filter": "off [fixed]",
	//                     "rx_vlan_stag_hw_parse": "off [fixed]",
	//                     "scatter_gather": "on",
	//                     "tcp_segmentation_offload": "on",
	//                     "tx_checksum_fcoe_crc": "off [fixed]",
	//                     "tx_checksum_ip_generic": "on",
	//                     "tx_checksum_ipv4": "off [fixed]",
	//                     "tx_checksum_ipv6": "off [fixed]",
	//                     "tx_checksum_sctp": "off [fixed]",
	//                     "tx_checksumming": "on",
	//                     "tx_fcoe_segmentation": "off [fixed]",
	//                     "tx_gre_csum_segmentation": "off [fixed]",
	//                     "tx_gre_segmentation": "off [fixed]",
	//                     "tx_gso_partial": "off [fixed]",
	//                     "tx_gso_robust": "off [fixed]",
	//                     "tx_ipip_segmentation": "off [fixed]",
	//                     "tx_lockless": "off [fixed]",
	//                     "tx_nocache_copy": "off",
	//                     "tx_scatter_gather": "on",
	//                     "tx_scatter_gather_fraglist": "off [fixed]",
	//                     "tx_sctp_segmentation": "off [fixed]",
	//                     "tx_sit_segmentation": "off [fixed]",
	//                     "tx_tcp6_segmentation": "off [fixed]",
	//                     "tx_tcp_ecn_segmentation": "off [fixed]",
	//                     "tx_tcp_mangleid_segmentation": "off",
	//                     "tx_tcp_segmentation": "on",
	//                     "tx_udp_tnl_csum_segmentation": "off [fixed]",
	//                     "tx_udp_tnl_segmentation": "off [fixed]",
	//                     "tx_vlan_offload": "on [fixed]",
	//                     "tx_vlan_stag_hw_insert": "off [fixed]",
	//                     "udp_fragmentation_offload": "off [fixed]",
	//                     "vlan_challenged": "off [fixed]"
	//                   },
	//                   "hw_timestamp_filters": [],
	//                   "ipv4": {
	//                     "address": "10.0.2.15",
	//                     "broadcast": "10.0.2.255",
	//                     "netmask": "255.255.255.0",
	//                     "network": "10.0.2.0"
	//                   },
	//                   "ipv6": [
	//                     {
	//                       "address": "fe80::5054:ff:fe4d:77d3",
	//                       "prefix": "64",
	//                       "scope": "link"
	//                     }
	//                   ],
	//                   "macaddress": "52:54:00:4d:77:d3",
	//                   "module": "e1000",
	//                   "mtu": 1500,
	//                   "pciid": "0000:00:03.0",
	//                   "promisc": false,
	//                   "speed": 1000,
	//                   "timestamping": [
	//                     "tx_software",
	//                     "rx_software",
	//                     "software"
	//                   ],
	//                   "type": "ether"
	//                 },
	//                 "ansible_eth1": {
	//                   "active": true,
	//                   "device": "eth1",
	//                   "features": {
	//                     "busy_poll": "off [fixed]",
	//                     "fcoe_mtu": "off [fixed]",
	//                     "generic_receive_offload": "on",
	//                     "generic_segmentation_offload": "on",
	//                     "highdma": "off [fixed]",
	//                     "hw_tc_offload": "off [fixed]",
	//                     "l2_fwd_offload": "off [fixed]",
	//                     "large_receive_offload": "off [fixed]",
	//                     "loopback": "off [fixed]",
	//                     "netns_local": "off [fixed]",
	//                     "ntuple_filters": "off [fixed]",
	//                     "receive_hashing": "off [fixed]",
	//                     "rx_all": "off",
	//                     "rx_checksumming": "off",
	//                     "rx_fcs": "off",
	//                     "rx_gro_hw": "off [fixed]",
	//                     "rx_udp_tunnel_port_offload": "off [fixed]",
	//                     "rx_vlan_filter": "on [fixed]",
	//                     "rx_vlan_offload": "on",
	//                     "rx_vlan_stag_filter": "off [fixed]",
	//                     "rx_vlan_stag_hw_parse": "off [fixed]",
	//                     "scatter_gather": "on",
	//                     "tcp_segmentation_offload": "on",
	//                     "tx_checksum_fcoe_crc": "off [fixed]",
	//                     "tx_checksum_ip_generic": "on",
	//                     "tx_checksum_ipv4": "off [fixed]",
	//                     "tx_checksum_ipv6": "off [fixed]",
	//                     "tx_checksum_sctp": "off [fixed]",
	//                     "tx_checksumming": "on",
	//                     "tx_fcoe_segmentation": "off [fixed]",
	//                     "tx_gre_csum_segmentation": "off [fixed]",
	//                     "tx_gre_segmentation": "off [fixed]",
	//                     "tx_gso_partial": "off [fixed]",
	//                     "tx_gso_robust": "off [fixed]",
	//                     "tx_ipip_segmentation": "off [fixed]",
	//                     "tx_lockless": "off [fixed]",
	//                     "tx_nocache_copy": "off",
	//                     "tx_scatter_gather": "on",
	//                     "tx_scatter_gather_fraglist": "off [fixed]",
	//                     "tx_sctp_segmentation": "off [fixed]",
	//                     "tx_sit_segmentation": "off [fixed]",
	//                     "tx_tcp6_segmentation": "off [fixed]",
	//                     "tx_tcp_ecn_segmentation": "off [fixed]",
	//                     "tx_tcp_mangleid_segmentation": "off",
	//                     "tx_tcp_segmentation": "on",
	//                     "tx_udp_tnl_csum_segmentation": "off [fixed]",
	//                     "tx_udp_tnl_segmentation": "off [fixed]",
	//                     "tx_vlan_offload": "on [fixed]",
	//                     "tx_vlan_stag_hw_insert": "off [fixed]",
	//                     "udp_fragmentation_offload": "off [fixed]",
	//                     "vlan_challenged": "off [fixed]"
	//                   },
	//                   "hw_timestamp_filters": [],
	//                   "ipv4": {
	//                     "address": "172.16.74.101",
	//                     "broadcast": "172.16.74.255",
	//                     "netmask": "255.255.255.0",
	//                     "network": "172.16.74.0"
	//                   },
	//                   "ipv6": [
	//                     {
	//                       "address": "fe80::a00:27ff:fef3:c25e",
	//                       "prefix": "64",
	//                       "scope": "link"
	//                     }
	//                   ],
	//                   "macaddress": "08:00:27:f3:c2:5e",
	//                   "module": "e1000",
	//                   "mtu": 1500,
	//                   "pciid": "0000:00:08.0",
	//                   "promisc": false,
	//                   "speed": 1000,
	//                   "timestamping": [
	//                     "tx_software",
	//                     "rx_software",
	//                     "software"
	//                   ],
	//                   "type": "ether"
	//                 },
	//                 "ansible_fibre_channel_wwn": [],
	//                 "ansible_fips": false,
	//                 "ansible_form_factor": "Other",
	//                 "ansible_fqdn": "managed11",
	//                 "ansible_hostname": "managed11",
	//                 "ansible_hostnqn": "",
	//                 "ansible_interfaces": [
	//                   "lo",
	//                   "eth1",
	//                   "eth0"
	//                 ],
	//                 "ansible_is_chroot": false,
	//                 "ansible_iscsi_iqn": "",
	//                 "ansible_kernel": "3.10.0-1127.el7.x86_64",
	//                 "ansible_kernel_version": "#1 SMP Tue Mar 31 23:36:51 UTC 2020",
	//                 "ansible_lo": {
	//                   "active": true,
	//                   "device": "lo",
	//                   "features": {
	//                     "busy_poll": "off [fixed]",
	//                     "fcoe_mtu": "off [fixed]",
	//                     "generic_receive_offload": "on",
	//                     "generic_segmentation_offload": "on",
	//                     "highdma": "on [fixed]",
	//                     "hw_tc_offload": "off [fixed]",
	//                     "l2_fwd_offload": "off [fixed]",
	//                     "large_receive_offload": "off [fixed]",
	//                     "loopback": "on [fixed]",
	//                     "netns_local": "on [fixed]",
	//                     "ntuple_filters": "off [fixed]",
	//                     "receive_hashing": "off [fixed]",
	//                     "rx_all": "off [fixed]",
	//                     "rx_checksumming": "on [fixed]",
	//                     "rx_fcs": "off [fixed]",
	//                     "rx_gro_hw": "off [fixed]",
	//                     "rx_udp_tunnel_port_offload": "off [fixed]",
	//                     "rx_vlan_filter": "off [fixed]",
	//                     "rx_vlan_offload": "off [fixed]",
	//                     "rx_vlan_stag_filter": "off [fixed]",
	//                     "rx_vlan_stag_hw_parse": "off [fixed]",
	//                     "scatter_gather": "on",
	//                     "tcp_segmentation_offload": "on",
	//                     "tx_checksum_fcoe_crc": "off [fixed]",
	//                     "tx_checksum_ip_generic": "on [fixed]",
	//                     "tx_checksum_ipv4": "off [fixed]",
	//                     "tx_checksum_ipv6": "off [fixed]",
	//                     "tx_checksum_sctp": "on [fixed]",
	//                     "tx_checksumming": "on",
	//                     "tx_fcoe_segmentation": "off [fixed]",
	//                     "tx_gre_csum_segmentation": "off [fixed]",
	//                     "tx_gre_segmentation": "off [fixed]",
	//                     "tx_gso_partial": "off [fixed]",
	//                     "tx_gso_robust": "off [fixed]",
	//                     "tx_ipip_segmentation": "off [fixed]",
	//                     "tx_lockless": "on [fixed]",
	//                     "tx_nocache_copy": "off [fixed]",
	//                     "tx_scatter_gather": "on [fixed]",
	//                     "tx_scatter_gather_fraglist": "on [fixed]",
	//                     "tx_sctp_segmentation": "on",
	//                     "tx_sit_segmentation": "off [fixed]",
	//                     "tx_tcp6_segmentation": "on",
	//                     "tx_tcp_ecn_segmentation": "on",
	//                     "tx_tcp_mangleid_segmentation": "on",
	//                     "tx_tcp_segmentation": "on",
	//                     "tx_udp_tnl_csum_segmentation": "off [fixed]",
	//                     "tx_udp_tnl_segmentation": "off [fixed]",
	//                     "tx_vlan_offload": "off [fixed]",
	//                     "tx_vlan_stag_hw_insert": "off [fixed]",
	//                     "udp_fragmentation_offload": "on",
	//                     "vlan_challenged": "on [fixed]"
	//                   },
	//                   "hw_timestamp_filters": [],
	//                   "ipv4": {
	//                     "address": "127.0.0.1",
	//                     "broadcast": "",
	//                     "netmask": "255.0.0.0",
	//                     "network": "127.0.0.0"
	//                   },
	//                   "ipv6": [
	//                     {
	//                       "address": "::1",
	//                       "prefix": "128",
	//                       "scope": "host"
	//                     }
	//                   ],
	//                   "mtu": 65536,
	//                   "promisc": false,
	//                   "timestamping": [
	//                     "rx_software",
	//                     "software"
	//                   ],
	//                   "type": "loopback"
	//                 },
	//                 "ansible_local": {},
	//                 "ansible_lsb": {},
	//                 "ansible_machine": "x86_64",
	//                 "ansible_machine_id": "f0a66ecc61bc74409cfa91adc9ce5a76",
	//                 "ansible_memfree_mb": 728,
	//                 "ansible_memory_mb": {
	//                   "nocache": {
	//                     "free": 859,
	//                     "used": 131
	//                   },
	//                   "real": {
	//                     "free": 728,
	//                     "total": 990,
	//                     "used": 262
	//                   },
	//                   "swap": {
	//                     "cached": 0,
	//                     "free": 2047,
	//                     "total": 2047,
	//                     "used": 0
	//                   }
	//                 },
	//                 "ansible_memtotal_mb": 990,
	//                 "ansible_mounts": [
	//                   {
	//                     "block_available": 9701555,
	//                     "block_size": 4096,
	//                     "block_total": 10480385,
	//                     "block_used": 778830,
	//                     "device": "/dev/sda1",
	//                     "fstype": "xfs",
	//                     "inode_available": 20943272,
	//                     "inode_total": 20971008,
	//                     "inode_used": 27736,
	//                     "mount": "/",
	//                     "options": "rw,seclabel,relatime,attr2,inode64,noquota",
	//                     "size_available": 39737569280,
	//                     "size_total": 42927656960,
	//                     "uuid": "1c419d6c-5064-4a2b-953c-05b2c67edb15"
	//                   }
	//                 ],
	//                 "ansible_nodename": "managed11",
	//                 "ansible_os_family": "RedHat",
	//                 "ansible_pkg_mgr": "yum",
	//                 "ansible_proc_cmdline": {
	//                   "BOOT_IMAGE": "/boot/vmlinuz-3.10.0-1127.el7.x86_64",
	//                   "LANG": "en_US.UTF-8",
	//                   "biosdevname": "0",
	//                   "console": [
	//                     "tty0",
	//                     "ttyS0,115200n8"
	//                   ],
	//                   "crashkernel": "auto",
	//                   "elevator": "noop",
	//                   "net.ifnames": "0",
	//                   "no_timer_check": true,
	//                   "ro": true,
	//                   "root": "UUID=1c419d6c-5064-4a2b-953c-05b2c67edb15"
	//                 },
	//                 "ansible_processor": [
	//                   "0",
	//                   "GenuineIntel",
	//                   "Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz"
	//                 ],
	//                 "ansible_processor_cores": 1,
	//                 "ansible_processor_count": 1,
	//                 "ansible_processor_nproc": 1,
	//                 "ansible_processor_threads_per_core": 1,
	//                 "ansible_processor_vcpus": 1,
	//                 "ansible_product_name": "VirtualBox",
	//                 "ansible_product_serial": "0",
	//                 "ansible_product_uuid": "F0A66ECC-61BC-7440-9CFA-91ADC9CE5A76",
	//                 "ansible_product_version": "1.2",
	//                 "ansible_python": {
	//                   "executable": "/usr/bin/python",
	//                   "has_sslcontext": true,
	//                   "type": "CPython",
	//                   "version": {
	//                     "major": 2,
	//                     "micro": 5,
	//                     "minor": 7,
	//                     "releaselevel": "final",
	//                     "serial": 0
	//                   },
	//                   "version_info": [
	//                     2,
	//                     7,
	//                     5,
	//                     "final",
	//                     0
	//                   ]
	//                 },
	//                 "ansible_python_version": "2.7.5",
	//                 "ansible_real_group_id": 0,
	//                 "ansible_real_user_id": 0,
	//                 "ansible_selinux": {
	//                   "config_mode": "enforcing",
	//                   "mode": "enforcing",
	//                   "policyvers": 31,
	//                   "status": "enabled",
	//                   "type": "targeted"
	//                 },
	//                 "ansible_selinux_python_present": true,
	//                 "ansible_service_mgr": "systemd",
	//                 "ansible_ssh_host_key_ecdsa_public": "AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBD5XdQEiALDjR22aIi58J1NhKq0wVMiISn2niIstEWDapgrwkrMmlShjfBpz0FOGskJWgVd/3ToWSDXyLS0Wu7w=",
	//                 "ansible_ssh_host_key_ecdsa_public_keytype": "ecdsa-sha2-nistp256",
	//                 "ansible_ssh_host_key_ed25519_public": "AAAAC3NzaC1lZDI1NTE5AAAAIO5I3Lhz/vRRRB3QWyeZ2ko63Eb17KpX4L/4phUWb0mU",
	//                 "ansible_ssh_host_key_ed25519_public_keytype": "ssh-ed25519",
	//                 "ansible_ssh_host_key_rsa_public": "AAAAB3NzaC1yc2EAAAADAQABAAABAQC8LNB7wRBOzT/yHR9f2NybMEgHDmFE24xxQKPdsAxoe1/u4mjpvjt99sOXu5ENTclkdPj06SXVK5NvGmlLD522Vkgj/r6imDnpUxU5b8RgANJl3jxWrTQ0SZ47B9FRa5RKv8zDh8m3eUJ6RZ3EJFPIzWhH5z3HPttdgA2HbyzfEB+jyCEvWBnv9OxWyR5xbkrG8uJLdpKr2AzO6zFULR3jBpOzlyJGJktL505clW9xwsDVCpD6YeUl1yGoZZK5bbqqueBkWDt5TxAClQLlF0mGm2Lrmnq5Dp5Klv1QfFr1NeplZvwZBXmbQakDYcAZpOlLXxXwwhXirAQ1McRWaENr",
	//                 "ansible_ssh_host_key_rsa_public_keytype": "ssh-rsa",
	//                 "ansible_swapfree_mb": 2047,
	//                 "ansible_swaptotal_mb": 2047,
	//                 "ansible_system": "Linux",
	//                 "ansible_system_capabilities": [
	//                   "cap_chown",
	//                   "cap_dac_override",
	//                   "cap_dac_read_search",
	//                   "cap_fowner",
	//                   "cap_fsetid",
	//                   "cap_kill",
	//                   "cap_setgid",
	//                   "cap_setuid",
	//                   "cap_setpcap",
	//                   "cap_linux_immutable",
	//                   "cap_net_bind_service",
	//                   "cap_net_broadcast",
	//                   "cap_net_admin",
	//                   "cap_net_raw",
	//                   "cap_ipc_lock",
	//                   "cap_ipc_owner",
	//                   "cap_sys_module",
	//                   "cap_sys_rawio",
	//                   "cap_sys_chroot",
	//                   "cap_sys_ptrace",
	//                   "cap_sys_pacct",
	//                   "cap_sys_admin",
	//                   "cap_sys_boot",
	//                   "cap_sys_nice",
	//                   "cap_sys_resource",
	//                   "cap_sys_time",
	//                   "cap_sys_tty_config",
	//                   "cap_mknod",
	//                   "cap_lease",
	//                   "cap_audit_write",
	//                   "cap_audit_control",
	//                   "cap_setfcap",
	//                   "cap_mac_override",
	//                   "cap_mac_admin",
	//                   "cap_syslog",
	//                   "35",
	//                   "36+ep"
	//                 ],
	//                 "ansible_system_capabilities_enforced": "True",
	//                 "ansible_system_vendor": "innotek GmbH",
	//                 "ansible_uptime_seconds": 4506,
	//                 "ansible_user_dir": "/root",
	//                 "ansible_user_gecos": "root",
	//                 "ansible_user_gid": 0,
	//                 "ansible_user_id": "root",
	//                 "ansible_user_shell": "/bin/bash",
	//                 "ansible_user_uid": 0,
	//                 "ansible_userspace_architecture": "x86_64",
	//                 "ansible_userspace_bits": "64",
	//                 "ansible_virtualization_role": "guest",
	//                 "ansible_virtualization_tech_guest": [
	//                   "virtualbox"
	//                 ],
	//                 "ansible_virtualization_tech_host": [],
	//                 "ansible_virtualization_type": "virtualbox",
	//                 "discovered_interpreter_python": "/usr/bin/python",
	//                 "gather_subset": [
	//                   "all"
	//                 ],
	//                 "module_setup": true
	//               },
	//               "changed": false,
	//               "deprecations": [],
	//               "warnings": []
	//             }
	//           },
	//           "task": {
	//             "duration": {
	//               "end": "2024-06-26T06:23:15.288297Z",
	//               "start": "2024-06-26T06:23:14.154015Z"
	//             },
	//             "id": "5254004d-77d3-22f8-e962-00000000000d",
	//             "name": "Gathering Facts"
	//           }
	//         },
	//         {
	//           "hosts": {
	//             "172.16.74.100": {
	//               "action": "debug",
	//               "changed": false,
	//               "msg": "All items completed",
	//               "results": [
	//                 {
	//                   "_ansible_item_label": "x86_64",
	//                   "_ansible_no_log": false,
	//                   "_ansible_verbose_always": true,
	//                   "ansible_loop_var": "item",
	//                   "changed": false,
	//                   "failed": false,
	//                   "item": "x86_64",
	//                   "msg": "x86_64"
	//                 },
	//                 {
	//                   "_ansible_item_label": "CentOS",
	//                   "_ansible_no_log": false,
	//                   "_ansible_verbose_always": true,
	//                   "ansible_loop_var": "item",
	//                   "changed": false,
	//                   "failed": false,
	//                   "item": "CentOS",
	//                   "msg": "CentOS"
	//                 },
	//                 {
	//                   "_ansible_item_label": 990,
	//                   "_ansible_no_log": false,
	//                   "_ansible_verbose_always": true,
	//                   "ansible_loop_var": "item",
	//                   "changed": false,
	//                   "failed": false,
	//                   "item": 990,
	//                   "msg": 990
	//                 },
	//                 {
	//                   "_ansible_item_label": "172.16.74.100",
	//                   "_ansible_no_log": false,
	//                   "_ansible_verbose_always": true,
	//                   "ansible_loop_var": "item",
	//                   "changed": false,
	//                   "failed": false,
	//                   "item": "172.16.74.100",
	//                   "msg": "172.16.74.100"
	//                 },
	//                 {
	//                   "_ansible_item_label": "10.0.2.15",
	//                   "_ansible_no_log": false,
	//                   "_ansible_verbose_always": true,
	//                   "ansible_loop_var": "item",
	//                   "changed": false,
	//                   "failed": false,
	//                   "item": "10.0.2.15",
	//                   "msg": "10.0.2.15"
	//                 },
	//                 {
	//                   "_ansible_item_label": "UTC",
	//                   "_ansible_no_log": false,
	//                   "_ansible_verbose_always": true,
	//                   "ansible_loop_var": "item",
	//                   "changed": false,
	//                   "failed": false,
	//                   "item": "UTC",
	//                   "msg": "UTC"
	//                 }
	//               ]
	//             },
	//             "172.16.74.101": {
	//               "action": "debug",
	//               "changed": false,
	//               "msg": "All items completed",
	//               "results": [
	//                 {
	//                   "_ansible_item_label": "x86_64",
	//                   "_ansible_no_log": false,
	//                   "_ansible_verbose_always": true,
	//                   "ansible_loop_var": "item",
	//                   "changed": false,
	//                   "failed": false,
	//                   "item": "x86_64",
	//                   "msg": "x86_64"
	//                 },
	//                 {
	//                   "_ansible_item_label": "CentOS",
	//                   "_ansible_no_log": false,
	//                   "_ansible_verbose_always": true,
	//                   "ansible_loop_var": "item",
	//                   "changed": false,
	//                   "failed": false,
	//                   "item": "CentOS",
	//                   "msg": "CentOS"
	//                 },
	//                 {
	//                   "_ansible_item_label": 990,
	//                   "_ansible_no_log": false,
	//                   "_ansible_verbose_always": true,
	//                   "ansible_loop_var": "item",
	//                   "changed": false,
	//                   "failed": false,
	//                   "item": 990,
	//                   "msg": 990
	//                 },
	//                 {
	//                   "_ansible_item_label": "172.16.74.101",
	//                   "_ansible_no_log": false,
	//                   "_ansible_verbose_always": true,
	//                   "ansible_loop_var": "item",
	//                   "changed": false,
	//                   "failed": false,
	//                   "item": "172.16.74.101",
	//                   "msg": "172.16.74.101"
	//                 },
	//                 {
	//                   "_ansible_item_label": "10.0.2.15",
	//                   "_ansible_no_log": false,
	//                   "_ansible_verbose_always": true,
	//                   "ansible_loop_var": "item",
	//                   "changed": false,
	//                   "failed": false,
	//                   "item": "10.0.2.15",
	//                   "msg": "10.0.2.15"
	//                 },
	//                 {
	//                   "_ansible_item_label": "UTC",
	//                   "_ansible_no_log": false,
	//                   "_ansible_verbose_always": true,
	//                   "ansible_loop_var": "item",
	//                   "changed": false,
	//                   "failed": false,
	//                   "item": "UTC",
	//                   "msg": "UTC"
	//                 }
	//               ]
	//             }
	//           },
	//           "task": {
	//             "duration": {
	//               "end": "2024-06-26T06:23:15.519469Z",
	//               "start": "2024-06-26T06:23:15.321888Z"
	//             },
	//             "id": "5254004d-77d3-22f8-e962-000000000009",
	//             "name": "Print with loop (items)"
	//           }
	//         }
	//       ]
	//     }
	//   ],
	//   "stats": {
	//     "172.16.74.100": {
	//       "changed": 0,
	//       "failures": 0,
	//       "ignored": 0,
	//       "ok": 2,
	//       "rescued": 0,
	//       "skipped": 0,
	//       "unreachable": 0
	//     },
	//     "172.16.74.101": {
	//       "changed": 0,
	//       "failures": 0,
	//       "ignored": 0,
	//       "ok": 2,
	//       "rescued": 0,
	//       "skipped": 0,
	//       "unreachable": 0
	//     }
	//   }
	// }`

	// output := ansible.Output(outputJsonData)
	// fmt.Printf("%+v", r)
	// output.Info()

	// gatherFactsExtendJsonData := `{
	// 	"DefaultAnsible":{
	// 		"type": "gather_facts",
	// 		"name": "게더팩트",
	// 		"account": "root",
	// 		"ips": ["172.16.74.100", "172.16.74.101"],
	// 		"password": "1123",
	// 		"playBook": ""
	// 	},
	// 	"options": {}
	// }`

	// atherFactsExtendAnsible := ansible.GetAnsibleFromFactory(gatherFactsExtendJsonData)
	// ansible.Excute(atherFactsExtendAnsible)

	// healthCheckExtendJsonData := `{
	// 	"DefaultAnsible":{
	// 		"type": "health_check",
	// 		"name": "헬스체크테스터",
	// 		"account": "root",
	// 		"ips": ["172.16.74.100", "172.16.74.101"],
	// 		"password": "1123",
	// 		"playBook": ""
	// 	},
	// 	"options": {}
	// }`

	// healthCheckExtendAnsible := ansible.GetAnsibleFromFactory(healthCheckExtendJsonData)
	// ansible.Excute(healthCheckExtendAnsible)

	// changePasswordExtendJsonData := `{
	// 	"DefaultAnsible":{
	// 		"type": "change_password",
	// 		"name": "비번변경테스터",
	// 		"account": "root",
	// 		"ips": ["172.16.74.100", "172.16.74.101"],
	// 		"password": "1123",
	// 		"playBook": ""
	// 	},
	// 	"options": {
	// 		"ACCOUNT": "root",
	// 		"CHANGE_PASSWORD": "1234"
	// 	}
	// }`

	// changePasswordExtendAnsible := ansible.GetAnsibleFromFactory(changePasswordExtendJsonData)
	// ansible.Excute(changePasswordExtendAnsible)

	// changePortExtendJsonData := `{
	// 	"DefaultAnsible":{
	// 		"type": "change_ssh_port",
	// 		"name": "포트변경 테스트",
	// 		"account": "root",
	// 		"ips": ["172.16.74.100", "172.16.74.101"],
	// 		"password": "1234",
	// 		"playBook": ""
	// 	},
	// 	"options": {
	// 		"NEW_PORT": "50000"
	// 	}
	// }`

	// changePortExtendAnsible := ansible.GetAnsibleFromFactory(changePortExtendJsonData)
	// ansible.Excute(changePortExtendAnsible)

	// changePasswordExtendJsonData2 := `{
	// 	"DefaultAnsible":{
	// 		"type": "change_password",
	// 		"name": "비번변경테스터",
	// 		"account": "root",
	// 		"ips": ["172.16.74.100", "172.16.74.101"],
	// 		"password": "1234",
	// 		"playBook": ""
	// 	},
	// 	"options": {
	// 		"ACCOUNT": "root",
	// 		"CHANGE_PASSWORD": "1123"
	// 	}
	// }`

	// changePasswordExtendAnsible2 := ansible.GetAnsibleFromFactory(changePasswordExtendJsonData2)
	// ansible.Excute(changePasswordExtendAnsible2)

	// http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusOK)
	// 	err := json.NewEncoder(w).Encode(output.Debug())
	// 	if err != nil {
	// 		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	// 		return
	// 	}
	// })

	log.Fatal(http.ListenAndServe(":8080", handlers.CorsMiddleware(router.NewRouter())))
}
