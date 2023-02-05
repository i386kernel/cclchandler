package main

import "time"

var KubeadmControlPlane struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Annotations struct {
			ClusterXK8SIoClonedFromGroupkind string `json:"cluster.x-k8s.io/cloned-from-groupkind"`
			ClusterXK8SIoClonedFromName      string `json:"cluster.x-k8s.io/cloned-from-name"`
		} `json:"annotations"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
		Finalizers        []string  `json:"finalizers"`
		Generation        int       `json:"generation"`
		Labels            struct {
			ClusterXK8SIoClusterName   string `json:"cluster.x-k8s.io/cluster-name"`
			TopologyClusterXK8SIoOwned string `json:"topology.cluster.x-k8s.io/owned"`
		} `json:"labels"`
		Name            string `json:"name"`
		Namespace       string `json:"namespace"`
		OwnerReferences []struct {
			APIVersion         string `json:"apiVersion"`
			BlockOwnerDeletion bool   `json:"blockOwnerDeletion"`
			Controller         bool   `json:"controller"`
			Kind               string `json:"kind"`
			Name               string `json:"name"`
			UID                string `json:"uid"`
		} `json:"ownerReferences"`
		ResourceVersion string `json:"resourceVersion"`
		UID             string `json:"uid"`
	} `json:"metadata"`
	Spec struct {
		KubeadmConfigSpec struct {
			ClusterConfiguration struct {
				APIServer struct {
					ExtraArgs struct {
						AdmissionControlConfigFile string `json:"admission-control-config-file"`
						CloudProvider              string `json:"cloud-provider"`
						TLSCipherSuites            string `json:"tls-cipher-suites"`
					} `json:"extraArgs"`
					ExtraVolumes []struct {
						HostPath  string `json:"hostPath"`
						MountPath string `json:"mountPath"`
						Name      string `json:"name"`
						PathType  string `json:"pathType"`
						ReadOnly  bool   `json:"readOnly"`
					} `json:"extraVolumes"`
					TimeoutForControlPlane string `json:"timeoutForControlPlane"`
				} `json:"apiServer"`
				ControllerManager struct {
					ExtraArgs struct {
						CloudProvider   string `json:"cloud-provider"`
						TLSCipherSuites string `json:"tls-cipher-suites"`
					} `json:"extraArgs"`
				} `json:"controllerManager"`
				DNS struct {
					ImageRepository string `json:"imageRepository"`
					ImageTag        string `json:"imageTag"`
				} `json:"dns"`
				Etcd struct {
					Local struct {
						DataDir   string `json:"dataDir"`
						ExtraArgs struct {
							CipherSuites string `json:"cipher-suites"`
						} `json:"extraArgs"`
						ImageRepository string `json:"imageRepository"`
						ImageTag        string `json:"imageTag"`
					} `json:"local"`
				} `json:"etcd"`
				ImageRepository string `json:"imageRepository"`
				Networking      struct {
				} `json:"networking"`
				Scheduler struct {
					ExtraArgs struct {
						TLSCipherSuites string `json:"tls-cipher-suites"`
					} `json:"extraArgs"`
				} `json:"scheduler"`
			} `json:"clusterConfiguration"`
			Files []struct {
				Content     string `json:"content"`
				Owner       string `json:"owner"`
				Path        string `json:"path"`
				Permissions string `json:"permissions"`
			} `json:"files"`
			Format            string `json:"format"`
			InitConfiguration struct {
				LocalAPIEndpoint struct {
				} `json:"localAPIEndpoint"`
				NodeRegistration struct {
					CriSocket        string `json:"criSocket"`
					KubeletExtraArgs struct {
						CloudProvider   string `json:"cloud-provider"`
						NodeLabels      string `json:"node-labels"`
						TLSCipherSuites string `json:"tls-cipher-suites"`
					} `json:"kubeletExtraArgs"`
					Name string `json:"name"`
				} `json:"nodeRegistration"`
			} `json:"initConfiguration"`
			JoinConfiguration struct {
				Discovery struct {
				} `json:"discovery"`
				NodeRegistration struct {
					CriSocket        string `json:"criSocket"`
					KubeletExtraArgs struct {
						CloudProvider   string `json:"cloud-provider"`
						NodeLabels      string `json:"node-labels"`
						TLSCipherSuites string `json:"tls-cipher-suites"`
					} `json:"kubeletExtraArgs"`
					Name string `json:"name"`
				} `json:"nodeRegistration"`
			} `json:"joinConfiguration"`
			PreKubeadmCommands       []string `json:"preKubeadmCommands"`
			UseExperimentalRetryJoin bool     `json:"useExperimentalRetryJoin"`
			Users                    []struct {
				Name              string   `json:"name"`
				SSHAuthorizedKeys []string `json:"sshAuthorizedKeys"`
				Sudo              string   `json:"sudo"`
			} `json:"users"`
		} `json:"kubeadmConfigSpec"`
		MachineTemplate struct {
			InfrastructureRef struct {
				APIVersion string `json:"apiVersion"`
				Kind       string `json:"kind"`
				Name       string `json:"name"`
				Namespace  string `json:"namespace"`
			} `json:"infrastructureRef"`
			Metadata struct {
				Annotations struct {
					RunTanzuVmwareComResolveOsImage string `json:"run.tanzu.vmware.com/resolve-os-image"`
				} `json:"annotations"`
				Labels struct {
					ClusterXK8SIoClusterName   string `json:"cluster.x-k8s.io/cluster-name"`
					TopologyClusterXK8SIoOwned string `json:"topology.cluster.x-k8s.io/owned"`
				} `json:"labels"`
			} `json:"metadata"`
		} `json:"machineTemplate"`
		Replicas      int `json:"replicas"`
		RolloutBefore struct {
			CertificatesExpiryDays int `json:"certificatesExpiryDays"`
		} `json:"rolloutBefore"`
		RolloutStrategy struct {
			RollingUpdate struct {
				MaxSurge int `json:"maxSurge"`
			} `json:"rollingUpdate"`
			Type string `json:"type"`
		} `json:"rolloutStrategy"`
		Version string `json:"version"`
	} `json:"spec"`
	Status struct {
		Conditions []struct {
			LastTransitionTime time.Time `json:"lastTransitionTime"`
			Status             string    `json:"status"`
			Type               string    `json:"type"`
		} `json:"conditions"`
		Initialized         bool   `json:"initialized"`
		ObservedGeneration  int    `json:"observedGeneration"`
		Ready               bool   `json:"ready"`
		ReadyReplicas       int    `json:"readyReplicas"`
		Replicas            int    `json:"replicas"`
		Selector            string `json:"selector"`
		UnavailableReplicas int    `json:"unavailableReplicas"`
		UpdatedReplicas     int    `json:"updatedReplicas"`
		Version             string `json:"version"`
	} `json:"status"`
}

var KubeadmConfigTemplate struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Annotations struct {
			ClusterXK8SIoClonedFromGroupkind string `json:"cluster.x-k8s.io/cloned-from-groupkind"`
			ClusterXK8SIoClonedFromName      string `json:"cluster.x-k8s.io/cloned-from-name"`
			KappK14SIoIdentity               string `json:"kapp.k14s.io/identity"`
			KappK14SIoOriginal               string `json:"kapp.k14s.io/original"`
			KappK14SIoOriginalDiffMd5        string `json:"kapp.k14s.io/original-diff-md5"`
		} `json:"annotations"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
		Generation        int       `json:"generation"`
		Labels            struct {
			ClusterXK8SIoClusterName            string `json:"cluster.x-k8s.io/cluster-name"`
			KappK14SIoApp                       string `json:"kapp.k14s.io/app"`
			KappK14SIoAssociation               string `json:"kapp.k14s.io/association"`
			TopologyClusterXK8SIoDeploymentName string `json:"topology.cluster.x-k8s.io/deployment-name"`
			TopologyClusterXK8SIoOwned          string `json:"topology.cluster.x-k8s.io/owned"`
		} `json:"labels"`
		Name            string `json:"name"`
		Namespace       string `json:"namespace"`
		OwnerReferences []struct {
			APIVersion string `json:"apiVersion"`
			Kind       string `json:"kind"`
			Name       string `json:"name"`
			UID        string `json:"uid"`
		} `json:"ownerReferences"`
		ResourceVersion string `json:"resourceVersion"`
		UID             string `json:"uid"`
	} `json:"metadata"`
	Spec struct {
		Template struct {
			Spec struct {
				Files []struct {
					Content     string `json:"content"`
					Owner       string `json:"owner"`
					Path        string `json:"path"`
					Permissions string `json:"permissions"`
				} `json:"files"`
				Format            string `json:"format"`
				JoinConfiguration struct {
					Discovery struct {
					} `json:"discovery"`
					NodeRegistration struct {
						CriSocket        string `json:"criSocket"`
						KubeletExtraArgs struct {
							CloudProvider   string `json:"cloud-provider"`
							NodeLabels      string `json:"node-labels"`
							TLSCipherSuites string `json:"tls-cipher-suites"`
						} `json:"kubeletExtraArgs"`
						Name string `json:"name"`
					} `json:"nodeRegistration"`
				} `json:"joinConfiguration"`
				PostKubeadmCommands      []string `json:"postKubeadmCommands"`
				PreKubeadmCommands       []string `json:"preKubeadmCommands"`
				UseExperimentalRetryJoin bool     `json:"useExperimentalRetryJoin"`
				Users                    []struct {
					Name              string   `json:"name"`
					SSHAuthorizedKeys []string `json:"sshAuthorizedKeys"`
					Sudo              string   `json:"sudo"`
				} `json:"users"`
			} `json:"spec"`
		} `json:"template"`
	} `json:"spec"`
}

