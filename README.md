# Managing Custom Certificates for TKG Cluster
This program will aid you in efficiently managing the lifecycle of custom certificates within a TKG Cluster by adding and deleting them to both the kubeadmcontrolplane and kubeadmconfigtemplate objects
## Features
- Add custom CA certificates to a TKG cluster
- Delete custom CA certificates from a TKG cluster
- Manage the life cycle of custom certificates for TKG cluster

## Prerequisites
Before you begin, you will need the following:

- Access to a TKG cluster
- The binary installed and configured on your machine

## Usage
The binary provides the following commands to manage custom certificates for a TKG cluster:

###Adding Custom Certificates
To add a custom certificate, use the following command:

```
./customcerthandler -a append --cert <path/to/certificate>
```
This will add the specified certificate to the trusted certificate store on the host machine and configure it for use with the TKG cluster.

### Deleting Custom Certificates

To delete a custom certificate, use the following command:
```
./customcerthandler -a delete --cert <path/to/certificate>
```
This will remove the specified certificate from the trusted certificate store on the host machine and remove it from the TKG cluster configuration.

## Conclusion
With this binary, you can easily manage the life cycle of custom certificates for your TKG cluster. This allows you to securely communicate with your cluster and maintain the security of your infrastructure.