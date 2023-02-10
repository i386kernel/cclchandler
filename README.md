# Managing Custom Certificates for TKG Cluster 
Custom Certificate Lifecycle Manager(cclcmgr) program will aid you in efficiently managing the lifecycle of custom certificates within a TKG Cluster by adding and deleting them to both the kubeadmcontrolplane and kubeadmconfigtemplate objects.
## Features
- Add custom CA certificates to a TKG cluster
- Delete custom CA certificates from a TKG cluster
- Manage the life cycle of custom certificates for TKG cluster

## Prerequisites
Before you begin, you will need the following:

- Access to a TKG cluster
- Make sure you are in Management Cluster Context
- The binary installed and configured on your machine

##  Download the binary file and set its execute permission to allow it to run
```
curl -o cclcmgr https://github.com/i386kernel/cclcmgr/releases/download/v1.0.0-beta/cclcmgr-v1.0.0-beta-linux-amd64
chmod +x cclcmgr
```

## Usage
The binary provides the following commands to manage custom certificates for a TKG cluster:

### Adding Custom Certificates
To add a custom certificate using existing certificate or pull from the server, use the following command:

```
# Append Certs to TKG with certificate

./cclcmgr --action append --cert <path/to/certificate>

-----

# Pulling the certificates from server

./cclcmgr --action append --server <Server Host Name>
```
This will add the specified certificate to the trusted certificate store on the host machine and configure it for use with the TKG cluster.

### Deleting existing Custom Certificates

To delete an existing custom certificate, use the following command:
```
./cclcmgr delete --cert <path/to/existing-certificate>
```
This will remove the specified certificate from the trusted certificate store on the host machine and remove it from the TKG cluster configuration.

## Conclusion
This binary streamlines the management of custom certificates within your TKG cluster by facilitating interaction with kubeadmcontrolplane, kubeadmconfigtemplates, and machinedeployments objects.