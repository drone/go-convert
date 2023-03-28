Package convert provides tooling to convert third party pipeline configurations to the Harness pipeline configuration format.


__Sample Usage__

Sample code to convert a Bitbucket pipeline to a Harness pipeline:

```Go
import "github.com/drone/go-convert/convert/bitbucket"
```

```Go
converter := bitbucket.New(
	bitbucket.WithDockerhub(c.dockerConn),
	bitbucket.WithKubernetes(c.kubeConn, c.kubeName),
)
converted, err := converter.ConvertFile("bitbucket-pipelines.yml")
if err != nil {
	log.Fatalln(err)
}
```

__Command Line__

This package provides command line tools for local development and debugging purposes. These command line tools are intentionally simple. For more robust command line tooling please use the [harness-convert](https://github.com/harness/harness-convert) project.

Compile the binary:

```
git clone https://github.com/drone/go-convert.git
cd go-convert
go build
```

__Bitbucket__

Convert a Bitbucket pipeline:

```
./go-convert bitbucket samples/bitbucket.yaml
```

Convert a Gitlab pipeline and print the before after:

```
./go-convert bitbucket --before-after samples/bitbucket.yaml
```

Convert a Bitbucket pipeline and downgrade to the Harness v0 format:

```
./go-convert bitbucket --downgrade samples/bitbucket.yaml
```

__Drone__

Convert a Drone pipeline:

```
./go-convert drone samples/drone.yaml
```

Convert a Drone pipeline and print the before after:

```
./go-convert drone --before-after samples/drone.yaml
```

Convert a Drone pipeline and downgrade to the Harness v0 format:

```
./go-convert drone --downgrade samples/drone.yaml
```

__Gitlab__

Convert a Gitlab pipeline:

```
./go-convert gitlab samples/gitlab.yaml
```

Convert a Gitlab pipeline and print the before after:

```
./go-convert gitlab --before-after samples/gitlab.yaml
```

Convert a Gitlab pipeline and downgrade to the Harness v0 format:

```
./go-convert gitlab --downgrade samples/gitlab.yaml
```

__Jenkins__

Convert a Jenkinsfile:

```
./go-convert jenkins --token=<chat-gpt-token> samples/Jenkinsfile
```

Convert a Jenkinsfile and downgrade to the Harness v0 format:

```
./go-convert jenkins --token=<chat-gpt-token> --downgrade samples/Jenkinsfile
```

__Syntax Highlighting__

The command line tools are compatble with [bat](https://github.com/sharkdp/bat) for syntax highlight.

```
./go-convert bitbucket --before-after samples/bitbucket.yaml | bat -l yaml
```
