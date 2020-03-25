package portfwd

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"

	"k8s.io/apimachinery/pkg/util/intstr"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	restclient "k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

type FoldyPortForwarder struct {
	config  *restclient.Config
	exit    map[string]chan<- struct{}
	l       sync.Mutex
	running int32
	Verbose bool
}

func RunPortForward(
	config *restclient.Config,
	serviceName string,
	namespace string,
	localPort int,
	remotePort int,
	stopChan <-chan struct{},
	verbose bool,
) error {
	cl, err := client.New(config, client.Options{})
	if err != nil {
		return err
	}
	service := &corev1.Service{}
	if err := cl.Get(
		context.TODO(),
		types.NamespacedName{
			Name:      serviceName,
			Namespace: namespace,
		},
		service,
	); err != nil {
		return err
	}
	var containerPort *intstr.IntOrString
	for _, svcPort := range service.Spec.Ports {
		if int(svcPort.Port) == remotePort {
			containerPort = &svcPort.TargetPort
		}
	}
	if containerPort == nil {
		return fmt.Errorf("unable to resolve containerPort for %v", remotePort)
	}
	selector := service.Spec.Selector
	pods := &corev1.PodList{}
	if err := cl.List(
		context.TODO(),
		pods,
		client.MatchingLabels(selector),
	); err != nil {
		return err
	}
	actualPort := containerPort.IntVal
	var podName string
	for _, pod := range pods.Items {
		if pod.Status.Phase == "Running" {
			// Found a running pod with the right labels
			podName = pod.ObjectMeta.Name

			// Resolve the containerPort if the service used a name
			if containerPort.Type == intstr.String {
				for _, container := range pod.Spec.Containers {
					for _, contPort := range container.Ports {
						if contPort.Name == containerPort.StrVal {
							actualPort = contPort.ContainerPort
							break
						}
					}
				}
			}
		}
	}
	if podName == "" {
		return fmt.Errorf("unable to resolve pod")
	}
	if actualPort == 0 {
		return fmt.Errorf("unable to resolve port")
	}
	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", namespace, podName)
	hostIP := strings.TrimLeft(config.Host, "htps:/")
	serverURL := url.URL{Scheme: "https", Path: path, Host: hostIP}
	roundTripper, upgrader, err := spdy.RoundTripperFor(config)
	if err != nil {
		return err
	}
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: roundTripper}, http.MethodPost, &serverURL)
	readyChan := make(chan struct{}, 1)
	out, errOut := new(bytes.Buffer), new(bytes.Buffer)
	forwarder, err := portforward.New(
		dialer,
		[]string{fmt.Sprintf("%d:%d", localPort, actualPort)},
		stopChan,
		readyChan,
		out,
		errOut)
	if err != nil {
		return err
	}
	if err := forwarder.ForwardPorts(); err != nil {
		return err
	}
	return nil
}

func NewFoldyPortForwarder(config *restclient.Config) (*FoldyPortForwarder, error) {
	return &FoldyPortForwarder{
		config:  config,
		running: 1,
		exit:    make(map[string]chan<- struct{}),
	}, nil
}

func (p *FoldyPortForwarder) Close() {
	p.l.Lock()
	defer p.l.Unlock()
	if p.exit == nil {
		return
	}
	atomic.StoreInt32(&p.running, 0)
	for _, exit := range p.exit {
		exit <- struct{}{}
	}
	p.exit = nil
}

func (p *FoldyPortForwarder) setStopChan(fullName string, stopChan chan<- struct{}) {
	p.l.Lock()
	if existing, ok := p.exit[fullName]; ok {
		// Don't forget to clean up
		close(existing)
	}
	p.exit[fullName] = stopChan
	p.l.Unlock()
}

func (p *FoldyPortForwarder) AddPort(serviceName string, namespace string, localPort int, remotePort int) {
	go func() {
		fullName := fmt.Sprintf("%s/%s/%d/%d", namespace, serviceName, localPort, remotePort)
		for {
			if atomic.LoadInt32(&p.running) == 0 {
				return
			}
			stopChan := make(chan struct{}, 1)
			p.setStopChan(fullName, stopChan)
			log.Printf("> kubectl port-forward -n %s svc/%s %d:%d", namespace, serviceName, localPort, remotePort)
			if err := RunPortForward(p.config, serviceName, namespace, localPort, remotePort, stopChan, p.Verbose); err != nil {
				log.Printf(err.Error())
			}
		}
	}()
}

func (p *FoldyPortForwarder) AddAllPorts() {
	p.AddPort("argocd-server", "argocd", 8080, 80)
	p.AddPort("foldy-ui", "foldy", 9000, 80)
}
