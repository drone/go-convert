Level of conversion support for GitLab [job keywords](https://docs.gitlab.com/ee/ci/yaml/#job-keywords) to [drone/spec](https://github.com/drone/spec).

| | Support level |
|-|-----------|
| 🟩 | Full |
| 🟨 | Partial |
| 🟥 | Unsupported |

## 🟥 [`after_script`](https://docs.gitlab.com/ee/ci/yaml/#after_script)

Issue [170](https://github.com/drone/go-convert/issues/170)

## 🟩 [`allow_failure`](https://docs.gitlab.com/ee/ci/yaml/#allow_failure)

<details>
  <summary>Example</summary>

Source
```yaml
job1:
  stage: test
  script:
    - execute_script_1

job2:
  stage: test
  script:
    - execute_script_2
  allow_failure: true

job3:
  stage: deploy
  script:
    - deploy_to_staging
  environment: staging
```

Converted
```yaml
stages:
- name: test
  spec:
    steps:
    - spec:
        steps:
        - name: job1
          spec:
            run: execute_script_1
          type: script
        - name: job2
          "on":
            failure:
              errors:
              - all
              type: ignore
          spec:
            run: execute_script_2
          type: script
      type: parallel
    - name: job3
      spec:
        run: deploy_to_staging
      type: script
  type: ci
version: 1
```

</details>

### 🟩 [`allow_failure:exit_codes`](https://docs.gitlab.com/ee/ci/yaml/#allow_failureexit_codes)

<details>
  <summary>Example</summary>

Source
```yaml
test_job_1:
  script:
    - echo "Run a script that results in exit code 1. This job fails."
    - exit 1
  allow_failure:
    exit_codes: 137

test_job_2:
  script:
    - echo "Run a script that results in exit code 137. This job is allowed to fail."
    - exit 137
  allow_failure:
    exit_codes:
      - 137
      - 255
```

Converted
```yaml
stages:
- name: test
  spec:
    steps:
    - spec:
        steps:
        - name: test_job_1
          "on":
            failure:
              errors:
              - all
              exit_codes:
              - "137"
              type: ignore
          spec:
            run: |-
              echo "Run a script that results in exit code 1. This job fails."
              exit 1
          type: script
        - name: test_job_2
          "on":
            failure:
              errors:
              - all
              exit_codes:
              - "137"
              - "255"
              type: ignore
          spec:
            run: |-
              echo "Run a script that results in exit code 137. This job is allowed to fail."
              exit 137
          type: script
      type: parallel
  type: ci
version: 1
```

</details>

## [`artifacts`](https://docs.gitlab.com/ee/ci/yaml/#artifacts)

## [`before_script`](https://docs.gitlab.com/ee/ci/yaml/#before_script)

## [`cache`](https://docs.gitlab.com/ee/ci/yaml/#cache)

### [`cache:paths`](https://docs.gitlab.com/ee/ci/yaml/#cachepaths)

### [`cache:key`](https://docs.gitlab.com/ee/ci/yaml/#cachekey)

### [`cache:untracked`](https://docs.gitlab.com/ee/ci/yaml/#cacheuntracked)

### [`cache:unprotect`](https://docs.gitlab.com/ee/ci/yaml/#cacheunprotect)

### [`cache:when`](https://docs.gitlab.com/ee/ci/yaml/#cachewhen)

### [`cache:policy`](https://docs.gitlab.com/ee/ci/yaml/#cachepolicy)

### [`cache:fallback_keys`](https://docs.gitlab.com/ee/ci/yaml/#cachefallback_keys)

## [`coverage`](https://docs.gitlab.com/ee/ci/yaml/#coverage)

## [`dast_configuration`](https://docs.gitlab.com/ee/ci/yaml/#dast_configuration)

## [`dependencies`](https://docs.gitlab.com/ee/ci/yaml/#dependencies)

## [`environment`](https://docs.gitlab.com/ee/ci/yaml/#environment)

### [`environment:name`](https://docs.gitlab.com/ee/ci/yaml/#environmentname)

### [`environment:url`](https://docs.gitlab.com/ee/ci/yaml/#environmenturl)

### [`environment:on_stop`](https://docs.gitlab.com/ee/ci/yaml/#environmenton_stop)

### [`environment:action`](https://docs.gitlab.com/ee/ci/yaml/#environmentaction)

### [`environment:auto_stop_in`](https://docs.gitlab.com/ee/ci/yaml/#environmentauto_stop_in)

### [`environment:kubernetes`](https://docs.gitlab.com/ee/ci/yaml/#environmentkubernetes)

### [`environment:deployment_tier`](https://docs.gitlab.com/ee/ci/yaml/#environmentdeployment_tier)

## [`extends`](https://docs.gitlab.com/ee/ci/yaml/#extends)

## [`hooks`](https://docs.gitlab.com/ee/ci/yaml/#hooks)

### [`hooks:pre_get_sources_script`](https://docs.gitlab.com/ee/ci/yaml/#hookspre_get_sources_script)

## [`id_tokens`](https://docs.gitlab.com/ee/ci/yaml/#id_tokens)

## [`image`](https://docs.gitlab.com/ee/ci/yaml/#image)

### [`image:name`](https://docs.gitlab.com/ee/ci/yaml/#imagename)

### [`image:entrypoint`](https://docs.gitlab.com/ee/ci/yaml/#imageentrypoint)

### [`image:pull_policy`](https://docs.gitlab.com/ee/ci/yaml/#imagepull_policy)

## [`inherit`](https://docs.gitlab.com/ee/ci/yaml/#inherit)

### [`inherit:default`](https://docs.gitlab.com/ee/ci/yaml/#inheritdefault)

### [`inherit:variables`](https://docs.gitlab.com/ee/ci/yaml/#inheritvariables)

## [`interruptible`](https://docs.gitlab.com/ee/ci/yaml/#interruptible)

## [`needs`](https://docs.gitlab.com/ee/ci/yaml/#needs)

### [`needs:artifacts`](https://docs.gitlab.com/ee/ci/yaml/#needsartifacts)

### [`needs:project`](https://docs.gitlab.com/ee/ci/yaml/#needsproject)

#### [`needs:pipeline:job`](https://docs.gitlab.com/ee/ci/yaml/#needspipelinejob)

### [`needs:optional`](https://docs.gitlab.com/ee/ci/yaml/#needsoptional)

### [`needs:pipeline`](https://docs.gitlab.com/ee/ci/yaml/#needspipeline)

#### [`needs:parallel:matrix`](https://docs.gitlab.com/ee/ci/yaml/#needsparallelmatrix)

## [`only / except`](https://docs.gitlab.com/ee/ci/yaml/#only--except)

### [`only:refs / except:refs`](https://docs.gitlab.com/ee/ci/yaml/#onlyrefs--exceptrefs)

### [`only:variables / except:variables`](https://docs.gitlab.com/ee/ci/yaml/#onlyvariables--exceptvariables)

### [`only:changes / except:changes`](https://docs.gitlab.com/ee/ci/yaml/#onlychanges--exceptchanges)

### [`only:kubernetes / except:kubernetes`](https://docs.gitlab.com/ee/ci/yaml/#onlykubernetes--exceptkubernetes)

## [`pages`](https://docs.gitlab.com/ee/ci/yaml/#pages)

### [`pages:publish`](https://docs.gitlab.com/ee/ci/yaml/#pagespublish)

## [`parallel`](https://docs.gitlab.com/ee/ci/yaml/#parallel)

### [`parallel:matrix`](https://docs.gitlab.com/ee/ci/yaml/#parallelmatrix)

## [`release`](https://docs.gitlab.com/ee/ci/yaml/#release)

### [`release:tag_name`](https://docs.gitlab.com/ee/ci/yaml/#releasetag_name)

### [`release:tag_message`](https://docs.gitlab.com/ee/ci/yaml/#releasetag_message)

### [`release:name`](https://docs.gitlab.com/ee/ci/yaml/#releasename)

### [`release:description`](https://docs.gitlab.com/ee/ci/yaml/#releasedescription)

### [`release:ref`](https://docs.gitlab.com/ee/ci/yaml/#releaseref)

### [`release:milestones`](https://docs.gitlab.com/ee/ci/yaml/#releasemilestones)

### [`release:released_at`](https://docs.gitlab.com/ee/ci/yaml/#releasereleased_at)

### [`release:assets:links`](https://docs.gitlab.com/ee/ci/yaml/#releaseassetslinks)

## [`resource_group`](https://docs.gitlab.com/ee/ci/yaml/#resource_group)

## [`retry`](https://docs.gitlab.com/ee/ci/yaml/#retry)

### [`retry:when`](https://docs.gitlab.com/ee/ci/yaml/#retrywhen)

## [`rules`](https://docs.gitlab.com/ee/ci/yaml/#rules)

### [`rules:if`](https://docs.gitlab.com/ee/ci/yaml/#rulesif)

### [`rules:changes`](https://docs.gitlab.com/ee/ci/yaml/#ruleschanges)

#### [`rules:changes:paths`](https://docs.gitlab.com/ee/ci/yaml/#ruleschangespaths)

#### [`rules:changes:compare_to`](https://docs.gitlab.com/ee/ci/yaml/#ruleschangescompare_to)

### [`rules:exists`](https://docs.gitlab.com/ee/ci/yaml/#rulesexists)

### [`rules:allow_failure`](https://docs.gitlab.com/ee/ci/yaml/#rulesallow_failure)

### [`rules:needs`](https://docs.gitlab.com/ee/ci/yaml/#rulesneeds)

### [`rules:variables`](https://docs.gitlab.com/ee/ci/yaml/#rulesvariables)

## [`script`](https://docs.gitlab.com/ee/ci/yaml/#script)

## [`secrets`](https://docs.gitlab.com/ee/ci/yaml/#secrets)

### [`secrets:vault`](https://docs.gitlab.com/ee/ci/yaml/#secretsvault)

### [`secrets:azure_key_vault`](https://docs.gitlab.com/ee/ci/yaml/#secretsazure_key_vault)

### [`secrets:file`](https://docs.gitlab.com/ee/ci/yaml/#secretsfile)

### [`secrets:token`](https://docs.gitlab.com/ee/ci/yaml/#secretstoken)

## [`services`](https://docs.gitlab.com/ee/ci/yaml/#services)

### [`service:pull_policy`](https://docs.gitlab.com/ee/ci/yaml/#servicepull_policy)

## [`stage`](https://docs.gitlab.com/ee/ci/yaml/#stage)

### [`stage: .pre`](https://docs.gitlab.com/ee/ci/yaml/#stage-pre)

### [`stage: .post`](https://docs.gitlab.com/ee/ci/yaml/#stage-post)

## [`tags`](https://docs.gitlab.com/ee/ci/yaml/#tags)

## [`timeout`](https://docs.gitlab.com/ee/ci/yaml/#timeout)

## [`trigger`](https://docs.gitlab.com/ee/ci/yaml/#trigger)

### [`trigger:include`](https://docs.gitlab.com/ee/ci/yaml/#triggerinclude)

### [`trigger:project`](https://docs.gitlab.com/ee/ci/yaml/#triggerproject)

### [`trigger:strategy`](https://docs.gitlab.com/ee/ci/yaml/#triggerstrategy)

### [`trigger:forward`](https://docs.gitlab.com/ee/ci/yaml/#triggerforward)

## [`variables`](https://docs.gitlab.com/ee/ci/yaml/#variables)

### [`variables:description`](https://docs.gitlab.com/ee/ci/yaml/#variablesdescription)

### [`variables:value`](https://docs.gitlab.com/ee/ci/yaml/#variablesvalue)

### [`variables:options`](https://docs.gitlab.com/ee/ci/yaml/#variablesoptions)

### [`variables:expand`](https://docs.gitlab.com/ee/ci/yaml/#variablesexpand)

## [`when`](https://docs.gitlab.com/ee/ci/yaml/#when)