var MachineDeployment struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Annotations struct {
			MachinedeploymentClustersXK8SIoRevision string `json:"machinedeployment.clusters.x-k8s.io/revision"`
		} `json:"annotations"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
		Finalizers        []string  `json:"finalizers"`
		Generation        int       `json:"generation"`
		Labels            struct {
			ClusterXK8SIoClusterName            string `json:"cluster.x-k8s.io/cluster-name"`
			TopologyClusterXK8SIoDeploymentName string `json:"topology.cluster.x-k8s.io/deployment-name"`
			TopologyClusterXK8SIoOwned          string `json:"topology.cluster.x-k8s.io/owned"`
		} `json:"labels"`
		Name            string `json:"name"`
		Namespace       string `json:"namespace"`
		OwnerReferences []struct {
			APIVersion string `json:"apiVersion"`
			Kind       string `json:"kind"`
			Name       string `json:"name"`
			UID        string `json:"uid"`
		} `json:"ownerReferences"`
		ResourceVersion string `json:"resourceVersion"`
		UID             string `json:"uid"`
	} `json:"metadata"`
	Spec struct {
		ClusterName             string `json:"clusterName"`
		MinReadySeconds         int    `json:"minReadySeconds"`
		ProgressDeadlineSeconds int    `json:"progressDeadlineSeconds"`
		Replicas                int    `json:"replicas"`
		RevisionHistoryLimit    int    `json:"revisionHistoryLimit"`
		Selector                struct {
			MatchLabels struct {
				ClusterXK8SIoClusterName            string `json:"cluster.x-k8s.io/cluster-name"`
				TopologyClusterXK8SIoDeploymentName string `json:"topology.cluster.x-k8s.io/deployment-name"`
				TopologyClusterXK8SIoOwned          string `json:"topology.cluster.x-k8s.io/owned"`
			} `json:"matchLabels"`
		} `json:"selector"`
		Strategy struct {
			RollingUpdate struct {
				MaxSurge       int `json:"maxSurge"`
				MaxUnavailable int `json:"maxUnavailable"`
			} `json:"rollingUpdate"`
			Type string `json:"type"`
		} `json:"strategy"`
		Template struct {
			Metadata struct {
				Annotations struct {
					Date                            string `json:"date"`
					RunTanzuVmwareComResolveOsImage string `json:"run.tanzu.vmware.com/resolve-os-image"`
				} `json:"annotations"`
				Labels struct {
					ClusterXK8SIoClusterName            string `json:"cluster.x-k8s.io/cluster-name"`
					TopologyClusterXK8SIoDeploymentName string `json:"topology.cluster.x-k8s.io/deployment-name"`
					TopologyClusterXK8SIoOwned          string `json:"topology.cluster.x-k8s.io/owned"`
				} `json:"labels"`
			} `json:"metadata"`
			Spec struct {
				Bootstrap struct {
					ConfigRef struct {
						APIVersion string `json:"apiVersion"`
						Kind       string `json:"kind"`
						Name       string `json:"name"`
						Namespace  string `json:"namespace"`
					} `json:"configRef"`
				} `json:"bootstrap"`
				ClusterName       string `json:"clusterName"`
				InfrastructureRef struct {
					APIVersion string `json:"apiVersion"`
					Kind       string `json:"kind"`
					Name       string `json:"name"`
					Namespace  string `json:"namespace"`
				} `json:"infrastructureRef"`
				Version string `json:"version"`
			} `json:"spec"`
		} `json:"template"`
	} `json:"spec"`
	Status struct {
		AvailableReplicas int `json:"availableReplicas"`
		Conditions        []struct {
			LastTransitionTime time.Time `json:"lastTransitionTime"`
			Status             string    `json:"status"`
			Type               string    `json:"type"`
		} `json:"conditions"`
		ObservedGeneration  int    `json:"observedGeneration"`
		Phase               string `json:"phase"`
		ReadyReplicas       int    `json:"readyReplicas"`
		Replicas            int    `json:"replicas"`
		Selector            string `json:"selector"`
		UnavailableReplicas int    `json:"unavailableReplicas"`
		UpdatedReplicas     int    `json:"updatedReplicas"`
	} `json:"status"`
}
