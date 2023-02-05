package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"io"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var kubeapiserver string
var kubeclient *http.Client
var certcontent string

func init() {
	fmt.Println("Checking for KubeConfig File, and Api Server Details...")
	getkubeclient(loadconfig())
}

func main() {
	app := cli.NewApp()
	app.Name = "customcerthandler"
	app.Usage = "TKG Custom Certificate Handler helps to manage the lifecycle of custom certificates in TKG Cluster"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:     "action, a",
			Usage:    "Select an action [append or delete] to execute, Either to Append Certs or Delete them",
			Required: true,
		},
		cli.StringFlag{
			Name:     "cert, c",
			Usage:    "provide a certificate, cert path. eg. ./tkg-custom-ca.crt",
			Required: true,
		},
	}

	app.Action = func(c *cli.Context) error {
		switch c.String("action") {
		case "append":
			appendCerts(c.String("cert"))
		case "delete":
			deleteCerts(c.String("cert"))
		default:
			fmt.Println("Invalid option")
			err := cli.ShowAppHelp(c)
			if err != nil {
				return err
			}
		}
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func appendCerts(cert string) {
	fileContents, err := os.ReadFile(cert)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	certcontent = string(fileContents)

	for _, kadm := range getkubeadmconfigTemplatesList(kubeclient) {
		appendKubeAdmCert(kubeclient, kadm)
	}
	for _, md := range getMachineDeployments(kubeclient) {
		fmt.Println("Applying MD", md)
		mergeMachineDeployments(kubeclient, md)
	}
}

func deleteCerts(cert string) {
	fileContents, err := os.ReadFile(cert)
	certcontent = string(fileContents)
	fmt.Println(certcontent)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	for _, kadm := range getkubeadmconfigTemplatesList(kubeclient) {
		deleteKubeAdmCerts(kubeclient, kadm)
	}
	for _, md := range getMachineDeployments(kubeclient) {
		fmt.Println("Applying MD", md)
		mergeMachineDeployments(kubeclient, md)
	}
}

// loadconfig
func loadconfig() *rest.Config {
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
	fmt.Println("Accessing kubeconfig from:", "'"+kubeconfigPath+"'")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading kubeconfig:", err)
		os.Exit(1)
	}
	kubeapiserver = config.Host
	fmt.Println("Here's the Current Host Details \n", kubeapiserver)
	return config
}

// getkubeclient creates a http client for kubernetes cluster in the current context
func getkubeclient(config *rest.Config) {
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(config.CAData)
	clientCert, err := tls.X509KeyPair(config.CertData, config.KeyData)
	if err != nil {
		log.Fatal(err)
	}
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      caCertPool,
			Certificates: []tls.Certificate{clientCert},
		},
	}
	kubeclient = &http.Client{Transport: transport}
}

func getkubeadmconfigTemplatesList(client *http.Client) []string {
	resp, err := client.Get(kubeapiserver + "/apis/bootstrap.cluster.x-k8s.io/v1beta1/namespaces/default/kubeadmconfigtemplates")
	if err != nil {
		log.Fatal("unable to retrieve with the given object", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "Unexpected status code:", resp.StatusCode)
		os.Exit(1)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading response body:", err)
		os.Exit(1)
	}
	var kadmList struct {
		Items []struct {
			Metadata struct {
				Name string `json:"name"`
			} `json:"metadata"`
		} `json:"items"`
	}
	if err := json.Unmarshal(body, &kadmList); err != nil {
		fmt.Fprintln(os.Stderr, "Error unmarshaling response:", err)
		os.Exit(1)
	}
	var kubeadmlist []string
	for _, kadm := range kadmList.Items {
		fmt.Println(kadm.Metadata.Name)
		kubeadmlist = append(kubeadmlist, kadm.Metadata.Name)
	}
	return kubeadmlist
}

// appendKubeAdmCert updates kubeadm object
func appendKubeAdmCert(client *http.Client, kadmdep string) {
	url := "/apis/bootstrap.cluster.x-k8s.io/v1beta1/namespaces/default/kubeadmconfigtemplates/" + kadmdep
	req, err := client.Get(kubeapiserver + url)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(req.Body)
	if req.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "Unexpected status code:", req.StatusCode)
		os.Exit(1)
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading response body:", err)
		os.Exit(1)
	}

	if err := json.Unmarshal(body, &KubeadmConfigTemplate); err != nil {
		fmt.Fprintln(os.Stderr, "Error unmarshaling response:", err)
		os.Exit(1)
	}
	newFile := struct {
		Content     string `json:"content"`
		Owner       string `json:"owner"`
		Path        string `json:"path"`
		Permissions string `json:"permissions"`
	}{
		Content:     certcontent,
		Owner:       "root",
		Path:        "/etc/ssl/certs/tkg-custom-ca.pem",
		Permissions: "0644",
	}

	KubeadmConfigTemplate.Spec.Template.Spec.Files = append(KubeadmConfigTemplate.Spec.Template.Spec.Files, newFile)
	KubeadmConfigTemplate.Spec.Template.Spec.PreKubeadmCommands = []string{"'! which rehash_ca_certificates.sh 2>/dev/null || rehash_ca_certificates.sh'", "'! which update-ca-certificates 2>/dev/null || (mv /etc/ssl/certs/tkg-custom-ca.pem /usr/local/share/ca-certificates/tkg-custom-ca.crt && update-ca-certificates)'"}
	data, err := json.Marshal(KubeadmConfigTemplate)

	request, err := http.NewRequest("PATCH", kubeapiserver+url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
	}
	request.Header = map[string][]string{"Content-type": {" application/merge-patch+json"}}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	bodyr, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bodyr))
}

