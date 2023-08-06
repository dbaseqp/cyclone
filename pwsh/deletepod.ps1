param(
[String] $Username,
[String] $Tag
)

#$cred = Import-CliXML -Path .\pwsh\vsphere_cred.xml
Connect-VIServer elsa.sdc.cpp -User bruharmycloner -Password '$w1ftCCDC!123'#-Credential $cred
Invoke-OrderSixtySix -Username $Username -Tag $Tag

