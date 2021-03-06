package s2csmall

const (
	//StatusComplete - a status flag for complete tasks
	StatusComplete = "complete"
	//StatusFailed - a status for failed tasks
	StatusFailed = "failed"
	//StatusProcessing - a status for in process items
	StatusProcessing = "processing"
	//StatusOutsourced - this is to indicate the the task tracking has been outsourced
	StatusOutsourced = "outsourced"
	//VCDTaskElementHrefMetaName - the name of the meta data field containing the href for the vcd task
	VCDTaskElementHrefMetaName = "vcd_task_element_href"
	//TaskActionUnDeploy --
	TaskActionUnDeploy = "undeploy_vapp"
	//TaskActionDeploy --
	TaskActionDeploy = "deploy_vapp"
	//TaskActionSelfDestruct --
	TaskActionSelfDestruct = "self_destruct"
	//SkuName2CSmall --
	SkuName2CSmall = "2c.small"
	//VCDUsernameField - name of the field in Procurement meta containing username for vcd
	VCDUsernameField = "username"
	//VCDPasswordField - name of the field in Procurement meta containing password for vcd
	VCDPasswordField = "password"
	//VCDAppIDField - name of the field in Procurement meta containing appid for vcd
	VCDAppIDField = "vapp_id"
	//VCDBaseURIField - name of the field in Procurement meta containing baseuri for vcd
	VCDBaseURIField = "base_uri"
	//VCDTemplateNameField - name of the field in Procurement meta containing template name for vcd
	VCDTemplateNameField = "template_name"
	//SubTaskIDField - name of the field in Task meta containing subtask id for vcd
	SubTaskIDField = "subtask_id"
	//LeaseExpiresFieldName --
	LeaseExpiresFieldName = "lease_expires"
	//InventoryIDFieldName ---
	InventoryIDFieldName = "inventory_id"
	//CredentialsFieldName ----
	CredentialsFieldName = "credentials"
	//VCDServiceName --
	VCDServiceName = "pezdispenser-2csmall-vcd-1"
)
