$cred = Import-Clixml -Path ./lib/creds/vsphere_cred.xml
Connect-VIServer -Server $env:vcenterurl -Credential $cred

# Create Resource Pools
if (!(Get-ResourcePool -Name $env:parentresourcepool)) {
    New-ResourcePool -Name $env:parentresourcepool -Location $env:cluster
}

if (!(Get-ResourcePool -Name $env:presettemplateresourcepool)) {
    New-ResourcePool -Name $env:presettemplateresourcepool -Location $env:parentresourcepool
}

if (!(Get-ResourcePool -Name $env:targetresourcepool)) {
    New-ResourcePool -Name $env:targetresourcepool -Location $env:parentresourcepool
}

# Create Roles
New-VIRole -Name KaminoUsers -Privilege (Get-VIPrivilege -Id System.Anonymous,System.Read,System.View,VApp.PowerOff,VApp.PowerOn,VirtualMachine.Interact.ConsoleInteract,VirtualMachine.Interact.PowerOff,VirtualMachine.Interact.PowerOn,VirtualMachine.Interact.Reset,VirtualMachine.State.RevertToSnapshot)
New-VIRole -Name KaminoUsersCustomPod -Privilege (Get-VIPrivilege -Id System.Anonymous,System.Read,System.View,VApp.PowerOff,VApp.PowerOn,VirtualMachine.Interact.ConsoleInteract,VirtualMachine.Interact.PowerOff,VirtualMachine.Interact.PowerOn,VirtualMachine.Interact.Reset,VirtualMachine.State.CreateSnapshot,VirtualMachine.State.RevertToSnapshot,VirtualMachine.State.RemoveSnapshot,VirtualMachine.State.RenameSnapshot)

# Create Folders
if (!(Get-Folder -Name $env:inventorylocation)) {
    New-Folder -Name $env:inventorylocation -Location $env:datacenter
}

# Set Permissions
$PermissionOptions = @{
    Role = (Get-VIRole -Name 'ReadOnly');
    Entity = (Get-ResourcePool -Name $env:parentresourcepool);
    Principal = ('Kamino.Labs\Kamino Users')
}
New-VIPermission @PermissionOptions | Out-Null