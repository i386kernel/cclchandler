package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	KUBEADMCONTROLPLANE   = "/apis/controlplane.cluster.x-k8s.io/v1beta1/namespaces/default/kubeadmcontrolplanes/"
	KUBEADMCONFIGTEMPLATE = "/apis/bootstrap.cluster.x-k8s.io/v1beta1/namespaces/default/kubeadmconfigtemplates/"
	MACHINEDEPLOYMENT     = "/apis/cluster.x-k8s.io/v1beta1/namespaces/default/machinedeployments/"
)

//go:embed overlay.yaml
var ob []byte

var kubeapiserver string
var kubeclient *http.Client
var kclient *rest.Config
var certcontent string

func main() {
	var dcert string
	var dserver string
	app := cli.App{
		Name:  "cclcmgr",
		Usage: "TKG Custom Certificate Manager helps to manage the lifecycle of custom certificates in TKG Cluster",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "action",
				Aliases:  []string{"a"},
				Usage:    "provide action to exectue",
				Required: true,
			},
			&cli.StringFlag{
				Name:        "cert",
				Aliases:     []string{"c"},
				Usage:       "provide a certificate, cert path. eg. ./tkg-custom-ca.crt",
				Required:    false,
				Destination: &dcert,
			},
			&cli.StringFlag{
				Name:        "server",
				Aliases:     []string{"s"},
				Usage:       "Provide Host to pull the certificate, eg. harbor.tanzu.prod",
				Required:    false,
				Destination: &dserver,
			},
		},
		Action: func(ctx *cli.Context) error {
			switch ctx.String("action") {
			case "append":
				if dcert == "" && dserver == "" {
					log.Fatal("execution failed, neither cert or host was provided")
				} else {
					setCertContentsandKubeClient(dcert, dserver)
					appendCerts()
				}
			case "delete":
				if dcert == "" {
					log.Fatal("deletion failed, certificate was not provided for deletion")
				} else {
					setCertContentsandKubeClient(dcert, dserver)
					deleteCerts()
				}
			}
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
	kclient = config
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

func getCertsFromHost(host string) {

	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	// Dial a TLS connection to the server
	conn, err := tls.Dial("tcp", host+":443", config)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Get the server's certificate chain
	certs := conn.ConnectionState().PeerCertificates

	var certextract [][]byte

	// Get certificates
	for _, cert := range certs {
		fmt.Println(cert.Subject)
		pemCert := pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		})
		certextract = append(certextract, pemCert)
	}
	certcontent = string(certextract[0])
}

func setCertContentsandKubeClient(cert, server string) {
	fmt.Println("Checking for KubeConfig File, and Api Server Details...")
	getkubeclient(loadconfig())
	switch {
	case cert != "" && server == "":
		fileContents, err := os.ReadFile(cert)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
		certcontent = string(fileContents)
	case cert == "" && server != "":
		getCertsFromHost(server)
	}
}

func writeFilesForFutureProvisioning() {
	err := os.WriteFile("/root/.config/tanzu/tkg/providers/ytt/03_customizations/overlay.yaml", ob, 0644)
	if err != nil {
		fmt.Println("unable to write embedded file to location", err)
	}
	err = os.WriteFile("/root/.config/tanzu/tkg/providers/ytt/03_customizations/tkg-custom-ca.pem", []byte(certcontent), 0644)
	if err != nil {
		fmt.Println("unable to write tkg-custom-ca.pem to location", err)
	}
}

// appendCerts Appends the cert to
func appendCerts() {
	fmt.Println("Appending Certs: \n", certcontent)
	writeFilesForFutureProvisioning()
	for _, kadm := range getkubeadmconfigTemplatesList(kubeclient) {
		appendKubeAdmCert(kubeclient, kadm)
	}
	for _, md := range getMachineDeployments(kubeclient) {
		fmt.Println("Applying MD", md)
		mergeMachineDeployments(kubeclient, md)
	}
}

func deleteCerts() {
	fmt.Println("Deleting Cert: \n", certcontent)
	for _, kadm := range getkubeadmconfigTemplatesList(kubeclient) {
		deleteKubeAdmConfigCerts(kubeclient, kadm)
	}
	for _, md := range getMachineDeployments(kubeclient) {
		fmt.Println("Applying MD", md)
		mergeMachineDeployments(kubeclient, md)
	}
}

// getkubeadmconfigTemplatesList returns all kubeadmconfigtemplatelist object names
func getkubeadmconfigTemplatesList(client *http.Client) []string {
	resp, err := client.Get(kubeapiserver + KUBEADMCONFIGTEMPLATE)
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
	var custfilename = "tkg-custom-ca" + time.Now().Format("Mon, 02-Jan-06 15:04:05 IST") + ".pem"
	url := KUBEADMCONFIGTEMPLATE + kadmdep
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
		Path:        "/etc/ssl/certs/" + custfilename,
		Permissions: "0644",
	}

	KubeadmConfigTemplate.Spec.Template.Spec.Files = append(KubeadmConfigTemplate.Spec.Template.Spec.Files, newFile)

	cfp := fmt.Sprintf("'! which update-ca-certificates 2>/dev/null || (mv /etc/ssl/certs/%s /usr/local/share/ca-certificates/%s && update-ca-certificates)'", custfilename, custfilename)

	KubeadmConfigTemplate.Spec.Template.Spec.PreKubeadmCommands = []string{"'! which rehash_ca_certificates.sh 2>/dev/null || rehash_ca_certificates.sh'", cfp}
	data, err := json.Marshal(KubeadmConfigTemplate)

	request, err := http.NewRequest("PATCH", kubeapiserver+url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
	request.Header = map[string][]string{"Content-type": {" application/merge-patch+json"}}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyr, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bodyr))
}

