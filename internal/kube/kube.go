package kube

import (
	"context"
	goerrors "errors"
	"fmt"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	ModeKubeConfig     = "kubeconfig"
	ModeServiceAccount = "serviceaccount"
)

type KubeClient struct {
	mode       string
	kubeConfig string
	ClientSet  *kubernetes.Clientset
}

func NewKubeClient(mode string, kubeConfig string) (*KubeClient, error) {
	var (
		config    *rest.Config
		clientset *kubernetes.Clientset
		err       error
	)

	switch mode {

	case ModeKubeConfig:
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			return nil, err
		}

		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			return nil, err
		}

		return &KubeClient{mode, kubeConfig, clientset}, nil

	case ModeServiceAccount:
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}

		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			return nil, err
		}

		return &KubeClient{mode, kubeConfig, clientset}, nil

	default:
		return nil, goerrors.New(fmt.Sprintf("Invalid mode '%s'", mode))

	}
}

func NewKubeClientFromCommand(cmd *cobra.Command) (*KubeClient, error) {
	kubeConfigFlag, _ := cmd.PersistentFlags().GetString("kubeconfig")
	serviceAccountFlag, _ := cmd.PersistentFlags().GetBool("service-account")

	if serviceAccountFlag {
		return NewKubeClient(ModeServiceAccount, "")
	} else {
		return NewKubeClient(ModeKubeConfig, kubeConfigFlag)
	}
}

func (k *KubeClient) GetNamespaceNames() ([]string, error) {
	var namespaces = make([]string, 0)

	if namespaceList, err := k.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{}); err == nil {
		for _, namespace := range namespaceList.Items {
			namespaces = append(namespaces, namespace.Name)
		}
	} else {
		return nil, err
	}

	return namespaces, nil
}

func (k *KubeClient) GetDeploymentNamesInNamespace(namespace string) ([]string, error) {
	var deployments, err = k.ClientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var deploymentNames = make([]string, 0)
	for _, deployment := range deployments.Items {
		deploymentNames = append(deploymentNames, deployment.Name)
	}
	return deploymentNames, nil
}

func (k *KubeClient) GetDeploymentsInNamespace(namespace string) (*v1.DeploymentList, error) {
	return k.ClientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
}

func (k *KubeClient) GetDeploymentInNamespace(namespace string, deploymentName string) (*v1.Deployment, bool, error) {
	var deployment, err = k.ClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return deployment, true, nil
}

func (k *KubeClient) GetDeploymentScale(namespace string, deploymentName string) (*autoscalingv1.Scale, bool, error) {
	var scale, err = k.ClientSet.AppsV1().Deployments(namespace).GetScale(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return scale, true, nil
}

func (k *KubeClient) SetDeploymentScale(namespace string, deploymentName string, count int32) (int32, bool, error) {
	var state, exists, err = k.GetDeploymentScale(namespace, deploymentName)
	if err != nil {
		return 0, false, err
	}

	if exists {
		state.Spec.Replicas = count
		newState, err := k.ClientSet.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), deploymentName, state, metav1.UpdateOptions{})
		if err != nil {
			return 0, true, err
		}
		return newState.Spec.Replicas, true, nil
	}
	return 0, false, nil
}

func (k *KubeClient) GetStatefulSetNamesInNamespace(namespace string) ([]string, error) {
	var deployments, err = k.ClientSet.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var StatefulSetNames = make([]string, 0)
	for _, deployment := range deployments.Items {
		StatefulSetNames = append(StatefulSetNames, deployment.Name)
	}
	return StatefulSetNames, nil
}

func (k *KubeClient) GetStatefulSetsInNamespace(namespace string) (*v1.StatefulSetList, error) {
	return k.ClientSet.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
}

func (k *KubeClient) GetStatefulSetInNamespace(namespace string, statefulSetName string) (*v1.StatefulSet, bool, error) {
	var statefulSet, err = k.ClientSet.AppsV1().StatefulSets(namespace).Get(context.TODO(), statefulSetName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return statefulSet, true, nil
}

func (k *KubeClient) GetStatefulSetScale(namespace string, statefulSetName string) (*autoscalingv1.Scale, bool, error) {
	var scale, err = k.ClientSet.AppsV1().StatefulSets(namespace).GetScale(context.TODO(), statefulSetName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return scale, true, nil
}

func (k *KubeClient) SetStatefulSetScale(namespace string, statefulSetName string, count int32) (int32, bool, error) {
	var state, exists, err = k.GetStatefulSetScale(namespace, statefulSetName)
	if err != nil {
		return 0, false, err
	}

	if exists {
		state.Spec.Replicas = count
		newState, err := k.ClientSet.AppsV1().StatefulSets(namespace).UpdateScale(context.TODO(), statefulSetName, state, metav1.UpdateOptions{})
		if err != nil {
			return 0, true, err
		}
		return newState.Spec.Replicas, true, nil
	}
	return 0, false, nil
}
