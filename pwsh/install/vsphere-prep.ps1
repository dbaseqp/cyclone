$cred = Import-Clixml -Path ./lib/creds/vsphere_cred.xml
Connect-VIServer -Server $env:vcenterurl -Credential $cred

# Create Resource Pools
if (!(Get-ResourcePool -Name $env:parentresourcepool)) {
    New-ResourcePool -Name $env:parentresourcepool -Location $env:cluster
}

if (!(Get-ResourcePool -Name $env:presettemplateresourcepool)) {
    New-ResourcePool -Name $env:presettemplateresourcepool -Location $env:parentresourcepool
    New-ResourcePool -Name TipocaTemplateExample -Location $env:presettemplateresourcepool
}

if (!(Get-ResourcePool -Name $env:targetresourcepool)) {
    New-ResourcePool -Name $env:targetresourcepool -Location $env:parentresourcepool
}

# Create Roles
New-VIRole -Name KaminoUsers -Privilege (Get-VIPrivilege -Id System.Anonymous,System.Read,System.View,VApp.PowerOff,VApp.PowerOn,VirtualMachine.Interact.ConsoleInteract,VirtualMachine.Interact.PowerOff,VirtualMachine.Interact.PowerOn,VirtualMachine.Interact.Reset,VirtualMachine.State.RevertToSnapshot)
New-VIRole -Name KaminoUsersCustomPod -Privilege (Get-VIPrivilege -Id System.Anonymous,System.Read,System.View,VApp.PowerOff,VApp.PowerOn,VirtualMachine.Interact.ConsoleInteract,VirtualMachine.Interact.PowerOff,VirtualMachine.Interact.PowerOn,VirtualMachine.Interact.Reset,VirtualMachine.State.CreateSnapshot,VirtualMachine.State.RevertToSnapshot,VirtualMachine.State.RemoveSnapshot,VirtualMachine.State.RenameSnapshot)

# Create Folders
if (!(Get-Folder -Name $env:inventorylocation)) {
    Get-Datacenter -Name $env:datacenter | Get-Folder -Name vm | New-Folder -Name $env:inventorylocation
}

if (!(Get-Folder -Name $env:templatelocation)) {
    Get-Datacenter -Name $env:datacenter | Get-Folder -Name vm | New-Folder -Name $env:templatelocation
}

Get-Folder -Name $env:templatelocation | New-Folder -Name Linux
Get-Folder -Name $env:templatelocation | New-Folder -Name Windows
Get-Folder -Name $env:templatelocation | New-Folder -Name Networking

# Set Permissions
$PermissionOptions = @{
    Role = (Get-VIRole -Name 'ReadOnly');
    Entity = (Get-ResourcePool -Name $env:parentresourcepool);
    Principal = ('Kamino.Labs\Kamino Users')
}
New-VIPermission @PermissionOptions | Out-Null