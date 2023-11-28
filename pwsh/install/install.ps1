Install-Module -Name VMware.PowerCLI -Force
Set-PowerCLIConfiguration -InvalidCertificateAction Ignore -Scope AllUsers -Confirm:0 
Set-PowerCLIConfiguration -ParticipateInCEIP:0 -Scope AllUsers -Confirm:0
[PSCredential]::New($env:vcenterusername, (ConvertTo-SecureString -String $env:vcenterpassword -AsPlainText -Force)) | Export-Clixml -Path ./lib/creds/vsphere_cred.xml