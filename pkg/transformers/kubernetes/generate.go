package kubernetes

import (
	"fmt"
	"net"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/golang/glog"

	"github.com/redhat-developer/henge/pkg/types"

	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/resource"
	"k8s.io/kubernetes/pkg/runtime"
	utilerrs "k8s.io/kubernetes/pkg/util/errors"
	"k8s.io/kubernetes/pkg/util/sets"

	"github.com/openshift/origin/pkg/generate/app"
	"github.com/openshift/origin/third_party/github.com/docker/libcompose/project"

	"github.com/redhat-developer/henge/pkg/utils"
)

// Generate accepts a set of Docker compose project paths and converts them in an
// Kubernetes List.
func Generate(vals *types.CmdValues) (*kapi.List, error) {
	var paths []string
	for _, file := range vals.Files {
		path, err := filepath.Abs(file)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	var bases []string
	for _, s := range paths {
		bases = append(bases, filepath.Dir(s))
	}

	context := &project.Context{
		ComposeFiles: paths,
	}
	p := project.NewProject(context)
	if err := p.Parse(); err != nil {
		return nil, err
	}
	if vals.Interactive {
		utils.AskForData(p.Configs)
	}
	list := &kapi.List{}

	serviceOrder := sets.NewString()
	warnings := make(map[string][]string)
	for k, v := range p.Configs {
		serviceOrder.Insert(k)
		warnUnusableComposeElements(k, v, warnings)
	}

	var errs []error

	// identify colocated components due to shared volumes
	joins := make(map[string]sets.String)
	volumesFrom := make(map[string][]string)
	for _, k := range serviceOrder.List() {
		if joins[k] == nil {
			joins[k] = sets.NewString(k)
		}
		v := p.Configs[k]
		if len(v.Build) != 0 {
			return nil, fmt.Errorf("Build is not currently supported for Kubernetes")
		}
		for _, from := range v.VolumesFrom {
			switch parts := strings.Split(from, ":"); len(parts) {
			case 1:
				joins[k].Insert(parts[0])
				volumesFrom[k] = append(volumesFrom[k], parts[0])
			case 2:
				target := parts[1]
				if parts[1] == "ro" || parts[1] == "rw" {
					target = parts[0]
				}
				joins[k].Insert(target)
				volumesFrom[k] = append(volumesFrom[k], target)
			case 3:
				joins[k].Insert(parts[1])
				volumesFrom[k] = append(volumesFrom[k], parts[1])
			}
		}
	}
	joinOrder := sets.NewString()
	for k := range joins {
		joinOrder.Insert(k)
	}
	var colocated []sets.String
	for _, k := range joinOrder.List() {
		set := joins[k]
		matched := -1
		for i, existing := range colocated {
			if set.Intersection(existing).Len() == 0 {
				continue
			}
			if matched != -1 {
				return nil, fmt.Errorf("%q belongs with %v, but %v also contains some overlapping elements", k, set, colocated[matched])
			}
			existing.Insert(set.List()...)
			matched = i
			continue
		}
		if matched == -1 {
			colocated = append(colocated, set)
		}
	}

	// identify service aliases
	aliases := make(map[string]sets.String)
	for _, v := range p.Configs {
		for _, s := range v.Links.Slice() {
			parts := strings.SplitN(s, ":", 2)
			if len(parts) != 2 || parts[0] == parts[1] {
				continue
			}
			set := aliases[parts[0]]
			if set == nil {
				set = sets.NewString()
				aliases[parts[0]] = set
			}
			set.Insert(parts[1])
		}
	}

	if len(errs) > 0 {
		return nil, utilerrs.NewAggregate(errs)
	}

	objects := []runtime.Object{}
	// create deployment groups
	for _, pod := range colocated {
		commonMounts := make(map[string]string)
		name := strings.Join(pod.List(), "_")

		allWarnings := sets.NewString()

		rc := kapi.ReplicationController{}
		rc.ObjectMeta.Name = name
		rc.Spec.Selector = make(map[string]string)
		rc.Spec.Selector["replicationcontroller"] = name

		rc.Spec.Template = &kapi.PodTemplateSpec{}
		rc.Spec.Template.ObjectMeta.Labels = make(map[string]string)
		rc.Spec.Template.ObjectMeta.Labels["replicationcontroller"] = name

		rc.ObjectMeta.Annotations = make(map[string]string)

		// TODO: is there number if replicas in docker-compose?
		rc.Spec.Replicas = 1

		for _, k := range pod.List() {
			// collect all warnings from all docker-compose servces for this rc
			for msg, services := range warnings {
				for _, service := range services {
					if service == k {
						allWarnings.Insert(msg)
					}
				}
			}
			v := p.Configs[k]
			glog.V(4).Infof("compose service: %#v", v)

			c := kapi.Container{}
			if len(v.Image) > 0 {
				c.Image = v.Image
			}

			if len(v.ContainerName) > 0 {
				c.Name = v.ContainerName
			} else {
				c.Name = k
			}
			for _, s := range v.Ports {
				container, _ := extractFirstPorts(s)
				if port, err := strconv.Atoi(container); err == nil {
					c.Ports = append(c.Ports, kapi.ContainerPort{ContainerPort: port})
				}
			}
			c.Args = v.Command.Slice()
			if len(v.Entrypoint.Slice()) > 0 {
				c.Command = v.Entrypoint.Slice()
			}
			if len(v.WorkingDir) > 0 {
				c.WorkingDir = v.WorkingDir
			}
			c.Env = append(c.Env, app.ParseEnvironment(v.Environment.Slice()...).List()...)
			if uid, err := strconv.Atoi(v.User); err == nil {
				uid64 := int64(uid)
				if c.SecurityContext == nil {
					c.SecurityContext = &kapi.SecurityContext{}
				}
				c.SecurityContext.RunAsUser = &uid64
			}
			c.TTY = v.Tty
			if v.StdinOpen {
				c.StdinOnce = true
				c.Stdin = true
			}
			if v.Privileged {
				if c.SecurityContext == nil {
					c.SecurityContext = &kapi.SecurityContext{}
				}
				t := true
				c.SecurityContext.Privileged = &t
			}
			if v.ReadOnly {
				if c.SecurityContext == nil {
					c.SecurityContext = &kapi.SecurityContext{}
				}
				t := true
				c.SecurityContext.ReadOnlyRootFilesystem = &t
			}
			if v.MemLimit > 0 {
				q := resource.NewQuantity(v.MemLimit, resource.DecimalSI)
				if c.Resources.Limits == nil {
					c.Resources.Limits = make(kapi.ResourceList)
				}
				c.Resources.Limits[kapi.ResourceMemory] = *q
			}

			if quota := v.CPUQuota; quota > 0 {
				if quota < 1000 {
					quota = 1000 // minQuotaPeriod
				}
				milliCPU := quota * 1000     // milliCPUtoCPU
				milliCPU = milliCPU / 100000 // quotaPeriod
				q := resource.NewMilliQuantity(milliCPU, resource.DecimalSI)
				if c.Resources.Limits == nil {
					c.Resources.Limits = make(kapi.ResourceList)
				}
				c.Resources.Limits[kapi.ResourceCPU] = *q
			}
			if shares := v.CPUShares; shares > 0 {
				if shares < 2 {
					shares = 2 // minShares
				}
				milliCPU := shares * 1000  // milliCPUtoCPU
				milliCPU = milliCPU / 1024 // sharesPerCPU
				q := resource.NewMilliQuantity(milliCPU, resource.DecimalSI)
				if c.Resources.Requests == nil {
					c.Resources.Requests = make(kapi.ResourceList)
				}
				c.Resources.Requests[kapi.ResourceCPU] = *q
			}

			mountPoints := make(map[string][]string)
			for _, s := range v.Volumes {
				switch parts := strings.SplitN(s, ":", 3); len(parts) {
				case 1:
					mountPoints[""] = append(mountPoints[""], parts[0])

				case 2:
					fallthrough
				default:
					mountPoints[parts[0]] = append(mountPoints[parts[0]], parts[1])
				}
			}
			for from, at := range mountPoints {
				name, ok := commonMounts[from]
				if !ok {
					name = fmt.Sprintf("dir-%d", len(commonMounts)+1)
					commonMounts[from] = name
				}
				for _, path := range at {
					c.VolumeMounts = append(c.VolumeMounts, kapi.VolumeMount{Name: name, MountPath: path})
				}
			}

			rc.Spec.Template.Spec.Containers = append(rc.Spec.Template.Spec.Containers, c)
			if len(allWarnings.List()) > 0 {
				rc.ObjectMeta.Annotations["app.henge/warnings"] = fmt.Sprintf("not all docker-compose fields were honored:\n* %s", strings.Join(allWarnings.List(), "\n* "))
			}

		}

		objects = append(objects, &rc)
	}

	if len(errs) > 0 {
		return nil, utilerrs.NewAggregate(errs)
	}

	// create services for each object with a name based on alias.
	containers := make(map[string]*kapi.Container)
	var services []*kapi.Service
	for _, obj := range objects {
		switch t := obj.(type) {
		case *kapi.ReplicationController:
			ports := app.UniqueContainerToServicePorts(app.AllContainerPorts(t.Spec.Template.Spec.Containers...))
			if len(ports) == 0 {
				continue
			}
			svc := app.GenerateService(t.ObjectMeta, t.Spec.Selector)
			if aliases[svc.Name].Len() == 1 {
				svc.Name = aliases[svc.Name].List()[0]
			}
			svc.Spec.Ports = ports
			services = append(services, svc)

			// take a reference to each container
			for i := range t.Spec.Template.Spec.Containers {
				c := &t.Spec.Template.Spec.Containers[i]
				containers[c.Name] = c
			}
		}
	}
	for _, svc := range services {
		objects = append(objects, svc)
	}

	// for each container that defines VolumesFrom, copy equivalent mounts.
	// TODO: ensure mount names are unique?
	for target, otherContainers := range volumesFrom {
		for _, from := range otherContainers {
			for _, volume := range containers[from].VolumeMounts {
				containers[target].VolumeMounts = append(containers[target].VolumeMounts, volume)
			}
		}
	}

	// add emptyDir volume for every volume defined in container
	for _, obj := range objects {
		switch o := obj.(type) {
		case *kapi.ReplicationController:
			for _, container := range o.Spec.Template.Spec.Containers {
				for _, volumeMount := range container.VolumeMounts {

					o.Spec.Template.Spec.Volumes = append(o.Spec.Template.Spec.Volumes,
						kapi.Volume{
							Name:         volumeMount.Name,
							VolumeSource: kapi.VolumeSource{EmptyDir: &kapi.EmptyDirVolumeSource{}},
						})
				}
			}
		}
	}

	list.Items = objects

	return list, nil
}

// extractFirstPorts converts a Docker compose port spec (CONTAINER, HOST:CONTAINER, or
// IP:HOST:CONTAINER) to the first container and host port in the range.  Host port will
// default to container port.
func extractFirstPorts(port string) (container, host string) {
	segments := strings.Split(port, ":")
	container = segments[len(segments)-1]
	container = rangeToPort(container)
	switch {
	case len(segments) == 3:
		host = rangeToPort(segments[1])
	case len(segments) == 2 && net.ParseIP(segments[0]) == nil:
		host = rangeToPort(segments[0])
	default:
		host = container
	}
	return container, host
}

func rangeToPort(s string) string {
	parts := strings.SplitN(s, "-", 2)
	return parts[0]
}

// warnUnusableComposeElements add warnings for unsupported elements in the provided service config
func warnUnusableComposeElements(k string, v *project.ServiceConfig, warnings map[string][]string) {
	fn := func(msg string) {
		warnings[msg] = append(warnings[msg], k)
	}
	if len(v.CapAdd) > 0 || len(v.CapDrop) > 0 {
		// TODO: we can support this
		fn("cap_add and cap_drop are not supported")
	}
	if len(v.CgroupParent) > 0 {
		fn("cgroup_parent is not supported")
	}
	if len(v.CPUSet) > 0 {
		fn("cpuset is not supported")
	}
	if len(v.Devices) > 0 {
		fn("devices are not supported")
	}
	if v.DNS.Len() > 0 || v.DNSSearch.Len() > 0 {
		fn("dns and dns_search are not supported")
	}
	if len(v.DomainName) > 0 {
		fn("domainname is not supported")
	}
	if len(v.Hostname) > 0 {
		fn("hostname is not supported")
	}
	if len(v.Labels.MapParts()) > 0 {
		fn("labels is ignored")
	}
	if len(v.Links.Slice()) > 0 {
		//fn("links are not supported, use services to talk to other pods")
		// TODO: display some sort of warning when linking will be inconsistent
	}
	if len(v.LogDriver) > 0 {
		fn("log_driver is not supported")
	}
	if len(v.MacAddress) > 0 {
		fn("mac_address is not supported")
	}
	if len(v.Net) > 0 {
		fn("net is not supported")
	}
	if len(v.Pid) > 0 {
		fn("pid is not supported")
	}
	if len(v.Uts) > 0 {
		fn("uts is not supported")
	}
	if len(v.Ipc) > 0 {
		fn("ipc is not supported")
	}
	if v.MemSwapLimit > 0 {
		fn("mem_swap_limit is not supported")
	}
	if len(v.Restart) > 0 {
		fn("restart is ignored - all pods are automatically restarted")
	}
	if len(v.SecurityOpt) > 0 {
		fn("security_opt is not supported")
	}
	if len(v.User) > 0 {
		if _, err := strconv.Atoi(v.User); err != nil {
			fn("setting user to a string is not supported - use numeric user value")
		}
	}
	if len(v.VolumeDriver) > 0 {
		fn("volume_driver is not supported")
	}
	if len(v.VolumesFrom) > 0 {
		fn("volumes_from is not supported")
		// TODO: use volumes from for colocated containers to automount volumes
	}
	if len(v.ExternalLinks) > 0 {
		fn("external_links are not supported - use services")
	}
	if len(v.LogOpt) > 0 {
		fn("log_opt is not supported")
	}
	if len(v.ExtraHosts) > 0 {
		fn("extra_hosts is not supported")
	}
	if len(v.Ulimits.Elements) > 0 {
		fn("ulimits is not supported")
	}
	// TODO: fields to handle
	// EnvFile       Stringorslice     `yaml:"env_file,omitempty"`
}
