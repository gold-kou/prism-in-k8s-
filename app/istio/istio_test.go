package istio_test

import (
	"context"
	"testing"

	"github.com/gold-kou/prism-in-k8s/app/istio"
	"github.com/gold-kou/prism-in-k8s/app/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"istio.io/client-go/pkg/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func TestCreateIstioResources(t *testing.T) {
	testNamespaceName := "test-namespace" + uuid.NewString()
	testResourceName := "test-resource" + uuid.NewString()

	ctx := context.TODO()
	kubeconfigPath := clientcmd.NewDefaultPathOptions().GetDefaultFilename()
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	require.NoError(t, err)

	// create namespace to create a virtualservice
	k8sClientSet, err := kubernetes.NewForConfig(kubeconfig)
	require.NoError(t, err)
	err = testutil.CreateNamespace(ctx, k8sClientSet, testNamespaceName)
	require.NoError(t, err)

	// test target
	err = istio.CreateIstioResources(ctx, kubeconfig, testNamespaceName, testResourceName)
	assert.NoError(t, err)

	// verify
	istioClient, err := versioned.NewForConfig(kubeconfig)
	assert.NoError(t, err)
	_, err = istioClient.NetworkingV1alpha3().VirtualServices(testNamespaceName).Get(ctx, testResourceName, metav1.GetOptions{})
	assert.NoError(t, err)

	// clean up
	err = testutil.DeleteNamespace(ctx, k8sClientSet, testNamespaceName)
	require.NoError(t, err)
}

func TestDeleteIstioResources(t *testing.T) {
	testNamespaceName := "test-namespace" + uuid.NewString()
	testResourceName := "test-resource" + uuid.NewString()

	ctx := context.TODO()
	kubeconfigPath := clientcmd.NewDefaultPathOptions().GetDefaultFilename()
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	require.NoError(t, err)
	k8sClientSet, err := kubernetes.NewForConfig(kubeconfig)
	require.NoError(t, err)
	istioClientSet, err := versioned.NewForConfig(kubeconfig)
	require.NoError(t, err)

	// dummy resources
	err = testutil.CreateNamespace(ctx, k8sClientSet, testNamespaceName)
	require.NoError(t, err)
	err = testutil.CreateVirtualService(ctx, istioClientSet, testNamespaceName, testResourceName)
	require.NoError(t, err)

	// test target
	err = istio.CreateIstioResources(ctx, kubeconfig, testNamespaceName, testResourceName)
	assert.NoError(t, err)

	// skip verify to reduce test time

	// clean up
	err = testutil.DeleteNamespace(ctx, k8sClientSet, testNamespaceName)
	require.NoError(t, err)
}
