# sapccmsget
Utility to get performance data from SAP CCMS monitoring tree element by full name. Used by zabbix-agent.

# example: 

Command 
>`[user@host1 ~]$ zabbix_agentd -c /etc/zabbix/zabbix_agentd.conf -t "sap.ccms.get[SID, 'SID\\host1_SID_00\\R3Services\\Dialog\\ResponseTime']"`

will produce output:
>`sap.ccms.get[SID, 'SID\host1_SID_00\R3Services\Dialog\ResponseTime'] [t|89]`