// deleteKubeAdmCerts deletes the existing certificates from kubeadmobjects
func deleteKubeAdmConfigCerts(client *http.Client, kadmdep string) {
	url := KUBEADMCONFIGTEMPLATE + kadmdep
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

	for i, v := range KubeadmConfigTemplate.Spec.Template.Spec.Files {
		if v.Content == certcontent {
			KubeadmConfigTemplate.Spec.Template.Spec.Files = append(KubeadmConfigTemplate.Spec.Template.Spec.Files[:i], KubeadmConfigTemplate.Spec.Template.Spec.Files[i+1:]...)
		}
	}
	//KubeadmConfigTemplate.Spec.Template.Spec.Files = KubeadmConfigTemplate.Spec.Template.Spec.Files[:0]
	KubeadmConfigTemplate.Spec.Template.Spec.PreKubeadmCommands = []string{"'! which rehash_ca_certificates.sh 2>/dev/null || rehash_ca_certificates.sh'",
		"'! which update-ca-certificates 2>/dev/null || (mv /etc/ssl/certs/tkg-custom-ca.pem /usr/local/share/ca-certificates/tkg-custom-ca.crt && update-ca-certificates)'"}
	data, err := json.Marshal(KubeadmConfigTemplate)

	request, err := http.NewRequest("PATCH", kubeapiserver+url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
	request.Header = map[string][]string{"Content-type": {" application/merge-patch+json"}}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyr, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bodyr))
}

// getMachineDeployments returns all the machinedpeloyments names
func getMachineDeployments(client *http.Client) []string {
	url := kubeapiserver + MACHINEDEPLOYMENT
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

// mergeMachineDeployments merges the newly created annotation with the current date and time
func mergeMachineDeployments(client *http.Client, mcdep string) {

	url := MACHINEDEPLOYMENT + mcdep
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

func createKappSecret() {
	clientset, err := kubernetes.NewForConfig(kclient)
	if err != nil {
		fmt.Println(err)
		return
	}
	certificate := []byte(certcontent)
	encodedCertificate := base64.StdEncoding.EncodeToString(certificate)

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kapp-controller-config",
			Namespace: "tkg-system",
		},
		Data: map[string][]byte{
			"certificate": []byte(encodedCertificate),
		},
		Type: v1.TLSCertKey,
	}

	result, err := clientset.CoreV1().Secrets("tkg-system").Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Secret %s/%s created\n", result.Namespace, result.Name)
}

// getkubeadmControlPlaneList returns all kubeadmcontrolplane object names
// future implementation if MUTABLE
func getkubeadmControlPlaneList(client *http.Client) []string {
	resp, err := client.Get(kubeapiserver + KUBEADMCONTROLPLANE)
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
	var kubeadmcplist []string
	for _, kadm := range kadmList.Items {
		fmt.Println(kadm.Metadata.Name)
		kubeadmcplist = append(kubeadmcplist, kadm.Metadata.Name)
	}
	fmt.Println(kubeadmcplist)
	return kubeadmcplist
}

// appendKubeAdmCPCert appends the provided certificate to kubeadmcontrolplane object
// future implementation if MUTABLE
func appendKubeAdmCPCert(client *http.Client, kadmcp string) {
	url := KUBEADMCONTROLPLANE + kadmcp
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

	if err := json.Unmarshal(body, &KubeadmControlPlane); err != nil {
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

	fmt.Println(newFile)

	KubeadmControlPlane.Spec.KubeadmConfigSpec.Files = append(KubeadmControlPlane.Spec.KubeadmConfigSpec.Files, newFile)
	KubeadmControlPlane.Spec.KubeadmConfigSpec.PreKubeadmCommands = []string{"'! which rehash_ca_certificates.sh 2>/dev/null || rehash_ca_certificates.sh'", "'! which update-ca-certificates 2>/dev/null || (mv /etc/ssl/certs/tkg-custom-ca.pem /usr/local/share/ca-certificates/tkg-custom-ca.crt && update-ca-certificates)'"}

	data, err := json.Marshal(KubeadmControlPlane)
	if err != nil {
		fmt.Println(err)
	}
	request, err := http.NewRequest("POST", kubeapiserver+url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
	}
	request.Header = map[string][]string{"Content-type": {" application/json"}}
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

// deleteKubeAdmCPCerts deletes all the certificates in kubeadmcontrolplane object
// future-implementation if MUTABLE
func deleteKubeAdmCPCerts(client *http.Client, kadmcp string) {
	url := KUBEADMCONTROLPLANE + kadmcp
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

	if err := json.Unmarshal(body, &KubeadmControlPlane); err != nil {
		fmt.Fprintln(os.Stderr, "Error unmarshaling response:", err)
		os.Exit(1)
	}

	KubeadmControlPlane.Spec.KubeadmConfigSpec.Files = KubeadmControlPlane.Spec.KubeadmConfigSpec.Files[:0]
	KubeadmControlPlane.Spec.KubeadmConfigSpec.PreKubeadmCommands = []string{"'! which rehash_ca_certificates.sh 2>/dev/null || rehash_ca_certificates.sh'", "'! which update-ca-certificates 2>/dev/null || (mv /etc/ssl/certs/tkg-custom-ca.pem /usr/local/share/ca-certificates/tkg-custom-ca.crt && update-ca-certificates)'"}
	data, err := json.Marshal(KubeadmControlPlane)

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
