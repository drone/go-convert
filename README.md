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
./go-convert bitbucket path/to/bitbucket-pipelines.yml
```

Convert a Gitlab pipeline and print the before after:

```
./go-convert bitbucket --before-after path/to/bitbucket-pipelines.yml
```

Convert a Bitbucket pipeline and downgrade to the Harness v0 format:

```
./go-convert bitbucket --downgrade path/to/bitbucket-pipelines.yml
```

__Gitlab__

Convert a Gitlab pipeline:

```
./go-convert gitlab path/to/.gitlab.yml
```

Convert a Gitlab pipeline and print the before after:

```
./go-convert gitlab --before-after path/to/.gitlab.yml
```

Convert a Gitlab pipeline and downgrade to the Harness v0 format:

```
./go-convert gitlab --downgrade path/to/bitbucket-pipelines.yml
```

