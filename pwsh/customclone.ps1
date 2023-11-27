param(
[String] $LabName,
[String] $Username,
[String[]] $VMsToClone,
[Boolean] $Natted,
[int] $PortGroup,
[String] $Target,
[String] $Domain,
[String] $WanPortGroup
)

$pg = [int] $PortGroup

$cred = Import-CliXML -Path .\lib\creds\vsphere_cred.xml
Connect-VIServer {vcenterfqdn} -Credential $cred

Invoke-CustomPod -LabName $LabName -Username $Username -Natted $Natted -Target $Target -Portgroup $pg -Domain $domain -WanPortGroup $WanPortGroup -VMsToClone $VMsToClone 