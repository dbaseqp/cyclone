param(
[String] $Template,
[String] $Username,
[String] $Port,
[String] $Target,
[String] $Domain,
[String] $WanPG
)

$pg = [int] $Port

$cred = Import-CliXML -Path .\pwsh\vsphere_cred.xml
Connect-VIServer elsa.sdc.cpp -Credential $cred

Invoke-WebClone -SourceResourcePool $Template -Target $Target -Portgroup $pg -Domain $domain -WanPortGroup $WanPG -Username $Username
