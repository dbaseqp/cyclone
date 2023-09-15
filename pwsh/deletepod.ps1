param(
[String] $Username,
[String] $Tag
)

$cred = Import-CliXML -Path $env:ProgramFiles\Kamino\lib\creds\vsphere_cred.xml
Connect-VIServer elsa.sdc.cpp -Credential $cred
Invoke-OrderSixtySix -Username $Username -Tag $Tag