// deleteKubeAdmCerts deletes the existing certificates from kubeadmobjects
func deleteKubeAdmCerts(client *http.Client, kadmdep string) {
	url := "/apis/bootstrap.cluster.x-k8s.io/v1beta1/namespaces/default/kubeadmconfigtemplates/" + kadmdep
	req, err := client.Get(kubeapiserver + url)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(req.Body)
	if req.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "Unexpected status code:", req.StatusCode)
		os.Exit(1)
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading response body:", err)
		os.Exit(1)
	}

	if err := json.Unmarshal(body, &KubeadmConfigTemplate); err != nil {
		fmt.Fprintln(os.Stderr, "Error unmarshaling response:", err)
		os.Exit(1)
	}

	KubeadmConfigTemplate.Spec.Template.Spec.Files = KubeadmConfigTemplate.Spec.Template.Spec.Files[:0]
	KubeadmConfigTemplate.Spec.Template.Spec.PreKubeadmCommands = []string{"'! which rehash_ca_certificates.sh 2>/dev/null || rehash_ca_certificates.sh'", "'! which update-ca-certificates 2>/dev/null || (mv /etc/ssl/certs/tkg-custom-ca.pem /usr/local/share/ca-certificates/tkg-custom-ca.crt && update-ca-certificates)'"}
	data, err := json.Marshal(KubeadmConfigTemplate)

	request, err := http.NewRequest("PATCH", kubeapiserver+url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
	}
	request.Header = map[string][]string{"Content-type": {" application/merge-patch+json"}}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	bodyr, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bodyr))
}

func getMachineDeployments(client *http.Client) []string {
	url := kubeapiserver + "/apis/cluster.x-k8s.io/v1beta1/namespaces/default/machinedeployments/"
	resp, err := client.Get(url)
	if err != nil {
		log.Fatal("unable to retrieve with the given object", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "Unexpected status code:", resp.StatusCode)
		os.Exit(1)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading response body:", err)
		os.Exit(1)
	}
	var mDepList struct {
		Items []struct {
			Metadata struct {
				Name string `json:"name"`
			} `json:"metadata"`
		} `json:"items"`
	}
	if err := json.Unmarshal(body, &mDepList); err != nil {
		_, err := fmt.Fprintln(os.Stderr, "Error unmarshaling response:", err)
		if err != nil {
			return nil
		}
		os.Exit(1)
	}
	var mdep []string
	for _, kadm := range mDepList.Items {
		fmt.Println(kadm.Metadata.Name)
		mdep = append(mdep, kadm.Metadata.Name)
	}
	return mdep
}

func mergeMachineDeployments(client *http.Client, mcdep string) {

	url := "/apis/cluster.x-k8s.io/v1beta1/namespaces/default/machinedeployments/" + mcdep
	req, err := client.Get(kubeapiserver + url)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(req.Body)
	if req.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "Unexpected status code:", req.StatusCode)
		os.Exit(1)
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading response body:", err)
		os.Exit(1)
	}
	if err := json.Unmarshal(body, &MachineDeployment); err != nil {
		fmt.Fprintln(os.Stderr, "Error unmarshaling response:", err)
		os.Exit(1)
	}

	getcurrenttime := time.Now().Format("Wed Feb 25 11:06:39 PST 2015")

	mdannotate := struct {
		Date                            string `yaml:"date"`
		RunTanzuVmwareComResolveOsImage string `yaml:"run.tanzu.vmware.com/resolve-os-image"`
	}{
		Date:                            getcurrenttime,
		RunTanzuVmwareComResolveOsImage: "run.tanzu.vmware.com/resolve-os-image",
	}

	MachineDeployment.Spec.Template.Metadata.Annotations = struct {
		Date                            string `json:"date"`
		RunTanzuVmwareComResolveOsImage string `json:"run.tanzu.vmware.com/resolve-os-image"`
	}(mdannotate)

	data, err := json.Marshal(MachineDeployment)
	fmt.Println(string(data))

	request, err := http.NewRequest("PATCH", kubeapiserver+url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
	}
	request.Header = map[string][]string{"Content-type": {"application/merge-patch+json"}}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("STATUS CODE: \n", resp.StatusCode)
	defer resp.Body.Close()
	bodyr, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bodyr))
}
