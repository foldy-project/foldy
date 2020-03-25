package installer

import (
	"os"

	"github.com/spf13/viper"
)

func init() {
	foldyRepoURL, ok := os.LookupEnv("FOLDY_GIT")
	if !ok {
		foldyRepoURL = "https://github.com/foldy-project/foldy.git"
	}

	AddComponent(&ApplicationComponent{
		Name:         "foldy",
		RepoURL:      foldyRepoURL,
		Path:         "charts/apps",
		Dependencies: []string{},
		CRDs: []string{
			// foldy
			"backends.app.foldy.dev",
			"datasets.app.foldy.dev",
			"experiments.app.foldy.dev",
			"models.app.foldy.dev",
			"transforms.app.foldy.dev",

			// argo
			"cronworkflows.argoproj.io",
			"workflows.argoproj.io",
			"workflowtemplates.argoproj.io",

			// argo-events
			"eventsources.argoproj.io",
			"gateways.argoproj.io",
			"sensors.argoproj.io",

			// traefik
			"ingressroutes.traefik.containo.us",
			"ingressroutetcps.traefik.containo.us",
			"middlewares.traefik.containo.us",
			"tlsoptions.traefik.containo.us",
			"traefikservices.traefik.containo.us",
		},
		ExtraRepos: []*Repository{{
			Name: "argo",
			URL:  "https://github.com/argoproj/argo-helm.git",
		}, {
			Name: "traefik",
			URL:  "https://github.com/containous/traefik-helm-chart.git",
		}, {
			Name: "jetstack",
			URL:  "https://charts.jetstack.io",
		}},
		PreInstall: func(s *Installer) error {
			ingressEnabled, _ := viper.Get("ingress.enabled").(bool)
			if ingressEnabled {
				if err := s.createNamespace("traefik"); err != nil {
					return err
				}
			}
			if err := s.createNamespace("argo"); err != nil {
				return err
			}
			if err := s.createNamespace("argo-events"); err != nil {
				return err
			}
			/*
				crds := []string{
					// cert-manager does not ship its CRDs with its chart

					// traefik uses the Helm v3+ crd/ directory which is not yet supported by argocd
					"https://raw.githubusercontent.com/foldy-project/foldy/master/charts/traefik/crds/ingressroute.yaml",
					"https://raw.githubusercontent.com/foldy-project/foldy/master/charts/traefik/crds/ingressroutetcp.yaml",
					"https://raw.githubusercontent.com/foldy-project/foldy/master/charts/traefik/crds/middlewares.yaml",
					"https://raw.githubusercontent.com/foldy-project/foldy/master/charts/traefik/crds/tlsoptions.yaml",
					"https://raw.githubusercontent.com/foldy-project/foldy/master/charts/traefik/crds/traefikservices.yaml",

					// the foldy monoapp itself utilizes some its dependencies' CRDs
					"https://raw.githubusercontent.com/argoproj/argo-helm/master/charts/argo-events/crds/event-source-crd.yml",
					"https://raw.githubusercontent.com/argoproj/argo-helm/master/charts/argo-events/crds/gateway-crd.yml",
					"https://raw.githubusercontent.com/argoproj/argo-helm/master/charts/argo-events/crds/sensor-crd.yml",
					"https://raw.githubusercontent.com/argoproj/argo-helm/master/charts/argo/crds/workflow-template-crd.yaml",
					"https://raw.githubusercontent.com/argoproj/argo-helm/master/charts/argo/crds/cron-workflow-crd.yaml",
					"https://raw.githubusercontent.com/argoproj/argo-helm/master/charts/argo/crds/workflow-crd.yaml",
				}
				dones := make([]chan error, len(crds), len(crds))
				for i, crd := range crds {
					done := make(chan error, 1)
					dones[i] = done
					go func(crd string, done chan<- error) {
						done <- s.exec("kubectl apply --validate=false -f %s", crd)
						close(done)
					}(crd, done)
				}*/
			return nil
		},
		PostUninstall: func(s *Installer) error {
			if err := s.AsyncDelete("namespace", []string{
				"traefik",
				"argo",
				"argo-events",
			}); err != nil {
				return err
			}
			return nil
		},
	})
}
