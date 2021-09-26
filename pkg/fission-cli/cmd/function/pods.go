/*
Copyright 2019 The Fission Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package function

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fission/fission/pkg/fission-cli/cliwrapper/cli"
	"github.com/fission/fission/pkg/fission-cli/cmd"
	flagkey "github.com/fission/fission/pkg/fission-cli/flag/key"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ListPodsSubCommand struct {
	cmd.CommandActioner
}

func ListPods(input cli.Input) error {
	return (&ListPodsSubCommand{}).do(input)
}

func (opts *ListPodsSubCommand) do(input cli.Input) error {

	m := &metav1.ObjectMeta{
		Name: input.String(flagkey.FnName),
		Labels: map[string]string{
			"executorType": input.String(flagkey.FnExecutorTypeWithoutDefault),
		},
	}

	pods, err := opts.Client().V1().Function().ListPods(m)
	if err != nil {
		return errors.Wrap(err, "error listing environments")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t\n", "NAME", "NAMESPACE", "STATUS", "EXECUTORTYPE", "MANAGED")
	for _, pod := range pods {
		labelList := pod.GetLabels()
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t\n", pod.ObjectMeta.Name, pod.ObjectMeta.Namespace, pod.Status.Phase, labelList["executorType"], labelList["managed"])
	}
	w.Flush()

	return nil
}
