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

This package provides command line tools for local development and debugging purposes. These command line tools are intentially simple. For more robust command line tooling please see the [harness-convert](https://github.com/harness/harness-convert) project.

Compile the binary:

```
git clone https://github.com/drone/go-convert.git
cd go-convert
go build
```

Convert a Bitbucket pipeline:

```
./go-convert gitlab path/to/bitbucket-pipelines.yml
```

Convert a Bitbucket pipeline and downgrade to the Harness v0 format:

```
./go-convert gitlab --downgrade path/to/bitbucket-pipelines.yml
```

Convert a Gitlab pipeline:

```
./go-convert gitlab path/to/.gitlab.yml
```

