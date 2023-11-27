param(
[String] $Username,
[String] $Tag
)

$cred = Import-CliXML -Path .\lib\creds\vsphere_cred.xml
Connect-VIServer {vcenterfqdn} -Credential $cred
Invoke-OrderSixtySix -Username $Username -Tag $Tag